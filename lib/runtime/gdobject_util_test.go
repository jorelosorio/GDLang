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
		{runtime.NewGDArray(nil, runtime.GDString("hola")), runtime.NewGDArray(nil, runtime.GDString("hola")), true},
		{runtime.NewGDArray(nil, runtime.GDZNil), runtime.NewGDArray(nil, runtime.GDString("Hola")), false},
		{runtime.NewGDArray(nil, runtime.GDZNil), runtime.NewGDArray(nil, runtime.GDZNil, runtime.GDZNil), false},
		{runtime.NewGDArray(nil, runtime.GDZNil), runtime.GDString("Hola"), false},
		{runtime.NewGDArray(nil, runtime.GDZNil, runtime.GDZNil), runtime.NewGDArray(nil, runtime.GDZNil, runtime.GDZNil), true},
		{runtime.NewGDArray(nil, runtime.GDString("hola")), runtime.NewGDArray(nil, runtime.GDString("hola")), true},
		{runtime.NewGDArray(nil, runtime.GDString("hola")), runtime.NewGDArray(nil, runtime.GDString("adios")), false},
		{runtime.NewGDArray(nil, runtime.GDString("hola")), runtime.NewGDArray(nil, runtime.GDString("hola"), runtime.GDString("adios")), false},
		{runtime.NewGDArray(nil, runtime.GDString("hola")), runtime.NewGDArray(nil, runtime.NewGDIntNumber(1)), false},
	}

	for _, tc := range testCases {
		if runtime.EqualObjects(tc.a, tc.b) != tc.want {
			t.Errorf("runtime.EqualObjects(%v, %v) = %v, want %v", tc.a, tc.b, !tc.want, tc.want)
		}
	}
}

