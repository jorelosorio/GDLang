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

type GDIRAGet struct {
	isNilSafe bool
	ident     runtime.GDIdent
	expr      GDIRNode
	GDIRBaseNode
}

func (i *GDIRAGet) BuildAssembly(padding string) string {
	return padding + fmt.Sprintf("%s%s %s %s", cpu.GetCPUInstName(cpu.AGet), tif(i.isNilSafe, " nilsafe", ""), i.expr.BuildAssembly(""), i.ident.ToString())
}

func (i *GDIRAGet) BuildBytecode(bytecode *bytes.Buffer, ctx *GDIRContext) error {
	ctx.AddMapping(bytecode, i.GetPosition())

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

func NewGDIRAIGet(ident runtime.GDIdent, isNilSafe bool, expr GDIRNode, node ast.Node) (*GDIRAGet, *GDIRObject) {
	return &GDIRAGet{isNilSafe, ident, expr, GDIRBaseNode{node}}, NewGDIRRegObject(cpu.RPop, node)
}
