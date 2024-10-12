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

package staticcheck

import (
	"errors"
	"gdlang/lib/builtin"
	"gdlang/lib/runtime"
	"gdlang/lib/tools"
	"gdlang/src/comn"
	"gdlang/src/gd/analysis"
	"gdlang/src/gd/ast"
	"gdlang/src/gd/scanner"
)

type Inference = ast.Inference

func NewInference(ident runtime.GDIdent, typ runtime.GDTypable) *Inference {
	return &Inference{Ident: ident, Type: typ, RuntimeIdent: nil}
}

type StaticStack = runtime.GDStack
type StaticSymbol = runtime.GDSymbol

type (
	ObjectEvaluator           = Evaluator[*Inference, *StaticStack]
	ObjectExpressionEvaluator = ExpressionEvaluator[*Inference, *StaticStack]
)

type StaticCheck struct {
	ObjectEvaluator           // Implements the Evaluator interface
	ObjectExpressionEvaluator // Embeds the evaluator process to evaluate the AST nodes
	tools.GDIdentGen
	*analysis.PackageDependenciesAnalyzer
}

func (t *StaticCheck) Check(stack *StaticStack) error {
	return t.EvalFileNodes(t.Nodes, stack)
}

func (t *StaticCheck) EvalAtom(a *ast.NodeLiteral, stack *StaticStack) (*Inference, error) {
	var obj runtime.GDObject
	switch a.Token {
	case scanner.STRING:
		obj = runtime.GDString(a.Lit)
	case scanner.CHAR:
		obj = runtime.GDChar(a.Lit[0])
	case scanner.INT:
		intVal, err := runtime.NewGDIntNumberFromString(a.Lit)
		if err != nil {
			return nil, comn.WrapSyntaxErr(err, a.Position)
		}

		obj = intVal
	case scanner.TRUE:
		obj = runtime.GDBool(true)
	case scanner.FALSE:
		obj = runtime.GDBool(false)
	case scanner.NIL:
		obj = runtime.GDZNil
	case scanner.IMAG:
		imgVal, err := runtime.NewGDComplexNumberFromString(a.Lit)
		if err != nil {
			return nil, comn.WrapSyntaxErr(err, a.Position)
		}

		obj = imgVal
	case scanner.FLOAT:
		floatVal, err := runtime.NewGDFloatNumberFromString(a.Lit)
		if err != nil {
			return nil, comn.WrapSyntaxErr(err, a.Position)
		}

		obj = floatVal
	default:
		panic("unexpected literal token: " + a.Lit)
	}

	// Set the inferred object to the AST node, to avoid re-evaluating the same object
	// in later stages of the static analysis.
	a.InferredObject = obj
	a.Inference = NewInference(nil, obj.GetType())

	return a.Inference, nil
}

// Mainly used to identify objects in the stack
func (t *StaticCheck) EvalIdent(i *ast.NodeIdent, stack *StaticStack) (*Inference, error) {
	ident := runtime.NewGDStrIdent(i.Lit)

	symbol, err := stack.GetSymbol(ident)
	if err != nil {
		return nil, comn.WrapFatalErr(err, i.GetPosition())
	}

	inference := symbol.Value.(*Inference)
	i.Inference = inference

	return inference, nil
}

func (t *StaticCheck) EvalLambda(l *ast.NodeLambda, stack *StaticStack) (*Inference, error) {
	lambdaStack := stack.NewStack(runtime.LambdaCtx)
	defer lambdaStack.Dispose()

	lambdaInf, err := t.evalLambda(l, lambdaStack)
	if err != nil {
		return nil, err
	}

	// Evaluate the block
	_, err = t.evalBlock(l.Block, lambdaStack)
	if err != nil {
		return nil, err
	}

	l.Inference = lambdaInf

	return lambdaInf, nil
}

