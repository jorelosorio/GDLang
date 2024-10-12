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

type GDStruct struct {
	Type  GDStructType
	stack *GDStack
}

func (gd *GDStruct) GetType() GDTypable    { return gd.Type }
func (gd *GDStruct) GetSubType() GDTypable { return nil }
func (gd *GDStruct) ToString() string {
	return "{" + JoinSlice(gd.Type, func(attrType *GDStructAttrType, _ int) string {
		symbol, err := gd.GetAttr(attrType.Ident)
		if err != nil {
			// NOTE: If attribute is not found, then panic!
			panic(err)
		}

		if symbol.Value != nil {
			return attrType.Ident.ToString() + ": " + ConvertObjectToString(symbol.Value)
		} else {
			return attrType.ToString()
		}
	}, ", ") + "}"
}
func (gd *GDStruct) CastToType(typ GDTypable) (GDObject, error) {
	switch typ := typ.(type) {
	case GDStringType:
		return GDString(gd.ToString()), nil
	case GDStructType:
		for _, attr := range gd.Type {
			typAttr, err := typ.GetAttrType(attr.Ident)
			if err != nil {
				return nil, TypeCastingWrongTypeWithHierarchyError(typ, gd.GetType(), err)
			}

			attrSymbol, err := gd.GetAttr(attr.Ident)
			if err != nil {
				return nil, err
			}

			attrObj := attrSymbol.Value.(GDObject)

			obj, err := attrObj.CastToType(typAttr)
			if err != nil {
				return nil, err
			}

			attrSymbol.Type = typAttr
			attrSymbol.Value = obj
		}

		gd.Type = typ

		return gd, nil
	}

	return nil, InvalidCastingWrongTypeErr(typ, gd.GetType())
}

func (gd *GDStruct) GetStack() *GDStack { return gd.stack }

func (gd *GDStruct) GetAttr(ident GDIdent) (*GDSymbol, error) {
	symbol, ok := gd.stack.Symbols[ident.GetRawValue()]
	if !ok {
		return nil, AttributeNotFoundErr(ident.ToString())
	}

	return symbol, nil
}

func (gd *GDStruct) SetAttr(ident GDIdent, typ GDTypable, object GDObject) (*GDSymbol, error) {
	symbol, err := gd.GetAttr(ident)
	if err != nil {
		return nil, err
	}

	err = symbol.SetObject(typ, object, gd.stack)
	if err != nil {
		return nil, err
	}

	return symbol, nil
}

func NewGDStruct(typ GDStructType, stack *GDStack) (*GDStruct, error) {
	structStack := stack.NewStack(StructCtx)

	for _, attr := range typ {
		// Struct attributes are not constants
		_, err := structStack.AddNewSymbol(attr.Ident, true, false, attr.Type, GDNilTypeRef, GDZNil)
		if err != nil {
			return nil, err
		}
	}

	return &GDStruct{typ, structStack}, nil
}

func QuickGDStruct(stack *GDStack, typ GDStructType, attrs ...GDObject) (*GDStruct, error) {
	structStack := stack.NewStack(StructCtx)

	if len(attrs) != len(typ) {
		panic("attributes count must be equal to structType count")
	}

	for i, attr := range typ {
		// Struct attributes are not constants
		_, err := structStack.AddNewSymbol(attr.Ident, true, false, attr.Type, attrs[i].GetType(), attrs[i])
		if err != nil {
			return nil, err
		}
	}

	return &GDStruct{typ, structStack}, nil
}
