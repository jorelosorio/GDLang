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

type GDIRObject struct {
	Type  runtime.GDTypable
	Value any
	GDIRBaseNode
}

func (o *GDIRObject) BuildAssembly(padding string) string {
	switch {
	case o.Type != nil && o.Value == nil:
		return IRTypeToString(o.Type)
	case o.Value != nil:
		switch val := o.Value.(type) {
		case runtime.GDObject:
			return IRObjectToString(val)
		case GDIRNode:
			node := val.BuildAssembly("")
			return fmt.Sprintf("(%s: %s)", IRTypeToString(o.Type), node)
		case []GDIRNode:
			nodes := runtime.JoinSlice(val, func(node GDIRNode, _ int) string {
				return node.BuildAssembly("")
			}, ", ")

			return fmt.Sprintf("(%s: [%s])", IRTypeToString(o.Type), nodes)
		}
	}

	return ""
}

func (o *GDIRObject) BuildBytecode(bytecode *bytes.Buffer, ctx *GDIRContext) error {
	ctx.AddMapping(bytecode, o.GetPosition())

	switch {
	case o.Type != nil && o.Value == nil:
		err := WriteType(bytecode, o.Type)
		if err != nil {
			return err
		}
	case o.Value != nil:
		switch obj := o.Value.(type) {
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
	}

	return nil
}

func NewGDIRObject(obj runtime.GDObject, node ast.Node) *GDIRObject {
	return &GDIRObject{nil, obj, GDIRBaseNode{node}}
}

// Used for special cases where the type wraps an object
// like the spreadable type.
func NewGDIRObjectWithType(typ runtime.GDTypable, obj interface{}, node ast.Node) *GDIRObject {
	return &GDIRObject{typ, obj, GDIRBaseNode{node}}
}

func NewGDIRRegObject(reg cpu.GDReg, node ast.Node) *GDIRObject {
	identTyp := runtime.NewGDObjectRefType(runtime.NewGDByteIdent(byte(reg)))
	return &GDIRObject{identTyp, nil, GDIRBaseNode{node}}
}

func NewGDIRObjectRef(ident runtime.GDIdent, node ast.Node) *GDIRObject {
	identTyp := runtime.NewGDObjectRefType(ident)
	return &GDIRObject{identTyp, nil, GDIRBaseNode{node}}
}

func NewGDIRObjects(typ runtime.GDTypable, values []GDIRNode) *GDIRObject {
	return &GDIRObject{typ, values, GDIRBaseNode{}}
}
