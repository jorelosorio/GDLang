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
)

type GDIRUse struct {
	ident   runtime.GDIdent
	mode    runtime.GDPackageMode
	imports []runtime.GDIdent
	GDIRBaseNode
}

func (u *GDIRUse) BuildAssembly(padding string) string {
	imports := runtime.JoinSlice(u.imports, func(importIdent runtime.GDIdent, _ int) string {
		return importIdent.ToString()
	}, ", ")
	return fmt.Sprintf("%s%s %s %s %s", padding, cpu.GetCPUInstName(cpu.Use), runtime.PackageModeMap[u.mode], u.ident.ToString(), imports)
}

func (u *GDIRUse) BuildBytecode(bytecode *bytes.Buffer, ctx *GDIRContext) error {
	ctx.AddMapping(bytecode, u.GetPosition())

	err := Write(bytecode, cpu.Use)
	if err != nil {
		return err
	}

	err = WriteByte(bytecode, byte(u.mode))
	if err != nil {
		return err
	}

	err = WriteIdent(bytecode, u.ident)
	if err != nil {
		return err
	}

	err = WriteByte(bytecode, byte(len(u.imports)))
	if err != nil {
		return err
	}

	for _, ident := range u.imports {
		err = WriteIdent(bytecode, ident)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewGDIRUse(mode runtime.GDPackageMode, ident runtime.GDIdent, imports []runtime.GDIdent) *GDIRUse {
	return &GDIRUse{ident, mode, imports, GDIRBaseNode{}}
}
