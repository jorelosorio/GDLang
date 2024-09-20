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
	"strings"
	"testing"
)

func TestNilLambdaCall(t *testing.T) {
	lambda := runtime.NewGDLambda(runtime.GDLambdaArgTypes{
		{aParamIdent, runtime.GDIntType},
	}, runtime.NewGDTupleType(), false, nil, func(_ *runtime.GDSymbolStack, args runtime.GDLambdaArgs) (runtime.GDObject, error) {
		return runtime.NewGDTuple(), nil
	})

	_, err := lambda.CheckArgValues(runtime.NewGDArray(runtime.NewGDIntNumber(1)))
	if err != nil {
		t.Errorf("Error calling lambda: %v", err)
	}
}

func TestBuildTypedObjectForLambda(t *testing.T) {
	lambda := runtime.NewGDLambda(runtime.GDLambdaArgTypes{
		{aParamIdent, runtime.GDIntType},
		{bParamIdent, runtime.GDStringType},
	}, runtime.NewGDTupleType(), false, nil, func(_ *runtime.GDSymbolStack, args runtime.GDLambdaArgs) (runtime.GDObject, error) {
		return runtime.NewGDTuple(), nil
	})

	args, err := lambda.CheckArgValues(runtime.NewGDArray(runtime.NewGDIntNumber(1), runtime.GDString("hello")))

	if err != nil {
		t.Errorf("Error building typed object args: %v", err)
		return
	}

	if len(args) != 2 {
		t.Errorf("Expected 2 args, got %v", len(args))
		return
	}

	if !runtime.EqualObjects(args.Get(aParamIdent), runtime.NewGDIntNumber(1)) {
		t.Errorf("Wrong type for arg `a`, got %v", args.Get(aParamIdent))
		return
	}

	if !runtime.EqualObjects(args.Get(bParamIdent), runtime.GDString("hello")) {
		t.Errorf("Wrong type for arg `b`, got %v", args.Get(bParamIdent).ToString())
		return
	}
}

func TestStructWithWrongParameterTypes(t *testing.T) {
	lambda := runtime.NewGDLambda(runtime.GDLambdaArgTypes{
		{aParamIdent, runtime.GDIntType},
		{bParamIdent, runtime.GDStringType},
	}, runtime.NewGDTupleType(), false, nil, func(_ *runtime.GDSymbolStack, args runtime.GDLambdaArgs) (runtime.GDObject, error) {
		return runtime.NewGDTuple(), nil
	})

	_, err := lambda.CheckArgValues(runtime.NewGDArray(runtime.GDString("hello"), runtime.NewGDIntNumber(1)))
	if err == nil {
		t.Errorf("Expected an error")
	}
}

func TestWithWrongNumberOfArgumentsLessThanAllowed(t *testing.T) {
	lambda := runtime.NewGDLambda(runtime.GDLambdaArgTypes{
		{aParamIdent, runtime.GDIntType},
		{bParamIdent, runtime.GDStringType},
	}, runtime.NewGDTupleType(runtime.GDIntType, runtime.GDStringType), false, nil, func(_ *runtime.GDSymbolStack, args runtime.GDLambdaArgs) (runtime.GDObject, error) {
		return runtime.NewGDTuple(), nil
	})

	_, err := lambda.CheckArgValues(runtime.NewGDArray(runtime.NewGDIntNumber(1)))
	if err == nil {
		t.Errorf("Expected an error")
	}
}

func TestWithWrongNumberOfArgumentsMoreThanAllowed(t *testing.T) {
	lambda := runtime.NewGDLambda(runtime.GDLambdaArgTypes{
		{aParamIdent, runtime.GDIntType},
		{bParamIdent, runtime.GDStringType},
	}, runtime.NewGDTupleType(runtime.GDIntType, runtime.GDStringType), false, nil, func(_ *runtime.GDSymbolStack, args runtime.GDLambdaArgs) (runtime.GDObject, error) {
		return runtime.NewGDTuple(), nil
	})

	_, err := lambda.CheckArgValues(runtime.NewGDArray(runtime.NewGDIntNumber(1), runtime.GDString("hello"), runtime.GDString("hello")))
	if err == nil {
		t.Errorf("Expected an error")
	}
}

func TestVariadicLambdaCallWithVariadicArguments(t *testing.T) {
	lambda := runtime.NewGDLambda(runtime.GDLambdaArgTypes{
		{aParamIdent, runtime.GDIntType},
		{bParamIdent, runtime.GDStringType},
	}, runtime.NewGDTupleType(runtime.GDIntType, runtime.GDStringType), true, nil, func(_ *runtime.GDSymbolStack, args runtime.GDLambdaArgs) (runtime.GDObject, error) {
		return runtime.NewGDTuple(), nil
	})

	args, err := lambda.CheckArgValues(runtime.NewGDArray(runtime.NewGDIntNumber(1), runtime.GDString("hello"), runtime.GDString("hello")))
	if err != nil {
		t.Errorf("Error calling lambda: %v", err)
	}

	if len(args) != 2 {
		t.Errorf("Expected 2 typed objects, got %v", len(args))
		return
	}

	if !runtime.EqualObjects(args.Get(aParamIdent), runtime.NewGDIntNumber(1)) {
		t.Errorf("Wrong type for arg `a`, got %v", args.Get(aParamIdent).ToString())
		return
	}

	if !runtime.EqualObjects(args.Get(bParamIdent), runtime.NewGDArray(runtime.GDString("hello"), runtime.GDString("hello"))) {
		t.Errorf("Wrong type for arg `b`, got %v", args.Get(bParamIdent).ToString())
		return
	}
}

