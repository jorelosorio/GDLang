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

func TestBasicFunctionType(t *testing.T) {
	funcType := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{aParamIdent, runtime.GDIntType},
		},
		runtime.NewGDTupleType(),
		false,
	)

	if funcType.GetCode() != runtime.GDLambdaTypeCode {
		t.Errorf("Expected to be a function type")
	}

	funcTypeStr := funcType.ToString()
	if funcTypeStr != "(a: int) => (untyped,)" {
		t.Errorf("Expected (a: int) => (untyped,), got %v", funcTypeStr)
	}

	sameFuncType := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{bParamIdent, runtime.GDIntType},
		},
		runtime.NewGDTupleType(),
		false,
	)

	err := runtime.EqualTypes(funcType, sameFuncType, nil)
	if err != nil {
		t.Errorf("Error comparing types: %v", err)
	}
}

func TestComplexFunctionType(t *testing.T) {
	funcType := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{aParamIdent, runtime.GDIntType},
			{bParamIdent, runtime.GDStringType},
		},
		runtime.NewGDTupleType(runtime.GDIntType, runtime.GDStringType),
		false,
	)

	if funcType.GetCode() != runtime.GDLambdaTypeCode {
		t.Errorf("Expected to be a function type")
	}

	funcTypeStr := funcType.ToString()
	if funcTypeStr != "(a: int, b: string) => (int, string)" {
		t.Errorf("Expected (a: int, b: string) => (int, string), got %v", funcTypeStr)
	}
}

func TestVariadicFunctionType(t *testing.T) {
	funcType := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{aParamIdent, runtime.GDIntType},
			{bParamIdent, runtime.GDStringType},
		},
		runtime.NewGDTupleType(runtime.GDIntType, runtime.GDStringType),
		true,
	)

	if funcType.GetCode() != runtime.GDLambdaTypeCode {
		t.Errorf("Expected to be a function type")
	}

	funcTypeStr := funcType.ToString()
	if funcTypeStr != "(a: int, b: string, ...) => (int, string)" {
		t.Errorf("Expected (a: int, b: string, ...) => (int, string), got %v", funcTypeStr)
	}
}

func TestFunctionWithNoArguments(t *testing.T) {
	funcType := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{},
		runtime.NewGDTupleType(runtime.GDIntType),
		false,
	)

	funcTypeStr := funcType.ToString()
	if funcTypeStr != "() => (int,)" {
		t.Errorf("Expected () => (int,) got %v", funcTypeStr)
	}
}

func TestVoidFunction(t *testing.T) {
	funcType := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{},
		runtime.NewGDTupleType(),
		false,
	)

	funcTypeStr := funcType.ToString()
	if funcTypeStr != "() => (untyped,)" {
		t.Errorf("Expected () => (untyped,), got %v", funcTypeStr)
	}
}

func TestCompareFunctionTypes(t *testing.T) {
	funcType1 := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{aParamIdent, runtime.GDIntType},
			{bParamIdent, runtime.GDStringType},
		},
		runtime.NewGDTupleType(runtime.GDIntType, runtime.GDStringType),
		false,
	)

	funcType2 := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{aParamIdent, runtime.GDIntType},
			{bParamIdent, runtime.GDStringType},
		},
		runtime.NewGDTupleType(runtime.GDIntType, runtime.GDStringType),
		false,
	)

	err := runtime.EqualTypes(funcType1, funcType2, nil)
	if err != nil {
		t.Errorf("Error comparing types: %v", err)
	}
}

func TestCompareFunctionTypesWithVariadic(t *testing.T) {
	funcType1 := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{aParamIdent, runtime.GDIntType},
		},
		runtime.NewGDTupleType(runtime.GDIntType, runtime.GDStringType),
		true,
	)

	funcType2 := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{aParamIdent, runtime.GDIntType},
		},
		runtime.NewGDTupleType(runtime.GDIntType, runtime.GDStringType),
		true,
	)

	err := runtime.EqualTypes(funcType1, funcType2, nil)
	if err != nil {
		t.Errorf("Error comparing types: %v", err)
	}
}

func TestCompareFunctionTypesWithDifferentVariadic(t *testing.T) {
	funcType1 := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{aParamIdent, runtime.GDIntType},
		},
		runtime.NewGDTupleType(runtime.GDIntType, runtime.GDStringType),
		true,
	)

	funcType2 := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{aParamIdent, runtime.GDIntType},
		},
		runtime.NewGDTupleType(runtime.GDIntType, runtime.GDStringType),
		false,
	)

	err := runtime.EqualTypes(funcType1, funcType2, nil)
	if err == nil {
		t.Errorf("Expected error comparing types")
	}
}

func TestCompareFunctionTypesWithDifferentArgs(t *testing.T) {
	funcType1 := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{aParamIdent, runtime.GDStringType},
			{bParamIdent, runtime.GDStringType},
		},
		runtime.NewGDTupleType(runtime.GDIntType, runtime.GDStringType),
		false,
	)

	funcType2 := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{aParamIdent, runtime.GDIntType},
			{bParamIdent, runtime.GDStringType},
		},
		runtime.NewGDTupleType(runtime.GDIntType, runtime.GDStringType),
		false,
	)

	err := runtime.EqualTypes(funcType1, funcType2, nil)
	if err == nil {
		t.Errorf("Expected error comparing types")
	}
}

func TestCompareFunctionTypesWithDifferentReturns(t *testing.T) {
	funcType1 := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{},
		runtime.NewGDTupleType(runtime.GDIntType, runtime.GDStringType),
		false,
	)

	funcType2 := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{},
		runtime.NewGDTupleType(runtime.GDIntType, runtime.GDStringType, runtime.GDStringType),
		false,
	)

	err := runtime.EqualTypes(funcType1, funcType2, nil)
	if err == nil {
		t.Errorf("Expected error comparing types")
	}
}
