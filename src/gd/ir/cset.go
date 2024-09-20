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

type GDIRICSet struct {
	idx  GDIRNode
	expr GDIRNode
	obj  GDIRNode
	GDIRBaseNode
}

func (i *GDIRICSet) BuildAssembly(padding string) string {
	return padding + fmt.Sprintf("%s %s %s %s", cpu.GetCPUInstName(cpu.CSet), i.idx.BuildAssembly(""), i.expr.BuildAssembly(""), i.obj.BuildAssembly(""))
}

func (i *GDIRICSet) BuildBytecode(bytecode *bytes.Buffer, ctx *GDIRContext) error {
	ctx.AddMapping(bytecode, i.pos)

	err := Write(bytecode, cpu.CSet)
	if err != nil {
		return err
	}

	err = i.idx.BuildBytecode(bytecode, ctx)
	if err != nil {
		return err
	}

	err = i.expr.BuildBytecode(bytecode, ctx)
	if err != nil {
		return err
	}

	err = i.obj.BuildBytecode(bytecode, ctx)
	if err != nil {
		return err
	}

	return nil
}

func NewGDIRISet(idx GDIRNode, expr GDIRNode, obj GDIRNode, pos scanner.Position) *GDIRICSet {
	return &GDIRICSet{idx, expr, obj, GDIRBaseNode{pos}}
}
