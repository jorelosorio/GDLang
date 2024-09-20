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

func TestArrayOfNumbers(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDIntNumber(1), runtime.NewGDIntNumber(2), runtime.NewGDIntNumber(3), runtime.NewGDFloatNumber(3.1))

	if array.ToString() != "[1, 2, 3, 3.1]" {
		t.Error("Wrong array string representation, got", array.ToString())
	}
}

func TestArrayOfStrings(t *testing.T) {
	array := runtime.NewGDArray(runtime.GDString("one"), runtime.GDString("two"), runtime.GDString("three"))

	if array.ToString() != `["one", "two", "three"]` {
		t.Error("Wrong array string representation, got", array.ToString())
	}
}

func TestArrayOfBools(t *testing.T) {
	array := runtime.NewGDArray(runtime.GDBool(true), runtime.GDBool(false), runtime.GDBool(true))
	if array.ToString() != "[true, false, true]" {
		t.Error("Wrong array string representation, got", array.ToString())
	}
}

func TestArrayOfChars(t *testing.T) {
	array := runtime.NewGDArray(runtime.GDChar('a'), runtime.GDChar('b'), runtime.GDChar('c'))
	if array.ToString() != `['a', 'b', 'c']` {
		t.Error("Wrong array string representation, got", array.ToString())
	}
}

func TestArrayOfArrays(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDArray(runtime.NewGDIntNumber(1), runtime.NewGDIntNumber(2)), runtime.NewGDArray(runtime.NewGDIntNumber(3), runtime.NewGDIntNumber(4)))

	if array.ToString() != `[[1, 2], [3, 4]]` {
		t.Error("Wrong array string representation, got", array.ToString())
	}
}

func TestArrayOfMixedTypes(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDIntNumber(1), runtime.GDChar('j'), runtime.GDString("two"), runtime.NewGDArray(runtime.NewGDIntNumber(3), runtime.GDString("four")))

	if array.ToString() != `[1, 'j', "two", [3, "four"]]` {
		t.Error("Wrong array string representation, got", array.ToString())
	}
}

func TestArrayGet(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDIntNumber(1))

	object, err := array.Get(0)
	if err != nil {
		t.Error("Error getting value")
	}

	if object.(runtime.GDInt8) != 1 {
		t.Errorf("Wrong value got %q but expected 1", object.ToString())
	}
}

func TestArraySet(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDIntNumber(1))

	err := array.Set(0, runtime.NewGDIntNumber(2), nil)
	if err != nil {
		t.Error("Error setting value")
	}

	value, err := array.Get(0)
	if err != nil {
		t.Error("Error getting value")
	}

	if value.(runtime.GDInt8) != 2 {
		t.Errorf("Wrong value got %q but expected 2", value.ToString())
	}
}

func TestArrayAppend(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDIntNumber(1))

	err := array.AddObject(runtime.NewGDIntNumber(2), nil)
	if err != nil {
		t.Errorf("Error when appending value, expected no error: %v", err)
	}

	if array.Length() != 2 {
		t.Errorf("Wrong length got %v but expected 2", array.Length())
	}

	value, err := array.Get(1)
	if err != nil {
		t.Error("Error getting value")
	}

	if value.(runtime.GDInt8) != 2 {
		t.Errorf("Wrong value got %q but expected 2", value.ToString())
	}
}

func TestArrayRemove(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDIntNumber(1), runtime.NewGDIntNumber(2), runtime.NewGDIntNumber(3))

	_, err := array.Remove(1)
	if err != nil && err.Error() == runtime.IndexOutOfBoundsErr.Error() {
		t.Errorf("Error when removing value, expected no error: %v", err)
	}

	if array.Length() != 2 {
		t.Errorf("Wrong length got %v but expected 2", array.Length())
	}

	value, err := array.Get(1)
	if err != nil && err.Error() == runtime.IndexOutOfBoundsErr.Error() {
		t.Errorf("Error when getting value after removing, expected no error: %v", err)
	}

	if value.(runtime.GDInt8) != 3 {
		t.Errorf("Wrong value got %q but expected 3", value.ToString())
	}
}

func TestArrayLength(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDIntNumber(1), runtime.NewGDIntNumber(2), runtime.NewGDIntNumber(3))

	if array.Length() != 3 {
		t.Errorf("Wrong length got %v but expected 3", array.Length())
	}
}

func TestArrayIsEmpty(t *testing.T) {
	array := runtime.NewGDArray()

	if !array.IsEmpty() {
		t.Error("Array should be empty")
	}
}

