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
	if computedType != runtime.GDUntypedType {
		t.Errorf("Computed type should be untyped but got %s", computedType.ToString())
	}
}

func TestComputeTypeDeterminationWithObjectTypes(t *testing.T) {
	computedType := runtime.ComputeTypeFromObjects([]runtime.GDObject{runtime.GDString("test"), runtime.GDString("test")})
	if computedType != runtime.GDStringType {
		t.Errorf("Computed type should be (string) but got %s", computedType.ToString())
	}
}

// Test compute type with different types

func TestComputeTypeWithDifferentTypes(t *testing.T) {
	for _, test := range []struct {
		types    []runtime.GDTypable
		expected runtime.GDTypable
	}{
		{[]runtime.GDTypable{runtime.GDStringType, runtime.GDStringType}, runtime.GDStringType},
		{[]runtime.GDTypable{runtime.GDStringType, runtime.GDBoolType}, runtime.NewGDUnionType(runtime.GDStringType, runtime.GDBoolType)},
		{[]runtime.GDTypable{runtime.NewGDArrayType(runtime.GDStringType), runtime.NewGDArrayType(runtime.GDBoolType)}, runtime.NewGDUnionType(runtime.NewGDArrayType(runtime.GDStringType), runtime.NewGDArrayType(runtime.GDBoolType))},
		{[]runtime.GDTypable{runtime.NewGDArrayType(runtime.GDStringType), runtime.NewGDArrayType(runtime.GDStringType)}, runtime.NewGDArrayType(runtime.GDStringType)},
		{[]runtime.GDTypable{
			runtime.NewGDArrayType(
				runtime.NewGDArrayType(
					runtime.NewGDArrayType(
						runtime.NewGDArrayType(runtime.GDStringType),
					),
				),
			), runtime.NewGDArrayType(
				runtime.NewGDArrayType(
					runtime.NewGDArrayType(
						runtime.NewGDArrayType(runtime.GDBoolType),
					),
				),
			),
		}, runtime.NewGDUnionType(
			runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringType)))),
			runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDBoolType)))),
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
		{runtime.GDStringType, runtime.GDStringType, nil, ""},
		{runtime.GDStringType, runtime.GDBoolType, nil, "types `string` and `bool` are not equal"},
		{runtime.NewGDArrayType(runtime.GDStringType), runtime.NewGDArrayType(runtime.GDStringType), nil, ""},
		{runtime.NewGDArrayType(runtime.GDStringType), runtime.NewGDArrayType(runtime.GDBoolType), nil, "types `[string]` and `[bool]` are not equal"},
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringType))), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringType))), nil, ""},
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringType))), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDBoolType))), nil, "types `[[[string]]]` and `[[[bool]]]` are not equal"},
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringType))), runtime.GDStringType, nil, "types `[[[string]]]` and `string` are not equal"},
		{runtime.GDStringType, runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringType))), nil, "types `string` and `[[[string]]]` are not equal"},
		{runtime.NewGDArrayType(runtime.GDStringType), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringType))), nil, "types `[string]` and `[[[string]]]` are not equal"},
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringType))), runtime.NewGDArrayType(runtime.GDStringType), nil, "types `[[[string]]]` and `[string]` are not equal"},
		{runtime.NewGDArrayType(runtime.GDAnyType), runtime.NewGDArrayType(runtime.GDStringType), nil, "types `[any]` and `[string]` are not equal"},
		{structWithAttrAAsInt, structWithAttrAAsInt, nil, ""},
		{structWithAttrAAsString, structWithAttrAAsInt, nil, "types `{a: string}` and `{a: int}` are not equal"},
		{structWithAttrsAStrBInt, structWithAttrsBIntAStr, nil, ""},
		{runtime.NewGDTuple(runtime.NewGDIntNumber(1), runtime.GDString("test")).GetType(), runtime.NewGDTupleType(runtime.GDIntType, runtime.GDStringType), nil, ""},
		{runtime.GDUntypedType, runtime.GDUntypedType, nil, ""},
		{runtime.GDUntypedType, runtime.GDStringType, nil, "types `untyped` and `string` are not equal"},
		{runtime.GDAnyType, runtime.GDStringType, nil, "types `any` and `string` are not equal"},
		{runtime.NewGDTupleType(runtime.GDIntType, runtime.NewGDTupleType(runtime.GDIntType, runtime.GDStringType)), runtime.NewGDTupleType(runtime.GDIntType, runtime.NewGDTupleType(runtime.GDIntType, runtime.GDStringType)), nil, ""},
		{runtime.NewGDTupleType(runtime.GDStringType, runtime.GDStringType), runtime.NewGDTupleType(runtime.GDStringType, runtime.GDStringType), nil, ""},
		{runtime.NewGDTupleType(runtime.GDIntType, runtime.NewGDTupleType(runtime.GDIntType, runtime.GDIntType)), runtime.NewGDTupleType(runtime.GDIntType, runtime.NewGDTupleType(runtime.GDIntType, runtime.GDStringType)), nil, "types `(int, (int, int))` and `(int, (int, string))` are not equal"},
		{runtime.NewGDUnionType(runtime.GDIntType, runtime.GDStringType), runtime.NewGDUnionType(runtime.GDIntType, runtime.GDStringType), nil, ""},
		{runtime.NewGDUnionType(runtime.GDStringType, runtime.GDIntType), runtime.NewGDUnionType(runtime.GDIntType, runtime.GDStringType), nil, ""},
		{runtime.GDStringType, runtime.NewGDUnionType(runtime.GDIntType, runtime.GDStringType), nil, "types `string` and `(int | string)` are not equal"},
		{runtime.NewGDUnionType(runtime.GDIntType, runtime.GDStringType), runtime.GDStringType, nil, ""},
	}

	TypeTests(t, tests, func(t *testing.T, test TypeTest) error {
		return runtime.EqualTypes(test.toType, test.fromType, nil)
	})
}

