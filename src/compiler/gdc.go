/*
 * Copyright (C) 2023 The GDLang Team.
 *
 * This file is part of GDLang.
 *
 * GDLang is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * GDLang is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with GDLang.  If not, see <http://www.gnu.org/licenses/>.
 */

package compiler

import (
	"bytes"
	"encoding/json"
	"gdlang/lib/builtin"
	"gdlang/lib/runtime"
	"gdlang/lib/tools"
	"gdlang/src/cpu"
	"gdlang/src/gd/analysis"
	"gdlang/src/gd/ast"
	"gdlang/src/gd/ir"
	"gdlang/src/gd/staticcheck"
	"os"
)

type (
	GDCmpEvaluator            = staticcheck.Evaluator[ir.GDIRNode, ir.GDIRStackNode]
	GDCompExpressionEvaluator = staticcheck.ExpressionEvaluator[ir.GDIRNode, ir.GDIRStackNode]
)

type ForCallback func(endLabel runtime.GDIdent, stack ir.GDIRStackNode) error

type GDCompiler struct {
	Ctx  *ir.GDIRContext
	Root *ir.GDIRBlock
	tools.GDIdentGen
	GDCmpEvaluator
	GDCompExpressionEvaluator
	*analysis.PackageDependenciesAnalyzer
}

func (c *GDCompiler) Compile(mainPkg string) error {
	stack := runtime.NewGDSymbolStack()
	err := builtin.ImportCoreBuiltins(stack)
	if err != nil {
		return err
	}

	err = c.Analyze(mainPkg, analysis.PackageDependenciesAnalyzerOptions{ShouldLookUpFromMain: true})
	if err != nil {
		return err
	}

	staticAnalyzer := staticcheck.NewStaticCheck(c.PackageDependenciesAnalyzer)
	err = staticAnalyzer.Check(stack)
	if err != nil {
		return err
	}

	stack.Dispose()

	for _, node := range c.Nodes {
		_, err := c.EvalNode(node, c.Root)
		if err != nil {
			return err
		}
	}

	ident := c.DeriveIdent(c.MainEntry)

	mainObj := runtime.NewGDIdObject(ident, runtime.GDZNil)
	irMain := ir.NewGDIRObject(mainObj, nil)

	inst, _ := ir.NewGDIRCall(irMain, ir.NewGDIRIterableObject(runtime.NewGDArrayType(runtime.GDAnyType), []ir.GDIRNode{}), nil)
	c.Root.AddNode(inst)

	return nil
}

func (c *GDCompiler) writeBytecode(outputFile string) error {
	buffer := &bytes.Buffer{}

	err := c.Root.BuildBytecode(buffer, c.Ctx)
	if err != nil {
		return err
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic("error closing file: " + err.Error())
		}
	}(file)

	_, err = file.Write(buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (c *GDCompiler) writeSourceMap(outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic("error closing file: " + err.Error())
		}
	}(file)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "")

	return encoder.Encode(c.Ctx.SrcMap)
}

func (c *GDCompiler) EvalAtom(a *ast.NodeLiteral, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	return ir.NewGDIRObject(a.InferredObject(), a), nil
}

func (c *GDCompiler) EvalIdent(i *ast.NodeIdent, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	ident := c.DeriveIdent(i)

	return ir.NewGDIRIdentObject(ident, i.InferredObject(), i), nil
}

func (c *GDCompiler) EvalLambda(l *ast.NodeLambda, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	lambdaType, ok := c.DeriveType(l).(*runtime.GDLambdaType)
	if !ok {
		panic("A lambda type must be defined")
	}

	lambda, reg := ir.NewGDIRLambda(lambdaType, l)

	block, err := c.evalBlock(l.Block, lambda)
	if err != nil {
		return nil, err
	}

	lambda.GDIRBlock = block.(*ir.GDIRBlock)

	stack.AddNode(lambda)

	return reg, nil
}

func (c *GDCompiler) EvalExprOp(e *ast.NodeExprOperation, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	l, err := c.EvalNode(e.L, stack)
	if err != nil {
		return nil, err
	}

	var r ir.GDIRNode
	if e.R != nil {
		r, err = c.EvalNode(e.R, stack)
		if err != nil {
			return nil, err
		}
	}

	inst, reg := ir.NewGDIROp(e.Op, l, r, e)

	stack.AddNode(inst)

	return reg, nil
}