func TestCastObject(t *testing.T) {
	stack := runtime.NewGDStack()
	defer stack.Dispose()

	userStruct, err := runtime.NewGDStruct(runtime.GDStructType{
		{runtime.NewGDStrIdent("name"), runtime.GDStringTypeRef},
		{runtime.NewGDStrIdent("age"), runtime.GDIntTypeRef},
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
		{runtime.GDString("hi"), runtime.GDStringTypeRef, ""},
		{runtime.GDString("1"), runtime.GDCharTypeRef, ""},
		{runtime.GDString("true"), runtime.GDBoolTypeRef, ""},
		{runtime.GDString("false"), runtime.GDBoolTypeRef, ""},
		{runtime.GDString("1"), runtime.GDBoolTypeRef, ""},
		{runtime.GDString("T"), runtime.GDIntTypeRef, "error trying to cast `T` into a `int`"},
		{runtime.GDString("1.0"), runtime.GDFloat32TypeRef, ""}, // Float32, Float64 are FloatType (Both are compatible)
		{runtime.GDString(runtime.GDFloat64(math.MaxFloat64).ToString()), runtime.GDFloat64TypeRef, ""},
		{runtime.GDString("1.0"), runtime.GDComplex64TypeRef, ""}, // Complex64, Complex128 are ComplexType (Both are compatible)
		{runtime.GDString(runtime.GDComplex128(complex(math.MaxFloat64, math.MaxFloat64)).ToString()), runtime.GDComplex128TypeRef, ""},
		{runtime.GDString("1.0"), runtime.GDIntTypeRef, "error trying to cast `1.0` into a `int`"},
		{runtime.GDString("1.0"), runtime.GDInt8TypeRef, "error trying to cast `1.0` into a `int`"}, // Internally it is an int8
		{runtime.NewGDArray(nil, runtime.GDString("hola")), runtime.GDStringTypeRef, ""},
		{runtime.NewGDArray(nil, runtime.GDString("hola")), runtime.GDIntTypeRef, "error while casting `[string]` to `int`"},
		// Char cases
		{runtime.GDChar('h'), runtime.GDCharTypeRef, ""},
		{runtime.GDChar('o'), runtime.GDInt8TypeRef, ""},
		{runtime.GDChar('l'), runtime.GDInt16TypeRef, ""},
		{runtime.GDChar('a'), runtime.GDIntTypeRef, ""},
		{runtime.GDChar('a'), runtime.GDStringTypeRef, ""},
		{runtime.GDChar('a'), runtime.GDBoolTypeRef, ""},
		{runtime.GDChar('a'), runtime.GDFloat32TypeRef, "error while casting `char` to `float32`"},
		{runtime.GDChar('a'), runtime.GDFloat64TypeRef, "error while casting `char` to `float64`"},
		{runtime.GDChar('a'), runtime.GDComplex64TypeRef, "error while casting `char` to `complex64`"},
		{runtime.GDChar('a'), runtime.GDComplex128TypeRef, "error while casting `char` to `complex128`"},
		{runtime.GDChar('l'), runtime.GDStringTypeRef, ""},
		// Int cases
		{runtime.GDInt(math.MaxInt32), runtime.GDIntTypeRef, ""},
		{runtime.GDInt(1), runtime.GDInt8TypeRef, ""},
		{runtime.GDInt(1), runtime.GDInt16TypeRef, ""},
		{runtime.GDInt(1), runtime.GDFloat32TypeRef, ""},
		{runtime.GDInt(1), runtime.GDFloat64TypeRef, ""},
		{runtime.GDInt(1), runtime.GDBoolTypeRef, ""},
		{runtime.GDInt(1), runtime.GDStringTypeRef, ""},
		{runtime.GDInt(1), runtime.GDComplex64TypeRef, ""},
		{runtime.GDInt(math.MaxInt64), runtime.GDComplex64TypeRef, ""}, // Real part is always 64 bits
		// Float cases
		{runtime.GDFloat32(1.0), runtime.GDFloat32TypeRef, ""},
		{runtime.GDFloat32(1.0), runtime.GDFloat64TypeRef, ""},
		// Float is either Float32 or Float64 (There is no case for FloatType)
		{runtime.GDFloat32(1.0), runtime.GDFloat32TypeRef, ""},
		{runtime.GDFloat32(math.MaxFloat32), runtime.GDIntTypeRef, ""},
		{runtime.GDFloat32(1.0), runtime.GDInt8TypeRef, ""},
		{runtime.GDFloat32(1.0), runtime.GDInt16TypeRef, ""},
		{runtime.GDFloat32(1.0), runtime.GDBoolTypeRef, ""},
		{runtime.GDFloat32(1.0), runtime.GDStringTypeRef, ""},
		{runtime.GDFloat32(1.0), runtime.GDComplex64TypeRef, ""},
		// Complex type is either Complex64 or Complex128 (There is no case for ComplexType)
		{runtime.GDFloat32(1.0), runtime.GDComplex64TypeRef, ""},
		// Complex cases
		{runtime.GDComplex64(1.0), runtime.GDComplex64TypeRef, ""},
		{runtime.GDComplex64(1.0), runtime.GDComplex128TypeRef, ""},
		{runtime.GDComplex64(math.MaxFloat32), runtime.GDIntTypeRef, ""},
		{runtime.GDComplex64(1.0), runtime.GDInt8TypeRef, ""},
		{runtime.GDComplex64(1.0), runtime.GDInt16TypeRef, ""},
		{runtime.GDComplex64(1.0), runtime.GDBoolTypeRef, ""},
		{runtime.GDComplex64(1.0), runtime.GDStringTypeRef, ""},
		{runtime.GDComplex64(1.0), runtime.GDFloat32TypeRef, ""},
		// Float requires 2 float32 values, but it uses only the real part
		{runtime.GDComplex64(math.MaxFloat32), runtime.GDFloat64TypeRef, ""},
		// Bool cases
		{runtime.GDBool(true), runtime.GDBoolTypeRef, ""},
		{runtime.GDBool(true), runtime.GDIntTypeRef, ""},
		{runtime.GDBool(true), runtime.GDInt8TypeRef, ""},
		{runtime.GDBool(true), runtime.GDInt16TypeRef, ""},
		{runtime.GDBool(true), runtime.GDFloat32TypeRef, ""},
		{runtime.GDBool(true), runtime.GDFloat64TypeRef, ""},
		{runtime.GDBool(true), runtime.GDComplex64TypeRef, ""},
		{runtime.GDBool(true), runtime.GDComplex128TypeRef, ""},
		// 1 fits in complex64
		{runtime.GDBool(true), runtime.GDComplex64TypeRef, ""},
		{runtime.GDBool(true), runtime.GDStringTypeRef, ""},
		{runtime.GDBool(true), runtime.GDCharTypeRef, ""},
		// Array cases
		{runtime.NewGDArray(nil, runtime.GDChar('h')), runtime.NewGDArrayType(runtime.GDStringTypeRef), ""},
		{runtime.NewGDArray(nil, runtime.GDString("hi")), runtime.NewGDArrayType(runtime.GDCharTypeRef), "error trying to cast `hi` into a `char`"},
		{runtime.NewGDArray(nil, runtime.GDString("h")), runtime.NewGDArrayType(runtime.GDCharTypeRef), ""},
		{runtime.NewGDArray(nil, runtime.NewGDArray(nil, runtime.GDString("h"))), runtime.NewGDArrayType(runtime.GDStringTypeRef), ""},
		{runtime.NewGDArray(nil, runtime.NewGDArray(nil, runtime.GDString("h"))), runtime.NewGDArrayType(runtime.GDIntTypeRef), "error while casting `[string]` to `int`"},
		{runtime.NewGDArray(nil), runtime.NewGDArrayType(runtime.GDIntTypeRef), ""},
		// Tuple cases
		{runtime.NewGDTuple(runtime.GDChar('h')), runtime.NewGDTupleType(runtime.GDStringTypeRef), ""},
		// Struct cases
		{userStruct, runtime.GDStructType{
			{runtime.NewGDStrIdent("name"), runtime.GDStringTypeRef},
			{runtime.NewGDStrIdent("age"), runtime.GDIntTypeRef},
		}, ""},
		{userStruct, runtime.GDStructType{
			{runtime.NewGDStrIdent("name"), runtime.GDStringTypeRef},
		}, "attribute `age`, not found"},
		{userStruct, runtime.GDStructType{
			{runtime.NewGDStrIdent("name"), runtime.GDStringTypeRef},
			{runtime.NewGDStrIdent("age"), runtime.GDStringTypeRef},
		}, ""},
	}

	for _, tc := range testCases {
		obj, err := tc.obj.CastToType(tc.toType)
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

		err = runtime.EqualTypes(objType, tc.toType, nil)
		if err != nil {
			t.Errorf("CastObject(%q) expect %q, but got %q", tc.obj.ToString(), tc.toType.ToString(), objType.ToString())
		}
	}
}
