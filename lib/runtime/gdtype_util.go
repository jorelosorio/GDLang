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

func ComputeTypeFromObjects(objects []GDObject, stack *GDStack) GDTypable {
	if len(objects) == 0 {
		return GDUntypedTypeRef
	}

	types := make([]GDTypable, 0)
	for _, obj := range objects {
		types = append(types, obj.GetType())
	}

	return ComputeTypeFromTypes(types)
}

func ComputeTypeFromTypes(types []GDTypable) GDTypable {
	if len(types) == 0 {
		return GDUntypedTypeRef
	}

	var computedType GDTypable
	if len(types) == 1 {
		return types[0]
	} else {
		uniqTypes := make([]GDTypable, 0)
		for _, objType := range types {
			isDuplicate := false
			for _, uniqType := range uniqTypes {
				if EqualTypes(objType, uniqType, nil) == nil {
					isDuplicate = true
					break
				}
			}
			if !isDuplicate {
				uniqTypes = append(uniqTypes, objType)
			}
		}

		if len(uniqTypes) == 1 {
			return uniqTypes[0]
		} else {
			computedType = NewGDUnionType(uniqTypes...)
		}
	}

	return computedType
}

func IsUntypedType(typ GDTypable) bool {
	switch typ := typ.(type) {
	case *GDArrayType:
		return IsUntypedType(typ.SubType)
	case GDTupleType:
		for _, t := range typ {
			if IsUntypedType(t) {
				return true
			}
		}

		return false
	}

	return typ == GDUntypedTypeRef
}

func CanBeAssign(left, right GDTypable, stack *GDStack) error {
	_, err := determineTypeCompatibility(left, right, true, stack)
	if err != nil {
		if err, isGDErr := err.(GDRuntimeErr); isGDErr {
			switch err.Code {
			case IncompatibleTypeCode:
				return WrongTypesErr(left, right)
			default:
				return err
			}
		}
	}

	return nil
}

func EqualTypes(left, right GDTypable, stack *GDStack) error {
	_, err := determineTypeCompatibility(left, right, false, stack)
	if err != nil {
		return TypesAreNotEqualErr(left, right)
	}

	return nil
}

// Rules:
// - Untyped is considered a weak unknown type, however, it can mutate to any other more strong type.
func InferType(toType, fromType GDTypable, stack *GDStack) (GDTypable, error) {
	typ, err := determineTypeCompatibility(toType, fromType, true, stack)
	if err != nil {
		if err, isGDErr := err.(GDRuntimeErr); isGDErr {
			switch err.Code {
			case IncompatibleTypeCode:
				return nil, WrongTypesErr(toType, fromType)
			default:
				return nil, err
			}
		}
	}

	return typ, nil
}