func (c *GDCompiler) EvalExpEllipsis(e *ast.NodeEllipsisExpr, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	obj, err := c.EvalNode(e.Expr, stack)
	if err != nil {
		return nil, err
	}

	spreadable := ir.NewGDIRObjectWithType(e.InferredType(), obj, e)

	return spreadable, nil
}

func (c *GDCompiler) EvalFunc(f *ast.NodeFunc, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	lambda, err := c.EvalLambda(f.NodeLambda, stack)
	if err != nil {
		return nil, err
	}

	ident := c.DeriveIdent(f)
	typ := c.DeriveType(f)

	disc := ir.NewGDIRDiscoverable(f.IsPub, true, ident, f)
	irFunc := ir.NewGDIRSet(disc, typ, lambda, f)

	stack.AddNode(irFunc)

	return nil, nil
}

func (c *GDCompiler) EvalTuple(nt *ast.NodeTuple, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	return c.collectNodes(nt.InferredType(), nt.Nodes, stack)
}

func (c *GDCompiler) EvalStruct(s *ast.NodeStruct, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	return c.collectNodes(s.InferredType(), s.Nodes, stack)
}

func (c *GDCompiler) EvalArray(a *ast.NodeArray, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	return c.collectNodes(a.InferredType(), a.Nodes, stack)
}

func (c *GDCompiler) EvalReturn(r *ast.NodeReturn, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	if r.Expr == nil {
		inst, reg := ir.NewGDIRRet(ir.NewGDIRObject(runtime.GDZNil, r), r)
		stack.AddNode(inst)

		return reg, nil
	}

	obj, err := c.EvalNode(r.Expr, stack)
	if err != nil {
		return nil, err
	}

	inst, reg := ir.NewGDIRRet(obj, r.Expr)
	stack.AddNode(inst)

	return reg, nil
}

func (c *GDCompiler) EvalIterIdxExpr(a *ast.NodeIterIdxExpr, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	expr, err := c.EvalNode(a.Expr, stack)
	if err != nil {
		return nil, err
	}

	idx, err := c.EvalNode(a.IdxExpr, stack)
	if err != nil {
		return nil, err
	}

	inst, reg := ir.NewGDIRIGet(idx, a.IsNilSafe, expr, a)
	stack.AddNode(inst)

	return reg, nil
}

func (c *GDCompiler) EvalCallExpr(call *ast.NodeCallExpr, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	expr, err := c.EvalNode(call.Expr, stack)
	if err != nil {
		return nil, err
	}

	args := make([]ir.GDIRNode, 0)
	for _, arg := range call.Args {
		arg, err := c.EvalNode(arg, stack)
		if err != nil {
			return nil, err
		}

		args = append(args, arg)
	}

	argsNode := ir.NewGDIRIterableObject(runtime.NewGDArrayType(runtime.GDAnyType), args)
	inst, reg := ir.NewGDIRCall(expr, argsNode, call.Expr)

	stack.AddNode(inst)

	return reg, nil
}

func (c *GDCompiler) EvalSafeDotExpr(s *ast.NodeSafeDotExpr, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	expr, err := c.EvalNode(s.Expr, stack)
	if err != nil {
		return nil, err
	}

	ident := c.DeriveIdent(s)

	inst, reg := ir.NewGDIRAIGet(ident, s.IsNilSafe, expr, s)
	stack.AddNode(inst)

	return reg, nil
}