func (t *StaticCheck) evalLambda(l *ast.NodeLambda, stack *StaticStack) (*Inference, error) {
	argTypes := make([]*runtime.GDLambdaArgType, len(l.Type.ArgTypes))
	funcArgsLen := len(l.Type.ArgTypes)
	if l.Type.IsVariadic {
		funcArgsLen--
		arg := l.Type.ArgTypes[funcArgsLen]

		// Functions arguments are always variables, not constants and not public
		// because they are only accessible within the function.
		symbol, err := t.addSymbol(arg.Key, false, false, runtime.NewGDArrayType(arg.Value), nil, stack)
		if err != nil {
			return nil, err
		}

		inference, isInf := symbol.Value.(*Inference)
		if !isInf {
			panic("expected an *Inference")
		}

		argTypes[funcArgsLen] = &runtime.GDLambdaArgType{Key: inference.RuntimeIdent, Value: inference.Type}
	}

	for i := 0; i < funcArgsLen; i++ {
		funcArg := l.Type.ArgTypes[i]

		symbol, err := t.addSymbol(funcArg.Key, false, false, funcArg.Value, nil, stack)
		if err != nil {
			return nil, err
		}

		inference, isInf := symbol.Value.(*Inference)
		if !isInf {
			panic("expected an *Inference")
		}

		argTypes[i] = &runtime.GDLambdaArgType{Key: inference.RuntimeIdent, Value: inference.Type}
	}

	inference := NewInference(nil, l.Type)
	runtimeInference := NewInference(nil, runtime.NewGDLambdaType(argTypes, l.Type.ReturnType, l.Type.IsVariadic))

	l.Inference = inference
	l.RuntimeInference = runtimeInference

	// Set the return type of the block
	l.Block.ReturnType = l.Type.ReturnType

	return inference, nil
}

func (t *StaticCheck) EvalExprOp(e *ast.NodeExprOperation, stack *StaticStack) (*Inference, error) {
	leftInf, err := t.EvalNode(e.L, stack)
	if err != nil {
		return nil, err
	}

	var rightInf *Inference
	if e.R != nil {
		rightInf, err = t.EvalNode(e.R, stack)
		if err != nil {
			return nil, err
		}
	}

	// Unary operation
	if e.R == nil {
		return leftInf, nil
	}

	evalTypesFromUnion := func(a runtime.GDTypable, b runtime.GDUnionType) ([]runtime.GDTypable, error) {
		types := make([]runtime.GDTypable, 0)
		for _, obj := range b {
			typ, err := runtime.TypeCheckExprOperation(e.Op, a, obj)
			if err != nil {
				return nil, err
			}

			types = append(types, typ)
		}

		return types, nil
	}

	evalTypesBetweenUnions := func(a, b runtime.GDUnionType) ([]runtime.GDTypable, error) {
		types := make([]runtime.GDTypable, 0)
		for _, a := range a {
			unionTypes, err := evalTypesFromUnion(a, b)
			if err != nil {
				return nil, comn.WrapFatalErr(err, e.R.GetPosition())
			}

			types = append(types, unionTypes...)
		}

		return types, nil
	}

	types := make([]runtime.GDTypable, 0)
	typeA, typeB := leftInf.Type, rightInf.Type
	switch typeA := typeA.(type) {
	case runtime.GDUnionType:
		switch typeB := typeB.(type) {
		case runtime.GDUnionType:
			unionTypes, err := evalTypesBetweenUnions(typeA, typeB)
			if err != nil {
				return nil, comn.WrapFatalErr(err, e.R.GetPosition())
			}

			types = append(types, unionTypes...)
		}
	default:
		switch TypeB := typeB.(type) {
		case runtime.GDUnionType:
			unionTypes, err := evalTypesFromUnion(typeA, TypeB)
			if err != nil {
				return nil, comn.WrapFatalErr(err, e.GetPosition())
			}

			types = append(types, unionTypes...)
		default:
			typ, err := runtime.TypeCheckExprOperation(e.Op, typeA, TypeB)
			if err != nil {
				return nil, comn.WrapFatalErr(err, e.GetPosition())
			}

			types = append(types, typ)
		}
	}

	return NewInference(nil, runtime.ComputeTypeFromTypes(types)), nil
}

func (t *StaticCheck) EvalExpEllipsis(e *ast.NodeEllipsisExpr, stack *StaticStack) (*Inference, error) {
	exprInf, err := t.EvalNode(e.Expr, stack)
	if err != nil {
		return nil, err
	}

	switch typ := exprInf.Type.(type) {
	case runtime.GDIterableCollectionType:
		spreadType := runtime.NewGDSpreadableType(typ)

		inference := NewInference(nil, spreadType)
		e.Inference = inference

		return inference, nil
	}

	return nil, comn.NewError(comn.InvalidSpreadableTypeErrCode, comn.InvalidArraySpreadExpressionErrorMsg, comn.FatalError, e.GetPosition(), nil)
}

