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
	"strings"
	"testing"
)

func TestObjectEqualsWithNumbers(t *testing.T) {
	testCases := []struct {
		a, b runtime.GDObject
		want bool
	}{
		{runtime.NewGDIntNumber(1), runtime.NewGDIntNumber(1), true},
		{runtime.NewGDIntNumber(1), runtime.NewGDIntNumber(2), false},
		{runtime.NewGDIntNumber(1), runtime.GDString("1"), false},
		{runtime.NewGDIntNumber(1), runtime.GDString("2"), false},
		{runtime.NewGDFloatNumber(2.2), runtime.NewGDFloatNumber(2.2), true},
		{runtime.GDString("Hola"), runtime.GDString("Hola"), true},
		{runtime.GDString("Hola"), runtime.GDString("HOLA"), false},
		{runtime.GDZNil, runtime.GDZNil, true},
	}

	for _, tc := range testCases {
		if runtime.EqualObjects(tc.a, tc.b) != tc.want {
			t.Errorf("runtime.EqualObjects(%v, %v) = %v, want %v", tc.a, tc.b, !tc.want, tc.want)
		}
	}
}

func TestObjectEqualsWithTuples(t *testing.T) {
	testCases := []struct {
		a, b runtime.GDObject
		want bool
	}{
		{runtime.NewGDTuple(runtime.GDString("hola")), runtime.NewGDTuple(runtime.GDString("hola")), true},
		{runtime.NewGDTuple(runtime.GDZNil), runtime.NewGDTuple(runtime.GDString("Hola")), false},
		{runtime.NewGDTuple(runtime.GDZNil), runtime.NewGDTuple(runtime.GDZNil, runtime.GDZNil), false},
		{runtime.NewGDTuple(runtime.GDZNil), runtime.GDString("Hola"), false},
		{runtime.NewGDTuple(runtime.GDZNil, runtime.GDZNil), runtime.NewGDTuple(runtime.GDZNil, runtime.GDZNil), true},
	}

	for _, tc := range testCases {
		if runtime.EqualObjects(tc.a, tc.b) != tc.want {
			t.Errorf("runtime.EqualObjects(%v, %v) = %v, want %v", tc.a, tc.b, !tc.want, tc.want)
		}
	}
}

func TestObjectEqualsWithArrays(t *testing.T) {
	testCases := []struct {
		a, b runtime.GDObject
		want bool
	}{
		{runtime.NewGDArray(runtime.GDString("hola")), runtime.NewGDArray(runtime.GDString("hola")), true},
		{runtime.NewGDArray(runtime.GDZNil), runtime.NewGDArray(runtime.GDString("Hola")), false},
		{runtime.NewGDArray(runtime.GDZNil), runtime.NewGDArray(runtime.GDZNil, runtime.GDZNil), false},
		{runtime.NewGDArray(runtime.GDZNil), runtime.GDString("Hola"), false},
		{runtime.NewGDArray(runtime.GDZNil, runtime.GDZNil), runtime.NewGDArray(runtime.GDZNil, runtime.GDZNil), true},
		{runtime.NewGDArray(runtime.GDString("hola")), runtime.NewGDArray(runtime.GDString("hola")), true},
		{runtime.NewGDArray(runtime.GDString("hola")), runtime.NewGDArray(runtime.GDString("adios")), false},
		{runtime.NewGDArray(runtime.GDString("hola")), runtime.NewGDArray(runtime.GDString("hola"), runtime.GDString("adios")), false},
		{runtime.NewGDArray(runtime.GDString("hola")), runtime.NewGDArray(runtime.NewGDIntNumber(1)), false},
	}

	for _, tc := range testCases {
		if runtime.EqualObjects(tc.a, tc.b) != tc.want {
			t.Errorf("runtime.EqualObjects(%v, %v) = %v, want %v", tc.a, tc.b, !tc.want, tc.want)
		}
	}
}

