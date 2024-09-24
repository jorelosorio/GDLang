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

package runtime_test

import (
	"gdlang/lib/runtime"
	"testing"
)

func TestIsInt(t *testing.T) {
	testCases := []struct {
		value any
		want  bool
	}{
		{1, true},
		{1.0, false},
		{runtime.GDInt(1), true},
		{runtime.GDInt8(1), true},
		{runtime.GDInt16(1), true},
		{"1", false},
	}

	for _, tc := range testCases {
		if runtime.IsInt(tc.value) != tc.want {
			t.Errorf("runtime.IsInt(%v) = %v, want %v", tc.value, !tc.want, tc.want)
		}
	}
}

func TestIsFloat(t *testing.T) {
	testCases := []struct {
		value interface{}
		want  bool
	}{
		{1, false},
		{1.0, true},
		{runtime.GDInt(1), false},
		{runtime.GDInt8(1), false},
		{runtime.GDInt16(1), false},
		{"1", false},
		{1.0, true},
		{runtime.GDFloat32(1.0), true},
		{runtime.GDFloat64(1.0), true},
	}

	for _, tc := range testCases {
		if runtime.IsFloat(tc.value) != tc.want {
			t.Errorf("runtime.IsFloat(%v) = %v, want %v", tc.value, !tc.want, tc.want)
		}
	}
}

func TestIsComplex(t *testing.T) {
	testCases := []struct {
		value interface{}
		want  bool
	}{
		{1, false},
		{1.0, false},
		{runtime.NewGDIntNumber(1), false},
		{runtime.NewGDIntNumber(1), false},
		{runtime.NewGDIntNumber(1), false},
		{"1", false},
		{runtime.NewGDFloatNumber(1.0), false},
		{runtime.NewGDFloatNumber(1.0), false},
		{runtime.NewGDFloatNumber(1.0), false},
		{runtime.NewGDComplexNumber(1.0), true},
		{runtime.NewGDComplexNumber(1.0), true},
	}

	for _, tc := range testCases {
		if runtime.IsComplex(tc.value) != tc.want {
			t.Errorf("runtime.IsComplex(%v) = %v, want %v", tc.value, !tc.want, tc.want)
		}
	}
}

func TestIsIntWithDifferentTypes(t *testing.T) {
	for i, test := range []struct {
		value any
		want  bool
	}{
		{uint(1), true}, // 1
		{uint8(1), true},
		{uint16(1), true},
		{uint32(1), true},
		{uint64(1), true},
		{1, true},
		{int8(1), true},
		{int16(1), true},
		{int32(1), true},
		{int64(1), true}, // 10
		{float32(1), false},
		{float64(1), false},
		{complex64(1), false},
		{complex128(1), false},
		{runtime.GDInt(1), true},
		{runtime.GDInt8(1), true},
		{runtime.GDInt16(1), true},
		{runtime.GDFloat32(1), false},
		{runtime.GDFloat64(1), false},
		{runtime.GDComplex64(1), false}, // 20
		{runtime.GDComplex128(1), false},
	} {
		resultValue := runtime.IsInt(test.value)
		if resultValue != test.want {
			t.Errorf("IsInt(%v) = %v, want %v for test case %d", test.value, resultValue, test.want, i+1)
		}
	}
}

func TestIsFloatWithDifferentTypes(t *testing.T) {
	for i, test := range []struct {
		value any
		want  bool
	}{
		{uint(1), false}, // 1
		{uint8(1), false},
		{uint16(1), false},
		{uint32(1), false},
		{uint64(1), false},
		{1, false},
		{int8(1), false},
		{int16(1), false},
		{int32(1), false},
		{int64(1), false}, // 10
		{float32(1), true},
		{float64(1), true},
		{complex64(1), false},
		{complex128(1), false},
		{runtime.GDInt(1), false},
		{runtime.GDInt8(1), false},
		{runtime.GDInt16(1), false},
		{runtime.GDFloat32(1), true},
		{runtime.GDFloat64(1), true},
		{runtime.GDComplex64(1), false}, // 20
		{runtime.GDComplex128(1), false},
	} {
		resultValue := runtime.IsFloat(test.value)
		if resultValue != test.want {
			t.Errorf("IsFloat(%v) = %v, want %v for test case %d", test.value, resultValue, test.want, i+1)
		}
	}
}

func TestIsComplexWithDifferentTypes(t *testing.T) {
	for i, test := range []struct {
		value any
		want  bool
	}{
		{uint(1), false}, // 1
		{uint8(1), false},
		{uint16(1), false},
		{uint32(1), false},
		{uint64(1), false},
		{1, false},
		{int8(1), false},
		{int16(1), false},
		{int32(1), false},
		{int64(1), false}, // 10
		{float32(1), false},
		{float64(1), false},
		{complex64(1), true},
		{complex128(1), true},
		{runtime.GDInt(1), false},
		{runtime.GDInt8(1), false},
		{runtime.GDInt16(1), false},
		{runtime.GDFloat32(1), false},
		{runtime.GDFloat64(1), false},
		{runtime.GDComplex64(1), true}, // 20
		{runtime.GDComplex128(1), true},
	} {
		resultValue := runtime.IsComplex(test.value)
		if resultValue != test.want {
			t.Errorf("IsComplex(%v) = %v, want %v for test case %d", test.value, resultValue, test.want, i+1)
		}
	}
}
