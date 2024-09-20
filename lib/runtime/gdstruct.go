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
	Attrs map[GDIdentType]*GDSymbol
	Type  GDStructType
	stack *GDSymbolStack
}

func (gd *GDStruct) GetType() GDTypable    { return gd.Type }
func (gd *GDStruct) GetSubType() GDTypable { return nil }
func (gd *GDStruct) ToString() string {
	return "{" + JoinSlice(gd.Type, func(attrType GDStructAttrType, _ int) string {
		attr := gd.Attrs[attrType.Ident]
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
		for ident, attr := range gd.Attrs {
			typAttr, err := typ.GetAttrType(ident)
			if err != nil {
				return nil, TypeCastingWrongTypeWithHierarchyError(typ, gd.GetType(), err)
			}

			obj, err := attr.Object.CastToType(typAttr, stack)
			if err != nil {
				return nil, err
			}

			attr.Type = typAttr
			attr.Object = obj
		}
		gd.Type = typ

		return gd, nil
	}

	return nil, InvalidCastingWrongTypeErr(typ, gd.GetType())
}

func (gd *GDStruct) GetStack() *GDSymbolStack { return gd.stack }

func (gd *GDStruct) GetAttr(ident GDIdentType) (*GDSymbol, error) {
	if attr, ok := gd.Attrs[ident]; ok {
		return attr, nil
	}

	return nil, AttributeNotFoundErr(ident.ToString())
}

func (gd *GDStruct) SetAttr(ident GDIdentType, object GDObject) error {
	attr, ok := gd.Attrs[ident]
	if !ok {
		return AttributeNotFoundErr(ident.ToString())
	}

	err := attr.SetObject(object, gd.stack)
	if err != nil {
		return err
	}

	return nil
}

func NewGDStruct(structType GDStructType, stack *GDSymbolStack) (*GDStruct, error) {
	structStack := stack.NewSymbolStack(StructCtx)

	seen := make(map[GDIdentType]bool)
	attrs := make(map[GDIdentType]*GDSymbol, len(structType))
	for _, attr := range structType {
		if seen[attr.Ident] {
			return nil, AttributeAlreadyExistsErr(attr.Ident.ToString())
		}
		seen[attr.Ident] = true
		// Struct attributes are not constants
		symbol, err := structStack.AddSymbol(attr.Ident, true, false, attr.Type, GDZNil)
		if err != nil {
			return nil, err
		}

		attrs[attr.Ident] = symbol
	}

	return &GDStruct{attrs, structType, structStack}, nil
}
