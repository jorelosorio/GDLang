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

// For

type NodeFor interface {
	SetEndLabel(runtime.GDIdent)
	GetEndLabel() runtime.GDIdent
}

// Nod For If

type NodeForIf struct {
	Sets  Node
	Conds []Node
	Block *NodeBlock

	// Labeling
	endLabel runtime.GDIdent

	BaseNode
}

func (f *NodeForIf) GetPosition() scanner.Position {
	return f.Block.GetPosition()
}
func (f *NodeForIf) Order() uint16 { return EquivalentOrder }

func (f *NodeForIf) SetEndLabel(endLabel runtime.GDIdent) { f.endLabel = endLabel }
func (f *NodeForIf) GetEndLabel() runtime.GDIdent         { return f.endLabel }

func NewNodeForIf(setObjs Node, ifConds []Node, block *NodeBlock) *NodeForIf {
	block.SetAsControlFlowBlock()
	nodeFor := &NodeForIf{setObjs, ifConds, block, nil, BaseNode{nodeType: NodeTypeFor}}
	block.SetParentNode(nodeFor)

	return nodeFor
}

// For In

type NodeForIn struct {
	Expr Node

	InferredIndex    *NodeSet
	InferredIterable *NodeSet

	*NodeForIf
}

func (f *NodeForIn) GetPosition() scanner.Position {
	return GetStartEndPosition([]Node{f.Sets, f.Expr, f.Block})
}
func (f *NodeForIn) Order() uint16 { return EquivalentOrder }

func NewNodeForIn(setObjs Node, expr Node, block *NodeBlock) Node {
	block.SetAsControlFlowBlock()
	nodeForIf := NewNodeForIf(
		setObjs,
		[]Node{},
		block,
	)
	block.SetParentNode(nodeForIf)

	forIn := &NodeForIn{expr, nil, nil, nodeForIf}
	nodeForIf.SetParentNode(forIn)

	return forIn
}

// Break

type NodeBreak struct {
	*NodeTokenInfo
	BaseNode
}

func (b *NodeBreak) GetPosition() scanner.Position { return b.Position }
func (b *NodeBreak) Order() uint16                 { return EquivalentOrder }

func NewNodeBreak(ident *NodeTokenInfo) *NodeBreak {
	return &NodeBreak{ident, BaseNode{nodeType: NodeTypeFor}}
}
