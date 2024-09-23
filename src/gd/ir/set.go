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

type GDIRSet struct {
	disc *GDIRDiscoverable
	typ  runtime.GDTypable
	expr GDIRNode
	GDIRBaseNode
}

func (s *GDIRSet) BuildAssembly(padding string) string {
	return padding + fmt.Sprintf("%s %s %s %s", cpu.GetCPUInstName(cpu.Set), s.disc.IR(""), IRTypeToString(s.typ), s.expr.BuildAssembly(padding))
}

func (s *GDIRSet) BuildBytecode(bytecode *bytes.Buffer, ctx *GDIRContext) error {
	ctx.AddMapping(bytecode, s.GetPosition())

	err := Write(bytecode, cpu.Set)
	if err != nil {
		return err
	}

	err = s.disc.Bytecode(bytecode, ctx)
	if err != nil {
		return err
	}

	err = WriteType(bytecode, s.typ)
	if err != nil {
		return err
	}

	err = s.expr.BuildBytecode(bytecode, ctx)
	if err != nil {
		return err
	}

	return nil
}

func NewGDIRSet(disc *GDIRDiscoverable, typ runtime.GDTypable, expr GDIRNode, node ast.Node) *GDIRSet {
	return &GDIRSet{disc, typ, expr, GDIRBaseNode{node}}
}
