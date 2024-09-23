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
	"fmt"
	"gdlang/lib/runtime"
	"gdlang/src/cpu"
	"gdlang/src/gd/ast"
)

type GDIRLabel struct {
	name runtime.GDIdent
	GDIRBaseNode
}

func (l *GDIRLabel) BuildAssembly(padding string) string {
	spacing := padding
	if len(padding)-3 > 0 {
		spacing = padding[:len(padding)-3]
	}
	return spacing + fmt.Sprintf("=> %s %s:", cpu.GetCPUInstName(cpu.Label), l.name.ToString())
}

func (l *GDIRLabel) BuildBytecode(bytecode *bytes.Buffer, ctx *GDIRContext) error {
	ctx.AddMapping(bytecode, l.GetPosition())

	ctx.AddLabel(bytecode, uint16(bytecode.Len()), l.name)

	return nil
}

func NewGDIRLabel(name runtime.GDIdent, node ast.Node) *GDIRLabel {
	return &GDIRLabel{name, GDIRBaseNode{node}}
}
