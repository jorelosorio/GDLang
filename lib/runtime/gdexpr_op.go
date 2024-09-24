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

package runtime

import (
	"math"
)

type ExprOperationType byte

const (
	ExprOperationUnaryPlus    ExprOperationType = iota // +
	ExprOperationUnaryMinus                            // -
	ExprOperationAdd                                   // +
	ExprOperationSubtract                              // -
	ExprOperationMultiply                              // *
	ExprOperationQuo                                   // /
	ExprOperationRem                                   // %
	ExprOperationGreater                               // >
	ExprOperationGreaterEqual                          // >=
	ExprOperationLess                                  // <
	ExprOperationLessEqual                             // <=
	ExprOperationEqual                                 // =
	ExprOperationNotEqual                              // !=
	ExprOperationAnd                                   // &&
	ExprOperationOr                                    // ||
	ExprOperationNot                                   // !
)

var ExprOperationMap = map[ExprOperationType]string{
	ExprOperationUnaryPlus:    "positive",
	ExprOperationUnaryMinus:   "negative",
	ExprOperationAdd:          "+",
	ExprOperationSubtract:     "-",
	ExprOperationMultiply:     "*",
	ExprOperationQuo:          "/",
	ExprOperationRem:          "%",
	ExprOperationGreater:      ">",
	ExprOperationGreaterEqual: ">=",
	ExprOperationLess:         "<",
	ExprOperationLessEqual:    "<=",
	ExprOperationEqual:        "==",
	ExprOperationNotEqual:     "!=",
	ExprOperationAnd:          "&&",
	ExprOperationOr:           "||",
	ExprOperationNot:          "!",
}

func IsUnaryOperation(op ExprOperationType) bool {
	return op == ExprOperationUnaryPlus || op == ExprOperationUnaryMinus || op == ExprOperationNot
}

func TypeCheckExprOperation(op ExprOperationType, a, b GDObject) (GDTypable, error) {
	isUnaryOp := IsUnaryOperation(op)
	if a == GDZNil {
		return GDNilType, nil
	} else if b == GDZNil && !isUnaryOp {
		return GDNilType, nil
	}

	switch {
	case IsString(a) || IsString(b):
		switch op {
		case ExprOperationAdd:
			return GDStringType, nil
		case ExprOperationGreater, ExprOperationGreaterEqual, ExprOperationLess, ExprOperationLessEqual, ExprOperationEqual, ExprOperationNotEqual:
			return GDBoolType, nil
		default:
			return nil, UnsupportedOperationBetweenTypesError(ExprOperationMap[op], a.GetType().ToString(), b.GetType().ToString())
		}
	case IsComplex(a) || IsComplex(b):
		if isUnaryOp {
			return GDComplexType, nil
		}

		switch op {
		case ExprOperationAdd, ExprOperationSubtract, ExprOperationMultiply, ExprOperationQuo:
			return GDComplexType, nil
		case ExprOperationGreater, ExprOperationGreaterEqual, ExprOperationLess, ExprOperationLessEqual, ExprOperationEqual, ExprOperationNotEqual:
			return GDBoolType, nil
		default:
			return nil, UnsupportedOperationBetweenTypesError(ExprOperationMap[op], a.GetType().ToString(), b.GetType().ToString())
		}
	case IsFloat(a) || IsFloat(b):
		if isUnaryOp {
			return GDFloatType, nil
		}

		switch op {
		case ExprOperationAdd, ExprOperationSubtract, ExprOperationMultiply, ExprOperationQuo, ExprOperationRem:
			return GDFloatType, nil
		case ExprOperationGreater, ExprOperationGreaterEqual, ExprOperationLess, ExprOperationLessEqual, ExprOperationEqual, ExprOperationNotEqual:
			return GDBoolType, nil
		}

		return nil, UnsupportedOperationBetweenTypesError(ExprOperationMap[op], a.GetType().ToString(), b.GetType().ToString())
	case IsInt(a) || IsInt(b):
		if isUnaryOp {
			return GDIntType, nil
		}

		switch op {
		case ExprOperationAdd, ExprOperationSubtract, ExprOperationMultiply, ExprOperationQuo, ExprOperationRem:
			return GDIntType, nil
		case ExprOperationGreater, ExprOperationGreaterEqual, ExprOperationLess, ExprOperationLessEqual, ExprOperationEqual, ExprOperationNotEqual:
			return GDBoolType, nil
		}
	case IsChar(a) || IsChar(b):
		switch op {
		case ExprOperationAdd:
			return GDStringType, nil
		case ExprOperationGreater, ExprOperationGreaterEqual, ExprOperationLess, ExprOperationLessEqual, ExprOperationEqual, ExprOperationNotEqual:
			return GDBoolType, nil
		default:
			return nil, UnsupportedOperationBetweenTypesError(ExprOperationMap[op], a.GetType().ToString(), b.GetType().ToString())
		}
	case IsBool(a) || IsBool(b):
		if isUnaryOp {
			return GDBoolType, nil
		}

		switch op {
		case ExprOperationAnd, ExprOperationOr, ExprOperationEqual, ExprOperationNotEqual:
			return GDBoolType, nil
		}
	}

	return nil, UnsupportedOperationBetweenTypesError(ExprOperationMap[op], a.GetType().ToString(), b.GetType().ToString())
}

