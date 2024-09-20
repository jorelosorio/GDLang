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

package ast

import (
	"gdlang/lib/runtime"
	"gdlang/src/gd/scanner"
)

// Lambda

type NodeLambda struct {
	Type  *runtime.GDLambdaType
	Block *NodeBlock
	BaseNode
}

func (l *NodeLambda) GetPosition() scanner.Position { return l.Block.GetPosition() }

func (l *NodeLambda) Order() uint16 { return EquivalentOrder }

func NewNodeLambda(funcType *runtime.GDLambdaType, block *NodeBlock) *NodeLambda {
	block.SetAsFuncBlock(funcType.ReturnType)

	nodeLambda := &NodeLambda{funcType, block, BaseNode{nodeType: NodeTypeLambda}}
	block.SetParentNode(nodeLambda)

	return nodeLambda
}

// Function

type NodeFunc struct {
	IsPub bool
	Ident *NodeIdent
	*NodeLambda
	BaseNode
}

func (f *NodeFunc) GetPosition() scanner.Position {
	return GetStartEndPosition([]Node{f.Ident})
}

func (f *NodeFunc) Order() uint16 {
	if f.IsPub {
		return PubFuncOrder
	}

	return PrivateFuncOrder
}

func NewNodeFunc(isPublic bool, ident *NodeIdent, funcType *runtime.GDLambdaType, block *NodeBlock) *NodeFunc {
	nodeFunc := &NodeFunc{isPublic, ident, NewNodeLambda(funcType, block), BaseNode{nodeType: NodeTypeFunc}}
	block.SetParentNode(nodeFunc)

	return nodeFunc
}

type NodeStructAttr struct {
	Ident *NodeIdent
	Expr  Node
	BaseNode
}

func (s *NodeStructAttr) GetPosition() scanner.Position {
	return GetStartEndPosition([]Node{s.Ident, s.Expr})
}
func (s *NodeStructAttr) Order() uint16 { return EquivalentOrder }

func NewNodeStructAttr(ident *NodeIdent, expr Node) *NodeStructAttr {
	return &NodeStructAttr{ident, expr, BaseNode{}}
}
