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

package ir

import (
	"bytes"
	"gdlang/src/gd/ast"
	"gdlang/src/gd/scanner"
)

type GDIRNode interface {
	BuildAssembly(string) string
	BuildBytecode(bytecode *bytes.Buffer, ctx *GDIRContext) error
	GetPosition() scanner.Position
}

type GDIRStackNode interface {
	GDIRNode
	AddNode(node ...GDIRNode)
}

type GDIRBaseNode struct {
	ast.Node
}

func (b *GDIRBaseNode) BuildAssembly(padding string) string                          { return "" }
func (b *GDIRBaseNode) BuildBytecode(bytecode *bytes.Buffer, ctx *GDIRContext) error { return nil }
func (b *GDIRBaseNode) GetPosition() scanner.Position {
	if b.Node != nil {
		return b.Node.GetPosition()
	}

	// There are some cases like in the test where an AST node is not specified,
	// as well, as some functions made by the compiler don't have a node.
	return scanner.ZeroPos
}
