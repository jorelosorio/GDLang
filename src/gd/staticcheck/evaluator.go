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
	"fmt"
	"gdlang/src/gd/ast"
)

// Evaluator is an interface that traverses the analyzed AST, and
// it is used for static analysis and code generation.
type Evaluator[T interface{}, E interface{}] interface {
	EvalAtom(a *ast.NodeLiteral, stack E) (T, error)
	EvalIdent(i *ast.NodeIdent, stack E) (T, error)
	EvalFunc(f *ast.NodeFunc, stack E) (T, error)
	EvalLambda(l *ast.NodeLambda, stack E) (T, error)
	EvalExprOp(e *ast.NodeExprOperation, stack E) (T, error)
	EvalExpEllipsis(e *ast.NodeEllipsisExpr, stack E) (T, error)
	EvalTuple(t *ast.NodeTuple, stack E) (T, error)
	EvalStruct(s *ast.NodeStruct, stack E) (T, error)
	EvalArray(a *ast.NodeArray, stack E) (T, error)
	EvalReturn(r *ast.NodeReturn, stack E) (T, error)
	EvalIterIdxExpr(a *ast.NodeIndexableExpr, stack E) (T, error)
	EvalCallExpr(c *ast.NodeCallExpr, stack E) (T, error)
	EvalSafeDotExpr(s *ast.NodeSafeDotExpr, stack E) (T, error)
	EvalSets(s *ast.NodeSets, stack E) (T, error)
	EvalSet(s *ast.NodeSet, stack E) (T, error)
	EvalUpdateSet(u *ast.NodeUpdateSet, stack E) (T, error)
	EvalLabel(l *ast.NodeLabel, stack E) (T, error)
	EvalIfElse(i *ast.NodeIfElse, stack E) (T, error)
	EvalTernaryIf(t *ast.NodeTernaryIf, stack E) (T, error)
	EvalForIn(f *ast.NodeForIn, stack E) (T, error)
	EvalForIf(f *ast.NodeForIf, stack E) (T, error)
	EvalCollectableOp(c *ast.NodeMutCollectionOp, stack E) (T, error)
	EvalTypeAlias(t *ast.NodeTypeAlias, stack E) (T, error)
	EvalCastExpr(c *ast.NodeCastExpr, stack E) (T, error)
	EvalPackage(p *ast.NodePackage, stack E) (T, error)
}

type ExpressionEvaluator[T interface{}, E interface{}] struct{ Evaluator[T, E] }

func (e *ExpressionEvaluator[T, E]) EvalFileNodes(nodes []ast.Node, stack E) error {
	for _, node := range nodes {
		_, err := e.EvalNode(node, stack)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *ExpressionEvaluator[T, E]) EvalNode(node ast.Node, stack E) (T, error) {
	var zeroT T
	switch node := node.(type) {
	case *ast.NodeLiteral:
		return e.EvalAtom(node, stack)
	case *ast.NodeIdent:
		return e.EvalIdent(node, stack)
	case *ast.NodeLambda:
		return e.EvalLambda(node, stack)
	case *ast.NodeExprOperation:
		return e.EvalExprOp(node, stack)
	case *ast.NodeEllipsisExpr:
		return e.EvalExpEllipsis(node, stack)
	case *ast.NodeFunc:
		return e.EvalFunc(node, stack)
	case *ast.NodeTuple:
		return e.EvalTuple(node, stack)
	case *ast.NodeStruct:
		return e.EvalStruct(node, stack)
	case *ast.NodeArray:
		return e.EvalArray(node, stack)
	case *ast.NodeReturn:
		return e.EvalReturn(node, stack)
	case *ast.NodeSafeDotExpr:
		return e.EvalSafeDotExpr(node, stack)
	case *ast.NodeSets:
		return e.EvalSets(node, stack)
	case *ast.NodeSet:
		return e.EvalSet(node, stack)
	case *ast.NodeCallExpr:
		return e.EvalCallExpr(node, stack)
	case *ast.NodeIndexableExpr:
		return e.EvalIterIdxExpr(node, stack)
	case *ast.NodeTypeAlias:
		return e.EvalTypeAlias(node, stack)
	case *ast.NodeCastExpr:
		return e.EvalCastExpr(node, stack)
	case *ast.NodeUpdateSet:
		return e.EvalUpdateSet(node, stack)
	case *ast.NodeLabel:
		return e.EvalLabel(node, stack)
	case *ast.NodeIfElse:
		return e.EvalIfElse(node, stack)
	case *ast.NodeTernaryIf:
		return e.EvalTernaryIf(node, stack)
	case *ast.NodeForIn:
		return e.EvalForIn(node, stack)
	case *ast.NodeForIf:
		return e.EvalForIf(node, stack)
	case *ast.NodeMutCollectionOp:
		return e.EvalCollectableOp(node, stack)
	case *ast.NodeBreak:
		return zeroT, nil
	case *ast.NodePackage:
		return e.EvalPackage(node, stack)
	}

	panic(fmt.Errorf("unhandled node type: %T", node))
}