func PerformExprOperation(op ExprOperationType, a, b GDObject) (GDObject, error) {
	isUnaryOp := IsUnaryOperation(op)
	if a == GDZNil {
		return GDZNil, nil
	} else if b == GDZNil && !isUnaryOp {
		return GDZNil, nil
	}

	switch {
	case IsString(a) || IsString(b):
		sA, err := ToString(a)
		if err != nil {
			return nil, err
		}

		sB, err := ToString(b)
		if err != nil {
			return nil, err
		}

		return performStringOp(op, sA, sB)
	case IsComplex(a) || IsComplex(b):
		nA, err := ToComplex(a)
		if err != nil {
			return nil, err
		}

		if isUnaryOp {
			return performComplexOp(op, nA, 0)
		}

		nB, err := ToComplex(b)
		if err != nil {
			return nil, err
		}

		return performComplexOp(op, nA, nB)
	case IsFloat(a) || IsFloat(b):
		nA, err := ToFloat(a)
		if err != nil {
			return nil, err
		}

		if isUnaryOp {
			return performFloatOp(op, nA, 0)
		}

		nB, err := ToFloat(b)
		if err != nil {
			return nil, err
		}

		return performFloatOp(op, nA, nB)
	case IsInt(a) || IsInt(b):
		nA, err := ToInt(a)
		if err != nil {
			return nil, err
		}

		if isUnaryOp {
			return performIntOp(op, nA, 0)
		}

		nB, err := ToInt(b)
		if err != nil {
			return nil, err
		}

		return performIntOp(op, nA, nB)
	case IsChar(a) || IsChar(b):
		cA, err := ToString(a)
		if err != nil {
			return nil, err
		}

		cB, err := ToString(b)
		if err != nil {
			return nil, err
		}

		return performStringOp(op, cA, cB)
	case IsBool(a) || IsBool(b):
		bA, err := ToBool(a)
		if err != nil {
			return nil, err
		}

		if isUnaryOp {
			return performLogicalOp(op, bA, false)
		}

		bB, err := ToBool(b)
		if err != nil {
			return nil, err
		}

		return performLogicalOp(op, bA, bB)
	}

	return nil, UnsupportedOperationBetweenTypesError(ExprOperationMap[op], a.GetType().ToString(), b.GetType().ToString())
}

func performLogicalOp(op ExprOperationType, a, b GDBool) (GDObject, error) {
	switch op {
	case ExprOperationNot:
		return !a, nil
	case ExprOperationAnd:
		return a && b, nil
	case ExprOperationOr:
		return a || b, nil
	case ExprOperationEqual:
		return GDBool(a == b), nil
	case ExprOperationNotEqual:
		return GDBool(a != b), nil
	}

	return nil, UnsupportedOperationBetweenTypesError(ExprOperationMap[op], a.GetType().ToString(), b.GetType().ToString())
}

func performStringOp(op ExprOperationType, a, b GDString) (GDObject, error) {
	switch op {
	case ExprOperationAdd:
		return a + b, nil
	case ExprOperationGreater:
		return GDBool(a > b), nil
	case ExprOperationGreaterEqual:
		return GDBool(a >= b), nil
	case ExprOperationLess:
		return GDBool(a < b), nil
	case ExprOperationLessEqual:
		return GDBool(a <= b), nil
	case ExprOperationEqual:
		return GDBool(a == b), nil
	case ExprOperationNotEqual:
		return GDBool(a != b), nil
	}

	return nil, UnsupportedOperationBetweenTypesError(ExprOperationMap[op], a.GetType().ToString(), b.GetType().ToString())
}

