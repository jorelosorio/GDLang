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

// Node order
const (
	TypeAliasOrder uint16 = iota
	SetObjectOrder
	UpdateObjectOrder
	StructOrder
	PrivateFuncOrder
	PubFuncOrder
	EquivalentOrder
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
	Order() uint16
	GetNodeType() BaseNodeType
	SetParentNode(Node)
	GetParentNode() Node
	GetParentNodeByType(BaseNodeType) Node
	Inferable
}

type Inferable interface {
	InferredType() runtime.GDTypable
	SetInferredType(runtime.GDTypable)

	RuntimeType() runtime.GDTypable
	SetRuntimeType(runtime.GDTypable)

	RuntimeIdent() runtime.GDIdent
	SetRuntimeIdent(runtime.GDIdent)

	InferredObject() runtime.GDObject
	SetInferredObject(runtime.GDObject)

	InferredIdent() runtime.GDIdent
	SetInferredIdent(runtime.GDIdent)
}

type BaseNode struct {
	inferredType  runtime.GDTypable
	runtimeType   runtime.GDTypable
	inferredObj   runtime.GDObject
	inferredIdent runtime.GDIdent
	runtimeIdent  runtime.GDIdent
	nodeType      BaseNodeType
	parent        Node
}

// Inferable interface implementation

func (n *BaseNode) InferredType() runtime.GDTypable       { return n.inferredType }
func (n *BaseNode) SetInferredType(typ runtime.GDTypable) { n.inferredType = typ }

func (n *BaseNode) RuntimeType() runtime.GDTypable       { return n.runtimeType }
func (n *BaseNode) SetRuntimeType(typ runtime.GDTypable) { n.runtimeType = typ }

func (n *BaseNode) RuntimeIdent() runtime.GDIdent         { return n.runtimeIdent }
func (n *BaseNode) SetRuntimeIdent(ident runtime.GDIdent) { n.runtimeIdent = ident }

func (n *BaseNode) InferredObject() runtime.GDObject       { return n.inferredObj }
func (n *BaseNode) SetInferredObject(obj runtime.GDObject) { n.inferredObj = obj }

func (n *BaseNode) InferredIdent() runtime.GDIdent         { return n.inferredIdent }
func (n *BaseNode) SetInferredIdent(ident runtime.GDIdent) { n.inferredIdent = ident }

// Node interface implementation

func (n *BaseNode) GetNodeType() BaseNodeType { return n.nodeType }
func (n *BaseNode) SetType(t BaseNodeType)    { n.nodeType = t }
func (n *BaseNode) SetParentNode(p Node)      { n.parent = p }
func (n *BaseNode) GetParentNode() Node       { return n.parent }

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
