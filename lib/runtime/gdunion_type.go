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

type GDUnionType []GDTypable

func (t GDUnionType) GetCode() GDTypableCode {
	return GDUnionTypeCode
}

func (t GDUnionType) ContainsType(typ GDTypable, stack *GDStack) bool {
	for _, existingTyp := range t {
		if err := EqualTypes(existingTyp, typ, stack); err == nil {
			return true
		}
	}
	return false
}

func (t GDUnionType) AppendType(typ GDTypable, stack *GDStack) GDUnionType {
	switch typ := typ.(type) {
	case GDUnionType:
		for _, typ := range typ {
			if !t.ContainsType(typ, stack) {
				t = append(t, typ)
			}
		}
	default:
		if !t.ContainsType(typ, stack) {
			return append(t, typ)
		}
	}

	return t
}

func (t GDUnionType) ToString() string {
	return "(" + JoinSlice(t, func(typ GDTypable, _ int) string {
		return typ.ToString()
	}, " | ") + ")"
}

func NewGDUnionType(fields ...GDTypable) GDUnionType {
	return fields
}
