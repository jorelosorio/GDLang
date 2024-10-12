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

type BaseNodeType byte

const (
	NodeTypeFor BaseNodeType = iota
	NodeTypeIf
	NodeTypeBlock
	NodeTypeLambda
	NodeTypeFunc
)

type Node interface {
	GetPosition() scanner.Position
	GetNodeType() BaseNodeType
	SetParentNode(Node)
	GetParentNode() Node
	GetParentNodeByType(BaseNodeType) Node
	GetInference() *Inference
	SetInference(*Inference)
}

type Inference struct {
	Ident        runtime.GDIdent
	Type         runtime.GDTypable
	RuntimeIdent runtime.GDIdent
}

func (i *Inference) ToString() string {
	return i.Ident.ToString()
}

type BaseNode struct {
	nodeType BaseNodeType
	parent   Node
	*Inference
}

// Node interface implementation

func (n *BaseNode) GetNodeType() BaseNodeType { return n.nodeType }
func (n *BaseNode) SetType(t BaseNodeType)    { n.nodeType = t }
func (n *BaseNode) SetParentNode(p Node)      { n.parent = p }
func (n *BaseNode) GetParentNode() Node       { return n.parent }
func (n *BaseNode) SetInference(i *Inference) { n.Inference = i }
func (n *BaseNode) GetInference() *Inference  { return n.Inference }

// A function that look up for the parents and stop where a node type is found.
func (n *BaseNode) GetParentNodeByType(t BaseNodeType) Node {
	if n.parent == nil {
		return nil
	}

	if n.parent.GetNodeType() == t {
		return n.parent
	}

	return n.parent.GetParentNodeByType(t)
}
