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

import (
	"math"
	"math/bits"
	"strconv"
)

// General function to create a number object
// It will create the smallest possible number object
// based on the input type.
func NewGDNumber(number any) (GDObject, error) {
	switch {
	case IsComplex(number):
		n, err := ToComplex(number)
		if err != nil {
			return nil, err
		}
		return NewGDComplexNumber(n), nil
	case IsFloat(number):
		n, err := ToFloat(number)
		if err != nil {
			return nil, err
		}
		return NewGDFloatNumber(n), nil
	case IsInt(number):
		n, err := ToInt(number)
		if err != nil {
			return nil, err
		}
		return NewGDIntNumber(n), nil
	}
	return nil, nil
}

func NewGDIntNumber(value GDInt) GDObject {
	bitSize := IntBitsLen(int(value))
	switch {
	case bitSize < 8:
		return GDInt8(value)
	case bitSize < 16:
		return GDInt16(value)
	}
	return value
}

func NewGDIntNumberFromString(number string) (GDObject, error) {
	intValue, err := strconv.ParseInt(number, 0, 64)
	if err != nil {
		return nil, InvalidCastingLitErr(number, GDIntType)
	}
	return NewGDIntNumber(GDInt(intValue)), nil
}

func NewGDFloatNumber(value GDFloat64) GDObject {
	if FloatBitsLen(float64(value)) == 32 {
		return GDFloat32(value)
	}
	return value
}

func NewGDFloatNumberFromString(number string) (GDObject, error) {
	floatValue, err := strconv.ParseFloat(number, 64)
	if err != nil {
		return nil, InvalidCastingLitErr(number, GDFloatType)
	}
	return NewGDFloatNumber(GDFloat64(floatValue)), nil
}

func NewGDComplexNumber(value GDComplex128) GDObject {
	if ComplexBitsLen(complex128(value)) == 64 {
		return GDComplex64(value)
	}
	return value
}

func NewGDComplexNumberFromString(number string) (GDObject, error) {
	complexValue, err := strconv.ParseComplex(number, 128)
	if err != nil {
		return nil, InvalidCastingLitErr(number, GDComplexType)
	}
	return NewGDComplexNumber(GDComplex128(complexValue)), nil
}

// Int

type GDInt int

func (gd GDInt) GetType() GDTypable    { return GDIntType }
func (gd GDInt) GetSubType() GDTypable { return nil }
func (gd GDInt) ToString() string      { return strconv.FormatInt(int64(gd), 10) }
func (gd GDInt) CastToType(typ GDTypable, stack *GDSymbolStack) (GDObject, error) {
	switch typ := typ.(type) {
	case GDType:
		switch typ {
		case GDIntType:
			return gd, nil
		case GDInt8Type:
			return GDInt8(int8(gd)), nil
		case GDInt16Type:
			return GDInt16(int16(gd)), nil
		case GDFloatType:
			return NewGDFloatNumber(GDFloat64(gd)), nil
		case GDFloat32Type:
			return GDFloat32(gd), nil
		case GDFloat64Type:
			return GDFloat64(gd), nil
		case GDComplexType:
			return NewGDComplexNumber(GDComplex128(complex(float64(gd), 0))), nil
		case GDComplex64Type:
			return GDComplex64(complex(float32(gd), 0)), nil
		case GDComplex128Type:
			return GDComplex128(complex(float64(gd), 0)), nil
		case GDStringType:
			return GDString(gd.ToString()), nil
		case GDCharType:
			return GDChar(rune(gd)), nil
		case GDBoolType:
			return GDBool(gd != 0), nil
		}
	}

	return nil, InvalidCastingWrongTypeErr(typ, gd.GetType())
}

type GDInt8 int8

