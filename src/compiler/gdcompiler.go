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
	"gdlang/src/analysis"
	"gdlang/src/cpu"
	"gdlang/src/gd/ast"
	"gdlang/src/gd/ir"
	"gdlang/src/gd/scanner"
	"os"
)

type (
	GDCmpEvaluator            = analysis.GDEvaluator[ir.GDIRNode, ir.GDIRStackNode]
	GDCompExpressionEvaluator = analysis.GDExpressionEvaluator[ir.GDIRNode, ir.GDIRStackNode]
)

type ForCallback func(endLabel runtime.GDIdentType, stack ir.GDIRStackNode) error

type GDCompiler struct {
	Ctx  *ir.GDIRContext
	Root *ir.GDIRBlock
	GDCmpEvaluator
	GDCompExpressionEvaluator
	*analysis.GDDepAnalyzer
}

func (c *GDCompiler) Compile(mainPkg string) error {
	stack := runtime.NewGDSymbolStack()
	err := builtin.Import(stack)
	if err != nil {
		return err
	}

	err = c.Build(mainPkg, analysis.DepAnalyzerOpt{MainAsEntryPoint: true})
	if err != nil {
		return err
	}

	staticAnalyzer := analysis.NewGDStaticAnalyzer(c.GDDepAnalyzer)
	err = staticAnalyzer.Check(stack)
	if err != nil {
		return err
	}

	stack.Dispose()

	for _, fileNode := range c.FileNodes {
		_, err := c.EvalNode(fileNode.Node, c.Root)
		if err != nil {
			return err
		}
	}

	main := ir.NewGDIRIdObject(runtime.GDStringIdentType("main"), runtime.GDZNil, scanner.ZeroPos)
	inst, _ := ir.NewGDIRCall(main, ir.NewGDIRIterableObject(runtime.NewGDArrayType(runtime.GDAnyType), []ir.GDIRNode{}, scanner.ZeroPos), scanner.ZeroPos)
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
	defer file.Close()

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
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "")

	return encoder.Encode(c.Ctx.SrcMap)
}

func (t *GDCompiler) EvalAtom(a *ast.NodeLiteral, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	return ir.NewGDIRObject(a.InferredObject(), a.GetPosition()), nil
}

func (t *GDCompiler) EvalIdent(i *ast.NodeIdent, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	return ir.NewGDIRObject(i.InferredObject(), i.GetPosition()), nil
}

func (t *GDCompiler) EvalLambda(l *ast.NodeLambda, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	lambda, reg := ir.NewGDIRLambda(l.Type, l.GetPosition())

	block, err := t.evalBlock(l.Block, lambda)
	if err != nil {
		return nil, err
	}

	lambda.GDIRBlock = block.(*ir.GDIRBlock)

	stack.AddNode(lambda)

	return reg, nil
}

func (t *GDCompiler) EvalExprOp(e *ast.NodeExprOperation, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	l, err := t.EvalNode(e.L, stack)
	if err != nil {
		return nil, err
	}

	var r ir.GDIRNode
	if e.R != nil {
		r, err = t.EvalNode(e.R, stack)
		if err != nil {
			return nil, err
		}
	}

	inst, reg := ir.NewGDIROp(e.Op, l, r, e.GetPosition())

	stack.AddNode(inst)

	return reg, nil
}

func (t *GDCompiler) EvalExpEllipsis(e *ast.NodeEllipsisExpr, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	obj, err := t.EvalNode(e.Expr, stack)
	if err != nil {
		return nil, err
	}

	spreadable := ir.NewGDIRObjectWithType(e.InferredType(), obj, e.GetPosition())

	return spreadable, nil
}

func (t *GDCompiler) EvalFunc(f *ast.NodeFunc, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	lambda, err := t.EvalLambda(f.NodeLambda, stack)
	if err != nil {
		return nil, err
	}

	ident := runtime.GDStringIdentType(f.Ident.Lit)
	disc := ir.NewGDIRDiscoverable(f.IsPub, true, ident)
	irFunc := ir.NewGDIRSet(disc, f.InferredType(), lambda, f.GetPosition())

	stack.AddNode(irFunc)

	return nil, nil
}

