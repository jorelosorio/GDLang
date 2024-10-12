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

import "strings"

var charScapeMap = strings.NewReplacer(
	"\\a", "\a",
	"\\b", "\b",
	"\\\\", "\\",
	"\\t", "\t",
	"\\n", "\n",
	"\\f", "\f",
	"\\r", "\r",
	"\\v", "\v",
	"\\'", "'",
	"\\\"", "\"",
)

type GDString string

func (gd GDString) Escape() string { return charScapeMap.Replace(string(gd)) }

// Object interface

func (gd GDString) GetType() GDTypable    { return GDStringTypeRef }
func (gd GDString) GetSubType() GDTypable { return nil }
func (gd GDString) ToString() string      { return string(gd) }
func (gd GDString) CastToType(typ GDTypable) (GDObject, error) {
	switch typ {
	case GDStringTypeRef:
		return gd, nil
	case GDBoolTypeRef:
		return NewGDBoolFromString(string(gd))
	case GDIntTypeRef, GDInt8TypeRef, GDInt16TypeRef:
		return NewGDIntNumberFromString(string(gd))
	case GDCharTypeRef:
		if len(gd) != 1 {
			return nil, InvalidCastingLitErr(gd.ToString(), GDCharTypeRef)
		}
		return GDChar(gd[0]), nil
	case GDFloatTypeRef, GDFloat32TypeRef, GDFloat64TypeRef:
		return NewGDFloatNumberFromString(string(gd))
	case GDComplexTypeRef, GDComplex64TypeRef, GDComplex128TypeRef:
		return NewGDComplexNumberFromString(string(gd))
	}

	return nil, InvalidCastingWrongTypeErr(typ, gd.GetType())
}

// Iterable interface

func (gd GDString) Length() int   { return len(gd) }
func (gd GDString) IsEmpty() bool { return len(gd) == 0 }
func (gd GDString) GetObjectAt(index int) (GDObject, error) {
	if err := gd.checkIndex(index); err != nil {
		return nil, err
	}

	return GDChar(gd[index]), nil
}
func (gd GDString) GetObjects() []GDObject {
	objects := make([]GDObject, len(gd))
	for i, c := range gd {
		objects[i] = GDChar(c)
	}
	return objects
}

// Iterable mutable interface

func (gd GDString) GetTypeAt(index int) GDTypable { return GDCharTypeRef }
func (gd GDString) GetIterableType() GDTypable    { return GDCharTypeRef }

func (gd GDString) checkIndex(index int) error {
	if index < 0 || index >= len(gd) {
		return IndexOutOfBoundsErr
	}
	return nil
}
