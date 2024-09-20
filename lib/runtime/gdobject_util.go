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

package runtime

func Unwrap(value GDObject) GDObject {
	switch value := value.(type) {
	case *GDIdObject:
		return Unwrap(value.GDObject)
	case *GDAttrIdObject:
		return Unwrap(value.GDObject)
	}

	return value
}

func EqualObjects(a, b GDObject) bool {
	switch a := a.(type) {
	case *GDTuple:
		if b, ok := b.(*GDTuple); ok {
			return equalArray(a.Objects, b.Objects)
		}
		return false
	// TODO: Add case for GDStruct, GDFunc
	case *GDArray:
		if b, ok := b.(*GDArray); ok {
			return equalArray(a.Objects, b.Objects)
		}
		return false
	}

	return a == b
}

func ObjectToStringForInternalData(object GDObject) string {
	switch object := object.(type) {
	case GDString:
		return Sprintf("\"%@\"", object)
	case GDChar:
		return Sprintf("'%@'", object)
	case GDObject:
		return object.ToString()
	}

	panic(NewGDRuntimeErr(UnsupportedTypeCode, "Unsupported type when converting to string"))
}

// This function requires that the type was already checked or inferred
// before calling it.
func TypeCoercion(obj GDObject, typ GDTypable, stack *GDSymbolStack) (GDObject, error) {
	typ, err := UnwrapIdentType(typ, stack)
	if err != nil {
		return nil, err
	}

	switch obj := obj.(type) {
	case *GDTuple:
		if typ, ok := typ.(GDTupleType); ok {
			obj.GDTupleType = typ
			return obj, nil
		}
	case *GDArray:
		if typ, ok := typ.(*GDArrayType); ok {
			obj.GDArrayType = typ
			return obj, nil
		}
	case *GDStruct:
		if typ, ok := typ.(GDStructType); ok {
			obj.Type = typ
			for _, attr := range typ {
				objAttr := obj.Attrs[attr.Ident]
				obj, err := TypeCoercion(objAttr.Object, attr.Type, stack)
				if err != nil {
					return nil, err
				}

				objAttr.Type = attr.Type
				objAttr.Object = obj
			}
		}
	}

	return obj, nil
}

func equalArray(a, b []GDObject) bool {
	if len(a) != len(b) {
		return false
	}

	for i, element := range a {
		if !EqualObjects(element, b[i]) {
			return false
		}
	}

	return true
}
