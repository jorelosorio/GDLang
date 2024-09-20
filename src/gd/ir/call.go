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

type GDIRCall struct {
	expr GDIRNode
	args GDIRNode
	GDIRBaseNode
}

func (c *GDIRCall) BuildAssembly(padding string) string {
	return padding + fmt.Sprintf("%s %s %s", cpu.GetCPUInstName(cpu.Call), c.expr.BuildAssembly(padding), c.args.BuildAssembly(""))
}

func (c *GDIRCall) BuildBytecode(bytecode *bytes.Buffer, ctx *GDIRContext) error {
	ctx.AddMapping(bytecode, c.pos)

	err := Write(bytecode, cpu.Call)
	if err != nil {
		return err
	}

	// Expr that defines the function
	err = c.expr.BuildBytecode(bytecode, ctx)
	if err != nil {
		return err
	}

	// Args
	err = c.args.BuildBytecode(bytecode, ctx)
	if err != nil {
		return err
	}

	return nil
}

func NewGDIRCall(expr GDIRNode, args GDIRNode, pos scanner.Position) (*GDIRCall, *GDIRObject) {
	return &GDIRCall{expr, args, GDIRBaseNode{pos}}, NewGDIRReg(cpu.RPop, pos)
}
