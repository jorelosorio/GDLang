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

type GDArrayType struct{ SubType GDTypable }

func (t *GDArrayType) GetCode() GDTypableCode { return GDArrayTypeCode }

func (t *GDArrayType) ToString() string {
	if t.SubType != nil {
		return "[" + t.SubType.ToString() + "]"
	}

	return GDTypeCodeMap[GDArrayTypeCode]
}

func (t *GDArrayType) GetTypeAt(index int) GDTypable { return t.SubType }
func (t *GDArrayType) GetIterableType() GDTypable    { return t.SubType }
func (t *GDArrayType) SetTypeAt(index int, typ GDTypable, stack *GDStack) error {
	return CanBeAssign(t.SubType, typ, stack)
}

func NewGDArrayType(subType GDTypable) *GDArrayType { return &GDArrayType{SubType: subType} }
func NewGDEmptyArrayType() *GDArrayType             { return &GDArrayType{SubType: GDUntypedTypeRef} }
func NewGDArrayTypeWithTypes(types []GDTypable, stack *GDStack) *GDArrayType {
	return &GDArrayType{SubType: ComputeTypeFromTypes(types)}
}