func (gd GDInt8) GetType() GDTypable    { return GDIntType }
func (gd GDInt8) GetSubType() GDTypable { return GDInt8Type }
func (gd GDInt8) ToString() string      { return strconv.FormatInt(int64(gd), 10) }
func (gd GDInt8) CastToType(typ GDTypable, stack *GDSymbolStack) (GDObject, error) {
	switch typ := typ.(type) {
	case GDType:
		switch typ {
		case GDIntType:
			return NewGDIntNumber(GDInt(int(gd))), nil
		case GDInt8Type:
			return gd, nil
		case GDInt16Type:
			return GDInt16(gd), nil
		case GDFloatType:
			return NewGDFloatNumber(GDFloat64(gd)), nil
		case GDFloat32Type:
			return GDFloat32(gd), nil
		case GDFloat64Type:
			return GDFloat64(gd), nil
		case GDComplexType:
			return NewGDComplexNumber(GDComplex128(complex(float64(gd), 0))), nil
		case GDComplex64Type:
			return GDComplex64(complex(float32(gd), 0)), nil
		case GDComplex128Type:
			return GDComplex128(complex(float64(gd), 0)), nil
		case GDStringType:
			return GDString(gd.ToString()), nil
		case GDCharType:
			return GDChar(gd), nil
		case GDBoolType:
			return GDBool(gd != 0), nil
		}
	}

	return nil, InvalidCastingWrongTypeErr(typ, gd.GetType())
}

type GDInt16 int16

func (gd GDInt16) GetType() GDTypable    { return GDIntType }
func (gd GDInt16) GetSubType() GDTypable { return GDInt16Type }
func (gd GDInt16) ToString() string      { return strconv.FormatInt(int64(gd), 10) }
func (gd GDInt16) CastToType(typ GDTypable, stack *GDSymbolStack) (GDObject, error) {
	switch typ := typ.(type) {
	case GDType:
		switch typ {
		case GDIntType:
			return NewGDIntNumber(GDInt(int(gd))), nil
		case GDInt8Type:
			return GDInt8(int8(gd)), nil
		case GDInt16Type:
			return gd, nil
		case GDFloatType:
			return NewGDFloatNumber(GDFloat64(gd)), nil
		case GDFloat32Type:
			return GDFloat32(gd), nil
		case GDFloat64Type:
			return GDFloat64(gd), nil
		case GDComplexType:
			return NewGDComplexNumber(GDComplex128(complex(float64(gd), 0))), nil
		case GDComplex64Type:
			return GDComplex64(complex(float32(gd), 0)), nil
		case GDComplex128Type:
			return GDComplex128(complex(float64(gd), 0)), nil
		case GDStringType:
			return GDString(gd.ToString()), nil
		case GDCharType:
			return GDChar(gd), nil
		case GDBoolType:
			return GDBool(gd != 0), nil
		}
	}

	return nil, InvalidCastingWrongTypeErr(typ, gd.GetType())
}

// Float

type GDFloat32 float32

func (gd GDFloat32) GetType() GDTypable    { return GDFloatType }
func (gd GDFloat32) GetSubType() GDTypable { return GDFloat32Type }
func (gd GDFloat32) ToString() string      { return strconv.FormatFloat(float64(gd), 'g', -1, 32) }
func (gd GDFloat32) CastToType(typ GDTypable, stack *GDSymbolStack) (GDObject, error) {
	switch typ := typ.(type) {
	case GDType:
		switch typ {
		case GDIntType:
			return NewGDIntNumber(GDInt(int(gd))), nil
		case GDInt8Type:
			return GDInt8(int8(gd)), nil
		case GDInt16Type:
			return GDInt16(int16(gd)), nil
		case GDFloatType:
			return NewGDFloatNumber(GDFloat64(gd)), nil
		case GDFloat32Type:
			return gd, nil
		case GDFloat64Type:
			return GDFloat64(gd), nil
		case GDComplexType:
			return NewGDComplexNumber(GDComplex128(complex(gd, 0))), nil
		case GDComplex64Type:
			return GDComplex64(complex(gd, 0)), nil
		case GDComplex128Type:
			return GDComplex128(complex(gd, 0)), nil
		case GDStringType:
			return GDString(gd.ToString()), nil
		case GDCharType:
			return GDChar(rune(gd)), nil
		case GDBoolType:
			return GDBool(gd != 0), nil
		}
	}

	return nil, InvalidCastingWrongTypeErr(typ, gd.GetType())
}

type GDFloat64 float64

