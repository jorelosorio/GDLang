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
	"gdlang/src/gd/ast"
	"gdlang/src/gd/scanner"
)

type GDIRCOp struct {
	op    ast.MutableCollectionOp
	left  GDIRNode
	right GDIRNode
	GDIRBaseNode
}

func (a *GDIRCOp) BuildAssembly(padding string) string {
	var inst cpu.GDInst
	switch a.op {
	case ast.MutableCollectionAddOp:
		inst = cpu.CAdd
	case ast.MutableCollectionRemoveOp:
		inst = cpu.CRemove
	}

	return padding + fmt.Sprintf("%s %s %s", cpu.GetCPUInstName(inst), a.left.BuildAssembly(""), a.right.BuildAssembly(""))
}

func (a *GDIRCOp) BuildBytecode(bytecode *bytes.Buffer, ctx *GDIRContext) error {
	ctx.AddMapping(bytecode, a.pos)

	var inst cpu.GDInst
	switch a.op {
	case ast.MutableCollectionAddOp:
		inst = cpu.CAdd
	case ast.MutableCollectionRemoveOp:
		inst = cpu.CRemove
	}

	err := Write(bytecode, inst)
	if err != nil {
		return err
	}

	err = a.left.BuildBytecode(bytecode, ctx)
	if err != nil {
		return err
	}

	err = a.right.BuildBytecode(bytecode, ctx)
	if err != nil {
		return err
	}

	return nil
}

func NewGDIRCOp(op ast.MutableCollectionOp, left GDIRNode, right GDIRNode, pos scanner.Position) (*GDIRCOp, *GDIRObject) {
	return &GDIRCOp{op, left, right, GDIRBaseNode{pos}}, NewGDIRReg(cpu.RPop, pos)
}
