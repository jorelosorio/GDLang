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

import "gdlang/src/gd/scanner"

// Shared expressions

type NodeSharedExpr struct {
	Expr             Node
	HasBeenProcessed bool
	BaseNode
}

func (e *NodeSharedExpr) GetPosition() scanner.Position { return e.Expr.GetPosition() }

func NewNodeSharedExpr(expr Node) *NodeSharedExpr {
	return &NodeSharedExpr{expr, false, BaseNode{}}
}

// Set object

type NodeSets struct {
	Nodes []Node
	BaseNode
}

func (d *NodeSets) GetPosition() scanner.Position { return GetStartEndPosition(d.Nodes) }

func NewNodeSets(nodes []Node) *NodeSets {
	return &NodeSets{nodes, BaseNode{}}
}

type NodeSet struct {
	IsPub         bool
	IsConst       bool
	IdentWithType *NodeIdentWithType
	Expr          Node
	Index         byte
	BaseNode
}

func (s *NodeSet) GetPosition() scanner.Position {
	return GetStartEndPosition([]Node{s.IdentWithType})
}

func NewNodeSet(isPub bool, isConst bool, identWithType *NodeIdentWithType, expr Node) *NodeSet {
	return &NodeSet{false, isConst, identWithType, expr, 0, BaseNode{}}
}

// Update object

type NodeUpdateSet struct {
	IdentExpr Node
	Expr      Node
	BaseNode
}

func (u *NodeUpdateSet) GetPosition() scanner.Position {
	return GetStartEndPosition([]Node{u.IdentExpr, u.Expr})
}

func NewNodeUpdateSet(identExpr Node, expr Node) *NodeUpdateSet {
	return &NodeUpdateSet{identExpr, expr, BaseNode{}}
}
