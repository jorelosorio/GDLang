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

type GDIRJump struct {
	label runtime.GDIdent
	GDIRBaseNode
}

func (j *GDIRJump) BuildAssembly(padding string) string {
	return padding + fmt.Sprintf("%s %s", cpu.GetCPUInstName(cpu.Jump), j.label.ToString())
}

func (j *GDIRJump) BuildBytecode(bytecode *bytes.Buffer, ctx *GDIRContext) error {
	ctx.AddMapping(bytecode, j.pos)

	err := Write(bytecode, cpu.Jump)
	if err != nil {
		return err
	}

	// Current offset
	offset := bytecode.Len()

	// Write space for the label offset
	err = WriteUInt16(bytecode, 0)
	if err != nil {
		return err
	}

	// Add the mark to wait for the label when is defined
	_ = ctx.AddMark(bytecode, offset, j.label)

	return nil
}

func NewGDIRJump(label runtime.GDIdent, pos scanner.Position) *GDIRJump {
	return &GDIRJump{label, GDIRBaseNode{pos}}
}