// Structure of a function node:
// func Ident(param: Type, ...) => Type { ... }
func (t *StaticCheck) EvalFunc(f *ast.NodeFunc, stack *StaticStack) (*Inference, error) {
	lambdaStack := stack.NewStack(runtime.LambdaCtx)
	defer lambdaStack.Dispose()

	lambdaInf, err := t.evalLambda(f.NodeLambda, lambdaStack)
	if err != nil {
		return nil, err
	}

	ident := runtime.NewGDStrIdent(f.Ident.Lit)
	symbol, err := t.addSymbol(ident, f.IsPub, true, lambdaInf.Type, nil, stack)
	if err != nil {
		return nil, comn.WrapFatalErr(err, f.Ident.Position)
	}

	// Evaluate the block
	_, err = t.evalBlock(f.NodeLambda.Block, lambdaStack)
	if err != nil {
		return nil, err
	}

	inference := symbol.Value.(*Inference)
	f.Inference = inference

	return inference, nil
}

func (t *StaticCheck) EvalTuple(tu *ast.NodeTuple, stack *StaticStack) (*Inference, error) {
	if len(tu.Nodes) == 0 {
		inference := NewInference(nil, runtime.NewGDTupleType())
		tu.Inference = inference
		return inference, nil
	}

	types := make([]runtime.GDTypable, len(tu.Nodes))
	for i, expr := range tu.Nodes {
		typeInf, err := t.EvalNode(expr, stack)
		if err != nil {
			return nil, err
		}

		types[i] = typeInf.Type
	}

	tupleType := runtime.NewGDTupleType(types...)

	inference := NewInference(nil, tupleType)
	tu.Inference = inference

	return inference, nil
}

func (t *StaticCheck) EvalStruct(s *ast.NodeStruct, stack *StaticStack) (*Inference, error) {
	attrTypes := make([]*runtime.GDStructAttrType, len(s.Nodes))
	for i, expr := range s.Nodes {
		switch expr := expr.(type) {
		case *ast.NodeStructAttr:
			exprInf, err := t.EvalNode(expr.Expr, stack)
			if err != nil {
				return nil, err
			}

			// Identifiers are always strings
			ident := runtime.NewGDStrIdent(expr.Ident.Lit)

			// Create a new attribute type
			attrTypes[i] = &runtime.GDStructAttrType{Ident: ident, Type: exprInf.Type}

			// Update the identifier of the inferred type
			exprInf.Ident = ident

			// Set the inferred type to the AST node
			expr.Inference = exprInf
		default:
			panic("expected a struct attribute")
		}
	}

	inference := NewInference(nil, runtime.NewGDStructType(attrTypes...))
	s.Inference = inference

	return inference, nil
}

func (t *StaticCheck) EvalArray(a *ast.NodeArray, stack *StaticStack) (*Inference, error) {
	if len(a.Nodes) == 0 {
		inference := NewInference(nil, runtime.NewGDEmptyArrayType())
		a.Inference = inference
		return inference, nil
	}

	types := make([]runtime.GDTypable, len(a.Nodes))
	for i, expr := range a.Nodes {
		exprInf, err := t.EvalNode(expr, stack)
		if err != nil {
			return nil, err
		}

		types[i] = exprInf.Type
	}

	inference := NewInference(nil, runtime.NewGDArrayTypeWithTypes(types, stack))

	a.Inference = inference

	return inference, nil
}

func (t *StaticCheck) EvalReturn(r *ast.NodeReturn, stack *StaticStack) (*Inference, error) {
	// Return with not expression
	// For example: `return`
	if r.Expr == nil {
		retInf := NewInference(nil, runtime.GDNilTypeRef)
		r.Inference = retInf
		return retInf, nil
	}

	exprInf, err := t.EvalNode(r.Expr, stack)
	if err != nil {
		return nil, err
	}

	inference := exprInf
	r.Inference = inference

	return inference, nil
}

