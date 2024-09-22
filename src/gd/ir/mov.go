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

type GDIRMov struct {
	target GDIRNode
	value  GDIRNode
	GDIRBaseNode
}

func (m *GDIRMov) BuildAssembly(padding string) string {
	return padding + fmt.Sprintf("%s %s %s", cpu.GetCPUInstName(cpu.Mov), m.target.BuildAssembly(""), m.value.BuildAssembly(""))
}

func (m *GDIRMov) BuildBytecode(bytecode *bytes.Buffer, ctx *GDIRContext) error {
	ctx.AddMapping(bytecode, m.pos)

	err := Write(bytecode, cpu.Mov)
	if err != nil {
		return err
	}

	err = m.target.BuildBytecode(bytecode, ctx)
	if err != nil {
		return err
	}

	err = m.value.BuildBytecode(bytecode, ctx)
	if err != nil {
		return err
	}

	return nil
}

func NewGDIRMov(target GDIRNode, value GDIRNode, pos scanner.Position) *GDIRMov {
	return &GDIRMov{target, value, GDIRBaseNode{pos}}
}
