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

type GDIRTypeAlias struct {
	disc *GDIRDiscoverable
	typ  runtime.GDTypable
	GDIRBaseNode
}

func (t *GDIRTypeAlias) BuildAssembly(padding string) string {
	return padding + fmt.Sprintf("%s %s %s", cpu.GetCPUInstName(cpu.TypeAlias), t.disc.IR(""), IRTypeToString(t.typ))
}

func (t *GDIRTypeAlias) BuildBytecode(bytecode *bytes.Buffer, ctx *GDIRContext) error {
	ctx.AddMapping(bytecode, t.GetPosition())

	err := Write(bytecode, cpu.TypeAlias)
	if err != nil {
		return err
	}

	err = t.disc.Bytecode(bytecode, ctx)
	if err != nil {
		return err
	}

	err = Write(bytecode, t.typ)
	if err != nil {
		return err
	}

	return nil
}

func NewGDIRTypeAlias(disc *GDIRDiscoverable, typ runtime.GDTypable, node ast.Node) *GDIRTypeAlias {
	return &GDIRTypeAlias{disc, typ, GDIRBaseNode{node}}
}
