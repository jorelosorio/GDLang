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

func (gd GDString) GetType() GDTypable    { return GDStringType }
func (gd GDString) GetSubType() GDTypable { return nil }
func (gd GDString) ToString() string      { return string(gd) }
func (gd GDString) CastToType(typ GDTypable, stack *GDSymbolStack) (GDObject, error) {
	switch typ {
	case GDStringType:
		return gd, nil
	case GDBoolType:
		return NewGDBoolFromString(string(gd))
	case GDIntType, GDInt8Type, GDInt16Type:
		return NewGDIntNumberFromString(string(gd))
	case GDCharType:
		if len(gd) != 1 {
			return nil, InvalidCastingLitErr(gd.ToString(), GDCharType)
		}
		return GDChar(gd[0]), nil
	case GDFloatType, GDFloat32Type, GDFloat64Type:
		return NewGDFloatNumberFromString(string(gd))
	case GDComplexType, GDComplex64Type, GDComplex128Type:
		return NewGDComplexNumberFromString(string(gd))
	}

	return nil, InvalidCastingWrongTypeErr(typ, gd.GetType())
}

// Iterable interface

func (gd GDString) Length() int   { return len(gd) }
func (gd GDString) IsEmpty() bool { return len(gd) == 0 }
func (gd GDString) Get(index int) (GDObject, error) {
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
func (gd GDString) GetTypes() ([]GDTypable, bool) { return []GDTypable{GDCharType}, true }
func (gd GDString) GetIterableType() GDTypable    { return GDCharType }

func (gd GDString) checkIndex(index int) error {
	if index < 0 || index >= len(gd) {
		return IndexOutOfBoundsErr
	}
	return nil
}
