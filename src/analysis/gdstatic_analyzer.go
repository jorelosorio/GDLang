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

package analysis

import (
	"gdlang/lib/runtime"
	"gdlang/src/comn"
	"gdlang/src/gd/ast"
	"gdlang/src/gd/scanner"
)

type (
	ObjectEvaluator           = GDEvaluator[runtime.GDObject, *runtime.GDSymbolStack]
	ObjectExpressionEvaluator = GDExpressionEvaluator[runtime.GDObject, *runtime.GDSymbolStack]
)

type GDStaticAnalyzer struct {
	ObjectEvaluator           // Implements the Evaluator interface
	ObjectExpressionEvaluator // Embeds the evaluator process to evaluate the AST nodes
	*GDDepAnalyzer
}

func (t *GDStaticAnalyzer) Check(stack *runtime.GDSymbolStack) error {
	return t.EvalFileNodes(t.FileNodes, stack)
}

func (t *GDStaticAnalyzer) EvalAtom(a *ast.NodeLiteral, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
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

	a.SetInferredObject(obj)

	return obj, nil
}

func (t *GDStaticAnalyzer) EvalIdent(i *ast.NodeIdent, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	ident := runtime.NewGDStringIdent(i.Lit)
	symbol, err := stack.GetSymbol(ident)
	if err != nil {
		return nil, comn.WrapFatalErr(err, i.GetPosition())
	}

	obj := runtime.NewGDIdObject(ident, symbol.Object)

	i.SetInferredObject(obj)

	return obj, nil
}

func (t *GDStaticAnalyzer) EvalLambda(l *ast.NodeLambda, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	lambdaStack := stack.NewSymbolStack(runtime.LambdaCtx)
	defer lambdaStack.Dispose()

	for _, arg := range l.Type.ArgTypes {
		// Functions arguments are always variables, not constants
		// and they are not public, because they are only accessible within the function
		// and they are not constants because they can be changed.
		_, err := lambdaStack.AddSymbol(arg.Key, false, false, arg.Value, nil)
		if err != nil {
			return nil, comn.WrapFatalErr(err, l.GetPosition())
		}
	}

	lambdaObj := runtime.NewGDLambdaWithType(l.Type, stack, nil)

	// Evaluate the block
	_, err := t.evalBlock(l.Block, lambdaStack)
	if err != nil {
		return nil, err
	}

	// Return type is the function type for lambda
	return lambdaObj, nil
}

func (t *GDStaticAnalyzer) EvalExprOp(e *ast.NodeExprOperation, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	leftObj, err := t.EvalNode(e.L, stack)
	if err != nil {
		return nil, err
	}

	var rightObj runtime.GDObject
	if e.R != nil {
		rightObj, err = t.EvalNode(e.R, stack)
		if err != nil {
			return nil, err
		}
	}

	// Unary operation
	if e.R == nil {
		return leftObj, nil
	}

	evalObjectsFromUnion := func(a runtime.GDObject, b *runtime.GDUnion) ([]runtime.GDObject, error) {
		objects := make([]runtime.GDObject, 0)
		for _, obj := range b.Objects {
			typ, err := runtime.TypeCheckExprOperation(e.Op, a, obj)
			if err != nil {
				return nil, err
			}

			obj, err = runtime.ZObjectForType(typ, stack)
			if err != nil {
				return nil, err
			}
			objects = append(objects, obj)
		}

		return objects, nil
	}

	evalObjectsBetweenUnions := func(a, b *runtime.GDUnion) ([]runtime.GDObject, error) {
		objects := []runtime.GDObject{}
		for _, a := range a.Objects {
			objs, err := evalObjectsFromUnion(a, b)
			if err != nil {
				return nil, comn.WrapFatalErr(err, e.R.GetPosition())
			}

			objects = append(objects, objs...)
		}

		return objects, nil
	}

	objects := make([]runtime.GDObject, 0)
	a, b := runtime.Unwrap(leftObj), runtime.Unwrap(rightObj)
	switch a := a.(type) {
	case *runtime.GDUnion:
		switch b := b.(type) {
		case *runtime.GDUnion:
			objs, err := evalObjectsBetweenUnions(a, b)
			if err != nil {
				return nil, comn.WrapFatalErr(err, e.R.GetPosition())
			}

			objects = append(objects, objs...)
		}
	default:
		switch b := b.(type) {
		case *runtime.GDUnion:
			objs, err := evalObjectsFromUnion(a, b)
			if err != nil {
				return nil, comn.WrapFatalErr(err, e.GetPosition())
			}

			objects = append(objects, objs...)
		default:
			typ, err := runtime.TypeCheckExprOperation(e.Op, a, b)
			if err != nil {
				return nil, comn.WrapFatalErr(err, e.GetPosition())
			}

			obj, err := runtime.ZObjectForType(typ, stack)
			if err != nil {
				return nil, comn.WrapFatalErr(err, e.GetPosition())
			}

			objects = append(objects, obj)
		}
	}

	typ := runtime.ComputeTypeFromObjects(objects)
	if union, isUnion := typ.(runtime.GDUnionType); isUnion {
		return runtime.NewGDUnion(union, objects...), nil
	}

	return objects[0], nil
}

