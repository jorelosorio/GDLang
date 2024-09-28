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
	"path"
)

// Token info

type NodeTokenInfo struct {
	scanner.Position
	Token scanner.Token
	Lit   string
	BaseNode
}

func (n *NodeTokenInfo) GetPosition() scanner.Position { return n.Position }

func NewNodeTokenInfo(token scanner.Token, pos scanner.Position, lit string) *NodeTokenInfo {
	return &NodeTokenInfo{pos, token, lit, BaseNode{}}
}

// File

type NodeFile struct {
	// Package declarations
	Packages []Node
	// Contains all the nodes statements
	Nodes []Node
	BaseNode
}

func (f *NodeFile) GetPosition() scanner.Position { return GetStartEndPosition(f.Nodes) }

func NewNodeFile(packages []Node, nodes []Node) *NodeFile {
	return &NodeFile{packages, nodes, BaseNode{}}
}

// Package

type NodePackageType byte

const (
	NodePackageSourced NodePackageType = iota
	NodePackageBuiltin
	NodePackageLib
)

type NodePackage struct {
	// Path of the package
	// Example: a.b.c
	PackagePath []Node
	// Public import references that are required from the package
	// Example: use a.b.c {object, ...}
	Imports []Node
	// Type of the package
	Type NodePackageType
	// Inferences
	InferredPath         string
	InferredAbsolutePath string
	BaseNode
}

func (p *NodePackage) GetPosition() scanner.Position { return GetStartEndPosition(p.PackagePath) }
func (p *NodePackage) GetName() string               { return p.PackagePath[len(p.PackagePath)-1].(*NodeIdent).Lit }
func (p *NodePackage) GetPath() string {
	var pkgPath string
	for _, ident := range p.PackagePath {
		pkgPath = path.Join(pkgPath, ident.(*NodeIdent).Lit)
	}

	return pkgPath
}

func NewNodePackage(packagePath []Node, imports []Node) *NodePackage {
	return &NodePackage{packagePath, imports, NodePackageSourced, "", "", BaseNode{}}
}

// Node for type definitions

type NodeTypeAlias struct {
	IsPub bool
	Ident *NodeIdent
	Type  runtime.GDTypable
	BaseNode
}

func (t *NodeTypeAlias) GetPosition() scanner.Position { return t.Ident.Position }

func NewNodeTypeAlias(isPub bool, ident *NodeIdent, identType runtime.GDTypable) *NodeTypeAlias {
	return &NodeTypeAlias{isPub, ident, identType, BaseNode{}}
}

// Ident with type

type NodeIdentWithType struct {
	Ident *NodeIdent
	Type  runtime.GDTypable
	BaseNode
}

func (t *NodeIdentWithType) GetPosition() scanner.Position { return t.Ident.Position }

func NewNodeIdentWithType(ident *NodeIdent, identType runtime.GDTypable) *NodeIdentWithType {
	return &NodeIdentWithType{ident, identType, BaseNode{}}
}

// Ident

type NodeIdent struct {
	*NodeTokenInfo
	BaseNode
}

func (i *NodeIdent) GetPosition() scanner.Position { return i.Position }

func NewNodeIdent(node *NodeTokenInfo) *NodeIdent {
	return &NodeIdent{node, BaseNode{}}
}

// Type

type NodeType struct {
	TypeTokenInfo *NodeTokenInfo
	BaseNode
}

func (t *NodeType) GetPosition() scanner.Position {
	return t.TypeTokenInfo.GetPosition()
}

func NewNodeType(typeTokenInfo *NodeTokenInfo) *NodeType {
	return &NodeType{typeTokenInfo, BaseNode{}}
}

// Nod cast expression

// Expr as Type
type NodeCastExpr struct {
	Expr Node
	Type runtime.GDTypable
	BaseNode
}

func (c *NodeCastExpr) GetPosition() scanner.Position { return c.Expr.GetPosition() }

func NewNodeCastExpr(expr Node, typ runtime.GDTypable) *NodeCastExpr {
	return &NodeCastExpr{expr, typ, BaseNode{}}
}

// Operator

type NodeExprOperation struct {
	Op   runtime.ExprOperationType
	L, R Node
	BaseNode
}

func (e *NodeExprOperation) GetPosition() scanner.Position {
	return GetStartEndPosition([]Node{e.L, e.R})
}