func TestVariadicArgsSendingEmptyArgs(t *testing.T) {
	lambda := runtime.NewGDLambda(runtime.GDLambdaArgTypes{
		{aParamIdent, runtime.GDIntType},
	}, runtime.NewGDTupleType(), true, nil, func(_ *runtime.GDSymbolStack, args runtime.GDLambdaArgs) (runtime.GDObject, error) {
		return runtime.NewGDTuple(), nil
	})

	args, err := lambda.CheckArgValues(runtime.NewGDArray())
	if err != nil {
		t.Errorf("Error calling lambda: %v", err)
		return
	}

	if len(args) != 1 {
		t.Errorf("Expected 1 typed objects, got %v", len(args))
		return
	}

	// Type should be compatible with empty arrays
	err = runtime.CanBeAssign(args.Get(aParamIdent).GetType(), runtime.NewGDArray().GetType(), nil)
	if err != nil {
		t.Errorf("Error checking if types are assignable: %v", err)
		return
	}
}

func TestSimpleLambdaCallReturnTypesWrong(t *testing.T) {
	lambda := runtime.NewGDLambda(runtime.GDLambdaArgTypes{
		{aParamIdent, runtime.GDIntType},
		{bParamIdent, runtime.GDStringType},
	}, runtime.NewGDTupleType(runtime.GDStringType, runtime.GDStringType), false, nil, func(_ *runtime.GDSymbolStack, args runtime.GDLambdaArgs) (runtime.GDObject, error) {
		return runtime.NewGDTuple(runtime.GDString("test"), runtime.NewGDIntNumber(2)), nil
	})

	_, err := lambda.Call(runtime.NewGDArray(runtime.NewGDIntNumber(1), runtime.GDString("hello")))
	errMsg := runtime.WrongTypesErr(runtime.NewGDTupleType(runtime.GDStringType, runtime.GDStringType), runtime.NewGDTupleType(runtime.GDStringType, runtime.GDIntType)).Error()
	if err != nil && !strings.Contains(err.Error(), errMsg) {
		t.Errorf("Expected error message to contain %q, got  %q", errMsg, err.Error())
		return
	}
}

func TestSimpleLambdaCallReturnTypes(t *testing.T) {
	lambda := runtime.NewGDLambda(runtime.GDLambdaArgTypes{
		{aParamIdent, runtime.GDIntType},
		{bParamIdent, runtime.GDStringType},
	}, runtime.NewGDTupleType(runtime.GDStringType, runtime.GDStringType), false, nil, func(_ *runtime.GDSymbolStack, args runtime.GDLambdaArgs) (runtime.GDObject, error) {
		return runtime.NewGDTuple(runtime.GDString("test"), runtime.GDString("test2")), nil
	})

	returns, err := lambda.Call(runtime.NewGDArray(runtime.NewGDIntNumber(1), runtime.GDString("hello")))
	if err != nil {
		t.Errorf("Error calling lambda: %v", err)
		return
	}

	if tuple, isTuple := returns.(*runtime.GDTuple); isTuple {
		if len(tuple.Objects) != 2 {
			t.Errorf("Expected 2 returns, got %v", len(tuple.Objects))
			return
		}
	} else {
		t.Errorf("Expected tuple, got %v", returns.GetType().ToString())
		return
	}

	err = runtime.EqualTypes(returns.GetType(), lambda.Type.ReturnType, nil)
	if err != nil {
		t.Errorf("Error comparing types: %v", err)
		return
	}
}

func TestSimpleLambdaCallNilReturnTypesVariadic(t *testing.T) {
	lambda := runtime.NewGDLambda(runtime.GDLambdaArgTypes{
		{aParamIdent, runtime.GDIntType},
		{bParamIdent, runtime.GDStringType},
	}, runtime.GDNilType, true, nil, func(_ *runtime.GDSymbolStack, args runtime.GDLambdaArgs) (runtime.GDObject, error) {
		return runtime.GDZNil, nil
	})

	nilReturn, err := lambda.Call(runtime.NewGDArray(runtime.NewGDIntNumber(1), runtime.GDString("hello")))
	if err != nil {
		t.Errorf("Error calling lambda: %v", err)
		return
	}

	if !runtime.EqualObjects(nilReturn, runtime.GDZNil) {
		t.Errorf("Expected string type, got %v", nilReturn.GetType().ToString())
		return
	}
}

func TestLambdaStringsButReturnsNils(t *testing.T) {
	lambda := runtime.NewGDLambda(runtime.GDLambdaArgTypes{
		{aParamIdent, runtime.GDIntType},
		{bParamIdent, runtime.GDStringType},
	}, runtime.NewGDTupleType(runtime.GDStringType, runtime.GDStringType), false, nil, func(_ *runtime.GDSymbolStack, args runtime.GDLambdaArgs) (runtime.GDObject, error) {
		return runtime.NewGDTuple(runtime.GDZNil, runtime.GDZNil), nil
	})

	returns, err1 := lambda.Call(runtime.NewGDArray(runtime.NewGDIntNumber(1), runtime.GDString("hello")))
	if err1 != nil {
		t.Errorf("Error calling lambda: %v", err1)
	}

	err := runtime.EqualTypes(returns.GetType(), runtime.NewGDTupleType(runtime.GDNilType, runtime.GDNilType), nil)
	if err != nil {
		t.Errorf("Error comparing types: %v", err)
	}
}
