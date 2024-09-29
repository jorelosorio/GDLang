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

func ZObjectForType(typ GDTypable, stack *GDSymbolStack) (GDObject, error) {
	switch typ.GetCode() {
	case GDUntypedTypeCode:
		return GDZUntyped, nil
	case GDIntTypeCode:
		return GDZInt, nil
	case GDFloatTypeCode:
		return GDZFloat, nil
	case GDNilTypeCode:
		return GDZNil, nil
	case GDComplexTypeCode:
		return GDZComplex, nil
	case GDStringTypeCode:
		return GDZString, nil
	case GDBoolTypeCode:
		return GDZBool, nil
	case GDAnyTypeCode:
		return GDZAny, nil
	case GDStructTypeCode:
		return NewGDStruct(typ.(GDStructType), stack)
	case GDRefTypeCode:
		symbol, err := stack.GetSymbol(typ.(GDIdent))
		if err != nil {
			return nil, err
		}

		return ZObjectForType(symbol.Type, stack)
	case GDCharTypeCode:
		return GDZChar, nil
	case GDTupleTypeCode:
		tupleType := typ.(GDTupleType)

		objs := make([]GDObject, len(tupleType))
		for i, field := range tupleType {
			obj, err := ZObjectForType(field, stack)
			if err != nil {
				return nil, err
			}

			objs[i] = obj
		}

		return &GDTuple{typ.(GDTupleType), objs}, nil
	case GDArrayTypeCode:
		return NewGDArrayWithTypeAndObjects(typ.(*GDArrayType), []GDObject{}), nil
	case GDUnionTypeCode:
		objs := make([]GDObject, 0)
		unionType := typ.(GDUnionType)
		for _, field := range unionType {
			obj, err := ZObjectForType(field, stack)
			if err != nil {
				return nil, err
			}

			objs = append(objs, obj)
		}

		return NewGDUnion(unionType, objs...), nil
	case GDLambdaTypeCode:
		return NewGDLambdaWithType(typ.(*GDLambdaType), stack, nil), nil
	}

	return nil, UnsupportedTypeErr(typ.ToString())
}
