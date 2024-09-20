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
	"math"
	"reflect"
	"strconv"
	"testing"
)

func TestIntBytes(t *testing.T) {
	for _, test := range []struct {
		value    int
		expected int
	}{
		{0, 0},
		{1, 1},
		{math.MaxInt8, 7},
		{math.MaxInt16, 15},
		{math.MaxInt32, 31},
		{math.MaxInt64, 63},
	} {
		if result := runtime.IntBitsLen(test.value); result != test.expected {
			t.Errorf("Expected %d, but got %d", test.expected, result)
		}
	}
}

func TestSubTypeNumberKind(t *testing.T) {
	for _, test := range []struct {
		expectedKindOfVariable reflect.Kind
		value                  runtime.GDObject
		expectedValue          any
	}{
		{reflect.Int8, runtime.NewGDIntNumber(math.MaxInt8), runtime.GDInt8(math.MaxInt8)},
		{reflect.Int16, runtime.NewGDIntNumber(math.MaxInt16), runtime.GDInt16(math.MaxInt16)},
		{reflect.Float32, runtime.NewGDFloatNumber(math.MaxInt8), runtime.GDFloat32(math.MaxInt8)},
		{reflect.Float64, runtime.NewGDFloatNumber(math.MaxFloat32), runtime.GDFloat64(math.MaxFloat32)},
	} {
		reflectTestValueKind := reflect.TypeOf(test.value).Kind()
		if reflectTestValueKind != test.expectedKindOfVariable {
			t.Errorf("Expected runtime.GDObject kind to be %v, but got %v", test.expectedKindOfVariable, reflectTestValueKind)
		}

		if test.value != test.expectedValue {
			t.Errorf("Expected runtime.GDObject value to be %v (%v), but got %v (%v)", test.expectedValue, reflect.TypeOf(test.expectedValue).Kind(), test.value, reflectTestValueKind)
		}
	}
}

func TestNumberToString(t *testing.T) {
	for _, test := range []struct {
		value    runtime.GDObject
		expected string
	}{
		{runtime.NewGDIntNumber(math.MinInt8), "-128"},
		{runtime.NewGDIntNumber(-10), "-10"},
		{runtime.NewGDFloatNumber(3.1), "3.1"},
		{runtime.NewGDIntNumber(3.0), "3"},
		{runtime.NewGDFloatNumber(3.01), "3.01"},
		{runtime.NewGDIntNumber(math.MinInt), strconv.FormatInt(math.MinInt, 10)},
		{runtime.NewGDFloatNumber(math.SmallestNonzeroFloat32), "1e-45"},
		{runtime.NewGDFloatNumber(runtime.GDFloat64(math.Pi)), strconv.FormatFloat(math.Pi, 'g', -1, 32)},
		{runtime.NewGDFloatNumber(-2.71828), "-2.71828"},
		{runtime.NewGDFloatNumber(runtime.GDFloat64(math.Inf(1))), "+Inf"},
		{runtime.NewGDFloatNumber(runtime.GDFloat64(math.Inf(-1))), "-Inf"},
		{runtime.NewGDFloatNumber(runtime.GDFloat64(math.NaN())), "NaN"},
		{runtime.NewGDComplexNumber(complex(0, 0)), "(0+0i)"},
		{runtime.NewGDFloatNumber(math.MaxFloat64), strconv.FormatFloat(math.MaxFloat64, 'g', -1, 64)},
	} {
		str := test.value.ToString()
		if str != test.expected {
			t.Errorf("Expected %s, but got %s", test.expected, str)
		}
	}
}

func TestNumberKind(t *testing.T) {
	for _, test := range []struct {
		number       runtime.GDObject
		expectedType reflect.Kind
	}{
		{runtime.NewGDFloatNumber(runtime.GDFloat64(math.NaN())), reflect.Float64},
		{runtime.NewGDFloatNumber(runtime.GDFloat64(math.Inf(1))), reflect.Float64},
		{runtime.NewGDFloatNumber(runtime.GDFloat64(math.Inf(-1))), reflect.Float64},
		{runtime.NewGDFloatNumber(math.MinInt32), reflect.Float32},
		{runtime.NewGDFloatNumber(math.MinInt64), reflect.Float32},
		{runtime.NewGDFloatNumber(math.MaxFloat32), reflect.Float64},
		{runtime.NewGDFloatNumber(math.MaxFloat64), reflect.Float64},
		{runtime.NewGDFloatNumber(-10.0), reflect.Float32},
		{runtime.NewGDIntNumber(0), reflect.Int8},
		{runtime.NewGDIntNumber(10), reflect.Int8},
		{runtime.NewGDIntNumber(-10), reflect.Int8},
		{runtime.NewGDFloatNumber(-2.71828), reflect.Float32},
	} {
		numberKind := reflect.TypeOf(test.number).Kind()
		if numberKind != test.expectedType {
			t.Errorf("Expected value kind to be %v, but got %v", test.expectedType, numberKind)
		}
	}
}

