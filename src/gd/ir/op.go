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

type GDIROp struct {
	op          runtime.ExprOperationType
	left, right GDIRNode
	GDIRBaseNode
}

func (e *GDIROp) BuildAssembly(padding string) string {
	rExpr := e.right
	if e.right == nil {
		rExpr = NewGDIRObject(runtime.GDZNil, e.pos)
	}

	return padding + fmt.Sprintf("%s %s %s %s", cpu.GetCPUInstName(cpu.Operation), runtime.ExprOperationMap[e.op], e.left.BuildAssembly(""), rExpr.BuildAssembly(""))
}

func (e *GDIROp) BuildBytecode(bytecode *bytes.Buffer, ctx *GDIRContext) error {
	ctx.AddMapping(bytecode, e.pos)

	err := Write(bytecode, cpu.Operation)
	if err != nil {
		return err
	}

	// Write operation type
	err = WriteByte(bytecode, byte(e.op))
	if err != nil {
		return err
	}

	if e.right != nil {
		err = e.right.BuildBytecode(bytecode, ctx)
		if err != nil {
			return err
		}
	} else {
		err = Write(bytecode, runtime.GDZNil)
		if err != nil {
			return err
		}
	}

	err = e.left.BuildBytecode(bytecode, ctx)
	if err != nil {
		return err
	}

	return nil
}

func NewGDIROp(op runtime.ExprOperationType, left, right GDIRNode, pos scanner.Position) (*GDIROp, *GDIRObject) {
	return &GDIROp{op, left, right, GDIRBaseNode{pos}}, NewGDIRReg(cpu.RPop, pos)
}