// Assignments

func TestCanBeAssigned(t *testing.T) {
	tests := []TypeTest{
		{runtime.GDStringType, runtime.GDStringType, nil, ""},
		{runtime.GDStringType, runtime.GDBoolType, nil, "expected `string` but got `bool`"},
		{runtime.NewGDArrayType(runtime.GDStringType), runtime.NewGDArrayType(runtime.GDStringType), nil, ""},
		{runtime.NewGDArrayType(runtime.GDStringType), runtime.NewGDArrayType(runtime.GDBoolType), nil, "expected `[string]` but got `[bool]`"},
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringType))), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringType))), nil, ""},
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringType))), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDBoolType))), nil, "expected `[[[string]]]` but got `[[[bool]]]`"},
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringType))), runtime.GDStringType, nil, "expected `[[[string]]]` but got `string`"},
		{runtime.GDStringType, runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringType))), nil, "expected `string` but got `[[[string]]]`"},
		{runtime.NewGDArrayType(runtime.GDStringType), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringType))), nil, "expected `[string]` but got `[[[string]]]`"},
		{runtime.GDAnyType, runtime.NewGDArrayType(runtime.GDStringType), nil, ""},
		{runtime.NewGDArrayType(runtime.GDStringType), runtime.GDAnyType, nil, "expected `[string]` but got `any`"},
		{runtime.NewGDArrayType(runtime.GDStringType), runtime.GDStringType, nil, "expected `[string]` but got `string`"},
		{runtime.NewGDTupleType(runtime.GDStringType), runtime.GDNilType, nil, ""},
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
		{runtime.NewGDArrayType(runtime.GDStringType), runtime.GDStringType, nil, ""},
		{runtime.NewGDArrayType(runtime.GDStringType), runtime.GDAnyType, nil, "expected `string` but got `any`"},
		{runtime.NewGDArrayType(runtime.GDAnyType), runtime.GDStringType, nil, ""},
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringType)), runtime.GDStringType, nil, "expected `[string]` but got `string`"},
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringType)), runtime.NewGDArrayType(runtime.GDStringType), nil, ""},
		{runtime.NewGDArrayType(runtime.GDAnyType), runtime.NewGDArrayType(runtime.NewGDUnionType(runtime.GDStringType, runtime.GDIntType)), nil, ""},
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDAnyType)), runtime.NewGDArrayType(runtime.NewGDUnionType(runtime.NewGDArrayType(runtime.GDStringType), runtime.NewGDArrayType(runtime.GDIntType))), nil, ""},
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
	stack := runtime.NewGDSymbolStack()

	structType1 := runtime.NewGDStructType(runtime.GDStructAttrType{aParamIdent, runtime.GDIntType})
	structType2 := runtime.NewGDStructType(runtime.GDStructAttrType{bParamIdent, runtime.GDIntType})
	func1 := runtime.NewGDLambdaType(runtime.GDLambdaArgTypes{}, runtime.GDAnyType, false)
	func2 := runtime.NewGDLambdaType(runtime.GDLambdaArgTypes{}, runtime.GDStringType, false)
	func3 := runtime.NewGDLambdaType(runtime.GDLambdaArgTypes{}, runtime.GDIntType, false)

	// Register a type alias
	typeAliasIdent := runtime.NewGDStringIdent("typ")
	_, err := stack.AddSymbol(typeAliasIdent, true, true, runtime.GDStringType, nil)
	if err != nil {
		t.Fatal(err)
	}

	tests := []TypeTest{
		// set a: any = int | float
		{runtime.GDAnyType, runtime.NewGDUnionType(runtime.GDIntType, runtime.GDFloatType), runtime.GDAnyType, ""},
		// set a: any = any
		{runtime.GDAnyType, runtime.GDAnyType, runtime.GDAnyType, ""},
		// set a: any = [int]
		{runtime.GDAnyType, runtime.NewGDArrayType(runtime.GDIntType), runtime.GDAnyType, ""},
		// set a: any = func() => any
		{runtime.GDAnyType, func1, runtime.GDAnyType, ""},
		// set a: string = any
		{runtime.GDStringType, runtime.GDAnyType, nil, "expected `string` but got `any`"},
		// set a: string = untyped
		{runtime.GDStringType, runtime.GDUntypedType, runtime.GDStringType, ""},
		// set a: string = nil
		{runtime.GDStringType, runtime.GDNilType, runtime.GDStringType, ""},
		// set a: [untyped] = untyped
		{runtime.NewGDEmptyArrayType(), runtime.GDUntypedType, runtime.NewGDEmptyArrayType(), ""},
		// set a: untyped = [untyped]
		{runtime.GDUntypedType, runtime.NewGDEmptyArrayType(), runtime.NewGDEmptyArrayType(), ""},
		// set a: [untyped] = [untyped]
		{runtime.NewGDEmptyArrayType(), runtime.NewGDEmptyArrayType(), runtime.NewGDEmptyArrayType(), ""},
		// set a: [untyped] = [[untyped]] (For this case `untyped` on the left side is a `[untyped]`)
		{runtime.NewGDEmptyArrayType(), runtime.NewGDArrayType(runtime.NewGDEmptyArrayType()), runtime.NewGDArrayType(runtime.NewGDEmptyArrayType()), ""},
		// set a: [[untyped]] = [untyped]
		{runtime.NewGDArrayType(runtime.NewGDEmptyArrayType()), runtime.NewGDEmptyArrayType(), runtime.NewGDArrayType(runtime.NewGDEmptyArrayType()), ""},
		// set a: [[string]] = [[untyped]]
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringType)), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDUntypedType)), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringType)), ""},
		// set a: [string] = [untyped]
		{runtime.NewGDArrayType(runtime.GDStringType), runtime.NewGDEmptyArrayType(), runtime.NewGDArrayType(runtime.GDStringType), ""},
		// set a: [untyped] = [string]
		{runtime.NewGDEmptyArrayType(), runtime.NewGDArrayType(runtime.GDStringType), runtime.NewGDArrayType(runtime.GDStringType), ""},
		// set a: [any] = nil
		{runtime.NewGDArrayType(runtime.GDAnyType), runtime.GDNilType, runtime.NewGDArrayType(runtime.GDAnyType), ""},
		// set a: [int] = [nil]
		{runtime.NewGDArrayType(runtime.GDIntType), runtime.NewGDArrayType(runtime.GDNilType), runtime.NewGDArrayType(runtime.GDIntType), ""},
		// set a: [string] = any
		{runtime.NewGDArrayType(runtime.GDStringType), runtime.GDAnyType, nil, "expected `[string]` but got `any`"},
		// set a: [[nil]] = untyped
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDNilType)), runtime.GDUntypedType, runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDNilType)), ""},
		// set a: untyped = [string]
		{runtime.GDUntypedType, runtime.NewGDArrayType(runtime.GDStringType), runtime.NewGDArrayType(runtime.GDStringType), ""},
		// set a: [?] = untyped
		{runtime.NewGDArrayType(runtime.GDNilType), runtime.GDUntypedType, runtime.NewGDArrayType(runtime.GDNilType), ""},
		{runtime.NewGDArrayType(runtime.GDIntType), runtime.GDUntypedType, runtime.NewGDArrayType(runtime.GDIntType), ""},
		{runtime.NewGDTupleType(runtime.GDIntType), runtime.GDUntypedType, runtime.NewGDTupleType(runtime.GDIntType), ""},
		{runtime.NewGDTupleType(runtime.GDNilType), runtime.GDUntypedType, runtime.NewGDTupleType(runtime.GDNilType), ""},
		// set a: [string] = untyped
		{runtime.NewGDArrayType(runtime.GDStringType), runtime.NewGDEmptyArrayType(), runtime.NewGDArrayType(runtime.GDStringType), ""},
		// set a: nil = untyped
		{runtime.GDNilType, runtime.GDUntypedType, runtime.GDNilType, ""},
		// set a: untyped = ?
		{runtime.GDUntypedType, runtime.GDNilType, runtime.GDUntypedType, ""},
		// set a: untyped = string
		{runtime.GDUntypedType, runtime.GDStringType, runtime.GDStringType, ""},
		// set a: int | string = int
		{runtime.NewGDUnionType(runtime.GDIntType, runtime.GDStringType), runtime.GDIntType, runtime.NewGDUnionType(runtime.GDIntType, runtime.GDStringType), ""},
		// set a: int | string = int | string
		{runtime.NewGDUnionType(runtime.GDIntType, runtime.GDStringType), runtime.NewGDUnionType(runtime.GDIntType, runtime.GDStringType), runtime.NewGDUnionType(runtime.GDIntType, runtime.GDStringType), ""},
		// set a: string | int = int | string
		{runtime.NewGDUnionType(runtime.GDStringType, runtime.GDIntType), runtime.NewGDUnionType(runtime.GDIntType, runtime.GDStringType), runtime.NewGDUnionType(runtime.GDStringType, runtime.GDIntType), ""},
		// set a: untyped = (untyped,)
		{runtime.GDUntypedType, runtime.NewGDTupleType(), runtime.NewGDTupleType(), ""},
		// set a: untyped = int | string
		{runtime.GDUntypedType, runtime.NewGDUnionType(runtime.GDIntType, runtime.GDStringType), runtime.NewGDUnionType(runtime.GDIntType, runtime.GDStringType), ""},
		// set a: (untyped,) = untyped
		{runtime.NewGDTupleType(), runtime.GDUntypedType, runtime.NewGDTupleType(), ""},
		// set a: (int,) = (untyped,)
		{runtime.NewGDTupleType(runtime.GDIntType), runtime.NewGDTupleType(), runtime.NewGDTupleType(runtime.GDIntType), ""},
		// set a: (untyped,) = (int,)
		{runtime.NewGDTupleType(), runtime.NewGDTupleType(runtime.GDIntType), runtime.NewGDTupleType(runtime.GDIntType), ""},
		// set a: ((int,),) = ((float,),)
		{runtime.NewGDTupleType(runtime.NewGDTupleType(runtime.GDIntType)), runtime.NewGDTupleType(runtime.NewGDTupleType(runtime.GDFloatType)), nil, "expected `((int,),)` but got `((float,),)`"},
		// set a: ((int,),) = ((untyped,),)
		{runtime.NewGDTupleType(runtime.NewGDTupleType(runtime.GDIntType)), runtime.NewGDTupleType(runtime.NewGDTupleType()), runtime.NewGDTupleType(runtime.NewGDTupleType(runtime.GDIntType)), ""},
		// set a: ((untyped,),) = ((int,),)
		{runtime.NewGDTupleType(runtime.NewGDTupleType()), runtime.NewGDTupleType(runtime.NewGDTupleType(runtime.GDIntType)), runtime.NewGDTupleType(runtime.NewGDTupleType(runtime.GDIntType)), ""},
		// set a: ((untyped,),) = ((untyped,),)
		{runtime.NewGDTupleType(runtime.NewGDTupleType()), runtime.NewGDTupleType(runtime.NewGDTupleType()), runtime.NewGDTupleType(runtime.NewGDTupleType()), ""},
		// set a: (untyped,) = ((untyped,),)
		{runtime.NewGDTupleType(), runtime.NewGDTupleType(runtime.NewGDTupleType()), runtime.NewGDTupleType(runtime.NewGDTupleType()), ""},
		// set a: int = int | string
		{runtime.GDIntType, runtime.NewGDUnionType(runtime.GDIntType, runtime.GDStringType), nil, "expected `int` but got `(int | string)`"},
		// set a: nil = int
		{runtime.GDNilType, runtime.GDIntType, nil, "expected `nil` but got `int`"},
		// set a: string = int
		{runtime.GDStringType, runtime.GDIntType, nil, "expected `string` but got `int`"},
		// set a: untyped = struct
		{runtime.GDUntypedType, structType1, structType1, ""},
		// set a: struct = untyped
		{structType1, runtime.GDUntypedType, structType1, ""},
		// set a: {a: untyped} = {a: int}
		{runtime.NewGDStructType(runtime.GDStructAttrType{runtime.NewGDStringIdent("a"), runtime.GDUntypedType}), runtime.NewGDStructType(runtime.GDStructAttrType{runtime.NewGDStringIdent("a"), runtime.GDIntType}), runtime.NewGDStructType(runtime.GDStructAttrType{runtime.NewGDStringIdent("a"), runtime.GDIntType}), ""},
		// set a: struct = {untyped}
		{structType1, runtime.NewGDStructType(), structType1, ""},
		// set a: struct{a: int} = struct{b: int}
		{structType1, structType2, nil, "expected `{a: int}` but got `{b: int}`"},
		// set a: func = untyped
		{func1, runtime.GDUntypedType, func1, ""},
		// set a: [string] = [int]
		{runtime.NewGDArrayType(runtime.GDStringType), runtime.NewGDArrayType(runtime.GDIntType), nil, "expected `[string]` but got `[int]`"},
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringType)), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDIntType)), nil, "expected `[[string]]` but got `[[int]]`"},
		{runtime.NewGDArrayType(runtime.GDNilType), runtime.NewGDArrayType(runtime.GDIntType), nil, "expected `[nil]` but got `[int]`"},
		{runtime.NewGDTupleType(runtime.GDIntType), runtime.NewGDArrayType(runtime.GDIntType), nil, "expected `(int,)` but got `[int]`"},
		// set a: [[int]] = [nil]
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDIntType)), runtime.NewGDArrayType(runtime.GDNilType), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDIntType)), ""},
		// set a: [[string]] | [[int]] = [[string]]
		{runtime.NewGDUnionType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringType)), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDIntType))), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringType)), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringType)), ""},
		// set a: [[any]] = [[int] | [string]]
		{runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDAnyType)), runtime.NewGDUnionType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDIntType)), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringType))), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDAnyType)), ""},
		// set a func2 = func3
		{func2, func3, nil, "expected `() => string` but got `() => int`"},
		// set a func2 = nil
		{func2, runtime.GDNilType, func2, ""},
		// set a func2 = untyped
		{func2, runtime.GDUntypedType, func2, ""},
		// set a func1 = func2
		{func1, func2, func1, "expected `() => any` but got `() => string`"},
		// set a func2 = func1
		{func2, func1, nil, "expected `() => string` but got `() => any`"},
		// string as type
		{runtime.NewStrRefType("typ"), runtime.GDStringType, runtime.NewStrRefType("typ"), ""},
		{runtime.NewStrRefType("typ"), runtime.GDIntType, nil, "expected `typ` but got `int`"},
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