func (t *GDStaticAnalyzer) EvalExpEllipsis(e *ast.NodeEllipsisExpr, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	obj, err := t.EvalNode(e.Expr, stack)
	if err != nil {
		return nil, err
	}

	switch spreadable := runtime.Unwrap(obj).(type) {
	case runtime.GDIterableCollection:
		spreadObj := runtime.NewGDSpreadable(spreadable)
		e.SetInferredType(spreadObj.GetType())
		e.SetInferredObject(spreadObj)

		return spreadObj, nil
	}

	return nil, comn.NErr(comn.InvalidSpreadableTypeErrCode, comn.InvalidArraySpreadExpressionErrorMsg, comn.FatalError, e.GetPosition(), nil)
}

// Structure of a function node:
// func Ident(param: Type, ...) => Type { ... }
func (t *GDStaticAnalyzer) EvalFunc(f *ast.NodeFunc, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	funcStack := stack.NewSymbolStack(runtime.FuncCtx)
	defer funcStack.Dispose()

	gdfunc := runtime.NewGDLambdaWithType(f.Type, stack, nil)

	_, err := stack.AddSymbol(runtime.NewGDStringIdent(f.Ident.Lit), f.IsPub, true, f.Type, gdfunc)
	if err != nil {
		return nil, comn.WrapFatalErr(err, f.Ident.Position)
	}

	addArgSymbol := func(arg runtime.GDLambdaArgType) error {
		obj, err := runtime.ZObjectForType(arg.Value, funcStack)
		if err != nil {
			return comn.WrapFatalErr(err, f.GetPosition())
		}

		_, err = funcStack.AddSymbol(arg.Key, false, false, arg.Value, obj)
		if err != nil {
			return comn.WrapFatalErr(err, f.GetPosition())
		}

		return nil
	}

	funcArgsLen := len(f.Type.ArgTypes)
	if f.Type.IsVariadic {
		funcArgsLen--
		funcArg := f.Type.ArgTypes[funcArgsLen]

		vArgType := runtime.NewGDArrayType(funcArg.Value)

		err := addArgSymbol(runtime.GDLambdaArgType{Key: funcArg.Key, Value: vArgType})
		if err != nil {
			return nil, err
		}
	}

	for i := 0; i < funcArgsLen; i++ {
		funcArg := f.Type.ArgTypes[i]
		err := addArgSymbol(funcArg)
		if err != nil {
			return nil, err
		}
	}

	// Evaluate the block
	_, err = t.evalBlock(f.Block, funcStack)
	if err != nil {
		return nil, err
	}

	f.SetInferredType(gdfunc.GetType())
	f.SetInferredObject(gdfunc)

	return gdfunc, nil
}

