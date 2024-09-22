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

type GDIRObject struct {
	Type runtime.GDTypable
	Obj  interface{}
	GDIRBaseNode
}

func (o *GDIRObject) BuildAssembly(padding string) string {
	switch value := o.Obj.(type) {
	case runtime.GDObject:
		return IRObjectWithTypeToString(value)
	case GDIRNode:
		node := value.BuildAssembly("")
		return fmt.Sprintf("(%s: %s)", IRTypeToString(o.Type), node)
	case []GDIRNode:
		nodes := runtime.JoinSlice(value, func(node GDIRNode, _ int) string {
			return node.BuildAssembly("")
		}, ", ")

		return fmt.Sprintf("(%s: [%s])", IRTypeToString(o.Type), nodes)
	}

	return padding
}

func (o *GDIRObject) BuildBytecode(bytecode *bytes.Buffer, ctx *GDIRContext) error {
	ctx.AddMapping(bytecode, o.pos)

	switch obj := o.Obj.(type) {
	case runtime.GDObject:
		err := Write(bytecode, obj)
		if err != nil {
			return err
		}

		return nil
	case GDIRNode:
		// Write the type
		err := WriteType(bytecode, o.Type)
		if err != nil {
			return err
		}

		err = obj.BuildBytecode(bytecode, ctx)
		if err != nil {
			return err
		}
	case []GDIRNode:
		// Write the type
		err := WriteType(bytecode, o.Type)
		if err != nil {
			return err
		}

		// Write the length
		err = WriteByte(bytecode, byte(len(obj)))
		if err != nil {
			return err
		}

		// Write the nodes
		for i := len(obj) - 1; i >= 0; i-- {
			node := obj[i]
			err := node.BuildBytecode(bytecode, ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func NewGDIRObject(obj runtime.GDObject, pos scanner.Position) *GDIRObject {
	return &GDIRObject{obj.GetType(), obj, GDIRBaseNode{pos}}
}

func NewGDIRRegObject(reg cpu.GDReg, pos scanner.Position) *GDIRObject {
	ident := runtime.NewGDObjRefType(runtime.NewGDByteIdent(byte(reg)))
	identObj := runtime.NewGDIdObject(ident, runtime.GDString(cpu.GetCPURegName(reg)))
	return NewGDIRObject(identObj, pos)
}

func NewGDIRIterableObject(typ runtime.GDTypable, values []GDIRNode, pos scanner.Position) *GDIRObject {
	return &GDIRObject{typ, values, GDIRBaseNode{pos}}
}

// Used for special cases where the type wraps an object
// like the spreadable type.
func NewGDIRObjectWithType(typ runtime.GDTypable, obj interface{}, pos scanner.Position) *GDIRObject {
	return &GDIRObject{typ, obj, GDIRBaseNode{pos}}
}