func TestCastObject(t *testing.T) {
	stack := runtime.NewRootGDSymbolStack()
	defer stack.Dispose()

	userStruct, err := runtime.NewGDStruct(runtime.GDStructType{
		{runtime.NewGDStringIdent("name"), runtime.GDStringType},
		{runtime.NewGDStringIdent("age"), runtime.GDIntType},
	}, stack)
	if err != nil {
		t.Fatalf("Error creating user struct: %v", err)
	}

	testCases := []struct {
		obj    runtime.GDObject
		toType runtime.GDTypable
		errMsg string
	}{
		// String cases
		{runtime.GDString("hi"), runtime.GDStringType, ""},
		{runtime.GDString("1"), runtime.GDCharType, ""},
		{runtime.GDString("true"), runtime.GDBoolType, ""},
		{runtime.GDString("false"), runtime.GDBoolType, ""},
		{runtime.GDString("1"), runtime.GDBoolType, ""},
		{runtime.GDString("T"), runtime.GDIntType, "error trying to cast `T` into a `int`"},
		{runtime.GDString("1.0"), runtime.GDFloat32Type, ""}, // Float32, Float64 are FloatType (Both are compatible)
		{runtime.GDString(runtime.GDFloat64(math.MaxFloat64).ToString()), runtime.GDFloat64Type, ""},
		{runtime.GDString("1.0"), runtime.GDComplex64Type, ""}, // Complex64, Complex128 are ComplexType (Both are compatible)
		{runtime.GDString(runtime.GDComplex128(complex(math.MaxFloat64, math.MaxFloat64)).ToString()), runtime.GDComplex128Type, ""},
		{runtime.GDString("1.0"), runtime.GDIntType, "error trying to cast `1.0` into a `int`"},
		{runtime.GDString("1.0"), runtime.GDInt8Type, "error trying to cast `1.0` into a `int`"}, // Internally it is a int8
		{runtime.NewGDArray(runtime.GDString("hola")), runtime.GDStringType, ""},
		{runtime.NewGDArray(runtime.GDString("hola")), runtime.GDIntType, "error while casting `[string]` to `int`"},
		// Char cases
		{runtime.GDChar('h'), runtime.GDCharType, ""},
		{runtime.GDChar('o'), runtime.GDInt8Type, ""},
		{runtime.GDChar('l'), runtime.GDInt16Type, ""},
		{runtime.GDChar('a'), runtime.GDIntType, ""},
		{runtime.GDChar('a'), runtime.GDStringType, ""},
		{runtime.GDChar('a'), runtime.GDBoolType, ""},
		{runtime.GDChar('a'), runtime.GDFloat32Type, "error while casting `char` to `float32`"},
		{runtime.GDChar('a'), runtime.GDFloat64Type, "error while casting `char` to `float64`"},
		{runtime.GDChar('a'), runtime.GDComplex64Type, "error while casting `char` to `complex64`"},
		{runtime.GDChar('a'), runtime.GDComplex128Type, "error while casting `char` to `complex128`"},
		{runtime.GDChar('l'), runtime.GDStringType, ""},
		// Int cases
		{runtime.GDInt(math.MaxInt32), runtime.GDIntType, ""},
		{runtime.GDInt(1), runtime.GDInt8Type, ""},
		{runtime.GDInt(1), runtime.GDInt16Type, ""},
		{runtime.GDInt(1), runtime.GDFloat32Type, ""},
		{runtime.GDInt(1), runtime.GDFloat64Type, ""},
		{runtime.GDInt(1), runtime.GDBoolType, ""},
		{runtime.GDInt(1), runtime.GDStringType, ""},
		{runtime.GDInt(1), runtime.GDComplex64Type, ""},
		{runtime.GDInt(math.MaxInt64), runtime.GDComplex64Type, ""}, // Real part is always 64 bits
		// Float cases
		{runtime.GDFloat32(1.0), runtime.GDFloat32Type, ""},
		{runtime.GDFloat32(1.0), runtime.GDFloat64Type, ""},
		// Float is either Float32 or Float64 (There is no case for FloatType)
		// {runtime.GDFloat32(1.0), runtime.GDFloatType, ""},
		{runtime.GDFloat32(math.MaxFloat32), runtime.GDIntType, ""},
		{runtime.GDFloat32(1.0), runtime.GDInt8Type, ""},
		{runtime.GDFloat32(1.0), runtime.GDInt16Type, ""},
		{runtime.GDFloat32(1.0), runtime.GDBoolType, ""},
		{runtime.GDFloat32(1.0), runtime.GDStringType, ""},
		{runtime.GDFloat32(1.0), runtime.GDComplex64Type, ""},
		// Complex type is either Complex64 or Complex128 (There is no case for ComplexType)
		// {runtime.GDFloat32(1.0), runtime.GDComplexType, ""},
		// Complex cases
		{runtime.GDComplex64(1.0), runtime.GDComplex64Type, ""},
		{runtime.GDComplex64(1.0), runtime.GDComplex128Type, ""},
		// {runtime.GDComplex64(1.0), runtime.GDComplexType, ""},
		{runtime.GDComplex64(math.MaxFloat32), runtime.GDIntType, ""},
		{runtime.GDComplex64(1.0), runtime.GDInt8Type, ""},
		{runtime.GDComplex64(1.0), runtime.GDInt16Type, ""},
		{runtime.GDComplex64(1.0), runtime.GDBoolType, ""},
		{runtime.GDComplex64(1.0), runtime.GDStringType, ""},
		{runtime.GDComplex64(1.0), runtime.GDFloat32Type, ""},
		// Float requires 2 float32 values, but it uses only the real part
		// {runtime.GDComplex64(math.MaxFloat32), runtime.GDFloatType, ""},
		// Bool cases
		{runtime.GDBool(true), runtime.GDBoolType, ""},
		{runtime.GDBool(true), runtime.GDIntType, ""},
		{runtime.GDBool(true), runtime.GDInt8Type, ""},
		{runtime.GDBool(true), runtime.GDInt16Type, ""},
		{runtime.GDBool(true), runtime.GDFloat32Type, ""},
		{runtime.GDBool(true), runtime.GDFloat64Type, ""},
		{runtime.GDBool(true), runtime.GDComplex64Type, ""},
		{runtime.GDBool(true), runtime.GDComplex128Type, ""},
		// 1 fits in complex64
		// {runtime.GDBool(true), runtime.GDComplexType, ""},
		{runtime.GDBool(true), runtime.GDStringType, ""},
		{runtime.GDBool(true), runtime.GDCharType, ""},
		// Array cases
		{runtime.NewGDArray(runtime.GDChar('h')), runtime.NewGDArrayType(runtime.GDStringType), ""},
		{runtime.NewGDArray(runtime.GDString("hi")), runtime.NewGDArrayType(runtime.GDCharType), "error trying to cast `hi` into a `char`"},
		{runtime.NewGDArray(runtime.GDString("h")), runtime.NewGDArrayType(runtime.GDCharType), ""},
		{runtime.NewGDArray(runtime.NewGDArray(runtime.GDString("h"))), runtime.NewGDArrayType(runtime.GDStringType), ""},
		{runtime.NewGDArray(runtime.NewGDArray(runtime.GDString("h"))), runtime.NewGDArrayType(runtime.GDIntType), "error while casting `[string]` to `int`"},
		{runtime.NewGDArray(), runtime.NewGDArrayType(runtime.GDIntType), ""},
		// Tuple cases
		{runtime.NewGDTuple(runtime.GDChar('h')), runtime.NewGDTupleType(runtime.GDStringType), ""},
		// Struct cases
		{userStruct, runtime.GDStructType{
			{runtime.NewGDStringIdent("name"), runtime.GDStringType},
			{runtime.NewGDStringIdent("age"), runtime.GDIntType},
		}, ""},
		{userStruct, runtime.GDStructType{
			{runtime.NewGDStringIdent("name"), runtime.GDStringType},
		}, "attribute `age`, not found"},
		{userStruct, runtime.GDStructType{
			{runtime.NewGDStringIdent("name"), runtime.GDStringType},
			{runtime.NewGDStringIdent("age"), runtime.GDStringType},
		}, ""},
	}

	for _, tc := range testCases {
		obj, err := tc.obj.CastToType(tc.toType, stack)
		if err != nil {
			if tc.errMsg != "" {
				if !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("Expected error to contain %q but got %q", tc.errMsg, err.Error())
				}
				continue
			} else {
				t.Errorf("And error occurred casting %q to %q: %q", tc.obj.ToString(), tc.toType.ToString(), err)
			}
		} else if tc.errMsg != "" {
			t.Errorf("An error was expected casting %q to %q but got no error", tc.obj.ToString(), tc.toType.ToString())
		}

		objType := obj.GetType()
		if obj.GetSubType() != nil {
			objType = obj.GetSubType()
		}

		err = runtime.EqualTypes(objType, tc.toType, stack)
		if err != nil {
			t.Errorf("CastObject(%q) expect %q, but got %q", tc.obj.ToString(), tc.toType.ToString(), objType.ToString())
		}
	}
}
