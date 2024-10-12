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

// Atom

type NodeLiteral struct {
	*NodeTokenInfo
	InferredObject runtime.GDObject
	*BaseNode
}

func (a *NodeLiteral) GetPosition() scanner.Position { return a.Position }

func NewNodeLiteral(token *NodeTokenInfo) *NodeLiteral {
	return &NodeLiteral{token, nil, &BaseNode{}}
}

func NewNodeIntLiteral(val int, pos scanner.Position) *NodeLiteral {
	lit := NewNodeLiteral(&NodeTokenInfo{
		Token:    scanner.INT,
		Position: pos,
		Lit:      runtime.GDInt(val).ToString(),
	})

	lit.InferredObject = runtime.GDInt(val)

	return lit
}

func NewNodeNilLiteral(pos scanner.Position) *NodeLiteral {
	lit := NewNodeLiteral(&NodeTokenInfo{
		Token:    scanner.NIL,
		Position: pos,
		Lit:      "nil",
	})

	lit.InferredObject = runtime.GDZNil

	return lit
}