func (t *GDStaticAnalyzer) EvalTuple(tu *ast.NodeTuple, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	if len(tu.Nodes) == 0 {
		tuple := runtime.NewGDTuple()
		tu.SetInferredType(tuple.GetType())
		tu.SetInferredObject(tuple)
		return tuple, nil
	}

	objects := make([]runtime.GDObject, len(tu.Nodes))
	for i, expr := range tu.Nodes {
		obj, err := t.EvalNode(expr, stack)
		if err != nil {
			return nil, err
		}

		objects[i] = obj
	}

	tuple := runtime.NewGDTuple(objects...)

	tu.SetInferredType(tuple.GetType())
	tu.SetInferredObject(tuple)

	return tuple, nil
}

func (t *GDStaticAnalyzer) EvalStruct(s *ast.NodeStruct, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	attrTypes := make([]runtime.GDStructAttrType, len(s.Nodes))
	objects := make([]runtime.GDObject, len(s.Nodes))
	for i, expr := range s.Nodes {
		switch expr := expr.(type) {
		case *ast.NodeStructAttr:
			obj, err := t.EvalNode(expr.Expr, stack)
			if err != nil {
				return nil, err
			}

			ident := runtime.NewGDStringIdent(expr.Ident.Lit)
			attrTypes[i] = runtime.GDStructAttrType{Ident: ident, Type: obj.GetType()}
			objects[i] = obj
		default:
			panic("expected a struct attribute")
		}
	}

	sType := runtime.NewGDStructType(attrTypes...)
	structObj, err := runtime.NewGDStruct(sType, stack)
	if err != nil {
		return nil, comn.WrapFatalErr(err, s.GetPosition())
	}

	for i, obj := range objects {
		err = structObj.SetAttr(attrTypes[i].Ident, obj)
		if err != nil {
			return nil, comn.WrapFatalErr(err, s.GetPosition())
		}
	}

	s.SetInferredType(sType)
	s.SetInferredObject(structObj)

	return structObj, nil
}

func (t *GDStaticAnalyzer) EvalArray(a *ast.NodeArray, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	if len(a.Nodes) == 0 {
		array := runtime.NewGDEmptyArray()
		a.SetInferredType(array.GetType())
		a.SetInferredObject(array)
		return array, nil
	}

	objects := make([]runtime.GDObject, len(a.Nodes))
	for i, expr := range a.Nodes {
		exprType, err := t.EvalNode(expr, stack)
		if err != nil {
			return nil, err
		}

		objects[i] = exprType
	}

	array, err := runtime.NewGDArrayWithObjects(objects, stack)
	if err != nil {
		return nil, comn.WrapFatalErr(err, a.GetPosition())
	}

	a.SetInferredType(array.GetType())
	a.SetInferredObject(array)

	return array, nil
}

func (t *GDStaticAnalyzer) EvalReturn(r *ast.NodeReturn, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	// Return with not expression
	// For example: `return`
	r.SetInferredObject(runtime.GDZNil)
	if r.Expr == nil {
		return r.InferredObject(), nil
	}

	obj, err := t.EvalNode(r.Expr, stack)
	if err != nil {
		return nil, err
	}

	r.SetInferredType(obj.GetType())
	r.SetInferredObject(obj)

	return obj, nil
}