func (t *GDCompiler) EvalTuple(nt *ast.NodeTuple, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	return t.collectNodes(nt.InferredType(), nt.Nodes, stack)
}

func (t *GDCompiler) EvalStruct(s *ast.NodeStruct, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	return t.collectNodes(s.InferredType(), s.Nodes, stack)
}

func (t *GDCompiler) EvalArray(a *ast.NodeArray, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	return t.collectNodes(a.InferredType(), a.Nodes, stack)
}

func (t *GDCompiler) EvalReturn(r *ast.NodeReturn, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	if r.Expr == nil {
		inst, reg := ir.NewGDIRRet(ir.NewGDIRObject(runtime.GDZNil, r.GetPosition()), r.GetPosition())
		stack.AddNode(inst)

		return reg, nil
	}

	obj, err := t.EvalNode(r.Expr, stack)
	if err != nil {
		return nil, err
	}

	inst, reg := ir.NewGDIRRet(obj, r.Expr.GetPosition())
	stack.AddNode(inst)

	return reg, nil
}

func (t *GDCompiler) EvalIterIdxExpr(a *ast.NodeIterIdxExpr, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	expr, err := t.EvalNode(a.Expr, stack)
	if err != nil {
		return nil, err
	}

	idx, err := t.EvalNode(a.IdxExpr, stack)
	if err != nil {
		return nil, err
	}

	inst, reg := ir.NewGDIRIGet(idx, a.IsNilSafe, expr, a.GetPosition())
	stack.AddNode(inst)

	return reg, nil
}

func (t *GDCompiler) EvalCallExpr(c *ast.NodeCallExpr, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	expr, err := t.EvalNode(c.Expr, stack)
	if err != nil {
		return nil, err
	}

	args := make([]ir.GDIRNode, 0)
	for _, arg := range c.Args {
		arg, err := t.EvalNode(arg, stack)
		if err != nil {
			return nil, err
		}

		args = append(args, arg)
	}

	argsNode := ir.NewGDIRIterableObject(runtime.NewGDArrayType(runtime.GDAnyType), args, c.GetPosition())
	inst, reg := ir.NewGDIRCall(expr, argsNode, c.Expr.GetPosition())

	stack.AddNode(inst)

	return reg, nil
}

func (t *GDCompiler) EvalSafeDotExpr(s *ast.NodeSafeDotExpr, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	expr, err := t.EvalNode(s.Expr, stack)
	if err != nil {
		return nil, err
	}

	inst, reg := ir.NewGDIRAIGet(s.InferredIdent(), s.IsNilSafe, expr, s.GetPosition())
	stack.AddNode(inst)

	return reg, nil
}