func TestIntOverflows(t *testing.T) {
	for _, test := range []struct {
		number       runtime.GDObject
		expectedType reflect.Kind
	}{
		{runtime.NewGDIntNumber(math.MaxInt8 + 1), reflect.Int16},
		{runtime.NewGDIntNumber(math.MaxInt16 + 1), reflect.Int},
		{runtime.NewGDIntNumber(math.MaxInt32 + 1), reflect.Int},
	} {
		numberKind := reflect.TypeOf(test.number).Kind()
		if numberKind != test.expectedType {
			t.Errorf("Expected runtime.GDObject value kind to be %v, but got %v", test.expectedType, numberKind)
		}
	}
}

func TestFloatOverflow(t *testing.T) {
	for _, test := range []struct {
		number       runtime.GDObject
		expectedType reflect.Kind
	}{
		{runtime.NewGDFloatNumber(math.MaxFloat32 + math.MaxFloat32), reflect.Float64},
	} {
		numberKind := reflect.TypeOf(test.number).Kind()
		if numberKind != test.expectedType {
			t.Errorf("Expected runtime.GDObject value kind to be %v, but got %v", test.expectedType, numberKind)
		}
	}
}

func TestNumberFloatFromString(t *testing.T) {
	for _, test := range []struct {
		value    string
		kind     reflect.Kind
		expected runtime.GDObject
	}{
		{"10", reflect.Float32, runtime.NewGDFloatNumber(10.0)},
		{"-10", reflect.Float32, runtime.NewGDFloatNumber(-10.0)},
		{"3.1416", reflect.Float32, runtime.NewGDFloatNumber(3.1416)},
		{"-3.1416", reflect.Float32, runtime.NewGDFloatNumber(-3.1416)},
		{"0", reflect.Float32, runtime.NewGDFloatNumber(.0)},
		{"-0", reflect.Int8, runtime.NewGDIntNumber(0)},
		{"1e-10", reflect.Float32, runtime.NewGDFloatNumber(1e-10)},
		{"-1e-10", reflect.Float32, runtime.NewGDFloatNumber(-1e-10)},
		{"1e10", reflect.Float32, runtime.NewGDFloatNumber(1e10)},
		{"-1e10", reflect.Float32, runtime.NewGDFloatNumber(-1e10)},
		{"1e+10", reflect.Float32, runtime.NewGDFloatNumber(1e10)},
		{"-1e+10", reflect.Float32, runtime.NewGDFloatNumber(-1e10)},
		{"1e+10", reflect.Float32, runtime.NewGDFloatNumber(1e10)},
		{"3.1", reflect.Float32, runtime.NewGDFloatNumber(3.1)},
	} {
		var result runtime.GDObject
		var err error
		switch test.kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			result, err = runtime.NewGDIntNumberFromString(test.value)
		case reflect.Float32, reflect.Float64:
			result, err = runtime.NewGDFloatNumberFromString(test.value)
		}

		if err != nil {
			t.Errorf("Error while parsing number: %v", err)
		}

		if result != test.expected {
			t.Errorf("Expected %v, but got %v", test.expected.GetType().ToString(), result.GetType().ToString())
		}
	}
}

func TestNumberIntFromString(t *testing.T) {
	for _, test := range []struct {
		value    string
		expected runtime.GDObject
	}{
		{"10", runtime.NewGDIntNumber(10)},
		{"-10", runtime.NewGDIntNumber(-10)},
		{"0", runtime.NewGDIntNumber(0)},
		{"-0", runtime.NewGDIntNumber(0)},
		{"1", runtime.NewGDIntNumber(1)},
		{"-1", runtime.NewGDIntNumber(-1)},
		{"9223372036854775807", runtime.NewGDIntNumber(math.MaxInt64)},
		{"-9223372036854775808", runtime.NewGDIntNumber(math.MinInt64)},
	} {
		result, err := runtime.NewGDIntNumberFromString(test.value)
		if err != nil {
			t.Errorf("Error while parsing number: %v", err)
		}

		if result != test.expected {
			t.Errorf("Expected %v, but got %v", test.expected, result)
		}
	}
}