func (t *StaticCheck) EvalIterIdxExpr(a *ast.NodeIndexableExpr, stack *StaticStack) (*Inference, error) {
	exprInf, err := t.EvalNode(a.Expr, stack)
	if err != nil {
		return nil, err
	}

	indexInf, err := t.EvalNode(a.IdxExpr, stack)
	if err != nil {
		return nil, err
	}

	if err := runtime.EqualTypes(indexInf.Type, runtime.GDIntTypeRef, stack); err != nil {
		return nil, comn.WrapFatalErr(err, a.IdxExpr.GetPosition())
	}

	if iterable, isIterable := exprInf.Type.(runtime.GDIterableCollectionType); isIterable {
		return NewInference(nil, iterable.GetIterableType()), nil
	} else {
		return nil, comn.WrapFatalErr(runtime.InvalidIterableTypeErr(exprInf.Type), a.Expr.GetPosition())
	}
}

func (t *StaticCheck) EvalCallExpr(c *ast.NodeCallExpr, stack *StaticStack) (*Inference, error) {
	exprInf, err := t.EvalNode(c.Expr, stack)
	if err != nil {
		return nil, err
	}

	lambdaType, isLambdaType := exprInf.Type.(*runtime.GDLambdaType)
	if !isLambdaType {
		return nil, comn.WrapFatalErr(runtime.InvalidCallableTypeErr(exprInf.Type), c.GetPosition())
	}

	err = lambdaType.CheckNumberOfArgs(uint(len(c.Args)))
	if err != nil {
		return nil, comn.WrapFatalErr(err, c.GetPosition())
	}

	for i := len(c.Args) - 1; i >= 0; i-- {
		arg := c.Args[i]
		argInf, err := t.EvalNode(arg, stack)
		if err != nil {
			return nil, err
		}

		err = lambdaType.CheckLambdaArgAtIndex(i, argInf.Type, stack)
		if err != nil {
			return nil, comn.WrapFatalErr(err, arg.GetPosition())
		}
	}

	inference := NewInference(nil, lambdaType.ReturnType)

	c.Inference = inference

	return inference, nil
}

func (t *StaticCheck) EvalSafeDotExpr(s *ast.NodeSafeDotExpr, stack *StaticStack) (*Inference, error) {
	switch identExpr := s.Ident.(type) {
	case *ast.NodeTokenInfo:
		idxExpr := ast.NewNodeIndexableExpr(s.IsNilSafe, s.Expr, ast.NewNodeLiteral(identExpr))
		return t.EvalIterIdxExpr(idxExpr, stack)
	case *ast.NodeIdent:
		exprInf, err := t.EvalNode(s.Expr, stack)
		if err != nil {
			return nil, err
		}

		attrIdent := runtime.NewGDStrIdent(identExpr.Lit)
		if exprInf.Type == runtime.GDNilTypeRef {
			if s.IsNilSafe {
				s.Inference = NewInference(attrIdent, runtime.GDNilTypeRef)
				return s.Inference, nil
			} else {
				return nil, comn.AnalysisErr(comn.NilAccessExceptionErrMsg, s.GetPosition())
			}
		}

		if attributable, isAttributable := exprInf.Type.(runtime.GDAttributableType); isAttributable {
			typ, err := attributable.GetAttrType(attrIdent)
			if err != nil {
				return nil, comn.WrapFatalErr(err, s.GetPosition())
			}

			inference := NewInference(attrIdent, typ)
			s.Inference = inference

			return inference, nil
		}

		return nil, comn.WrapFatalErr(runtime.InvalidAttributableTypeErr(exprInf.Type), s.Expr.GetPosition())
	}

	return nil, nil
}