func (t *GDCompiler) EvalSets(s *ast.NodeSets, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	for _, set := range s.Nodes {
		_, err := t.EvalNode(set, stack)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (t *GDCompiler) resolveNodeSetExpr(set *ast.NodeSet, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	var exprObj ir.GDIRNode
	if set.Expr != nil {
		switch expr := set.Expr.(type) {
		case *ast.NodeSharedExpr:
			sharedExprReg := ir.NewGDIRReg(cpu.Ra, set.GetPosition())

			// Check if shared expression has been processed
			if !expr.HasBeenProcessed {
				obj, err := t.EvalNode(expr.Expr, stack)
				if err != nil {
					return nil, err
				}

				inst := ir.NewGDIRMov(sharedExprReg, obj, set.GetPosition())
				stack.AddNode(inst)
			}

			idxValue := runtime.NewGDIntNumber(runtime.GDInt(set.Index))
			idxNode := ir.NewGDIRObject(idxValue, set.GetPosition())

			inst, igetObj := ir.NewGDIRIGet(idxNode, true, sharedExprReg, set.GetPosition())
			stack.AddNode(inst)

			exprObj = igetObj
		default:
			obj, err := t.EvalNode(expr, stack)
			if err != nil {
				return nil, err
			}

			exprObj = obj
		}
	}

	return exprObj, nil
}

func (t *GDCompiler) EvalSet(s *ast.NodeSet, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	exprObj, err := t.resolveNodeSetExpr(s, stack)
	if err != nil {
		return nil, err
	}

	disc := ir.NewGDIRDiscoverable(s.IsPub, s.IsConst, s.InferredIdent())
	set := ir.NewGDIRSet(disc, s.InferredType(), exprObj, s.GetPosition())
	stack.AddNode(set)

	return nil, nil
}

func (t *GDCompiler) EvalUpdateSet(u *ast.NodeUpdateSet, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	assignObj, err := t.EvalNode(u.Expr, stack)
	if err != nil {
		return nil, err
	}

	switch identExpr := u.IdentExpr.(type) {
	case *ast.NodeIterIdxExpr:
		expr, err := t.EvalNode(identExpr.Expr, stack)
		if err != nil {
			return nil, err
		}

		idx, err := t.EvalNode(identExpr.IdxExpr, stack)
		if err != nil {
			return nil, err
		}

		inst := ir.NewGDIRISet(idx, expr, assignObj, u.GetPosition())
		stack.AddNode(inst)
	default:
		exp, err := t.EvalNode(identExpr, stack)
		if err != nil {
			return nil, err
		}

		inst := ir.NewGDIRMov(exp, assignObj, u.GetPosition())
		stack.AddNode(inst)
	}

	return nil, nil
}

func (t *GDCompiler) EvalLabel(l *ast.NodeLabel, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	return nil, nil
}

func (t *GDCompiler) EvalIfElse(i *ast.NodeIfElse, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	endLabel := ir.NewStrIdent()
	e, isElse := i.Else.(*ast.NodeIf)

	ifNodes := append([]ast.Node{i.If}, i.ElseIf...)
	ifNodeCount := len(ifNodes)

	for _, node := range ifNodes {
		if i, ok := node.(*ast.NodeIf); ok {
			ident := ir.NewStrIdent()
			i.Ident = ident

			for _, node := range i.Conds {
				cond, err := t.EvalNode(node, stack)
				if err != nil {
					return nil, err
				}

				cmpJump := ir.NewGDIRCompJump(cond, ir.NewGDIRObject(runtime.GDBool(true), i.GetPosition()), ident, i.GetPosition())
				stack.AddNode(cmpJump)
			}
		}
	}

	if isElse {
		block, err := t.evalBlock(e.Block, stack)
		if err != nil {
			return nil, err
		}
		stack.AddNode(block)
	}

	// Jump to the end label
	stack.AddNode(ir.NewGDIRJump(endLabel, i.GetPosition()))

	for i, node := range ifNodes {
		if nIf, ok := node.(*ast.NodeIf); ok {
			block, err := t.evalBlock(nIf.Block, stack)
			if err != nil {
				return nil, err
			}

			stack.AddNode(ir.NewGDIRLabel(nIf.Ident, nIf.GetPosition()))
			stack.AddNode(block)

			// Jumps to the end label if is not the last node
			if i != ifNodeCount-1 {
				stack.AddNode(ir.NewGDIRJump(endLabel, nIf.GetPosition()))
			}
		}
	}

	stack.AddNode(ir.NewGDIRLabel(endLabel, i.GetPosition()))

	return nil, nil
}

func (t *GDCompiler) EvalTernaryIf(tif *ast.NodeTernaryIf, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	ifNode, err := t.EvalNode(tif.Expr, stack)
	if err != nil {
		return nil, err
	}

	thenObj, err := t.EvalNode(tif.Then, stack)
	if err != nil {
		return nil, err
	}

	elseObj, err := t.EvalNode(tif.Else, stack)
	if err != nil {
		return nil, err
	}

	inst, reg := ir.NewGDIRTIf(ifNode, thenObj, elseObj, tif.GetPosition())
	stack.AddNode(inst)

	return reg, nil
}

func (t *GDCompiler) EvalForIn(f *ast.NodeForIn, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	// Register where the iterable
	ra := ir.NewGDIRReg(cpu.Ra, f.Expr.GetPosition())

	// Register for the iterable index
	ri := ir.NewGDIRReg(cpu.Ri, f.GetPosition())

	return t.evalFor(
		f.NodeForIf,
		stack,
		func(_ runtime.GDIdentType, stack ir.GDIRStackNode) error {
			return nil
		},
		func(endLabel runtime.GDIdentType, stack ir.GDIRStackNode) error {
			expr, err := t.EvalNode(f.Expr, stack)
			if err != nil {
				return err
			}

			// Iterable expression
			stack.AddNode(ir.NewGDIRMov(ra, expr, f.Expr.GetPosition()))
			stack.AddNode(ir.NewGDIRMov(ri, ir.NewGDIRObject(runtime.GDInt(0), f.Sets.GetPosition()), f.Expr.GetPosition()))

			return nil
		},
		func(endLabel runtime.GDIdentType, stack ir.GDIRStackNode) error {
			// Compare jump condition:
			// if (ri < len(ra)) == false => goto endLabel
			lenInst, lenReg := ir.NewGDIRLen(ra, f.Expr.GetPosition())
			opInst, opReg := ir.NewGDIROp(runtime.ExprOperationLess, ri, lenReg, f.GetPosition())
			stack.AddNode(
				lenInst,
				opInst,
				ir.NewGDIRCompJump(opReg, ir.NewGDIRObject(runtime.GDBool(false), f.GetPosition()), endLabel, f.GetPosition()),
			)

			// Get the value of the index
			inst, reg := ir.NewGDIRIGet(ir.NewGDIRReg(cpu.Ri, f.Expr.GetPosition()), true, ir.NewGDIRReg(cpu.Ra, f.GetPosition()), f.Expr.GetPosition())
			stack.AddNode(inst)

			if f.InferredIterable != nil {
				// Iterable reference
				iterableIdent := ir.NewGDIRIdObject(f.InferredIterable.InferredIdent(), runtime.GDZNil, f.InferredIterable.GetPosition())
				stack.AddNode(
					ir.NewGDIRMov(iterableIdent, reg, f.InferredIterable.GetPosition()),
				)
			}

			if f.InferredIndex != nil {
				// Index reference
				indexIdent := ir.NewGDIRIdObject(f.InferredIndex.InferredIdent(), runtime.GDZNil, f.InferredIndex.GetPosition())
				stack.AddNode(
					ir.NewGDIRMov(indexIdent, ir.NewGDIRReg(cpu.Ri, f.GetPosition()), f.InferredIndex.GetPosition()),
				)
			}

			// Increment the index
			addInst, addReg := ir.NewGDIROp(runtime.ExprOperationAdd, ri, ir.NewGDIRObject(runtime.GDInt(1), f.GetPosition()), f.GetPosition())
			stack.AddNode(
				addInst,
				ir.NewGDIRMov(ri, addReg, f.GetPosition()),
			)

			return nil
		},
	)
}

func (t *GDCompiler) EvalForIf(f *ast.NodeForIf, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	return t.evalFor(f, stack, nil, nil, nil)
}

func (t *GDCompiler) EvalCollectableOp(c *ast.NodeMutCollectionOp, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	exprLObj, err := t.EvalNode(c.L, stack)
	if err != nil {
		return nil, err
	}

	exprRObj, err := t.EvalNode(c.R, stack)
	if err != nil {
		return nil, err
	}

	inst, reg := ir.NewGDIRCOp(c.Op, exprLObj, exprRObj, c.GetPosition())
	stack.AddNode(inst)

	return reg, nil
}

func (t *GDCompiler) EvalTypeAlias(ta *ast.NodeTypeAlias, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	disc := ir.NewGDIRDiscoverable(ta.IsPub, true, runtime.GDStringIdentType(ta.Ident.Lit))

	alias := ir.NewGDIRTypeAlias(disc, ta.Type, ta.GetPosition())
	stack.AddNode(alias)

	return nil, nil
}

func (t *GDCompiler) EvalCastExpr(c *ast.NodeCastExpr, stack ir.GDIRStackNode) (ir.GDIRNode, error) {
	exprObj, err := t.EvalNode(c.Expr, stack)
	if err != nil {
		return nil, err
	}

	switch obj := exprObj.(type) {
	case *ir.GDIRObject:
		switch obj.Obj.(type) {
		case *runtime.GDIdObject:
			inst, reg := ir.NewGDIRCastObject(c.Type, exprObj, c.GetPosition())
			stack.AddNode(inst)

			return reg, nil
		}

		irObj := ir.NewGDIRObject(c.InferredObject(), c.GetPosition())
		return irObj, nil
	default:
		return obj, nil
	}
}

func (t *GDCompiler) evalBlock(b *ast.NodeBlock, _ ir.GDIRStackNode) (ir.GDIRNode, error) {
	block := ir.NewGDIRBlock()

	for _, node := range b.Nodes {
		_, err := t.EvalNode(node, block)
		if err != nil {
			return nil, err
		}

		switch node := node.(type) {
		case *ast.NodeBreak:
			parent := node.GetParentNodeByType(ast.NodeTypeFor)
			if parent != nil {
				if nodeFor, isNodeFor := parent.(ast.NodeFor); isNodeFor {
					block.AddNode(ir.NewGDIRJump(nodeFor.GetEndLabel(), node.GetPosition()))
				}
			}
		}
	}

	return block, nil
}

func (t *GDCompiler) evalFor(f *ast.NodeForIf, stack ir.GDIRStackNode, forStart, preSets, inLoop ForCallback) (ir.GDIRNode, error) {
	loopLabel, endLabel := ir.NewStrIdent(), ir.NewStrIdent()
	f.SetEndLabel(endLabel)

	block := ir.NewGDIRBlock()

	if forStart != nil {
		err := forStart(f.GetEndLabel(), block)
		if err != nil {
			return nil, err
		}
	}

	if sets, ok := f.Sets.(*ast.NodeSets); ok {
		_, err := t.EvalSets(sets, block)
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
	block.AddNode(ir.NewGDIRLabel(loopLabel, f.GetPosition()))

	// Loop conditions
	if inLoop != nil {
		err := inLoop(endLabel, block)
		if err != nil {
			return nil, err
		}
	}

	// Evaluate the condition
	if f.Conds != nil {
		for _, node := range f.Conds {
			cond, err := t.EvalNode(node, block)
			if err != nil {
				return nil, err
			}

			block.AddNode(ir.NewGDIRCompJump(cond, ir.NewGDIRObject(runtime.GDBool(false), node.GetPosition()), endLabel, node.GetPosition()))
		}
	}

	// for block
	loopBody, err := t.evalBlock(f.Block, block)
	if err != nil {
		return nil, err
	}
	block.AddNode(loopBody)

	// Jump to the loop label
	block.AddNode(ir.NewGDIRJump(loopLabel, f.GetPosition()))

	// And end of the loop
	block.AddNode(ir.NewGDIRLabel(endLabel, f.GetPosition()))

	// Add the block
	stack.AddNode(block)

	return block, nil
}

func (t *GDCompiler) collectNodes(typ runtime.GDTypable, nodes []ast.Node, stack ir.GDIRStackNode) (*ir.GDIRObject, error) {
	gdNodes := make([]ir.GDIRNode, 0)
	for _, node := range nodes {
		var expr ast.Node
		switch node := node.(type) {
		case *ast.NodeStructAttr:
			expr = node.Expr
		default:
			expr = node
		}

		gdNode, err := t.EvalNode(expr, stack)
		if err != nil {
			return nil, err
		}

		gdNodes = append(gdNodes, gdNode)
	}

	return ir.NewGDIRIterableObject(typ, gdNodes, ast.GetStartEndPosition(nodes)), nil
}

func NewGDCompiler() *GDCompiler {
	byteCode := &GDCompiler{Ctx: ir.NewGDIRContext(), Root: ir.NewGDIRBlock(), GDDepAnalyzer: analysis.NDepAnalyzerProc()}
	byteCode.GDCompExpressionEvaluator = GDCompExpressionEvaluator{GDEvaluator: byteCode}

	return byteCode
}