func NewNodeExprOperation(op runtime.ExprOperationType, l, r Node) *NodeExprOperation {
	return &NodeExprOperation{op, l, r, BaseNode{}}
}

// Ellipsis ...

type NodeEllipsisExpr struct {
	Expr Node
	BaseNode
}

func (e *NodeEllipsisExpr) GetPosition() scanner.Position { return e.Expr.GetPosition() }

func NewNodeEllipsisExpr(exp Node) Node {
	return &NodeEllipsisExpr{exp, BaseNode{}}
}

// Struct

type NodeStruct struct {
	Nodes []Node
	BaseNode
}

func (s *NodeStruct) GetPosition() scanner.Position { return GetStartEndPosition(s.Nodes) }

func NewNodeStruct(nodes ...Node) *NodeStruct {
	return &NodeStruct{nodes, BaseNode{}}
}

// Tuple

type NodeTuple struct {
	Nodes []Node
	BaseNode
}

func (t *NodeTuple) GetPosition() scanner.Position { return GetStartEndPosition(t.Nodes) }

func NewNodeTuple(exprs ...Node) *NodeTuple {
	return &NodeTuple{exprs, BaseNode{}}
}

// Collectable

type MutableCollectionOp byte

const (
	MutableCollectionAddOp    MutableCollectionOp = iota // Add a new element to the collection (<<)
	MutableCollectionRemoveOp                            // Remove an element from the collectable (>>)
)

// L >> R or L << R
type NodeMutCollectionOp struct {
	Op MutableCollectionOp
	L  Node
	R  Node
	BaseNode
}

func (c *NodeMutCollectionOp) GetPosition() scanner.Position {
	return GetStartEndPosition([]Node{c.L, c.R})
}

func NewNodeMutCollectionOp(op MutableCollectionOp, l, r Node) *NodeMutCollectionOp {
	return &NodeMutCollectionOp{op, l, r, BaseNode{}}
}

// Array

type NodeArray struct {
	// It is used to get the position of the array
	// when the array is empty.
	exprStart, exprEnd Node
	Nodes              []Node
	BaseNode
}

func (a *NodeArray) GetPosition() scanner.Position {
	if len(a.Nodes) == 0 {
		return GetStartEndPosition([]Node{a.exprStart, a.exprEnd})
	}

	return GetStartEndPosition(a.Nodes)
}

func NewNodeArray(exprStart, exprEnd Node, nodes []Node) *NodeArray {
	return &NodeArray{exprStart, exprEnd, nodes, BaseNode{}}
}

// Label

type NodeLabel struct {
	Ident NodeIdent
	BaseNode
}

func (l *NodeLabel) GetPosition() scanner.Position { return l.Ident.Position }

func NewNodeLabel(ident NodeIdent) *NodeLabel {
	return &NodeLabel{ident, BaseNode{}}
}

// Primary expressions

// expr(args)
type NodeCallExpr struct {
	Expr Node
	Args []Node
	BaseNode
}

func (s *NodeCallExpr) GetPosition() scanner.Position { return s.Expr.GetPosition() }

func NewNodeCallExpr(expr Node, args []Node) *NodeCallExpr {
	return &NodeCallExpr{expr, args, BaseNode{}}
}

// expr.?[IdxExpr]
type NodeIterIdxExpr struct {
	IsNilSafe bool
	Expr      Node
	IdxExpr   Node
	BaseNode
}

func (s *NodeIterIdxExpr) GetPosition() scanner.Position {
	return GetStartEndPosition([]Node{s.Expr, s.IdxExpr})
}

func NewNodeIterIdxExpr(isNilSafe bool, expr, idxExpr Node) *NodeIterIdxExpr {
	return &NodeIterIdxExpr{isNilSafe, expr, idxExpr, BaseNode{}}
}

// expr.?ident
type NodeSafeDotExpr struct {
	Expr      Node
	Ident     Node
	IsNilSafe bool
	BaseNode
}

func (s *NodeSafeDotExpr) GetPosition() scanner.Position { return s.Expr.GetPosition() }

func NewNodeSafeDotExpr(node Node, isNilSafe bool, ident Node) *NodeSafeDotExpr {
	return &NodeSafeDotExpr{node, ident, isNilSafe, BaseNode{}}
}
