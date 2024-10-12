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
	"fmt"
	"gdlang/lib/runtime"
	"strings"
	"testing"
)

type TypeTest struct {
	toType, fromType runtime.GDTypable
	expected         runtime.GDTypable
	errMsg           string
}

func (t TypeTest) ToString() string {
	return t.toType.ToString() + " & " + t.fromType.ToString()
}

// Computed type determination with object types

func TestComputeTypeDeterminationWithNonObjectTypes(t *testing.T) {
	computedType := runtime.ComputeTypeFromTypes([]runtime.GDTypable{})
	if computedType != runtime.GDUntypedTypeRef {
		t.Errorf("Computed type should be untyped but got %s", computedType.ToString())
	}
}

func TestComputeTypeDeterminationWithObjectTypes(t *testing.T) {
	computedType := runtime.ComputeTypeFromObjects([]runtime.GDObject{runtime.GDString("test"), runtime.GDString("test")}, nil)
	if computedType != runtime.GDStringTypeRef {
		t.Errorf("Computed type should be (string) but got %s", computedType.ToString())
	}
}

// Test compute type with different types

func TestComputeTypeWithDifferentTypes(t *testing.T) {
	for _, test := range []struct {
		types    []runtime.GDTypable
		expected runtime.GDTypable
	}{
		{[]runtime.GDTypable{runtime.GDStringTypeRef, runtime.GDStringTypeRef}, runtime.GDStringTypeRef},
		{[]runtime.GDTypable{runtime.GDStringTypeRef, runtime.GDBoolTypeRef}, runtime.NewGDUnionType(runtime.GDStringTypeRef, runtime.GDBoolTypeRef)},
		{[]runtime.GDTypable{runtime.NewGDArrayType(runtime.GDStringTypeRef), runtime.NewGDArrayType(runtime.GDBoolTypeRef)}, runtime.NewGDUnionType(runtime.NewGDArrayType(runtime.GDStringTypeRef), runtime.NewGDArrayType(runtime.GDBoolTypeRef))},
		{[]runtime.GDTypable{runtime.NewGDArrayType(runtime.GDStringTypeRef), runtime.NewGDArrayType(runtime.GDStringTypeRef)}, runtime.NewGDArrayType(runtime.GDStringTypeRef)},
		{[]runtime.GDTypable{
			runtime.NewGDArrayType(
				runtime.NewGDArrayType(
					runtime.NewGDArrayType(
						runtime.NewGDArrayType(runtime.GDStringTypeRef),
					),
				),
			), runtime.NewGDArrayType(
				runtime.NewGDArrayType(
					runtime.NewGDArrayType(
						runtime.NewGDArrayType(runtime.GDBoolTypeRef),
					),
				),
			),
		}, runtime.NewGDUnionType(
			runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringTypeRef)))),
			runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDBoolTypeRef)))),
		)},
	} {
		t.Run(test.expected.ToString(), func(t *testing.T) {
			computedType := runtime.ComputeTypeFromTypes(test.types)
			if err := runtime.EqualTypes(computedType, test.expected, nil); err != nil {
				t.Errorf("Expected %q but got %q", test.expected.ToString(), computedType.ToString())
				return
			}
		})
	}
}

// Equal types

