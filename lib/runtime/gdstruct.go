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
	Attrs []GDKeyValue[GDIdent, *GDSymbol]
	Type  GDStructType
	stack *GDSymbolStack
}

func (gd *GDStruct) GetType() GDTypable    { return gd.Type }
func (gd *GDStruct) GetSubType() GDTypable { return nil }
func (gd *GDStruct) ToString() string {
	return "{" + JoinSlice(gd.Type, func(attrType GDStructAttrType, _ int) string {
		attr, err := gd.GetAttr(attrType.Ident)
		if err != nil {
			// NOTE: If attribute is not found, then panic!
			panic(err)
		}

		if attr.Object != nil {
			return attrType.Ident.ToString() + ": " + ObjectToStringForInternalData(attr.Object)
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
		for _, attr := range gd.Attrs {
			typAttr, err := typ.GetAttrType(attr.Key)
			if err != nil {
				return nil, TypeCastingWrongTypeWithHierarchyError(typ, gd.GetType(), err)
			}

			obj, err := attr.Value.Object.CastToType(typAttr, stack)
			if err != nil {
				return nil, err
			}

			attr.Value.Type = typAttr
			attr.Value.Object = obj
		}
		gd.Type = typ

		return gd, nil
	}

	return nil, InvalidCastingWrongTypeErr(typ, gd.GetType())
}

func (gd *GDStruct) GetStack() *GDSymbolStack { return gd.stack }

func (gd *GDStruct) GetAttr(ident GDIdent) (*GDSymbol, error) {
	for _, attr := range gd.Attrs {
		if attr.Key.GetRawValue() == ident.GetRawValue() {
			return attr.Value, nil
		}
	}

	return nil, AttributeNotFoundErr(ident.ToString())
}

func (gd *GDStruct) SetAttr(ident GDIdent, object GDObject) error {
	symbol, error := gd.GetAttr(ident)
	if error != nil {
		return error
	}

	err := symbol.SetObject(object, gd.stack)
	if err != nil {
		return err
	}

	return nil
}

func NewGDStruct(structType GDStructType, stack *GDSymbolStack) (*GDStruct, error) {
	structStack := stack.NewSymbolStack(StructCtx)

	seen := make(map[any]bool)
	attrs := make([]GDKeyValue[GDIdent, *GDSymbol], len(structType))
	for i, attr := range structType {
		if seen[attr.Ident.GetRawValue()] {
			return nil, AttributeAlreadyExistsErr(attr.Ident.ToString())
		}
		seen[attr.Ident.GetRawValue()] = true

		// Struct attributes are not constants
		symbol, err := structStack.AddSymbol(attr.Ident, true, false, attr.Type, GDZNil)
		if err != nil {
			return nil, err
		}

		attrs[i] = GDKeyValue[GDIdent, *GDSymbol]{attr.Ident, symbol}
	}

	return &GDStruct{attrs, structType, structStack}, nil
}
