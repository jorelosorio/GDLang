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

// Block

type NodeBlockType byte

const (
	// A block that is used for lambdas.
	// It allows the use of return statements.
	// For instance:
	// pub func main() {
	//     set b = func () => int {
	//         return 1
	//     }
	// }
	FuncBlockType NodeBlockType = iota
	// A block that is used for control flow statements
	// with no declared return type.
	ControlFlowBlockType
)

type NodeBlock struct {
	Type       NodeBlockType
	Nodes      []Node
	ReturnType runtime.GDTypable
	*BaseNode
}

func (n *NodeBlock) GetPosition() scanner.Position { return GetStartEndPosition(n.Nodes) }

func (n *NodeBlock) SetAsControlFlowBlock() {
	n.Type = ControlFlowBlockType
	n.ReturnType = nil
}

func (n *NodeBlock) SetAsFuncBlock(retType runtime.GDTypable) {
	n.Type = FuncBlockType
	n.ReturnType = retType
}

func NewNodeBlock(nodes []Node) *NodeBlock {
	nodeBlock := &NodeBlock{Nodes: nodes, BaseNode: &BaseNode{nodeType: NodeTypeBlock}}
	for _, node := range nodes {
		node.SetParentNode(nodeBlock)
	}

	return nodeBlock
}