func TestEqualTypes(t *testing.T) {
	tests := []TypeTest{
		{runtime.GDStringTypeRef, runtime.GDStringTypeRef, nil, ""},
		{runtime.GDStringTypeRef, runtime.GDBoolTypeRef, nil, "types `string` and `bool` are not equal"},
		{runtime.NewGDArrayType(runtime.GDStringTypeRef), runtime.NewGDArrayType(runtime.GDStringTypeRef), nil, ""},
		{runtime.NewGDArrayType(runtime.GDStringTypeRef), runtime.NewGDArrayType(runtime.GDBoolTypeRef), nil, "types `[string]` and `[bool]` are not equal"},
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringTypeRef))), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringTypeRef))), nil, ""},
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringTypeRef))), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDBoolTypeRef))), nil, "types `[[[string]]]` and `[[[bool]]]` are not equal"},
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringTypeRef))), runtime.GDStringTypeRef, nil, "types `[[[string]]]` and `string` are not equal"},
		{runtime.GDStringTypeRef, runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringTypeRef))), nil, "types `string` and `[[[string]]]` are not equal"},
		{runtime.NewGDArrayType(runtime.GDStringTypeRef), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringTypeRef))), nil, "types `[string]` and `[[[string]]]` are not equal"},
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringTypeRef))), runtime.NewGDArrayType(runtime.GDStringTypeRef), nil, "types `[[[string]]]` and `[string]` are not equal"},
		{runtime.NewGDArrayType(runtime.GDAnyTypeRef), runtime.NewGDArrayType(runtime.GDStringTypeRef), nil, "types `[any]` and `[string]` are not equal"},
		{structWithAttrAAsInt, structWithAttrAAsInt, nil, ""},
		{structWithAttrAAsString, structWithAttrAAsInt, nil, "types `{a: string}` and `{a: int}` are not equal"},
		{structWithAttrsAStrBInt, structWithAttrsBIntAStr, nil, ""},
		{runtime.NewGDTuple(runtime.NewGDIntNumber(1), runtime.GDString("test")).GetType(), runtime.NewGDTupleType(runtime.GDIntTypeRef, runtime.GDStringTypeRef), nil, ""},
		{runtime.GDUntypedTypeRef, runtime.GDUntypedTypeRef, nil, ""},
		{runtime.GDUntypedTypeRef, runtime.GDStringTypeRef, nil, "types `untyped` and `string` are not equal"},
		{runtime.GDAnyTypeRef, runtime.GDStringTypeRef, nil, "types `any` and `string` are not equal"},
		{runtime.NewGDTupleType(runtime.GDIntTypeRef, runtime.NewGDTupleType(runtime.GDIntTypeRef, runtime.GDStringTypeRef)), runtime.NewGDTupleType(runtime.GDIntTypeRef, runtime.NewGDTupleType(runtime.GDIntTypeRef, runtime.GDStringTypeRef)), nil, ""},
		{runtime.NewGDTupleType(runtime.GDStringTypeRef, runtime.GDStringTypeRef), runtime.NewGDTupleType(runtime.GDStringTypeRef, runtime.GDStringTypeRef), nil, ""},
		{runtime.NewGDTupleType(runtime.GDIntTypeRef, runtime.NewGDTupleType(runtime.GDIntTypeRef, runtime.GDIntTypeRef)), runtime.NewGDTupleType(runtime.GDIntTypeRef, runtime.NewGDTupleType(runtime.GDIntTypeRef, runtime.GDStringTypeRef)), nil, "types `(int, (int, int))` and `(int, (int, string))` are not equal"},
		{runtime.NewGDUnionType(runtime.GDIntTypeRef, runtime.GDStringTypeRef), runtime.NewGDUnionType(runtime.GDIntTypeRef, runtime.GDStringTypeRef), nil, ""},
		{runtime.NewGDUnionType(runtime.GDStringTypeRef, runtime.GDIntTypeRef), runtime.NewGDUnionType(runtime.GDIntTypeRef, runtime.GDStringTypeRef), nil, ""},
		{runtime.GDStringTypeRef, runtime.NewGDUnionType(runtime.GDIntTypeRef, runtime.GDStringTypeRef), nil, "types `string` and `(int | string)` are not equal"},
		{runtime.NewGDUnionType(runtime.GDIntTypeRef, runtime.GDStringTypeRef), runtime.GDStringTypeRef, nil, ""},
	}

	TypeTests(t, tests, func(t *testing.T, test TypeTest) error {
		return runtime.EqualTypes(test.toType, test.fromType, nil)
	})
}

// Assignments