func performIntOp(op ExprOperationType, a, b GDInt) (GDObject, error) {
	switch op {
	case ExprOperationUnaryPlus:
		return NewGDIntNumber(+a), nil
	case ExprOperationUnaryMinus:
		return NewGDIntNumber(-a), nil
	case ExprOperationAdd:
		return NewGDIntNumber(a + b), nil
	case ExprOperationSubtract:
		return NewGDIntNumber(a - b), nil
	case ExprOperationMultiply:
		return NewGDIntNumber(a * b), nil
	case ExprOperationQuo:
		if b == 0 {
			return nil, DivByZeroErr
		}

		return NewGDIntNumber(a / b), nil
	case ExprOperationRem:
		if b == 0 {
			return nil, DivByZeroErr
		}

		return NewGDIntNumber(a % b), nil
	case ExprOperationGreater:
		return GDBool(a > b), nil
	case ExprOperationGreaterEqual:
		return GDBool(a >= b), nil
	case ExprOperationLess:
		return GDBool(a < b), nil
	case ExprOperationLessEqual:
		return GDBool(a <= b), nil
	case ExprOperationEqual:
		return GDBool(a == b), nil
	case ExprOperationNotEqual:
		return GDBool(a != b), nil
	}

	return nil, UnsupportedOperationBetweenTypesError(ExprOperationMap[op], a.GetType().ToString(), b.GetType().ToString())
}

func performFloatOp(op ExprOperationType, a, b GDFloat64) (GDObject, error) {
	switch op {
	case ExprOperationUnaryPlus:
		return NewGDFloatNumber(+a), nil
	case ExprOperationUnaryMinus:
		return NewGDFloatNumber(-a), nil
	case ExprOperationAdd:
		return NewGDFloatNumber(a + b), nil
	case ExprOperationSubtract:
		return NewGDFloatNumber(a - b), nil
	case ExprOperationMultiply:
		return NewGDFloatNumber(a * b), nil
	case ExprOperationQuo:
		if b == 0 {
			return nil, DivByZeroErr
		}

		return NewGDFloatNumber(a / b), nil
	case ExprOperationRem:
		if b == 0 {
			return nil, DivByZeroErr
		}

		return NewGDFloatNumber(GDFloat64(math.Mod(float64(a), float64(b)))), nil
	case ExprOperationGreater:
		return GDBool(a > b), nil
	case ExprOperationGreaterEqual:
		return GDBool(a >= b), nil
	case ExprOperationLess:
		return GDBool(a < b), nil
	case ExprOperationLessEqual:
		return GDBool(a <= b), nil
	case ExprOperationEqual:
		return GDBool(a == b), nil
	case ExprOperationNotEqual:
		return GDBool(a != b), nil
	}

	return nil, UnsupportedOperationBetweenTypesError(ExprOperationMap[op], a.GetType().ToString(), b.GetType().ToString())
}

func performComplexOp(op ExprOperationType, a, b GDComplex128) (GDObject, error) {
	switch op {
	case ExprOperationUnaryPlus:
		return NewGDComplexNumber(+a), nil
	case ExprOperationUnaryMinus:
		return NewGDComplexNumber(-a), nil
	case ExprOperationAdd:
		return NewGDComplexNumber(a + b), nil
	case ExprOperationSubtract:
		return NewGDComplexNumber(a - b), nil
	case ExprOperationMultiply:
		return NewGDComplexNumber(a * b), nil
	case ExprOperationQuo:
		if b == 0 {
			return nil, DivByZeroErr
		}

		return NewGDComplexNumber(a / b), nil
	case ExprOperationRem:
		return nil, UnsupportedOperationErr(ExprOperationMap[op])
	case ExprOperationGreater:
		return GDBool(real(a) > real(b) && imag(a) > imag(b)), nil
	case ExprOperationGreaterEqual:
		return GDBool(real(a) >= real(b) && imag(a) >= imag(b)), nil
	case ExprOperationLess:
		return GDBool(real(a) < real(b) && imag(a) < imag(b)), nil
	case ExprOperationLessEqual:
		return GDBool(real(a) <= real(b) && imag(a) <= imag(b)), nil
	case ExprOperationEqual:
		return GDBool(a == b), nil
	case ExprOperationNotEqual:
		return GDBool(a != b), nil
	}

	return nil, UnsupportedOperationBetweenTypesError(ExprOperationMap[op], a.GetType().ToString(), b.GetType().ToString())
}
