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
	"gdlang/src/cpu"
	"gdlang/src/gd/scanner"
)

type GDIRILen struct {
	expr GDIRNode
	GDIRBaseNode
}

func (l *GDIRILen) BuildAssembly(padding string) string {
	return padding + fmt.Sprintf("%s %s", cpu.GetCPUInstName(cpu.ILen), l.expr.BuildAssembly(""))
}

func (l *GDIRILen) BuildBytecode(bytecode *bytes.Buffer, ctx *GDIRContext) error {
	ctx.AddMapping(bytecode, l.pos)

	err := Write(bytecode, cpu.ILen)
	if err != nil {
		return err
	}

	err = l.expr.BuildBytecode(bytecode, ctx)
	if err != nil {
		return err
	}

	return nil
}

func NewGDIRLen(expr GDIRNode, pos scanner.Position) (*GDIRILen, *GDIRObject) {
	return &GDIRILen{expr, GDIRBaseNode{pos}}, NewGDIRRegObject(cpu.RPop, pos)
}