func TestCanBeAssigned(t *testing.T) {
	tests := []TypeTest{
		{runtime.GDStringTypeRef, runtime.GDStringTypeRef, nil, ""},
		{runtime.GDStringTypeRef, runtime.GDBoolTypeRef, nil, "expected `string` but got `bool`"},
		{runtime.NewGDArrayType(runtime.GDStringTypeRef), runtime.NewGDArrayType(runtime.GDStringTypeRef), nil, ""},
		{runtime.NewGDArrayType(runtime.GDStringTypeRef), runtime.NewGDArrayType(runtime.GDBoolTypeRef), nil, "expected `[string]` but got `[bool]`"},
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringTypeRef))), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringTypeRef))), nil, ""},
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringTypeRef))), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDBoolTypeRef))), nil, "expected `[[[string]]]` but got `[[[bool]]]`"},
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringTypeRef))), runtime.GDStringTypeRef, nil, "expected `[[[string]]]` but got `string`"},
		{runtime.GDStringTypeRef, runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringTypeRef))), nil, "expected `string` but got `[[[string]]]`"},
		{runtime.NewGDArrayType(runtime.GDStringTypeRef), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringTypeRef))), nil, "expected `[string]` but got `[[[string]]]`"},
		{runtime.GDAnyTypeRef, runtime.NewGDArrayType(runtime.GDStringTypeRef), nil, ""},
		{runtime.NewGDArrayType(runtime.GDStringTypeRef), runtime.GDAnyTypeRef, nil, "expected `[string]` but got `any`"},
		{runtime.NewGDArrayType(runtime.GDStringTypeRef), runtime.GDStringTypeRef, nil, "expected `[string]` but got `string`"},
		{runtime.NewGDTupleType(runtime.GDStringTypeRef), runtime.GDNilTypeRef, nil, ""},
	}

	TypeTests(t, tests, func(t *testing.T, test TypeTest) error {
		err := runtime.CanBeAssign(test.toType, test.fromType, nil)
		if err != nil {
			return err
		}

		return nil
	})
}

func TestCanBeAppendedTo(t *testing.T) {
	tests := []TypeTest{
		{runtime.NewGDArrayType(runtime.GDStringTypeRef), runtime.GDStringTypeRef, nil, ""},
		{runtime.NewGDArrayType(runtime.GDStringTypeRef), runtime.GDAnyTypeRef, nil, "expected `string` but got `any`"},
		{runtime.NewGDArrayType(runtime.GDAnyTypeRef), runtime.GDStringTypeRef, nil, ""},
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringTypeRef)), runtime.GDStringTypeRef, nil, "expected `[string]` but got `string`"},
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringTypeRef)), runtime.NewGDArrayType(runtime.GDStringTypeRef), nil, ""},
		{runtime.NewGDArrayType(runtime.GDAnyTypeRef), runtime.NewGDArrayType(runtime.NewGDUnionType(runtime.GDStringTypeRef, runtime.GDIntTypeRef)), nil, ""},
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDAnyTypeRef)), runtime.NewGDArrayType(runtime.NewGDUnionType(runtime.NewGDArrayType(runtime.GDStringTypeRef), runtime.NewGDArrayType(runtime.GDIntTypeRef))), nil, ""},
	}

	TypeTests(t, tests, func(t *testing.T, test TypeTest) error {
		err := runtime.CanBeAssign(test.toType.(*runtime.GDArrayType).SubType, test.fromType, nil)
		if err != nil {
			return err
		}

		return nil
	})
}

// Infer type