func determineTypeCompatibility(toType, fromType GDTypable, isAssignmentNeeded bool, stack *GDStack) (GDTypable, error) {
	if toRefType, isRefType := toType.(GDTypeRefType); isRefType && stack != nil {
		symbol, err := stack.GetSymbol(toRefType.GDIdent)
		if err != nil {
			return nil, err
		}

		toType = symbol.Type
	}

	if fromRefType, isRefType := fromType.(GDTypeRefType); isRefType && stack != nil {
		symbol, err := stack.GetSymbol(fromRefType.GDIdent)
		if err != nil {
			return nil, err
		}

		fromType = symbol.Type
	}

	switch toType := toType.(type) {
	// It should fall into the GDTypeRefType case
	//If the types are not exhaustively checked in the stack,
	// they must be compared against the names.
	case GDTypeRefType:
		if typ, ok := fromType.(GDTypeRefType); ok {
			if toType.GDIdent.GetRawValue() == typ.GDIdent.GetRawValue() {
				return toType, nil
			}
		}
	case *GDTypeAliasType:
		if typ, ok := fromType.(*GDTypeAliasType); ok {
			fromType = typ.GDTypable
		}

		_, err := determineTypeCompatibility(toType.GDTypable, fromType, isAssignmentNeeded, stack)
		if err != nil {
			return nil, err
		}

		return toType, nil
	case *GDArrayType:
		if fromType, ok := fromType.(*GDArrayType); ok {
			typ, err := determineTypeCompatibility(toType.SubType, fromType.SubType, isAssignmentNeeded, stack)
			if err != nil {
				return nil, err
			}

			return NewGDArrayType(typ), nil
		}
	// Union types do not have untyped types
	case GDUnionType:
		if fromTypeUnion, isUnion := fromType.(GDUnionType); isUnion {
			if len(toType) != len(fromTypeUnion) {
				return nil, WrongTypesErr(toType, fromType)
			}

			for _, typ := range fromTypeUnion {
				if !toType.ContainsType(typ, stack) {
					return nil, WrongTypesErr(toType, fromType)
				}
			}

			return toType, nil
		} else if toType.ContainsType(fromType, stack) {
			return toType, nil
		}
	case GDTupleType:
		if fromTypeTuple, ok := fromType.(GDTupleType); ok {
			if len(toType) != len(fromTypeTuple) {
				return nil, WrongTypesErr(toType, fromType)
			}

			for i, typ := range toType {
				typ, err := determineTypeCompatibility(typ, fromTypeTuple[i], isAssignmentNeeded, stack)
				if err != nil {
					return nil, err
				}

				toType[i] = typ
			}

			return toType, nil
		}
	case GDStructType:
		if fromType, ok := fromType.(GDStructType); ok {
			if len(fromType) == 0 {
				return toType, nil
			}

			if len(toType) != len(fromType) {
				return nil, WrongTypesErr(toType, fromType)
			}

			structAttrTypes := make([]*GDStructAttrType, len(toType))
			for i, fromAttr := range fromType {
				attrType, err := toType.GetAttrType(fromAttr.Ident)
				if err != nil {
					return nil, WrongTypesErr(toType, fromType)
				}

				typ, err := determineTypeCompatibility(attrType, fromAttr.Type, isAssignmentNeeded, stack)
				if err != nil {
					return nil, err
				}

				structAttrTypes[i] = &GDStructAttrType{Ident: fromAttr.Ident, Type: typ}
			}

			return NewGDStructType(structAttrTypes...), nil
		}
	// For functions, is enough to check both types are equal
	// there is no need to check the arguments and return types with untyped types
	case *GDLambdaType:
		if fromType, ok := fromType.(*GDLambdaType); ok {
			if len(toType.ArgTypes) != len(fromType.ArgTypes) {
				return nil, WrongTypesErr(toType, fromType)
			}

			if toType.IsVariadic != fromType.IsVariadic {
				return nil, WrongTypesErr(toType, fromType)
			}

			if _, err := determineTypeCompatibility(toType.ReturnType, fromType.ReturnType, false, stack); err != nil {
				return nil, WrongTypesErr(toType, fromType)
			}

			for i, argType := range toType.ArgTypes {
				if _, err := determineTypeCompatibility(argType.Value, fromType.ArgTypes[i].Value, false, stack); err != nil {
					return nil, WrongTypesErr(toType, fromType)
				}
			}

			return toType, nil
		}
	}

	switch fromType := fromType.(type) {
	case GDUnionType:
		cTypes := make([]GDTypable, len(fromType))
		for i, cType := range fromType {
			typ, err := determineTypeCompatibility(toType, cType, isAssignmentNeeded, stack)
			if err != nil {
				return nil, err
			}

			cTypes[i] = typ
		}

		return NewGDUnionType(cTypes...), nil
	}

	if isAssignmentNeeded {
		// Untyped behaves as the type `any`, but with weak inference,
		// that means it mutates to the other type, even if the other type is untyped.
		if toType == GDUntypedTypeRef {
			if fromType != GDNilTypeRef {
				return fromType, nil // Mutate to the fromType
			} else {
				return toType, nil
			}
		}

		if fromType == GDUntypedTypeRef {
			return toType, nil // Mutate to the toType
		}

		// Any and nil cases
		if toType == GDAnyTypeRef || fromType == GDNilTypeRef {
			return toType, nil
		}
	}

	// Are them equal?
	if toType.GetCode() == fromType.GetCode() {
		return toType, nil
	}

	return nil, WrongTypesErr(toType, fromType)
}
