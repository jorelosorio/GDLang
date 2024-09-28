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
	stack *GDSymbolStack
}

func (gd *GDStruct) GetType() GDTypable    { return gd.Type }
func (gd *GDStruct) GetSubType() GDTypable { return nil }
func (gd *GDStruct) ToString() string {
	return "{" + JoinSlice(gd.Type, func(attrType GDStructAttrType, _ int) string {
		symbol, err := gd.GetAttr(attrType.Ident)
		if err != nil {
			// NOTE: If attribute is not found, then panic!
			panic(err)
		}

		if symbol.Object != nil {
			return attrType.Ident.ToString() + ": " + ObjectToStringForInternalData(symbol.Object)
		} else {
			return attrType.ToString()
		}
	}, ", ") + "}"
}
func (gd *GDStruct) CastToType(typ GDTypable, stack *GDSymbolStack) (GDObject, error) {
	switch typ := typ.(type) {
	case GDType:
		switch typ {
		case GDStringType:
			return GDString(gd.ToString()), nil
		}
	case GDStructType:
		for _, attr := range gd.Type {
			typAttr, err := typ.GetAttrType(attr.Ident)
			if err != nil {
				return nil, TypeCastingWrongTypeWithHierarchyError(typ, gd.GetType(), err)
			}

			symbol, err := gd.GetAttr(attr.Ident)
			if err != nil {
				return nil, err
			}

			obj, err := symbol.Object.CastToType(typAttr, stack)
			if err != nil {
				return nil, err
			}

			symbol.Type = typAttr
			symbol.Object = obj
		}
		gd.Type = typ

		return gd, nil
	}

	return nil, InvalidCastingWrongTypeErr(typ, gd.GetType())
}

func (gd *GDStruct) GetStack() *GDSymbolStack { return gd.stack }

func (gd *GDStruct) GetAttr(ident GDIdent) (*GDSymbol, error) {
	symbol, err := gd.stack.GetLocalSymbol(ident)
	if err != nil {
		return nil, AttributeNotFoundErr(ident.ToString())
	}

	return symbol, nil
}

func (gd *GDStruct) SetAttr(ident GDIdent, object GDObject) (*GDSymbol, error) {
	symbol, err := gd.GetAttr(ident)
	if err != nil {
		return nil, err
	}

	err = symbol.SetObject(object, gd.stack)
	if err != nil {
		return nil, err
	}

	return symbol, nil
}

func NewGDStruct(structType GDStructType, stack *GDSymbolStack) (*GDStruct, error) {
	structStack := stack.NewSymbolStack(StructCtx)

	for _, attr := range structType {
		// Struct attributes are not constants
		_, err := structStack.AddSymbol(attr.Ident, true, false, attr.Type, GDZNil)
		if err != nil {
			return nil, err
		}
	}

	return &GDStruct{structType, structStack}, nil
}
