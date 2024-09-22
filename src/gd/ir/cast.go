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

type GDIRCastObject struct {
	Type runtime.GDTypable
	Expr GDIRNode
	GDIRBaseNode
}

func (c *GDIRCastObject) BuildAssembly(padding string) string {
	return padding + fmt.Sprintf("%s %s %s", cpu.GetCPUInstName(cpu.CastObj), IRTypeToString(c.Type), c.Expr.BuildAssembly(""))
}

func (c *GDIRCastObject) BuildBytecode(bytecode *bytes.Buffer, ctx *GDIRContext) error {
	ctx.AddMapping(bytecode, c.pos)

	err := Write(bytecode, cpu.CastObj, c.Type)
	if err != nil {
		return err
	}

	err = c.Expr.BuildBytecode(bytecode, ctx)
	if err != nil {
		return err
	}

	return nil
}

func NewGDIRCastObject(typ runtime.GDTypable, expr GDIRNode, pos scanner.Position) (*GDIRCastObject, *GDIRObject) {
	return &GDIRCastObject{typ, expr, GDIRBaseNode{pos}}, NewGDIRRegObject(cpu.RPop, pos)
}
