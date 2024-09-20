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

type GDIRLambda struct {
	typ *runtime.GDLambdaType
	pos scanner.Position
	*GDIRBlock
}

func (l *GDIRLambda) BuildAssembly(padding string) string {
	return padding + fmt.Sprintf("lambda %s\n%s", IRTypeToString(l.typ), l.GDIRBlock.BuildAssembly(padding))
}

func (l *GDIRLambda) BuildBytecode(bytecode *bytes.Buffer, ctx *GDIRContext) error {
	// Write type
	err := Write(bytecode, cpu.Lambda, l.typ)
	if err != nil {
		return err
	}

	// Write block
	err = l.GDIRBlock.BuildBytecode(bytecode, ctx)
	if err != nil {
		return err
	}

	return nil
}

func NewGDIRLambda(typ *runtime.GDLambdaType, pos scanner.Position) (*GDIRLambda, *GDIRObject) {
	return &GDIRLambda{typ, pos, NewGDIRBlock()}, NewGDIRReg(cpu.RPop, scanner.ZeroPos)
}
