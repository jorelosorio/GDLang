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

// Comparators

func IsNumber(value any) bool {
	return IsInt(value) || IsFloat(value) || IsComplex(value)
}

func IsInt(value any) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	case GDInt, GDInt8, GDInt16:
		return true
	}

	return false
}

func IsFloat(value any) bool {
	switch value.(type) {
	case float32, float64:
		return true
	case GDFloat32, GDFloat64:
		return true
	}

	return false
}

func IsComplex(value any) bool {
	switch value.(type) {
	case complex64, complex128:
		return true
	case GDComplex64, GDComplex128:
		return true
	}

	return false
}

func ToInt(value any) (GDInt, error) {
	switch value := value.(type) {
	case GDInt:
		return value, nil
	case GDInt8:
		return GDInt(value), nil
	case GDInt16:
		return GDInt(value), nil
	case GDFloat32:
		return GDInt(value), nil
	case GDFloat64:
		return GDInt(value), nil
	case GDComplex64:
		return GDInt(real(value)), nil
	case GDComplex128:
		return GDInt(real(value)), nil
	case GDObject:
		return GDInt(0), InvalidCastingWrongTypeErr(GDIntType, value.GetType())
	default:
		return GDInt(0), InvalidCastingExpectedTypeErr(GDIntType)
	}
}

func ToFloat(value any) (GDFloat64, error) {
	switch value := value.(type) {
	case GDInt:
		return GDFloat64(value), nil
	case GDInt8:
		return GDFloat64(value), nil
	case GDInt16:
		return GDFloat64(value), nil
	case GDFloat32:
		return GDFloat64(value), nil
	case GDFloat64:
		return value, nil
	case GDComplex64:
		return GDFloat64(real(value)), nil
	case GDComplex128:
		return GDFloat64(real(value)), nil
	case GDObject:
		return GDFloat64(0), InvalidCastingWrongTypeErr(GDFloatType, value.GetType())
	default:
		return GDFloat64(0), InvalidCastingExpectedTypeErr(GDFloatType)
	}
}

func ToComplex(value any) (GDComplex128, error) {
	switch value := value.(type) {
	case GDInt:
		return GDComplex128(complex(float64(value), 0)), nil
	case GDInt8:
		return GDComplex128(complex(float64(value), 0)), nil
	case GDInt16:
		return GDComplex128(complex(float64(value), 0)), nil
	case GDFloat32:
		return GDComplex128(complex(value, 0)), nil
	case GDFloat64:
		return GDComplex128(complex(value, 0)), nil
	case GDComplex64:
		return GDComplex128(value), nil
	case GDComplex128:
		return value, nil
	case GDObject:
		return GDComplex128(0), InvalidCastingWrongTypeErr(GDComplexType, value.GetType())
	default:
		return GDComplex128(0), InvalidCastingExpectedTypeErr(GDComplexType)
	}
}
