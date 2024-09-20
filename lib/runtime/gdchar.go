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

type GDChar rune

func (gd GDChar) GetType() GDTypable    { return GDCharType }
func (gd GDChar) GetSubType() GDTypable { return nil }
func (gd GDChar) ToString() string      { return string(gd) }
func (gd GDChar) CastToType(typ GDTypable, stack *GDSymbolStack) (GDObject, error) {
	switch typ := typ.(type) {
	case GDType:
		switch typ {
		case GDCharType:
			return gd, nil
		case GDBoolType:
			return GDBool(gd != 0), nil
		case GDStringType:
			return GDString(gd.ToString()), nil
		case GDIntType:
			return GDInt(int(gd)), nil
		case GDInt8Type:
			return GDInt8(int8(gd)), nil
		case GDInt16Type:
			return GDInt16(int16(gd)), nil
		}
	}

	return nil, InvalidCastingWrongTypeErr(typ, gd.GetType())
}

func GDCharFromString(value string) (GDChar, error) {
	if len(value) != 1 {
		return GDChar(0), InvalidCastingLitErr(value, GDCharType)
	}

	return GDChar(rune(value[0])), nil
}