func (c *GDCompiler) EvalSets(s *ast.NodeSets, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	for _, set := range s.Nodes {
		_, err := c.EvalNode(set, stack)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (c *GDCompiler) resolveNodeSetExpr(set *ast.NodeSet, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	var exprObj ir.GDIRNode
	if set.Expr != nil {
		switch expr := set.Expr.(type) {
		case *ast.NodeSharedExpr:
			sharedExprReg := ir.NewGDIRRegObject(cpu.Ra, set)

			// Check if shared expression has been processed
			if !expr.HasBeenProcessed {
				obj, err := c.EvalNode(expr.Expr, stack)
				if err != nil {
					return nil, err
				}

				inst := ir.NewGDIRMov(sharedExprReg, obj, set)
				stack.AddNode(inst)
			}

			idxValue := runtime.NewGDIntNumber(runtime.GDInt(set.Index))
			idxNode := ir.NewGDIRObject(idxValue, set)

			inst, igetObj := ir.NewGDIRIGet(idxNode, true, sharedExprReg, set)
			stack.AddNode(inst)

			exprObj = igetObj
		default:
			obj, err := c.EvalNode(expr, stack)
			if err != nil {
				return nil, err
			}

			exprObj = obj
		}
	}

	return exprObj, nil
}

func (c *GDCompiler) EvalSet(s *ast.NodeSet, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	exprObj, err := c.resolveNodeSetExpr(s, stack)
	if err != nil {
		return nil, err
	}

	ident := c.DeriveIdent(s)
	typ := c.DeriveType(s)

	disc := ir.NewGDIRDiscoverable(s.IsPub, s.IsConst, ident, s)
	set := ir.NewGDIRSet(disc, typ, exprObj, s)
	stack.AddNode(set)

	return nil, nil
}

func (c *GDCompiler) EvalUpdateSet(u *ast.NodeUpdateSet, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	assignObj, err := c.EvalNode(u.Expr, stack)
	if err != nil {
		return nil, err
	}

	switch identExpr := u.IdentExpr.(type) {
	case *ast.NodeIterIdxExpr:
		expr, err := c.EvalNode(identExpr.Expr, stack)
		if err != nil {
			return nil, err
		}

		idx, err := c.EvalNode(identExpr.IdxExpr, stack)
		if err != nil {
			return nil, err
		}

		inst := ir.NewGDIRISet(idx, expr, assignObj, u)
		stack.AddNode(inst)
	default:
		exp, err := c.EvalNode(identExpr, stack)
		if err != nil {
			return nil, err
		}

		inst := ir.NewGDIRMov(exp, assignObj, u)
		stack.AddNode(inst)
	}

	return nil, nil
}

func (c *GDCompiler) EvalLabel(l *ast.NodeLabel, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	return nil, nil
}

func (c *GDCompiler) EvalIfElse(i *ast.NodeIfElse, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	endLabel := c.NewIdent()
	e, isElse := i.Else.(*ast.NodeIf)

	ifNodes := append([]ast.Node{i.If}, i.ElseIf...)
	ifNodeCount := len(ifNodes)

	for _, node := range ifNodes {
		if i, ok := node.(*ast.NodeIf); ok {
			ident := c.NewIdent()
			i.Ident = ident

			for _, node := range i.Conditions {
				cond, err := c.EvalNode(node, stack)
				if err != nil {
					return nil, err
				}

				cmpJump := ir.NewGDIRCompJump(cond, ir.NewGDIRObject(runtime.GDBool(true), i), ident, i)
				stack.AddNode(cmpJump)
			}
		}
	}

	if isElse {
		block, err := c.evalBlock(e.Block, stack)
		if err != nil {
			return nil, err
		}
		stack.AddNode(block)
	}

	// Jump to the end label
	stack.AddNode(ir.NewGDIRJump(endLabel, i))

	for i, node := range ifNodes {
		if nIf, ok := node.(*ast.NodeIf); ok {
			block, err := c.evalBlock(nIf.Block, stack)
			if err != nil {
				return nil, err
			}

			stack.AddNode(ir.NewGDIRLabel(nIf.Ident, nIf))
			stack.AddNode(block)

			// Jumps to the end label if is not the last node
			if i != ifNodeCount-1 {
				stack.AddNode(ir.NewGDIRJump(endLabel, nIf))
			}
		}
	}

	stack.AddNode(ir.NewGDIRLabel(endLabel, i))

	return nil, nil
}

func (c *GDCompiler) EvalTernaryIf(tif *ast.NodeTernaryIf, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	ifNode, err := c.EvalNode(tif.Expr, stack)
	if err != nil {
		return nil, err
	}

	thenObj, err := c.EvalNode(tif.Then, stack)
	if err != nil {
		return nil, err
	}

	elseObj, err := c.EvalNode(tif.Else, stack)
	if err != nil {
		return nil, err
	}

	inst, reg := ir.NewGDIRTIf(ifNode, thenObj, elseObj, tif)
	stack.AddNode(inst)

	return reg, nil
}

func (c *GDCompiler) EvalForIn(f *ast.NodeForIn, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	// Register where the iterable
	ra := ir.NewGDIRRegObject(cpu.Ra, f.Expr)

	// Register for the iterable index
	ri := ir.NewGDIRRegObject(cpu.Ri, f)

	return c.evalFor(
		f.NodeForIf,
		stack,
		func(_ runtime.GDIdent, stack ir.GDIRStackNode) error {
			return nil
		},
		func(endLabel runtime.GDIdent, stack ir.GDIRStackNode) error {
			expr, err := c.EvalNode(f.Expr, stack)
			if err != nil {
				return err
			}

			// Iterable expression
			stack.AddNode(ir.NewGDIRMov(ra, expr, f.Expr))
			stack.AddNode(ir.NewGDIRMov(ri, ir.NewGDIRObject(runtime.GDInt(0), f.Sets), f.Expr))

			return nil
		},
		func(endLabel runtime.GDIdent, stack ir.GDIRStackNode) error {
			// Compare jump condition:
			// if (ri < len(ra)) == false => goto endLabel
			lenInst, lenReg := ir.NewGDIRLen(ra, f.Expr)
			opInst, opReg := ir.NewGDIROp(runtime.ExprOperationLess, ri, lenReg, f)
			stack.AddNode(
				lenInst,
				opInst,
				ir.NewGDIRCompJump(opReg, ir.NewGDIRObject(runtime.GDBool(false), f), endLabel, f),
			)

			// Get the value of the index
			inst, reg := ir.NewGDIRIGet(ir.NewGDIRRegObject(cpu.Ri, f.Expr), true, ir.NewGDIRRegObject(cpu.Ra, f), f.Expr)
			stack.AddNode(inst)

			if f.InferredIterable != nil {
				// Iterable reference
				ident := c.DeriveIdent(f.InferredIterable)
				identObj := runtime.NewGDIdObject(ident, runtime.GDZNil)
				iterableIdent := ir.NewGDIRObject(identObj, f.InferredIterable)
				stack.AddNode(
					ir.NewGDIRMov(iterableIdent, reg, f.InferredIterable),
				)
			}

			if f.InferredIndex != nil {
				// Index reference
				ident := c.DeriveIdent(f.InferredIndex)
				identObj := runtime.NewGDIdObject(ident, runtime.GDZNil)
				indexIdent := ir.NewGDIRObject(identObj, f.InferredIndex)
				stack.AddNode(
					ir.NewGDIRMov(indexIdent, ir.NewGDIRRegObject(cpu.Ri, f), f.InferredIndex),
				)
			}

			// Increment the index
			addInst, addReg := ir.NewGDIROp(runtime.ExprOperationAdd, ri, ir.NewGDIRObject(runtime.GDInt(1), f), f)
			stack.AddNode(
				addInst,
				ir.NewGDIRMov(ri, addReg, f),
			)

			return nil
		},
	)
}

func (c *GDCompiler) EvalForIf(f *ast.NodeForIf, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	return c.evalFor(f, stack, nil, nil, nil)
}

func (c *GDCompiler) EvalCollectableOp(collectable *ast.NodeMutCollectionOp, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	exprLObj, err := c.EvalNode(collectable.L, stack)
	if err != nil {
		return nil, err
	}

	exprRObj, err := c.EvalNode(collectable.R, stack)
	if err != nil {
		return nil, err
	}

	inst, reg := ir.NewGDIRCOp(collectable.Op, exprLObj, exprRObj, collectable)
	stack.AddNode(inst)

	return reg, nil
}

func (c *GDCompiler) EvalTypeAlias(ta *ast.NodeTypeAlias, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	disc := ir.NewGDIRDiscoverable(ta.IsPub, true, runtime.NewGDStringIdent(ta.Ident.Lit), ta)

	alias := ir.NewGDIRTypeAlias(disc, ta.Type, ta)
	stack.AddNode(alias)

	return nil, nil
}

func (c *GDCompiler) EvalCastExpr(cast *ast.NodeCastExpr, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	exprObj, err := c.EvalNode(cast.Expr, stack)
	if err != nil {
		return nil, err
	}

	switch obj := exprObj.(type) {
	case *ir.GDIRObject:
		switch obj.Obj.(type) {
		case *runtime.GDIdObject:
			inst, reg := ir.NewGDIRCastObject(cast.Type, exprObj, cast)
			stack.AddNode(inst)

			return reg, nil
		}

		irObj := ir.NewGDIRObject(cast.InferredObject(), cast)
		return irObj, nil
	default:
		return obj, nil
	}
}

func (c *GDCompiler) EvalPackage(p *ast.NodePackage, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	switch p.InferredMode {
	case runtime.PackageModeBuiltin, runtime.PackageModeSource:
		ident := runtime.NewGDStringIdent(p.InferredPath)
		importNodes := make([]runtime.GDIdent, len(p.Imports))
		for i, node := range p.Imports {
			identNode, isIdentNode := node.(*ast.NodeIdent)
			if !isIdentNode {
				panic("Invalid node type: expected *ast.NodeIdent")
			}

			importNodes[i] = runtime.NewGDStringIdent(identNode.Lit)
		}

		irUse := ir.NewGDIRUse(p.InferredMode, ident, importNodes)
		stack.AddNode(irUse)

		return irUse, nil
	}

	return nil, nil
}

func (c *GDCompiler) evalBlock(b *ast.NodeBlock, _ ir.GDIRStackNode) (ir.GDIRNode, error) {
	block := ir.NewGDIRBlock()

	for _, node := range b.Nodes {
		_, err := c.EvalNode(node, block)
		if err != nil {
			return nil, err
		}

		switch node := node.(type) {
		case *ast.NodeBreak:
			parent := node.GetParentNodeByType(ast.NodeTypeFor)
			if parent != nil {
				if nodeFor, isNodeFor := parent.(ast.NodeFor); isNodeFor {
					block.AddNode(ir.NewGDIRJump(nodeFor.GetEndLabel(), node))
				}
			}
		}
	}

	return block, nil
}

func (c *GDCompiler) evalFor(f *ast.NodeForIf, stack ir.GDIRStackNode, forStart, preSets, inLoop ForCallback) (ir.GDIRNode, error) {
	loopLabel, endLabel := c.NewIdent(), c.NewIdent()
	f.SetEndLabel(endLabel)

	block := ir.NewGDIRBlock()

	if forStart != nil {
		err := forStart(f.GetEndLabel(), block)
		if err != nil {
			return nil, err
		}
	}

	if sets, ok := f.Sets.(*ast.NodeSets); ok {
		_, err := c.EvalSets(sets, block)
		if err != nil {
			return nil, err
		}
	}

	// Pre loop sets
	if preSets != nil {
		err := preSets(endLabel, block)
		if err != nil {
			return nil, err
		}
	}

	// Label for the loop
	block.AddNode(ir.NewGDIRLabel(loopLabel, f))

	// Loop conditions
	if inLoop != nil {
		err := inLoop(endLabel, block)
		if err != nil {
			return nil, err
		}
	}

	// Evaluate the condition
	if f.Conditions != nil {
		for _, node := range f.Conditions {
			cond, err := c.EvalNode(node, block)
			if err != nil {
				return nil, err
			}

			block.AddNode(ir.NewGDIRCompJump(cond, ir.NewGDIRObject(runtime.GDBool(false), node), endLabel, node))
		}
	}

	// for block
	loopBody, err := c.evalBlock(f.Block, block)
	if err != nil {
		return nil, err
	}
	block.AddNode(loopBody)

	// Jump to the loop label
	block.AddNode(ir.NewGDIRJump(loopLabel, f))

	// And end of the loop
	block.AddNode(ir.NewGDIRLabel(endLabel, f))

	// Add the block
	stack.AddNode(block)

	return block, nil
}

func (c *GDCompiler) collectNodes(typ runtime.GDTypable, nodes []ast.Node, stack ir.GDIRStackNode) (*ir.GDIRObject, error) {
	gdNodes := make([]ir.GDIRNode, 0)
	for _, node := range nodes {
		var expr ast.Node
		switch node := node.(type) {
		case *ast.NodeStructAttr:
			expr = node.Expr
		default:
			expr = node
		}

		gdNode, err := c.EvalNode(expr, stack)
		if err != nil {
			return nil, err
		}

		gdNodes = append(gdNodes, gdNode)
	}

	return ir.NewGDIRIterableObject(typ, gdNodes), nil
}

func NewGDCompiler() *GDCompiler {
	byteCode := &GDCompiler{Ctx: ir.NewGDIRContext(), Root: ir.NewGDIRBlock(), PackageDependenciesAnalyzer: analysis.NewPackageDependenciesAnalyzer(), GDIdentGen: tools.NewGDIdentStringGen()}
	byteCode.GDCompExpressionEvaluator = GDCompExpressionEvaluator{Evaluator: byteCode}

	return byteCode
}
