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
)

func IRObjectWithTypeToString(obj runtime.GDObject) string {
	switch obj := obj.(type) {
	case *runtime.GDIdObject:
		return IRObjectTypeToString(obj)
	default:
		return fmt.Sprintf("(%s: %s)", IRObjectTypeToString(obj), IRObjectWithoutTypeToString(obj))
	}
}

func IRTypeToString(typ runtime.GDTypable) string {
	return formatTypeToString(typ)
}

func IRObjectTypeToString(obj runtime.GDObject) string {
	if typ := obj.GetSubType(); typ != nil {
		return formatTypeToString(typ)
	}

	return formatTypeToString(obj.GetType())
}

func formatTypeToString(typ runtime.GDTypable) string {
	typeCode := runtime.GDTypeCodeMap[typ.GetCode()]
	switch typ := typ.(type) {
	case runtime.GDRefType:
		return fmt.Sprintf("(%s: %s)", typeCode, typ.ToString())
	case runtime.GDObjRefType:
		return fmt.Sprintf("(%s: %s)", typeCode, typ.ToString())
	default:
		return typ.ToString()
	}
}

func IRObjectWithoutTypeToString(obj runtime.GDObject) string {
	switch obj := obj.(type) {
	case *runtime.GDArray:
		return fmt.Sprintf("[%s]", runtime.JoinSlice(obj.Objects, func(object runtime.GDObject, _ int) string {
			return IRObjectWithTypeToString(object)
		}, ", "))
	case *runtime.GDTuple:
		return fmt.Sprintf("(%s)", runtime.JoinSlice(obj.Objects, func(object runtime.GDObject, _ int) string {
			return IRObjectWithTypeToString(object)
		}, ", "))
	case *runtime.GDStruct:
		strObjs := make([]string, len(obj.Type))
		for i, attrType := range obj.Type {
			symbol, err := obj.GetAttr(attrType.Ident)
			if err != nil {
				panic(err)
			}

			strObjs[i] = IRObjectWithTypeToString(symbol.Object)
		}
		return fmt.Sprintf("{%s}", runtime.JoinSlice(strObjs, func(str string, _ int) string {
			return str
		}, ", "))
	case *runtime.GDSpreadable:
		return IRObjectWithTypeToString(obj.Iterable) + "..."
	case runtime.GDObject:
		return runtime.ObjectToStringForInternalData(obj)
	}

	panic(runtime.NewGDRuntimeErr(runtime.UnsupportedTypeCode, "Unsupported type when converting to string"))
}
