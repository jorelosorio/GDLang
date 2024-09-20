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

type GDStructAttrType struct {
	Ident GDIdentType
	Type  GDTypable
}

func (t GDStructAttrType) GetCode() GDTypableCode { return GDNilTypeCode }
func (t GDStructAttrType) ToString() string       { return t.Ident.ToString() + ": " + t.Type.ToString() }

type GDStructType []GDStructAttrType

func (t GDStructType) GetCode() GDTypableCode { return GDStructTypeCode }

func (t GDStructType) ToString() string {
	return "{" + JoinSlice(t, func(attr GDStructAttrType, _ int) string {
		return attr.ToString()
	}, ", ") + "}"
}

func (t GDStructType) GetAttrType(key GDIdentType) (GDTypable, error) {
	for _, attr := range t {
		if attr.Ident == key {
			return attr.Type, nil
		}
	}

	return nil, AttributeNotFoundErr(key.ToString())
}

// An empty struct type is a struct type with no attributes
// This is useful to represent an empty struct
func NewGDStructType(attrs ...GDStructAttrType) GDStructType {
	return GDStructType(attrs)
}
