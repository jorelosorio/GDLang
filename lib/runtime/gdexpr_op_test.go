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
	"testing"
)

func TestNumberAddition(t *testing.T) {
	for _, test := range []struct {
		a, b     runtime.GDObject
		expected runtime.GDObject
	}{
		{runtime.NewGDIntNumber(10), runtime.NewGDIntNumber(5), runtime.NewGDIntNumber(15)},
		{runtime.NewGDIntNumber(10), runtime.NewGDIntNumber(-5), runtime.NewGDIntNumber(5)},
		{runtime.NewGDIntNumber(10), runtime.NewGDIntNumber(0), runtime.NewGDIntNumber(10)},
		{runtime.NewGDIntNumber(10), runtime.NewGDFloatNumber(math.Pi), runtime.NewGDFloatNumber(13.141593)},
		{runtime.NewGDIntNumber(10), runtime.NewGDFloatNumber(runtime.GDFloat64(math.Inf(1))), runtime.NewGDFloatNumber(runtime.GDFloat64(math.Inf(1)))},
		{runtime.NewGDIntNumber(10), runtime.NewGDFloatNumber(runtime.GDFloat64(math.Inf(-1))), runtime.NewGDFloatNumber(runtime.GDFloat64(math.Inf(-1)))},
	} {
		result, err := runtime.PerformExprOperation(runtime.ExprOperationAdd, test.a, test.b)
		if err != nil {
			t.Errorf("Error while performing operation: %v", err)
		}

		if result != test.expected {
			t.Errorf("Expected %v, but got %v", test.expected, result)
		}
	}
}

func TestNumberSubtraction(t *testing.T) {
	for _, test := range []struct {
		a, b     runtime.GDObject
		expected runtime.GDObject
	}{
		{runtime.NewGDIntNumber(10), runtime.NewGDIntNumber(5), runtime.NewGDIntNumber(5)},
		{runtime.NewGDIntNumber(10), runtime.NewGDIntNumber(-5), runtime.NewGDIntNumber(15)},
		{runtime.NewGDIntNumber(10), runtime.NewGDIntNumber(0), runtime.NewGDIntNumber(10)},
		{runtime.NewGDIntNumber(10), runtime.NewGDFloatNumber(runtime.GDFloat64(math.Pi)), runtime.NewGDFloatNumber(6.858407)},
		{runtime.NewGDIntNumber(10), runtime.NewGDFloatNumber(runtime.GDFloat64(math.Inf(1))), runtime.NewGDFloatNumber(runtime.GDFloat64(math.Inf(-1)))},
		{runtime.NewGDIntNumber(10), runtime.NewGDFloatNumber(runtime.GDFloat64(math.Inf(-1))), runtime.NewGDFloatNumber(runtime.GDFloat64(math.Inf(1)))},
	} {
		result, err := runtime.PerformExprOperation(runtime.ExprOperationSubtract, test.a, test.b)
		if err != nil {
			t.Errorf("Error while performing operation: %v", err)
		}
		if result != test.expected {
			t.Errorf("Expected %v, but got %v", test.expected, result)
		}
	}
}

func TestNumberMultiplication(t *testing.T) {
	for _, test := range []struct {
		a, b     runtime.GDObject
		expected runtime.GDObject
	}{
		{runtime.NewGDIntNumber(10), runtime.NewGDIntNumber(5), runtime.NewGDIntNumber(50)},
		{runtime.NewGDIntNumber(10), runtime.NewGDIntNumber(-5), runtime.NewGDIntNumber(-50)},
		{runtime.NewGDIntNumber(10), runtime.NewGDIntNumber(0), runtime.NewGDIntNumber(0)},
		{runtime.NewGDIntNumber(10), runtime.NewGDFloatNumber(3.1416), runtime.NewGDFloatNumber(31.415998)},
		{runtime.NewGDIntNumber(10), runtime.NewGDFloatNumber(runtime.GDFloat64(math.Inf(1))), runtime.NewGDFloatNumber(runtime.GDFloat64(math.Inf(1)))},
		{runtime.NewGDIntNumber(10), runtime.NewGDFloatNumber(runtime.GDFloat64(math.Inf(-1))), runtime.NewGDFloatNumber(runtime.GDFloat64(math.Inf(-1)))},
	} {
		result, err := runtime.PerformExprOperation(runtime.ExprOperationMultiply, test.a, test.b)
		if err != nil {
			t.Errorf("Error while performing operation: %v", err)
		}

		if result != test.expected {
			t.Errorf("Expected %v, but got %v", test.expected, result)
		}
	}
}

