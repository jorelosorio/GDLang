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

// Nod ExprIf

type NodeTernaryIf struct {
	Expr Node
	Then Node
	Else Node
	BaseNode
}

func (e *NodeTernaryIf) GetPosition() scanner.Position {
	return GetStartEndPosition([]Node{e.Expr, e.Then, e.Else})
}
func (e *NodeTernaryIf) Order() uint16 { return EquivalentOrder }

func NewNodeTernaryIf(expr Node, then Node, elseNode Node) *NodeTernaryIf {
	return &NodeTernaryIf{expr, then, elseNode, BaseNode{nodeType: NodeTypeIf}}
}

// If

type NodeIf struct {
	Conditions []Node
	Block      *NodeBlock
	Ident      runtime.GDIdent
	BaseNode
}

func (i *NodeIf) GetPosition() scanner.Position {
	return GetStartEndPosition(append(i.Conditions, i.Block))
}
func (i *NodeIf) Order() uint16 { return EquivalentOrder }

func NewNodeIf(ifConditions []Node, block *NodeBlock) *NodeIf {
	nodeIf := &NodeIf{ifConditions, block, nil, BaseNode{nodeType: NodeTypeIf}}
	// Set the block as a control flow block if it is not nil.
	// Nil block it is used of ternary if and conditions validation
	if block != nil {
		block.SetParentNode(nodeIf)
		block.SetAsControlFlowBlock()
	}

	return nodeIf
}

type NodeIfElse struct {
	If     Node
	ElseIf []Node
	Else   Node
	BaseNode
}

func (i *NodeIfElse) GetPosition() scanner.Position {
	return GetStartEndPosition([]Node{i.If, i.Else})
}
func (i *NodeIfElse) Order() uint16 { return EquivalentOrder }

func NewNodeIfElse(nodIf Node, nodIfs []Node, nodeElse Node) *NodeIfElse {
	nodeElseIf := &NodeIfElse{nodIf, nodIfs, nodeElse, BaseNode{nodeType: NodeTypeIf}}
	// Set parent node for the main if
	nodIf.SetParentNode(nodeElseIf)
	// Set parent node for the else if
	for _, nodeIf := range nodIfs {
		nodeIf.SetParentNode(nodeElseIf)
	}
	// Set parent node for the else
	if nodeElse != nil {
		nodeElse.SetParentNode(nodeElseIf)
	}
	return nodeElseIf
}