func TestTypeInference(t *testing.T) {
	stack := runtime.NewGDStack()

	structType1 := runtime.NewGDStructType(&runtime.GDStructAttrType{aParamIdent, runtime.GDIntTypeRef})
	structType2 := runtime.NewGDStructType(&runtime.GDStructAttrType{bParamIdent, runtime.GDIntTypeRef})
	func1 := runtime.NewGDLambdaType(runtime.GDLambdaArgTypes{}, runtime.GDAnyTypeRef, false)
	func2 := runtime.NewGDLambdaType(runtime.GDLambdaArgTypes{}, runtime.GDStringTypeRef, false)
	func3 := runtime.NewGDLambdaType(runtime.GDLambdaArgTypes{}, runtime.GDIntTypeRef, false)

	// Register a type alias
	typeAliasIdent := runtime.NewGDStrIdent("typ")
	_, err := stack.AddNewSymbol(typeAliasIdent, true, true, runtime.NewGDTypeAliasType(runtime.NewGDStrIdent("typ"), runtime.GDStringTypeRef), nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	tests := []TypeTest{
		// set a: any = int | float
		{runtime.GDAnyTypeRef, runtime.NewGDUnionType(runtime.GDIntTypeRef, runtime.GDFloatTypeRef), runtime.GDAnyTypeRef, ""},
		// set a: any = any
		{runtime.GDAnyTypeRef, runtime.GDAnyTypeRef, runtime.GDAnyTypeRef, ""},
		// set a: any = [int]
		{runtime.GDAnyTypeRef, runtime.NewGDArrayType(runtime.GDIntTypeRef), runtime.GDAnyTypeRef, ""},
		// set a: any = func() => any
		{runtime.GDAnyTypeRef, func1, runtime.GDAnyTypeRef, ""},
		// set a: string = any
		{runtime.GDStringTypeRef, runtime.GDAnyTypeRef, nil, "expected `string` but got `any`"},
		// set a: string = untyped
		{runtime.GDStringTypeRef, runtime.GDUntypedTypeRef, runtime.GDStringTypeRef, ""},
		// set a: string = nil
		{runtime.GDStringTypeRef, runtime.GDNilTypeRef, runtime.GDStringTypeRef, ""},
		// set a: [untyped] = untyped
		{runtime.NewGDEmptyArrayType(), runtime.GDUntypedTypeRef, runtime.NewGDEmptyArrayType(), ""},
		// set a: untyped = [untyped]
		{runtime.GDUntypedTypeRef, runtime.NewGDEmptyArrayType(), runtime.NewGDEmptyArrayType(), ""},
		// set a: [untyped] = [untyped]
		{runtime.NewGDEmptyArrayType(), runtime.NewGDEmptyArrayType(), runtime.NewGDEmptyArrayType(), ""},
		// set a: [untyped] = [[untyped]] (For this case `untyped` on the left side is a `[untyped]`)
		{runtime.NewGDEmptyArrayType(), runtime.NewGDArrayType(runtime.NewGDEmptyArrayType()), runtime.NewGDArrayType(runtime.NewGDEmptyArrayType()), ""},
		// set a: [[untyped]] = [untyped]
		{runtime.NewGDArrayType(runtime.NewGDEmptyArrayType()), runtime.NewGDEmptyArrayType(), runtime.NewGDArrayType(runtime.NewGDEmptyArrayType()), ""},
		// set a: [[string]] = [[untyped]]
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringTypeRef)), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDUntypedTypeRef)), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringTypeRef)), ""},
		// set a: [string] = [untyped]
		{runtime.NewGDArrayType(runtime.GDStringTypeRef), runtime.NewGDEmptyArrayType(), runtime.NewGDArrayType(runtime.GDStringTypeRef), ""},
		// set a: [untyped] = [string]
		{runtime.NewGDEmptyArrayType(), runtime.NewGDArrayType(runtime.GDStringTypeRef), runtime.NewGDArrayType(runtime.GDStringTypeRef), ""},
		// set a: [any] = nil
		{runtime.NewGDArrayType(runtime.GDAnyTypeRef), runtime.GDNilTypeRef, runtime.NewGDArrayType(runtime.GDAnyTypeRef), ""},
		// set a: [int] = [nil]
		{runtime.NewGDArrayType(runtime.GDIntTypeRef), runtime.NewGDArrayType(runtime.GDNilTypeRef), runtime.NewGDArrayType(runtime.GDIntTypeRef), ""},
		// set a: [string] = any
		{runtime.NewGDArrayType(runtime.GDStringTypeRef), runtime.GDAnyTypeRef, nil, "expected `[string]` but got `any`"},
		// set a: [[nil]] = untyped
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDNilTypeRef)), runtime.GDUntypedTypeRef, runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDNilTypeRef)), ""},
		// set a: untyped = [string]
		{runtime.GDUntypedTypeRef, runtime.NewGDArrayType(runtime.GDStringTypeRef), runtime.NewGDArrayType(runtime.GDStringTypeRef), ""},
		// set a: [?] = untyped
		{runtime.NewGDArrayType(runtime.GDNilTypeRef), runtime.GDUntypedTypeRef, runtime.NewGDArrayType(runtime.GDNilTypeRef), ""},
		{runtime.NewGDArrayType(runtime.GDIntTypeRef), runtime.GDUntypedTypeRef, runtime.NewGDArrayType(runtime.GDIntTypeRef), ""},
		{runtime.NewGDTupleType(runtime.GDIntTypeRef), runtime.GDUntypedTypeRef, runtime.NewGDTupleType(runtime.GDIntTypeRef), ""},
		{runtime.NewGDTupleType(runtime.GDNilTypeRef), runtime.GDUntypedTypeRef, runtime.NewGDTupleType(runtime.GDNilTypeRef), ""},
		// set a: [string] = untyped
		{runtime.NewGDArrayType(runtime.GDStringTypeRef), runtime.NewGDEmptyArrayType(), runtime.NewGDArrayType(runtime.GDStringTypeRef), ""},
		// set a: nil = untyped
		{runtime.GDNilTypeRef, runtime.GDUntypedTypeRef, runtime.GDNilTypeRef, ""},
		// set a: untyped = ?
		{runtime.GDUntypedTypeRef, runtime.GDNilTypeRef, runtime.GDUntypedTypeRef, ""},
		// set a: untyped = string
		{runtime.GDUntypedTypeRef, runtime.GDStringTypeRef, runtime.GDStringTypeRef, ""},
		// set a: int | string = int
		{runtime.NewGDUnionType(runtime.GDIntTypeRef, runtime.GDStringTypeRef), runtime.GDIntTypeRef, runtime.NewGDUnionType(runtime.GDIntTypeRef, runtime.GDStringTypeRef), ""},
		// set a: int | string = int | string
		{runtime.NewGDUnionType(runtime.GDIntTypeRef, runtime.GDStringTypeRef), runtime.NewGDUnionType(runtime.GDIntTypeRef, runtime.GDStringTypeRef), runtime.NewGDUnionType(runtime.GDIntTypeRef, runtime.GDStringTypeRef), ""},
		// set a: string | int = int | string
		{runtime.NewGDUnionType(runtime.GDStringTypeRef, runtime.GDIntTypeRef), runtime.NewGDUnionType(runtime.GDIntTypeRef, runtime.GDStringTypeRef), runtime.NewGDUnionType(runtime.GDStringTypeRef, runtime.GDIntTypeRef), ""},
		// set a: untyped = (untyped,)
		{runtime.GDUntypedTypeRef, runtime.NewGDTupleType(), runtime.NewGDTupleType(), ""},
		// set a: untyped = int | string
		{runtime.GDUntypedTypeRef, runtime.NewGDUnionType(runtime.GDIntTypeRef, runtime.GDStringTypeRef), runtime.NewGDUnionType(runtime.GDIntTypeRef, runtime.GDStringTypeRef), ""},
		// set a: (untyped,) = untyped
		{runtime.NewGDTupleType(), runtime.GDUntypedTypeRef, runtime.NewGDTupleType(), ""},
		// set a: (int,) = (untyped,)
		{runtime.NewGDTupleType(runtime.GDIntTypeRef), runtime.NewGDTupleType(), runtime.NewGDTupleType(runtime.GDIntTypeRef), ""},
		// set a: (untyped,) = (int,)
		{runtime.NewGDTupleType(), runtime.NewGDTupleType(runtime.GDIntTypeRef), runtime.NewGDTupleType(runtime.GDIntTypeRef), ""},
		// set a: ((int,),) = ((float,),)
		{runtime.NewGDTupleType(runtime.NewGDTupleType(runtime.GDIntTypeRef)), runtime.NewGDTupleType(runtime.NewGDTupleType(runtime.GDFloatTypeRef)), nil, "expected `((int,),)` but got `((float,),)`"},
		// set a: ((int,),) = ((untyped,),)
		{runtime.NewGDTupleType(runtime.NewGDTupleType(runtime.GDIntTypeRef)), runtime.NewGDTupleType(runtime.NewGDTupleType()), runtime.NewGDTupleType(runtime.NewGDTupleType(runtime.GDIntTypeRef)), ""},
		// set a: ((untyped,),) = ((int,),)
		{runtime.NewGDTupleType(runtime.NewGDTupleType()), runtime.NewGDTupleType(runtime.NewGDTupleType(runtime.GDIntTypeRef)), runtime.NewGDTupleType(runtime.NewGDTupleType(runtime.GDIntTypeRef)), ""},
		// set a: ((untyped,),) = ((untyped,),)
		{runtime.NewGDTupleType(runtime.NewGDTupleType()), runtime.NewGDTupleType(runtime.NewGDTupleType()), runtime.NewGDTupleType(runtime.NewGDTupleType()), ""},
		// set a: (untyped,) = ((untyped,),)
		{runtime.NewGDTupleType(), runtime.NewGDTupleType(runtime.NewGDTupleType()), runtime.NewGDTupleType(runtime.NewGDTupleType()), ""},
		// set a: int = int | string
		{runtime.GDIntTypeRef, runtime.NewGDUnionType(runtime.GDIntTypeRef, runtime.GDStringTypeRef), nil, "expected `int` but got `(int | string)`"},
		// set a: nil = int
		{runtime.GDNilTypeRef, runtime.GDIntTypeRef, nil, "expected `nil` but got `int`"},
		// set a: string = int
		{runtime.GDStringTypeRef, runtime.GDIntTypeRef, nil, "expected `string` but got `int`"},
		// set a: untyped = struct
		{runtime.GDUntypedTypeRef, structType1, structType1, ""},
		// set a: struct = untyped
		{structType1, runtime.GDUntypedTypeRef, structType1, ""},
		// set a: {a: untyped} = {a: int}
		{runtime.NewGDStructType(&runtime.GDStructAttrType{runtime.NewGDStrIdent("a"), runtime.GDUntypedTypeRef}), runtime.NewGDStructType(&runtime.GDStructAttrType{runtime.NewGDStrIdent("a"), runtime.GDIntTypeRef}), runtime.NewGDStructType(&runtime.GDStructAttrType{runtime.NewGDStrIdent("a"), runtime.GDIntTypeRef}), ""},
		// set a: struct = {untyped}
		{structType1, runtime.NewGDStructType(), structType1, ""},
		// set a: struct{a: int} = struct{b: int}
		{structType1, structType2, nil, "expected `{a: int}` but got `{b: int}`"},
		// set a: func = untyped
		{func1, runtime.GDUntypedTypeRef, func1, ""},
		// set a: [string] = [int]
		{runtime.NewGDArrayType(runtime.GDStringTypeRef), runtime.NewGDArrayType(runtime.GDIntTypeRef), nil, "expected `[string]` but got `[int]`"},
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringTypeRef)), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDIntTypeRef)), nil, "expected `[[string]]` but got `[[int]]`"},
		{runtime.NewGDArrayType(runtime.GDNilTypeRef), runtime.NewGDArrayType(runtime.GDIntTypeRef), nil, "expected `[nil]` but got `[int]`"},
		{runtime.NewGDTupleType(runtime.GDIntTypeRef), runtime.NewGDArrayType(runtime.GDIntTypeRef), nil, "expected `(int,)` but got `[int]`"},
		// set a: [[int]] = [nil]
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDIntTypeRef)), runtime.NewGDArrayType(runtime.GDNilTypeRef), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDIntTypeRef)), ""},
		// set a: [[string]] | [[int]] = [[string]]
		{runtime.NewGDUnionType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringTypeRef)), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDIntTypeRef))), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringTypeRef)), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringTypeRef)), ""},
		// set a: [[any]] = [[int] | [string]]
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDAnyTypeRef)), runtime.NewGDUnionType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDIntTypeRef)), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringTypeRef))), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDAnyTypeRef)), ""},
		// set a func2 = func3
		{func2, func3, nil, "expected `func() => string` but got `func() => int`"},
		// set a func2 = nil
		{func2, runtime.GDNilTypeRef, func2, ""},
		// set a func2 = untyped
		{func2, runtime.GDUntypedTypeRef, func2, ""},
		// set a func1 = func2
		{func1, func2, func1, "expected `func() => any` but got `func() => string`"},
		// set a func2 = func1
		{func2, func1, nil, "expected `func() => string` but got `func() => any`"},
		// string as type
		{runtime.NewGDStrTypeRefType("typ"), runtime.GDStringTypeRef, runtime.NewGDStrTypeRefType("typ"), ""},
		{runtime.NewGDStrTypeRefType("typ"), runtime.GDIntTypeRef, nil, "expected `typ` but got `int`"},
	}

	TypeTests(t, tests, func(t *testing.T, test TypeTest) error {
		nType, err := runtime.InferType(test.toType, test.fromType, stack)
		if err != nil {
			return err
		}

		if test.expected != nil {
			if err := runtime.EqualTypes(nType, test.expected, stack); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("there was no error but found expected type nil")
		}

		return nil
	})
}

func TypeTests(t *testing.T, tests []TypeTest, callback func(t *testing.T, test TypeTest) error) {
	for _, test := range tests {
		if !t.Run(test.ToString(), func(t *testing.T) {
			err := callback(t, test)
			if err != nil {
				if test.errMsg != "" {
					if !strings.Contains(err.Error(), test.errMsg) {
						t.Errorf("Expected error message to contain %q but got %q", test.errMsg, err.Error())
					}
					return
				}

				t.Errorf("Expected no error but got %q", err)
				return
			} else if test.errMsg != "" {
				t.Errorf("Expected error message %q but got no error", test.errMsg)
			}
		}) {
			break
		}
	}
}