func TestNumberDivision(t *testing.T) {
	for _, test := range []struct {
		a, b     runtime.GDObject
		expected runtime.GDObject
	}{
		{runtime.NewGDIntNumber(10), runtime.NewGDIntNumber(5), runtime.NewGDIntNumber(2)},
		{runtime.NewGDIntNumber(10), runtime.NewGDIntNumber(-5), runtime.NewGDIntNumber(-2)},
		{runtime.NewGDIntNumber(10), runtime.NewGDFloatNumber(3.1416), runtime.NewGDFloatNumber(3.1830916)},
	} {
		result, err := runtime.PerformExprOperation(runtime.ExprOperationQuo, test.a, test.b)
		if err != nil {
			t.Errorf("Error while performing operation: %v", err)
		}

		if result != test.expected {
			t.Errorf("Expected %v, but got %v", test.expected, result)
		}
	}
}

func TestNumberModulus(t *testing.T) {
	for _, test := range []struct {
		a, b     runtime.GDObject
		expected runtime.GDObject
	}{
		{runtime.NewGDIntNumber(10), runtime.NewGDIntNumber(5), runtime.NewGDIntNumber(0)},
		{runtime.NewGDIntNumber(10), runtime.NewGDIntNumber(-5), runtime.NewGDIntNumber(0)},
	} {
		result, err := runtime.PerformExprOperation(runtime.ExprOperationRem, test.a, test.b)
		if err != nil {
			t.Errorf("Error while performing operation: %v", err)
		}

		if result != test.expected {
			t.Errorf("Expected %v, but got %v", test.expected, result)
		}
	}
}

func TestNumberOperationModulusWithWrongTypes(t *testing.T) {
	_, err := runtime.PerformExprOperation(runtime.ExprOperationRem, runtime.NewGDIntNumber(10), runtime.NewGDIntNumber(complex(0, 0)))
	if err == nil {
		t.Errorf("Expected error performing modulus between int and complex, but got nil")
	}

	_, err = runtime.PerformExprOperation(runtime.ExprOperationRem, runtime.NewGDIntNumber(10), runtime.NewGDIntNumber(0))
	if err == nil {
		t.Errorf("Expected division by zero error, but got nil")
	}
}

func TestZeroDivisionError(t *testing.T) {
	_, err := runtime.PerformExprOperation(runtime.ExprOperationQuo, runtime.NewGDIntNumber(10), runtime.NewGDIntNumber(complex(0, 0)))
	if err == nil {
		t.Errorf("Expected division by zero error, but got nil")
	}

	_, err = runtime.PerformExprOperation(runtime.ExprOperationQuo, runtime.NewGDIntNumber(10.0), runtime.NewGDIntNumber(0.0))
	if err == nil {
		t.Errorf("Expected division by zero error, but got nil")
	}

	_, err = runtime.PerformExprOperation(runtime.ExprOperationQuo, runtime.NewGDIntNumber(10), runtime.NewGDIntNumber(0))
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}
}

func TestModulusByZeroError(t *testing.T) {
	_, err := runtime.PerformExprOperation(runtime.ExprOperationRem, runtime.NewGDIntNumber(10), runtime.NewGDIntNumber(complex(0, 0)))
	if err == nil {
		t.Errorf("Expected division by zero error, but got nil")
	}

	_, err = runtime.PerformExprOperation(runtime.ExprOperationRem, runtime.NewGDIntNumber(10.0), runtime.NewGDIntNumber(0.0))
	if err == nil {
		t.Errorf("Expected division by zero error, but got nil")
	}

	_, err = runtime.PerformExprOperation(runtime.ExprOperationRem, runtime.NewGDIntNumber(10), runtime.NewGDIntNumber(0))
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}
}

func TestStringAddition(t *testing.T) {
	for _, test := range []struct {
		a, b     runtime.GDObject
		expected runtime.GDObject
	}{
		{runtime.GDString("Hello"), runtime.GDString("World"), runtime.GDString("HelloWorld")},
		{runtime.GDString("Hello"), runtime.NewGDIntNumber(5), runtime.GDString("Hello5")},
		{runtime.GDString("Hello"), runtime.NewGDFloatNumber(3.1416), runtime.GDString("Hello3.1416")},
	} {
		result, err := runtime.PerformExprOperation(runtime.ExprOperationAdd, test.a, test.b)
		if err != nil {
			t.Errorf("Error while performing operation: %v", err)
		}

		if result != test.expected {
			t.Errorf("Expected %v, but got %v", test.expected, result)
		}
	}
}

func TestWrongTypesAddition(t *testing.T) {
	result, err := runtime.PerformExprOperation(runtime.ExprOperationAdd, runtime.GDChar('b'), runtime.GDChar('a'))
	if err != nil {
		t.Errorf("Error while performing operation: %v", err)
	}

	if !runtime.EqualObjects(result, runtime.GDString("ba")) {
		t.Errorf("Expected %q, but got %v", "ba", result)
	}
}