func TestGetIndexOutOfBounds(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDIntNumber(1))

	_, err := array.Get(1)
	if err != runtime.IndexOutOfBoundsErr {
		t.Error("Expected error when getting index out of bounds")
	}
}

func TestSetIndexOutOfBounds(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDIntNumber(1))

	err := array.Set(1, runtime.NewGDIntNumber(2), nil)
	if err != runtime.IndexOutOfBoundsErr {
		t.Error("Expected error when setting index out of bounds")
	}
}

func TestRemoveIndexOutOfBounds(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDIntNumber(1))

	_, err := array.Remove(1)
	if err != runtime.IndexOutOfBoundsErr {
		t.Error("Expected error when removing index out of bounds")
	}
}

func TestAppendIndexOutOfBounds(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDIntNumber(1))

	err := array.AddObject(runtime.NewGDIntNumber(2), nil)
	if err != nil {
		t.Errorf("Error when appending value, expected no error: %v", err)
	}

	_, err = array.Get(2)
	if err != runtime.IndexOutOfBoundsErr {
		t.Error("Expected error when getting index out of bounds")
	}
}

func TestGDArrayWithAutoTypeForString(t *testing.T) {
	array := runtime.NewGDArray(runtime.GDString("one"), runtime.GDString("two"), runtime.GDString("three"))

	err := runtime.EqualTypes(array.GetType(), runtime.NewGDArrayType(runtime.GDStringType), nil)
	if err != nil {
		t.Error(err)
	}
}

func TestGDArrayWithAutoTypeForNumber(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDIntNumber(1), runtime.NewGDIntNumber(2), runtime.NewGDIntNumber(3))

	err := runtime.EqualTypes(array.GetType(), runtime.NewGDArrayType(runtime.GDIntType), nil)
	if err != nil {
		t.Error(err)
	}
}

func TestGDArrayWithMixedAnySubArrayTypes(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDArray(runtime.NewGDIntNumber(1), runtime.NewGDIntNumber(2)), runtime.NewGDArray(runtime.GDString("three"), runtime.GDString("four")))
	err := runtime.EqualTypes(array.GetType(), runtime.NewGDArrayType(runtime.NewGDUnionType(runtime.NewGDArrayType(runtime.GDIntType), runtime.NewGDArrayType(runtime.GDStringType))), nil)
	if err != nil {
		t.Error(err)
	}
}

func TestGDArrayWithMixedSubArrayTypes(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDArray(runtime.NewGDIntNumber(1), runtime.NewGDIntNumber(2)), runtime.GDString("three"), runtime.GDString("four"))
	err := runtime.EqualTypes(array.GetType(), runtime.NewGDArrayType(runtime.NewGDUnionType(runtime.NewGDArrayType(runtime.GDIntType), runtime.GDStringType)), nil)
	if err != nil {
		t.Error(err)
	}
}

func TestGDArrayWithAutoTypeForArray(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDArray(runtime.NewGDIntNumber(1), runtime.NewGDIntNumber(2)), runtime.NewGDArray(runtime.NewGDIntNumber(3), runtime.NewGDIntNumber(4)))
	err := runtime.EqualTypes(array.GetType(), runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDIntType)), nil)
	if err != nil {
		t.Error(err)
	}
}

func TestGDArrayWithAutoTypeForMixedArrayTypes(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDArray(runtime.NewGDIntNumber(1), runtime.GDString("two")), runtime.NewGDArray(runtime.GDString("three"), runtime.GDBool(true)))
	err := runtime.EqualTypes(array.GetType(), runtime.NewGDArrayType(runtime.NewGDUnionType(runtime.NewGDArrayType(runtime.NewGDUnionType(runtime.GDIntType, runtime.GDStringType)), runtime.NewGDArrayType(runtime.NewGDUnionType(runtime.GDBoolType, runtime.GDStringType)))), nil)
	if err != nil {
		t.Error(err)
	}
}

func TestGDArrayWithAutoTypeForEmptyArray(t *testing.T) {
	array := runtime.NewGDArray()
	err := runtime.EqualTypes(array.GetType(), runtime.NewGDArrayType(runtime.GDUntypedType), nil)
	if err != nil {
		t.Error(err)
	}
}

// Array append operations with type checking

