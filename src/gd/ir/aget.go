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
	"gdlang/src/gd/scanner"
)

type GDIRAGet struct {
	ident     runtime.GDIdentType
	isNilSafe bool
	expr      GDIRNode
	GDIRBaseNode
}

func (i *GDIRAGet) BuildAssembly(padding string) string {
	return padding + fmt.Sprintf("%s %s %s", cpu.GetCPUInstName(cpu.AGet), i.expr.BuildAssembly(""), IRTypeToString(i.ident))
}

func (i *GDIRAGet) BuildBytecode(bytecode *bytes.Buffer, ctx *GDIRContext) error {
	ctx.AddMapping(bytecode, i.pos)

	err := Write(bytecode, cpu.AGet, i.isNilSafe)
	if err != nil {
		return err
	}

	err = i.expr.BuildBytecode(bytecode, ctx)
	if err != nil {
		return err
	}

	err = Write(bytecode, i.ident)
	if err != nil {
		return err
	}

	return nil
}

func NewGDIRAIGet(ident runtime.GDIdentType, isNilSafe bool, expr GDIRNode, pos scanner.Position) (*GDIRAGet, *GDIRObject) {
	return &GDIRAGet{ident, isNilSafe, expr, GDIRBaseNode{pos}}, NewGDIRReg(cpu.RPop, pos)
}