func (t *GDStaticAnalyzer) EvalIterIdxExpr(a *ast.NodeIterIdxExpr, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	exprObj, err := t.EvalNode(a.Expr, stack)
	if err != nil {
		return nil, err
	}

	indexObj, err := t.EvalNode(a.IdxExpr, stack)
	if err != nil {
		return nil, err
	}

	if err := runtime.EqualTypes(indexObj.GetType(), runtime.GDIntType, stack); err != nil {
		return nil, comn.WrapFatalErr(err, a.IdxExpr.GetPosition())
	}

	if iterable, isIterable := runtime.Unwrap(exprObj).(runtime.GDIterableCollection); isIterable {
		obj, err := runtime.ZObjectForType(iterable.GetIterableType(), stack)
		if err != nil {
			return nil, comn.WrapFatalErr(err, a.GetPosition())
		}

		return obj, nil
	} else {
		return nil, comn.WrapFatalErr(runtime.InvalidIterableTypeErr(exprObj.GetType()), a.Expr.GetPosition())
	}
}

func (t *GDStaticAnalyzer) EvalCallExpr(c *ast.NodeCallExpr, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	exprObj, err := t.EvalNode(c.Expr, stack)
	if err != nil {
		return nil, err
	}

	exprType := exprObj.GetType()
	funcType, isFuncType := exprType.(*runtime.GDLambdaType)
	if !isFuncType {
		return nil, comn.WrapFatalErr(runtime.InvalidCallableTypeErr(exprType), c.GetPosition())
	}

	err = funcType.CheckNumberOfArgs(uint(len(c.Args)))
	if err != nil {
		return nil, comn.WrapFatalErr(err, c.GetPosition())
	}

	for i := len(c.Args) - 1; i >= 0; i-- {
		arg := c.Args[i]
		argObj, err := t.EvalNode(arg, stack)
		if err != nil {
			return nil, err
		}

		err = funcType.CheckArgAtIndex(i, argObj.GetType(), stack)
		if err != nil {
			return nil, comn.WrapFatalErr(err, arg.GetPosition())
		}
	}

	obj, err := runtime.ZObjectForType(funcType.ReturnType, stack)
	if err != nil {
		return nil, comn.WrapFatalErr(err, c.GetPosition())
	}

	return obj, nil
}

func (t *GDStaticAnalyzer) EvalSafeDotExpr(s *ast.NodeSafeDotExpr, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	switch identExpr := s.Ident.(type) {
	case *ast.NodeTokenInfo:
		idxExpr := ast.NewNodeIterIdxExpr(s.IsNilSafe, s.Expr, ast.NewNodeLiteral(identExpr))
		return t.EvalIterIdxExpr(idxExpr, stack)
	case *ast.NodeIdent:
		obj, err := t.EvalNode(s.Expr, stack)
		if err != nil {
			return nil, err
		}

		attrIdent := runtime.NewGDStringIdent(identExpr.Lit)
		s.SetInferredIdent(attrIdent)

		obj = runtime.Unwrap(obj)
		if obj == runtime.GDZNil {
			if s.IsNilSafe {
				return runtime.GDZNil, nil
			} else {
				return nil, comn.AnalysisErr(comn.NilAccessExceptionErrMsg, s.GetPosition())
			}
		}

		if attributable, isAttributable := obj.(runtime.GDAttributable); isAttributable {
			attrObj, err := attributable.GetAttr(attrIdent)
			if err != nil {
				return nil, comn.WrapFatalErr(err, s.GetPosition())
			}

			zObj, err := runtime.ZObjectForType(attrObj.Type, stack)
			if err != nil {
				return nil, comn.WrapFatalErr(err, s.GetPosition())
			}

			return runtime.NewGDAttrIdObject(attrIdent, zObj, attributable), nil
		}

		return nil, comn.WrapFatalErr(runtime.InvalidAttributableTypeErr(obj.GetType()), s.GetPosition())
	}

	return nil, nil
}