func TestAppendNumberToArrayOfNumbers(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDIntNumber(1), runtime.NewGDIntNumber(2), runtime.NewGDIntNumber(3))
	err := array.AddObject(runtime.NewGDIntNumber(4), nil)
	if err != nil {
		t.Errorf("Error when appending number to array of numbers, expected no error: %v", err)
	}

	if array.Length() != 4 {
		t.Errorf("Wrong length got %v but expected 4", array.Length())
	}
}

func TestAppendStringToArrayOfStrings(t *testing.T) {
	array := runtime.NewGDArray(runtime.GDString("one"), runtime.GDString("two"), runtime.GDString("three"))
	err := array.AddObject(runtime.GDString("four"), nil)
	if err != nil {
		t.Errorf("Error when appending string to array of strings, expected no error: %v", err)
	}

	if array.Length() != 4 {
		t.Errorf("Wrong length got %v but expected 4", array.Length())
	}
}

func TestAppendArrayToArrayOfArrays(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDArray(runtime.NewGDIntNumber(1), runtime.NewGDIntNumber(2)), runtime.NewGDArray(runtime.NewGDIntNumber(3), runtime.NewGDIntNumber(4)))
	err := array.AddObjects([]runtime.GDObject{runtime.NewGDArray(runtime.NewGDIntNumber(5), runtime.NewGDIntNumber(6))}, nil)

	if err != nil {
		t.Errorf("Error when appending array to array of arrays, expected no error: %v", err)
	}

	if array.Length() != 3 {
		t.Errorf("Wrong length got %v but expected 3", array.Length())
	}
}

func TestAppendMixedTypesToArrayOfNumbers(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDIntNumber(1), runtime.NewGDIntNumber(2), runtime.NewGDIntNumber(3))
	err := array.AddObject(runtime.GDString("four"), nil)
	errMsg := runtime.WrongTypesErr(array.GDArrayType.SubType, runtime.GDStringType).Error()
	if err != nil && err.Error() != errMsg {
		t.Errorf("Error when appending string to array of numbers, expected error: %q but got %q", errMsg, err)
	}

	if array.Length() != 3 {
		t.Errorf("Wrong length got %v but expected 3", array.Length())
	}
}

func TestAppendMixedTypesToArrayOfStrings(t *testing.T) {
	array := runtime.NewGDArray(runtime.GDString("one"), runtime.GDString("two"), runtime.GDString("three"))
	err := array.AddObject(runtime.NewGDIntNumber(4), nil)
	if err != nil && err.Error() != runtime.WrongTypesErr(array.GDArrayType.SubType, runtime.GDIntType).Error() {
		t.Errorf("Error when appending int to array of strings, expected error: %v", err)
	}

	if array.Length() != 3 {
		t.Errorf("Wrong length got %v but expected 3", array.Length())
	}
}

func TestAppendMixedTypesToArrayOfArrays(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDArray(runtime.NewGDIntNumber(1), runtime.NewGDIntNumber(2)), runtime.NewGDArray(runtime.NewGDIntNumber(3), runtime.NewGDIntNumber(4)))
	err := array.AddObject(runtime.GDString("five"), nil)
	if err != nil && err.Error() != runtime.WrongTypesErr(array.GDArrayType.SubType, runtime.GDStringType).Error() {
		t.Errorf("Error when appending string to array of arrays, expected error: %v", err)
	}

	if array.Length() != 2 {
		t.Errorf("Wrong length got %v but expected 2", array.Length())
	}
}

func TestAppendMixedTypesToArrayOfMixedTypes(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDIntNumber(1), runtime.GDString("two"), runtime.NewGDArray(runtime.NewGDIntNumber(3), runtime.GDString("four")))
	err := array.AddObject(runtime.GDBool(true), nil)

	if err == nil {
		t.Errorf("Error when appending boolean to array of mixed types, expected error: %v", err)
	}
}

func TestAppendMixedTypesToArrayOfAnyWithOneLevel(t *testing.T) {
	array := runtime.NewGDArrayWithType(runtime.GDAnyType)
	err := array.AddObjects([]runtime.GDObject{runtime.GDBool(true), runtime.GDString("two"), runtime.NewGDIntNumber(3)}, nil)
	if err != nil {
		t.Errorf("Error when appending boolean to array of mixed types, expected no error: %v", err)
	}

	if array.Length() != 3 {
		t.Errorf("Wrong length got %v but expected 3", array.Length())
	}
}