func (gd GDFloat64) GetType() GDTypable    { return GDFloatType }
func (gd GDFloat64) GetSubType() GDTypable { return GDFloat64Type }
func (gd GDFloat64) ToString() string      { return strconv.FormatFloat(float64(gd), 'g', -1, 64) }
func (gd GDFloat64) CastToType(typ GDTypable, stack *GDSymbolStack) (GDObject, error) {
	switch typ := typ.(type) {
	case GDType:
		switch typ {
		case GDIntType:
			return NewGDIntNumber(GDInt(int(gd))), nil
		case GDInt8Type:
			return GDInt8(int8(gd)), nil
		case GDInt16Type:
			return GDInt16(int16(gd)), nil
		case GDFloatType:
			return NewGDFloatNumber(gd), nil
		case GDFloat32Type:
			return GDFloat32(gd), nil
		case GDFloat64Type:
			return gd, nil
		case GDComplexType:
			return NewGDComplexNumber(GDComplex128(complex(gd, 0))), nil
		case GDComplex64Type:
			return GDComplex64(complex(float32(gd), 0)), nil
		case GDComplex128Type:
			return GDComplex128(complex(gd, 0)), nil
		case GDStringType:
			return GDString(gd.ToString()), nil
		case GDCharType:
			return GDChar(rune(gd)), nil
		case GDBoolType:
			return GDBool(gd != 0), nil
		}
	}

	return nil, InvalidCastingWrongTypeErr(typ, gd.GetType())
}

// Complex

type GDComplex64 complex64

func (gd GDComplex64) GetType() GDTypable    { return GDComplexType }
func (gd GDComplex64) GetSubType() GDTypable { return GDComplex64Type }
func (gd GDComplex64) ToString() string      { return strconv.FormatComplex(complex128(gd), 'g', -1, 64) }
func (gd GDComplex64) CastToType(typ GDTypable, stack *GDSymbolStack) (GDObject, error) {
	switch typ := typ.(type) {
	case GDType:
		switch typ {
		case GDIntType:
			return NewGDIntNumber(GDInt(real(gd))), nil
		case GDInt8Type:
			return GDInt8(int8(real(gd))), nil
		case GDInt16Type:
			return GDInt16(int16(real(gd))), nil
		case GDFloatType:
			return NewGDFloatNumber(GDFloat64(real(gd))), nil
		case GDFloat32Type:
			return GDFloat32(real(gd)), nil
		case GDFloat64Type:
			return GDFloat64(real(gd)), nil
		case GDComplexType:
			return NewGDComplexNumber(GDComplex128(gd)), nil
		case GDComplex64Type:
			return gd, nil
		case GDComplex128Type:
			return GDComplex128(gd), nil
		case GDStringType:
			return GDString(gd.ToString()), nil
		case GDBoolType:
			return GDBool(gd != 0), nil
		}
	}

	return nil, InvalidCastingWrongTypeErr(typ, gd.GetType())
}

type GDComplex128 complex128

func (gd GDComplex128) GetType() GDTypable    { return GDComplexType }
func (gd GDComplex128) GetSubType() GDTypable { return GDComplex128Type }
func (gd GDComplex128) ToString() string      { return strconv.FormatComplex(complex128(gd), 'g', -1, 128) }
func (gd GDComplex128) CastToType(typ GDTypable, stack *GDSymbolStack) (GDObject, error) {
	switch typ := typ.(type) {
	case GDType:
		switch typ {
		case GDIntType:
			return NewGDIntNumber(GDInt(real(gd))), nil
		case GDInt8Type:
			return GDInt8(int8(real(gd))), nil
		case GDInt16Type:
			return GDInt16(int16(real(gd))), nil
		case GDFloatType:
			return NewGDFloatNumber(GDFloat64(real(gd))), nil
		case GDFloat32Type:
			return GDFloat32(real(gd)), nil
		case GDFloat64Type:
			return GDFloat64(real(gd)), nil
		case GDComplexType:
			return NewGDComplexNumber(gd), nil
		case GDComplex64Type:
			return GDComplex64(complex64(gd)), nil
		case GDComplex128Type:
			return gd, nil
		case GDStringType:
			return GDString(gd.ToString()), nil
		case GDBoolType:
			return GDBool(gd != 0), nil
		}
	}

	return nil, InvalidCastingWrongTypeErr(typ, gd.GetType())
}

func IntBitsLen(value int) int {
	// Len internally uses unit size of the arch
	if value < 0 {
		return bits.Len(uint(-value))
	}
	return bits.Len(uint(value))
}

func FloatBitsLen(value float64) int {
	// Constraint to 32 if arch does not support 64-bit floats
	if bits.UintSize == 32 {
		return 32
	}
	if value < math.MaxFloat32 && value > -math.MaxFloat32 {
		return 32
	}
	return 64
}

func ComplexBitsLen(value complex128) int {
	if FloatBitsLen(real(value)) == 32 && FloatBitsLen(imag(value)) == 32 {
		return 64
	}
	return 128
}
