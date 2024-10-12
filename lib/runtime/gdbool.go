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

import "strconv"

type GDBool bool

func (gd GDBool) GetType() GDTypable    { return GDBoolTypeRef }
func (gd GDBool) GetSubType() GDTypable { return nil }
func (gd GDBool) ToString() string {
	if gd {
		return "true"
	}
	return "false"
}
func (gd GDBool) CastToType(typ GDTypable) (GDObject, error) {
	switch typ := typ.(type) {
	case GDStringType:
		return GDString(gd.ToString()), nil
	case GDType:
		switch typ {
		case GDBoolTypeRef:
			return gd, nil
		case GDIntTypeRef:
			if gd {
				return GDInt(1), nil
			}

			return GDInt(0), nil
		case GDInt8TypeRef:
			if gd {
				return GDInt8(1), nil
			}
		case GDInt16TypeRef:
			if gd {
				return GDInt16(1), nil
			}
		case GDFloat32TypeRef:
			if gd {
				return GDFloat32(1), nil
			}

			return GDFloat32(0), nil
		case GDFloat64TypeRef:
			if gd {
				return GDFloat64(1), nil
			}

			return GDFloat64(0), nil
		case GDComplexTypeRef:
			if gd {
				return GDComplex64(1), nil
			}

			return GDComplex64(0), nil
		case GDComplex64TypeRef:
			if gd {
				return GDComplex64(1), nil
			}

			return GDComplex64(0), nil
		case GDComplex128TypeRef:
			if gd {
				return GDComplex128(1), nil
			}

			return GDComplex128(0), nil
		case GDCharTypeRef:
			if gd {
				return GDChar('T'), nil
			}

			return GDChar('F'), nil
		}
	}

	return nil, InvalidCastingWrongTypeErr(typ, gd.GetType())
}

func NewGDBoolFromString(value string) (GDBool, error) {
	boolean, err := strconv.ParseBool(value)
	if err != nil {
		return false, InvalidCastingLitErr(value, GDBoolTypeRef)
	}

	return GDBool(boolean), nil
}
