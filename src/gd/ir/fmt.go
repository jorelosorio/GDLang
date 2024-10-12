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
	"fmt"
	"gdlang/lib/runtime"
	"gdlang/src/cpu"
)

func IRObjectWithTypeToString(typ runtime.GDTypable, obj runtime.GDObject) string {
	if obj == nil {
		return IRObjectTypeToString(obj)
	}

	return fmt.Sprintf("(%s: %s)", IRObjectTypeToString(obj), IRObjectWithoutTypeToString(obj))
}

func IRObjectToString(obj runtime.GDObject) string {
	return fmt.Sprintf("(%s: %s)", IRObjectTypeToString(obj), IRObjectWithoutTypeToString(obj))
}

func IRObjectTypeToString(obj runtime.GDObject) string {
	if typ := obj.GetSubType(); typ != nil {
		return typ.ToString()
	}

	return IRTypeToString(obj.GetType())
}

func IRObjectWithoutTypeToString(obj runtime.GDObject) string {
	switch obj := obj.(type) {
	case *runtime.GDArray:
		return fmt.Sprintf("[%s]", runtime.JoinSlice(obj.GetObjects(), func(obj runtime.GDObject, _ int) string {
			return IRObjectToString(obj)
		}, ", "))
	case *runtime.GDTuple:
		return fmt.Sprintf("(%s)", runtime.JoinSlice(obj.GetObjects(), func(obj runtime.GDObject, _ int) string {
			return IRObjectToString(obj)
		}, ", "))
	case *runtime.GDStruct:
		strObjs := make([]string, len(obj.Type))
		for i, attrType := range obj.Type {
			symbol, err := obj.GetAttr(attrType.Ident)
			if err != nil {
				panic(err)
			}

			strObjs[i] = IRObjectWithTypeToString(symbol.Type, symbol.Value.(runtime.GDObject))
		}
		return fmt.Sprintf("{%s}", runtime.JoinSlice(strObjs, func(str string, _ int) string {
			return str
		}, ", "))
	case *runtime.GDSpreadable:
		return IRObjectWithTypeToString(obj.GetType(), obj.Iterable) + "..."
	case runtime.GDObject:
		return runtime.ConvertObjectToString(obj)
	}

	panic(runtime.NewGDRuntimeErr(runtime.UnsupportedTypeCode, "Unsupported type when converting to string"))
}

func IRTypeToString(typ runtime.GDTypable) string {
	return formatTypeToString(typ)
}

func formatTypeToString(typ runtime.GDTypable) string {
	typeCode := runtime.GDTypeCodeMap[typ.GetCode()]
	switch typ := typ.(type) {
	case runtime.GDTypeRefType:
		return fmt.Sprintf("%s:<%s>", typeCode, typ.ToString())
	case runtime.GDObjectRefType:
		switch typ.GDIdent.GetMode() {
		case runtime.GDByteIdentMode:
			rawValue, ok := typ.GDIdent.GetRawValue().(byte)
			if !ok {
				panic("GDIdent.GetRawValue() is not of type byte")
			}

			return fmt.Sprintf("%s:<%s>", typeCode, cpu.GetCPURegName(cpu.GDReg(rawValue)))
		default:
			return fmt.Sprintf("%s<%s>", typeCode, typ.ToString())
		}
	case *runtime.GDTypeAliasType:
		return fmt.Sprintf("%s:<%s>", typeCode, typ.ToString())
	case *runtime.GDArrayType:
		return fmt.Sprintf("[%s]", formatTypeToString(typ.SubType))
	case runtime.GDTupleType:
		return fmt.Sprintf("(%s)", runtime.JoinSlice(typ, func(typ runtime.GDTypable, _ int) string {
			return formatTypeToString(typ)
		}, ", "))
	case runtime.GDStructType:
		return fmt.Sprintf("{%s}", runtime.JoinSlice(typ, func(typ *runtime.GDStructAttrType, _ int) string {
			return fmt.Sprintf("%s: %s", typ.Ident.ToString(), formatTypeToString(typ.Type))
		}, ", "))
	case *runtime.GDLambdaType:
		return fmt.Sprintf("func(%s) => %s", runtime.JoinSlice(typ.ArgTypes, func(arg *runtime.GDLambdaArgType, _ int) string {
			return fmt.Sprintf("%s: %s", arg.Key.ToString(), formatTypeToString(arg.Value))
		}, ", "), formatTypeToString(typ.ReturnType))
	case runtime.GDSpreadableType:
		return fmt.Sprintf("%s...", formatTypeToString(typ.GetIterableType()))
	default:
		return typ.ToString()
	}
}