func (t *StaticCheck) EvalSets(s *ast.NodeSets, stack *StaticStack) (*Inference, error) {
	for _, node := range s.Nodes {
		_, err := t.EvalNode(node, stack)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (t *StaticCheck) checkNodeSetExpr(set *ast.NodeSet, stack *StaticStack) (*Inference, error) {
	var exprType runtime.GDTypable = runtime.GDNilTypeRef
	if set.Expr != nil {
		switch expr := set.Expr.(type) {
		case *ast.NodeSharedExpr:
			var sharedType runtime.GDTypable
			if expr.Inference == nil {
				exprInf, err := t.EvalNode(expr.Expr, stack)
				if err != nil {
					return nil, err
				}

				sharedType = exprInf.Type

				expr.Inference = NewInference(nil, sharedType)
			} else {
				sharedType = expr.Inference.Type
			}

			if iterable, itIterable := sharedType.(runtime.GDIterableCollectionType); itIterable {
				exprType = iterable.GetTypeAt(int(set.Index))
			} else {
				return nil, comn.WrapFatalErr(runtime.InvalidIterableTypeErr(sharedType), set.GetPosition())
			}
		default:
			exprInf, err := t.EvalNode(expr, stack)
			if err != nil {
				return nil, err
			}

			exprType = exprInf.Type
		}
	} else {
		exprType = runtime.GDNilTypeRef
		set.Expr = ast.NewNodeNilLiteral(set.GetPosition())
	}

	inference := NewInference(nil, exprType)

	return inference, nil
}

func (t *StaticCheck) EvalSet(s *ast.NodeSet, stack *StaticStack) (*Inference, error) {
	exprInf, err := t.checkNodeSetExpr(s, stack)
	if err != nil {
		return nil, err
	}

	ident := runtime.NewGDStrIdent(s.IdentWithType.Ident.Lit)

	inferredType, err := runtime.InferType(s.IdentWithType.Type, exprInf.Type, stack)
	if err != nil {
		return nil, comn.WrapFatalErr(err, s.GetPosition())
	}

	symbol, err := t.addSymbol(ident, s.IsPub, s.IsConst, inferredType, nil, stack)
	if err != nil {
		var err runtime.GDRuntimeErr
		switch {
		case errors.As(err, &err):
			switch err.Code {
			case runtime.DuplicatedObjectCreationCode:
				return nil, comn.WrapFatalErr(err, s.IdentWithType.Ident.GetPosition())
			case runtime.SetObjectWrongTypeErrCode, runtime.IncompatibleTypeCode:
				if s.Expr != nil {
					return nil, comn.WrapFatalErr(err, s.Expr.GetPosition())
				}

				return nil, comn.WrapFatalErr(err, s.IdentWithType.Ident.GetPosition())
			default:
				panic("unhandled default case")
			}
		}
		return nil, comn.WrapFatalErr(err, s.GetPosition())
	}

	inference := symbol.Value.(*Inference)
	s.Inference = inference

	return inference, nil
}

func (t *StaticCheck) EvalUpdateSet(u *ast.NodeUpdateSet, stack *StaticStack) (*Inference, error) {
	rhsExprInf, err := t.EvalNode(u.RhsExpr, stack)
	if err != nil {
		return nil, err
	}

	switch lhsExpr := u.LhsExpr.(type) {
	// Expr[IdxExpr] = RhsExpr
	case *ast.NodeIndexableExpr:
		// Indexable expression type
		exprInf, err := t.EvalNode(lhsExpr.Expr, stack)
		if err != nil {
			return nil, err
		}

		// Index type
		indexInf, err := t.EvalNode(lhsExpr.IdxExpr, stack)
		if err != nil {
			return nil, err
		}

		if indexInf.Type != runtime.GDIntTypeRef {
			return nil, comn.WrapFatalErr(runtime.WrongTypesErr(runtime.GDIntTypeRef, indexInf.Type), lhsExpr.IdxExpr.GetPosition())
		}

		// Collectable
		mutable, isMutableType := exprInf.Type.(runtime.GDMutableCollectionType)
		if isMutableType {
			iterableType := mutable.GetIterableType()
			err := runtime.CanBeAssign(iterableType, rhsExprInf.Type, stack)
			if err != nil {
				return nil, comn.WrapFatalErr(runtime.WrongTypesErr(iterableType, rhsExprInf.Type), u.RhsExpr.GetPosition())
			}
		} else {
			return nil, comn.WrapFatalErr(runtime.InvalidMutableCollectionTypeErr(exprInf.Type), lhsExpr.Expr.GetPosition())
		}
	// Expr.Ident = RhsExpr
	case *ast.NodeSafeDotExpr:
		exprInf, err := t.EvalNode(lhsExpr.Expr, stack)
		if err != nil {
			return nil, err
		}
		attributable, isAttributable := exprInf.Type.(runtime.GDAttributableType)
		if !isAttributable {
			return nil, comn.WrapFatalErr(runtime.InvalidAttributableTypeErr(exprInf.Type), lhsExpr.Expr.GetPosition())
		}

		nodeIdent, isNodeIdent := lhsExpr.Ident.(*ast.NodeIdent)
		if !isNodeIdent {
			panic("expected an *ast.NodeIdent")
		}

		ident := runtime.NewGDStrIdent(nodeIdent.Lit)
		err = attributable.SetAttrType(ident, rhsExprInf.Type, stack)
		if err != nil {
			return nil, comn.WrapFatalErr(err, lhsExpr.Ident.GetPosition())
		}

		// Attributable type do not obfuscate the ids
		lhsExpr.Ident.SetInference(NewInference(ident, nil))
	// identExpr = assignObj
	default:
		lhsExprInf, err := t.EvalNode(lhsExpr, stack)
		if err != nil {
			return nil, err
		}

		symbol, err := stack.GetSymbol(lhsExprInf.Ident)
		if err != nil {
			return nil, comn.WrapFatalErr(err, lhsExpr.GetPosition())
		}

		err = symbol.SetType(rhsExprInf.Type, stack)
		if err != nil {
			return nil, comn.WrapFatalErr(err, u.RhsExpr.GetPosition())
		}
	}

	return nil, nil
}

func (t *StaticCheck) EvalLabel(l *ast.NodeLabel, stack *StaticStack) (*Inference, error) {
	return nil, nil
}

func (t *StaticCheck) EvalIfElse(i *ast.NodeIfElse, stack *StaticStack) (*Inference, error) {
	evalIfNode := func(ifNode ast.Node, stack *StaticStack) (*Inference, error) {
		if ifNode, isIfNode := ifNode.(*ast.NodeIf); isIfNode {
			_, err := t.evalIfNode(ifNode, stack)
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	}

	_, err := evalIfNode(i.If, stack)
	if err != nil {
		return nil, err
	}

	for _, ifNode := range i.ElseIf {
		_, err := evalIfNode(ifNode, stack)
		if err != nil {
			return nil, err
		}
	}

	_, err = evalIfNode(i.Else, stack)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *StaticCheck) EvalTernaryIf(tIf *ast.NodeTernaryIf, stack *StaticStack) (*Inference, error) {
	ifInf, err := t.EvalNode(tIf.Expr, stack)
	if err != nil {
		return nil, err
	}

	err = runtime.EqualTypes(runtime.GDBoolTypeRef, ifInf.Type, stack)
	if err != nil {
		return nil, comn.WrapFatalErr(err, tIf.Expr.GetPosition())
	}

	thenInf, err := t.EvalNode(tIf.Then, stack)
	if err != nil {
		return nil, err
	}

	elseInf, err := t.EvalNode(tIf.Else, stack)
	if err != nil {
		return nil, err
	}

	typ := runtime.ComputeTypeFromTypes([]runtime.GDTypable{thenInf.Type, elseInf.Type})

	inference := NewInference(nil, typ)

	tIf.Inference = inference

	return inference, nil
}

func (t *StaticCheck) EvalForIn(f *ast.NodeForIn, stack *StaticStack) (*Inference, error) {
	// Create a new stack for the for loop
	forStack := stack.NewStack(runtime.ForCtx)
	defer forStack.Dispose()

	// Evaluate the iterable expression inside the for loop
	exprInf, err := t.EvalNode(f.Expr, forStack)
	if err != nil {
		return nil, err
	}

	iterable, isIterable := exprInf.Type.(runtime.GDIterableCollectionType)
	if !isIterable {
		return nil, comn.WrapFatalErr(runtime.InvalidIterableTypeErr(exprInf.Type), f.Expr.GetPosition())
	}

	// Type check for the sets
	nodeSets, isSets := f.Sets.(*ast.NodeSets)
	if !isSets {
		panic("expected an *ast.NodeSets")
	}

	// Resolve the sets, objects are assigned to the symbols stack
	_, err = t.EvalNode(nodeSets, forStack)
	if err != nil {
		return nil, err
	}

	updateSymbolType := func(set ast.Node, assignType runtime.GDTypable) (*ast.NodeSet, error) {
		if set, isSet := set.(*ast.NodeSet); isSet {
			symbol, err := forStack.GetSymbol(set.Inference.Ident)
			if err != nil {
				return nil, err
			}

			// It is expected that the type of the symbol is the same as
			// the required type, and it can be either an Int or the iterable type
			err = symbol.SetType(assignType, stack)
			if err != nil {
				return nil, err
			}

			// Update the set according to the iterable type
			set.Inference.Type = symbol.Type

			return set, nil
		} else {
			panic("expected an *ast.NodeSet")
		}
	}

	sets := nodeSets.Nodes

	// Check set types, it must allow int and the iterable type
	sListLen := len(sets)
	switch {
	// for set (0) in iterable
	case sListLen == 1:
		iterSet, err := updateSymbolType(sets[0], iterable.GetIterableType())
		if err != nil {
			return nil, comn.WrapFatalErr(err, sets[0].GetPosition())
		}

		f.InferredIterable = iterSet
	// for set (0), (1) in iterable
	case sListLen > 1:
		idxSet, err := updateSymbolType(sets[0], runtime.GDIntTypeRef)
		if err != nil {
			return nil, comn.WrapFatalErr(err, sets[0].GetPosition())
		}
		f.InferredIndex = idxSet

		iterSet, err := updateSymbolType(sets[1], iterable.GetIterableType())
		if err != nil {
			return nil, comn.WrapFatalErr(err, sets[1].GetPosition())
		}

		f.InferredIterable = iterSet
	}

	_, err = t.evalBlock(f.NodeForIf.Block, forStack)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *StaticCheck) EvalForIf(f *ast.NodeForIf, stack *StaticStack) (*Inference, error) {
	// Create a new stack for the for loop
	forStack := stack.NewStack(runtime.ForCtx)
	defer forStack.Dispose()

	if f.Sets != nil {
		_, err := t.EvalNode(f.Sets, forStack)
		if err != nil {
			return nil, err
		}
	}

	if f.Conditions != nil {
		err := t.checkIfConditions(f.Conditions, forStack)
		if err != nil {
			return nil, err
		}
	}

	_, err := t.evalBlock(f.Block, forStack)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *StaticCheck) EvalCollectableOp(c *ast.NodeMutCollectionOp, stack *StaticStack) (*Inference, error) {
	exprLInf, err := t.EvalNode(c.L, stack)
	if err != nil {
		return nil, err
	}

	exprRInf, err := t.EvalNode(c.R, stack)
	if err != nil {
		return nil, err
	}

	if mutCollection, isMutCollection := exprLInf.Type.(runtime.GDMutableCollectionType); isMutCollection {
		switch c.Op {
		case ast.MutableCollectionAddOp:
			if ident, ok := c.L.(*ast.NodeIdent); ok {
				symbolId := runtime.NewGDStrIdent(ident.Lit)
				symbol, err := stack.GetSymbol(symbolId)
				if err != nil {
					return nil, comn.WrapFatalErr(err, ident.GetPosition())
				}

				err = runtime.CanBeAssign(mutCollection.GetIterableType(), exprRInf.Type, stack)
				if err != nil {
					return nil, comn.WrapFatalErr(err, c.GetPosition())
				}

				err = symbol.SetType(runtime.NewGDArrayType(exprRInf.Type), stack)
				if err != nil {
					return nil, comn.WrapFatalErr(err, c.GetPosition())
				}

				return exprRInf, nil
			}

			err := runtime.CanBeAssign(mutCollection.GetIterableType(), exprRInf.Type, stack)
			if err != nil {
				return nil, comn.WrapFatalErr(err, c.GetPosition())
			}

			return exprRInf, nil
		case ast.MutableCollectionRemoveOp:
			err := runtime.EqualTypes(runtime.GDIntTypeRef, exprRInf.Type, stack)
			if err != nil {
				return nil, comn.WrapFatalErr(err, c.GetPosition())
			}

			err = runtime.CanBeAssign(mutCollection.GetIterableType(), exprRInf.Type, stack)
			if err != nil {
				return nil, comn.WrapFatalErr(err, c.GetPosition())
			}

			return NewInference(nil, mutCollection.GetIterableType()), nil
		}
	}

	return nil, comn.WrapFatalErr(runtime.InvalidMutableCollectionTypeErr(exprLInf.Type), c.GetPosition())
}

func (t *StaticCheck) EvalTypeAlias(ta *ast.NodeTypeAlias, stack *StaticStack) (*Inference, error) {
	ident := runtime.NewGDStrIdent(ta.Ident.Lit)
	symbol, err := t.addSymbol(ident, ta.IsPub, true, ta.Type, nil, stack)
	if err != nil {
		return nil, comn.WrapFatalErr(err, ta.GetPosition())
	}

	inference := symbol.Value.(*Inference)

	ta.Inference = inference

	return nil, nil
}

func (t *StaticCheck) EvalCastExpr(c *ast.NodeCastExpr, stack *StaticStack) (*Inference, error) {
	exprInf, err := t.EvalNode(c.Expr, stack)
	if err != nil {
		return nil, err
	}

	err = runtime.CheckCastToType(exprInf.Type, c.Type)
	if err != nil {
		return nil, comn.WrapFatalErr(err, c.GetPosition())
	}

	c.Inference = exprInf

	return c.Inference, nil
}

// Register a new package in the symbol stack
// NOTE: Package must exist before evaluation
// Those checks are performed during the dependency analysis
func (t *StaticCheck) EvalPackage(p *ast.NodePackage, stack *StaticStack) (*Inference, error) {
	switch p.InferredMode {
	case runtime.PackageModeBuiltin:
		if pkg, found := builtin.Packages[p.InferredPath]; found {
			for _, node := range p.Imports {
				identNode, isIdentNode := node.(*ast.NodeIdent)
				if !isIdentNode {
					panic("Invalid node type: expected *ast.NodeIdent")
				}

				ident := runtime.NewGDStrIdent(identNode.Lit)
				if symbol, err := pkg.GetMember(ident); err == nil {
					_, err := t.addSymbol(ident, symbol.IsPub, symbol.IsConst, symbol.Type, nil, stack)
					if err != nil {
						return nil, comn.WrapFatalErr(err, p.GetPosition())
					}
				}
			}
		}
	}

	return nil, nil
}

func (t *StaticCheck) evalIfNode(i *ast.NodeIf, stack *StaticStack) (*Inference, error) {
	err := t.checkIfConditions(i.Conditions, stack)
	if err != nil {
		return nil, err
	}

	// Block evaluation
	return t.evalBlock(i.Block, stack)
}

// Checks that all the conditions are of type bool
// and returns an error if any of them is not.
// if condition, ... { ... }
func (t *StaticCheck) checkIfConditions(conditions []ast.Node, stack *StaticStack) error {
	for _, cond := range conditions {
		condInf, err := t.EvalNode(cond, stack)
		if err != nil {
			return err
		}

		err = runtime.EqualTypes(condInf.Type, runtime.GDBoolTypeRef, stack)
		if err != nil {
			return comn.WrapFatalErr(err, cond.GetPosition())
		}
	}

	return nil
}

func (t *StaticCheck) evalBlock(b *ast.NodeBlock, stack *StaticStack) (*Inference, error) {
	blockStack := stack.NewStack(runtime.BlockCtx)
	defer blockStack.Dispose()

	for _, node := range b.Nodes {
		nodeInf, err := t.EvalNode(node, blockStack)
		if err != nil {
			return nil, err
		}

		switch node := node.(type) {
		case *ast.NodeBreak:
			if b.Type == ast.FuncBlockType {
				return nil, comn.CompilerErr(comn.MisplacedBreakErrMsg, node.GetPosition())
			}
		case *ast.NodeReturn:
			// If not a flow control block, then return type is expected
			if b.Type != ast.ControlFlowBlockType {
				inferredType, err := runtime.InferType(b.ReturnType, nodeInf.Type, stack)
				if err != nil {
					return nil, comn.WrapFatalErr(err, node.GetPosition())
				}

				b.Inference = NewInference(nil, inferredType)

				continue
			}
		}
	}

	return nil, nil
}

func (t *StaticCheck) addSymbol(ident runtime.GDIdent, isPub, isConst bool, typ, assignType runtime.GDTypable, stack *StaticStack) (*StaticSymbol, error) {
	symbol, err := stack.AddNewSymbol(ident, isPub, isConst, typ, assignType, nil)
	if err != nil {
		return nil, err
	}

	// Set internal identifier for the symbol
	staticInf := NewInference(ident, symbol.Type)

	// It is used only internally to identify the symbol to map with the runtime identifier
	staticInf.RuntimeIdent = t.NewIdent()

	// Set the inference object to the symbol
	symbol.Value = staticInf

	return symbol, nil
}

func NewStaticCheck(analyzer *analysis.PackageDependenciesAnalyzer) *StaticCheck {
	staticCheck := &StaticCheck{PackageDependenciesAnalyzer: analyzer, GDIdentGen: NewIdentGenerator()}
	staticCheck.ObjectExpressionEvaluator = ObjectExpressionEvaluator{staticCheck}

	return staticCheck
}