func (t *GDStaticAnalyzer) EvalSets(s *ast.NodeSets, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	for _, node := range s.Nodes {
		_, err := t.EvalNode(node, stack)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (t *GDStaticAnalyzer) checkNodeSetExpr(set *ast.NodeSet, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	var exprObj runtime.GDObject = runtime.GDZNil
	if set.Expr != nil {
		switch expr := set.Expr.(type) {
		case *ast.NodeSharedExpr:
			var sharedObj runtime.GDObject
			if expr.InferredObject() == nil {
				obj, err := t.EvalNode(expr.Expr, stack)
				if err != nil {
					return nil, err
				}

				sharedObj = obj
			} else {
				sharedObj = expr.InferredObject()
			}

			if iterable, ok := runtime.Unwrap(sharedObj).(runtime.GDIterableCollection); ok {
				obj, err := iterable.Get(int(set.Index))
				if err != nil {
					return nil, comn.WrapFatalErr(err, set.GetPosition())
				}

				exprObj = obj
			} else {
				return nil, comn.WrapFatalErr(runtime.InvalidIterableTypeErr(exprObj.GetType()), set.GetPosition())
			}
		default:
			obj, err := t.EvalNode(expr, stack)
			if err != nil {
				return nil, err
			}

			exprObj = obj
		}
	} else {
		exprObj = runtime.GDZNil
		set.Expr = ast.NewNodeNilLiteral(set.GetPosition())
	}

	return exprObj, nil
}

func (t *GDStaticAnalyzer) EvalSet(s *ast.NodeSet, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	exprObj, err := t.checkNodeSetExpr(s, stack)
	if err != nil {
		return nil, err
	}

	ident := runtime.NewGDStringIdent(s.IdentWithType.Ident.Lit)

	inferredType, err := runtime.InferType(s.IdentWithType.Type, exprObj.GetType(), stack)
	if err != nil {
		return nil, comn.WrapFatalErr(err, s.GetPosition())
	}

	exprObj, err = runtime.TypeCoercion(exprObj, inferredType, stack)
	if err != nil {
		return nil, comn.WrapFatalErr(err, s.GetPosition())
	}

	if s.Expr != nil {
		switch s.Expr.(type) {
		case *ast.NodeSharedExpr:
		default:
			s.Expr.SetInferredType(inferredType)
			s.Expr.SetInferredObject(exprObj)
		}
	}

	s.SetInferredIdent(ident)
	s.SetInferredType(inferredType)
	s.SetInferredObject(exprObj)

	_, err = stack.AddSymbol(ident, s.IsPub, s.IsConst, inferredType, exprObj)
	if err != nil {
		switch err := err.(type) {
		case runtime.GDRuntimeErr:
			switch err.Code {
			case runtime.DuplicatedObjectCreationCode:
				return nil, comn.WrapFatalErr(err, s.IdentWithType.Ident.GetPosition())
			case runtime.SetObjectWrongTypeErrCode, runtime.IncompatibleTypeCode:
				if s.Expr != nil {
					return nil, comn.WrapFatalErr(err, s.Expr.GetPosition())
				}

				return nil, comn.WrapFatalErr(err, s.IdentWithType.Ident.GetPosition())
			}
		}
		return nil, comn.WrapFatalErr(err, s.GetPosition())
	}

	return runtime.NewGDIdObject(ident, exprObj), nil
}

func (t *GDStaticAnalyzer) EvalUpdateSet(u *ast.NodeUpdateSet, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	assignObj, err := t.EvalNode(u.Expr, stack)
	if err != nil {
		return nil, err
	}

	switch identExpr := u.IdentExpr.(type) {
	// identExpr.Expr[identExpr.IdxExpr] = assignType
	case *ast.NodeIterIdxExpr:
		expressionObj, err := t.EvalNode(identExpr.Expr, stack)
		if err != nil {
			return nil, err
		}

		// Index type
		indexObj, err := t.EvalNode(identExpr.IdxExpr, stack)
		if err != nil {
			return nil, err
		}

		if indexObj.GetType() != runtime.GDIntType {
			return nil, comn.WrapFatalErr(runtime.WrongTypesErr(runtime.GDIntType, indexObj.GetType()), identExpr.IdxExpr.GetPosition())
		}

		// Collectable
		mutable, isMutableType := runtime.Unwrap(expressionObj).(runtime.GDMutableCollection)
		if isMutableType {
			err := runtime.CanBeAssign(mutable.GetIterableType(), assignObj.GetType(), stack)
			if err != nil {
				return nil, comn.WrapFatalErr(runtime.WrongTypesErr(expressionObj.GetType(), assignObj.GetType()), u.Expr.GetPosition())
			}

			// Store the inferred type
			u.SetInferredObject(expressionObj)
		} else {
			return nil, comn.WrapFatalErr(runtime.InvalidMutableCollectionTypeErr(expressionObj.GetType()), identExpr.Expr.GetPosition())
		}
	// identExpr = assignObj
	default:
		idxExpr, err := t.EvalNode(identExpr, stack)
		if err != nil {
			return nil, err
		}

		switch expr := idxExpr.(type) {
		case *runtime.GDIdObject:
			symbol, err := stack.GetSymbol(expr.Ident)
			if err != nil {
				return nil, comn.WrapFatalErr(err, identExpr.GetPosition())
			}

			err = symbol.SetObject(assignObj, stack)
			if err != nil {
				return nil, comn.WrapFatalErr(err, u.Expr.GetPosition())
			}
		case *runtime.GDAttrIdObject:
			err := expr.SetAttr(expr.Ident, assignObj)
			if err != nil {
				return nil, comn.WrapFatalErr(err, u.Expr.GetPosition())
			}
		}
	}

	return nil, nil
}

func (t *GDStaticAnalyzer) EvalLabel(l *ast.NodeLabel, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	return nil, nil
}

func (t *GDStaticAnalyzer) EvalIfElse(i *ast.NodeIfElse, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	evalIfNode := func(ifNode ast.Node, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
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

func (t *GDStaticAnalyzer) EvalTernaryIf(tIf *ast.NodeTernaryIf, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	ifObj, err := t.EvalNode(tIf.Expr, stack)
	if err != nil {
		return nil, err
	}

	err = runtime.EqualTypes(runtime.GDBoolType, ifObj.GetType(), stack)
	if err != nil {
		return nil, comn.WrapFatalErr(err, tIf.Expr.GetPosition())
	}

	thenObj, err := t.EvalNode(tIf.Then, stack)
	if err != nil {
		return nil, err
	}

	elseObj, err := t.EvalNode(tIf.Else, stack)
	if err != nil {
		return nil, err
	}

	typ := runtime.ComputeTypeFromObjects([]runtime.GDObject{thenObj, elseObj})

	obj, err := runtime.ZObjectForType(typ, stack)
	if err != nil {
		return nil, comn.WrapFatalErr(err, tIf.GetPosition())
	}

	tIf.SetInferredObject(obj)

	return obj, nil
}

func (t *GDStaticAnalyzer) EvalForIn(f *ast.NodeForIn, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	// Create a new stack for the for loop
	forStack := stack.NewSymbolStack(runtime.ForCtx)
	defer forStack.Dispose()

	// Evaluate the iterable expression inside the for loop
	exprObj, err := t.EvalNode(f.Expr, forStack)
	if err != nil {
		return nil, err
	}

	f.SetInferredType(exprObj.GetType())

	iterable, isIterable := runtime.Unwrap(exprObj).(runtime.GDIterableCollection)
	if !isIterable {
		return nil, comn.WrapFatalErr(runtime.InvalidIterableTypeErr(exprObj.GetType()), f.Expr.GetPosition())
	}

	// Type check for the sets
	nodeSets, isSets := f.Sets.(*ast.NodeSets)
	if !isSets {
		panic("expected a NodeSets")
	}

	// Resolve the sets, objects are assigned to the symbols stack
	_, err = t.EvalNode(nodeSets, forStack)
	if err != nil {
		return nil, err
	}

	updateSymbolType := func(set ast.Node, toType runtime.GDTypable, obj runtime.GDObject) (*ast.NodeSet, error) {
		if set, isSet := set.(*ast.NodeSet); isSet {
			ident := set.InferredIdent()
			symbol, err := forStack.GetSymbol(ident)
			if err != nil {
				return nil, err
			}

			// It is expected that the type of the symbol is the same as
			// the required type and it can be either a Int or the iterable type
			err = symbol.SetObject(obj, forStack)
			if err != nil {
				return nil, err
			}

			// Update the set according to the iterable type
			set.SetInferredType(toType)
			set.SetInferredObject(obj)

			return set, nil
		} else {
			panic("expected a NodeSet")
		}
	}

	sets := nodeSets.Nodes

	// An iterable object is created for the iterable type
	iterableZObj, err := runtime.ZObjectForType(iterable.GetIterableType(), forStack)
	if err != nil {
		return nil, comn.WrapFatalErr(err, f.Expr.GetPosition())
	}

	// Check set types, it must allow int and the iterable type
	sListLen := len(sets)
	switch {
	case sListLen == 1:
		set, err := updateSymbolType(sets[0], iterable.GetIterableType(), iterableZObj)
		if err != nil {
			return nil, comn.WrapFatalErr(err, sets[0].GetPosition())
		}

		f.InferredIterable = set
	case sListLen > 1:
		idxSet, err := updateSymbolType(sets[0], runtime.GDIntType, iterableZObj)
		if err != nil {
			return nil, comn.WrapFatalErr(err, sets[0].GetPosition())
		}
		f.InferredIndex = idxSet

		iterSet, err := updateSymbolType(sets[1], iterable.GetIterableType(), runtime.GDInt8(0))
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

func (t *GDStaticAnalyzer) EvalForIf(f *ast.NodeForIf, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	// Create a new stack for the for loop
	forStack := stack.NewSymbolStack(runtime.ForCtx)
	defer forStack.Dispose()

	if f.Sets != nil {
		_, err := t.EvalNode(f.Sets, forStack)
		if err != nil {
			return nil, err
		}
	}

	if f.Conds != nil {
		err := t.checkIfConditions(f.Conds, forStack)
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

func (t *GDStaticAnalyzer) EvalCollectableOp(c *ast.NodeMutCollectionOp, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	exprLObj, err := t.EvalNode(c.L, stack)
	if err != nil {
		return nil, err
	}

	exprRObj, err := t.EvalNode(c.R, stack)
	if err != nil {
		return nil, err
	}

	if mutCollection, isMutCollection := runtime.Unwrap(exprLObj).(runtime.GDMutableCollection); isMutCollection {
		switch c.Op {
		case ast.MutableCollectionAddOp:
			if ident, ok := c.L.(*ast.NodeIdent); ok {
				symbolId := runtime.NewGDStringIdent(ident.Lit)
				symbol, err := stack.GetSymbol(symbolId)
				if err != nil {
					return nil, comn.WrapFatalErr(err, ident.GetPosition())
				}

				err = runtime.CanBeAssign(mutCollection.GetIterableType(), exprRObj.GetType(), stack)
				if err != nil {
					return nil, comn.WrapFatalErr(err, c.GetPosition())
				}

				err = symbol.SetType(runtime.NewGDArrayType(exprRObj.GetType()), stack)
				if err != nil {
					return nil, comn.WrapFatalErr(err, c.GetPosition())
				}

				return exprRObj, nil
			}

			err := runtime.CanBeAssign(mutCollection.GetIterableType(), exprRObj.GetType(), stack)
			if err != nil {
				return nil, comn.WrapFatalErr(err, c.GetPosition())
			}

			return exprRObj, nil
		case ast.MutableCollectionRemoveOp:
			err := runtime.EqualTypes(runtime.GDIntType, exprRObj.GetType(), stack)
			if err != nil {
				return nil, comn.WrapFatalErr(err, c.GetPosition())
			}

			err = runtime.CanBeAssign(mutCollection.GetIterableType(), exprRObj.GetType(), stack)
			if err != nil {
				return nil, comn.WrapFatalErr(err, c.GetPosition())
			}

			zIterObj, err := runtime.ZObjectForType(mutCollection.GetIterableType(), stack)
			if err != nil {
				return nil, comn.WrapFatalErr(err, c.GetPosition())
			}

			return zIterObj, nil
		}
	}

	return nil, comn.WrapFatalErr(runtime.InvalidMutableCollectionTypeErr(exprLObj.GetType()), c.GetPosition())
}

func (t *GDStaticAnalyzer) EvalTypeAlias(ta *ast.NodeTypeAlias, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	ident := runtime.NewGDStringIdent(ta.Ident.Lit)
	_, err := stack.AddSymbol(ident, ta.IsPub, true, ta.Type, nil)
	if err != nil {
		return nil, comn.WrapFatalErr(err, ta.GetPosition())
	}

	return nil, nil
}

func (t *GDStaticAnalyzer) EvalCastExpr(c *ast.NodeCastExpr, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	exprObj, err := t.EvalNode(c.Expr, stack)
	if err != nil {
		return nil, err
	}

	castObj, err := exprObj.CastToType(c.Type, stack)
	if err != nil {
		return nil, comn.WrapFatalErr(err, c.GetPosition())
	}

	c.SetInferredType(c.Type)
	c.SetInferredObject(castObj)

	return castObj, nil
}

func (t *GDStaticAnalyzer) evalIfNode(i *ast.NodeIf, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	err := t.checkIfConditions(i.Conds, stack)
	if err != nil {
		return nil, err
	}

	// Block evaluation
	return t.evalBlock(i.Block, stack)
}

// Checks that all the conditions are of type bool
// and returns an error if any of them is not.
// if condition, ... { ... }
func (t *GDStaticAnalyzer) checkIfConditions(conditions []ast.Node, stack *runtime.GDSymbolStack) error {
	for _, cond := range conditions {
		condObj, err := t.EvalNode(cond, stack)
		if err != nil {
			return err
		}

		err = runtime.EqualTypes(condObj.GetType(), runtime.GDBoolType, stack)
		if err != nil {
			return comn.WrapFatalErr(err, cond.GetPosition())
		}
	}

	return nil
}

func (t *GDStaticAnalyzer) evalBlock(b *ast.NodeBlock, stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	funcStack := stack.NewSymbolStack(runtime.BlockCtx)
	defer funcStack.Dispose()

	for _, node := range b.Nodes {
		obj, err := t.EvalNode(node, funcStack)
		if err != nil {
			return nil, err
		}

		switch node := node.(type) {
		case *ast.NodeBreak:
			if b.Type == ast.FuncBlockType {
				return nil, comn.CompilerErr(comn.MisplassedBreakErrMsg, node.GetPosition())
			}
		case *ast.NodeReturn:
			// If not a flow control block, then return type is expected
			if b.Type != ast.ControlFlowBlockType {
				inferredType, err := runtime.InferType(b.ReturnType, obj.GetType(), stack)
				if err != nil {
					return nil, comn.WrapFatalErr(err, node.GetPosition())
				}

				node.SetInferredType(inferredType)
				if node.Expr != nil {
					node.Expr.SetInferredType(inferredType)
				}

				node.SetInferredObject(obj)

				continue
			}
		}
	}

	return nil, nil
}

func NewGDStaticAnalyzer(depAnalyzer *GDDepAnalyzer) *GDStaticAnalyzer {
	staticAnalyzer := &GDStaticAnalyzer{GDDepAnalyzer: depAnalyzer}
	staticAnalyzer.ObjectExpressionEvaluator = ObjectExpressionEvaluator{staticAnalyzer}

	return staticAnalyzer
}