func TestAppendMixedTypesToArrayOfAnyWithTwoLevel(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDArrayWithType(runtime.GDAnyType))
	err := array.AddObjects([]runtime.GDObject{runtime.NewGDArray(runtime.GDBool(true), runtime.GDString("two"), runtime.NewGDIntNumber(3))}, nil)
	if err != nil {
		t.Errorf("Error when appending boolean to array of mixed types, expected no error: %v", err)
	}

	if array.Length() != 2 {
		t.Errorf("Wrong length got %v but expected 3", array.Length())
	}
}

func TestAppendMixedTypesToArrayOfAnyWithThreeLevel(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDArray(runtime.NewGDArrayWithType(runtime.GDAnyType)))
	err := array.AddObjects([]runtime.GDObject{runtime.NewGDArray(runtime.NewGDArray(runtime.GDBool(true), runtime.GDString("two"), runtime.NewGDIntNumber(3)))}, nil)
	if err != nil {
		t.Errorf("Error when appending boolean to array of mixed types, expected no error: %v", err)
	}

	if array.Length() != 2 {
		t.Errorf("Wrong length got %v but expected 3", array.Length())
	}
}

func TestCreateArrayWithValidType(t *testing.T) {
	array, err := runtime.NewGDArrayWithSubTypeAndObjects(runtime.GDIntType, []runtime.GDObject{runtime.NewGDIntNumber(1)}, nil)
	if err != nil {
		t.Errorf("Error creating array with allowed type, expected no error: %v", err)
	}

	err = runtime.EqualTypes(array.GetType(), runtime.NewGDArrayType(runtime.GDIntType), nil)
	if err != nil {
		t.Error(err)
		return
	}

	if array.Length() != 1 {
		t.Errorf("Wrong length got %v but expected 1", array.Length())
	}
}

func TestAppendWrongObjectTypeOnArrayWithTypeOnInitialization(t *testing.T) {
	_, err := runtime.NewGDArrayWithSubTypeAndObjects(runtime.GDIntType, []runtime.GDObject{runtime.GDString("two")}, nil)
	if err != nil && err.Error() != runtime.WrongTypesErr(runtime.GDIntType, runtime.GDStringType).Error() {
		t.Errorf("Error when appending string to array of numbers, expected error: %v", err)
	}
}

func TestAppendWrongObjectTypeOnArrayWithTypeOnInitializationWithAnyType(t *testing.T) {
	array, err := runtime.NewGDArrayWithSubTypeAndObjects(runtime.GDAnyType, []runtime.GDObject{runtime.GDString("two")}, nil)
	if err != nil {
		t.Errorf("Error when appending string to array of any, expected no error: %v", err)
	}

	err = runtime.EqualTypes(array.GetType(), runtime.NewGDArrayType(runtime.GDAnyType), nil)
	if err != nil {
		t.Error(err)
		return
	}

	if array.Length() != 1 {
		t.Errorf("Wrong length got %v but expected 1", array.Length())
	}
}

func TestAppendWrongObjectTypeOnArrayWithTypeOnInitializationWithAnyTypeAndTwoLevel(t *testing.T) {
	array, err := runtime.NewGDArrayWithSubTypeAndObjects(runtime.GDAnyType, []runtime.GDObject{runtime.NewGDArray(runtime.NewGDArray(runtime.GDString("two")))}, nil)
	if err != nil {
		t.Errorf("Error when appending array to array of any, expected no error: %v", err)
	}

	err = runtime.EqualTypes(array.GetType(), runtime.NewGDArrayType(runtime.GDAnyType), nil)
	if err != nil {
		t.Error(err)
		return
	}

	if array.Length() != 1 {
		t.Errorf("Wrong length got %v but expected 1", array.Length())
	}
}

func TestComputeArrayType(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDIntNumber(1), runtime.NewGDIntNumber(2), runtime.NewGDIntNumber(3))

	err := runtime.EqualTypes(runtime.NewGDArrayType(runtime.ComputeTypeFromObjects(array.Objects)), runtime.NewGDArrayType(runtime.GDIntType), nil)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestComputeArrayTypeWithMixedTypes(t *testing.T) {
	array := runtime.NewGDArray(runtime.NewGDIntNumber(1), runtime.GDString("two"), runtime.GDBool(true))

	aT := runtime.ComputeTypeFromObjects(array.Objects)
	err := runtime.EqualTypes(runtime.NewGDArrayType(aT), runtime.NewGDArrayType(runtime.NewGDUnionType(runtime.GDIntType, runtime.GDStringType, runtime.GDBoolType)), nil)
	if err != nil {
		t.Error(err)
		return
	}
}
