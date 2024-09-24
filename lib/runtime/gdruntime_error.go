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

const (
	IncompatibleTypeCode = iota
	InvalidArgumentTypeCode
	UnsupportedOperationCode
	UnsupportedTypeCode
	AttrNameAlreadyExistsCode
	AttrNotFoundCode
	SetObjectWrongTypeErrCode
	ObjectNotFoundErrCode
	DuplicatedObjectCreationCode
	DivByZeroCode
	IndexOutOfBoundsCode
	VariadicFnRequiresArgsCode
	FuncMissingArgsCode
	InvalidCharConversionCode
	NoFunctionCallbackErrCode
	RuntimeErrorCode
)

var (
	DivByZeroErr          = NewGDRuntimeErr(DivByZeroCode, "division by zero")
	IndexOutOfBoundsErr   = NewGDRuntimeErr(IndexOutOfBoundsCode, "index out of bounds")
	NoFunctionCallbackErr = NewGDRuntimeErr(NoFunctionCallbackErrCode, "no function callback")
)

type GDRuntimeErr struct {
	Code   int
	Msg    string
	parent *GDRuntimeErr
}

func (e GDRuntimeErr) Error() string {
	return e.Msg
}

func NewGDRuntimeErr(code int, msg string) GDRuntimeErr {
	return GDRuntimeErr{code, msg, nil}
}

func InvalidCallableTypeErr(got GDTypable) GDRuntimeErr {
	return NewGDRuntimeErr(IncompatibleTypeCode, Sprintf("invalid callable type: `%@`", got))
}

func InvalidCastingWrongTypeErr(expected, got GDTypable) GDRuntimeErr {
	return NewGDRuntimeErr(IncompatibleTypeCode, Sprintf("error while casting `%@` to `%@`", got, expected))
}

func TypeCastingWrongTypeWithHierarchyError(expected, got GDTypable, err error) GDRuntimeErr {
	return NewGDRuntimeErr(IncompatibleTypeCode, Sprintf("error while casting `%@` to `%@` ^ %@", got, expected, err.Error()))
}

func InvalidCastingExpectedTypeErr(expected GDTypable) GDRuntimeErr {
	return NewGDRuntimeErr(IncompatibleTypeCode, Sprintf("a value of type `%@` is expected", expected))
}

func InvalidCastingLitErr(lit string, typ GDTypable) GDRuntimeErr {
	return NewGDRuntimeErr(InvalidCharConversionCode, Sprintf("error trying to cast `%@` into a `%@`", lit, typ))
}

func WrongTypesErr(expected, got GDTypable) GDRuntimeErr {
	return NewGDRuntimeErr(IncompatibleTypeCode, Sprintf("expected `%@` but got `%@`", expected, got))
}

func InvalidIterableTypeErr(got GDTypable) GDRuntimeErr {
	return NewGDRuntimeErr(IncompatibleTypeCode, Sprintf("invalid iterable type: `%@`", got))
}

func InvalidMutableCollectionTypeErr(got GDTypable) GDRuntimeErr {
	return NewGDRuntimeErr(IncompatibleTypeCode, Sprintf("invalid collectable type: `%@`", got))
}

func InvalidAttributableTypeErr(got GDTypable) GDRuntimeErr {
	return NewGDRuntimeErr(IncompatibleTypeCode, Sprintf("invalid attributable type: `%@`", got))
}

func TypesAreNotEqualErr(a, b GDTypable) GDRuntimeErr {
	return NewGDRuntimeErr(IncompatibleTypeCode, Sprintf("types `%@` and `%@` are not equal", a, b))
}

func InvalidArgumentTypeErr(argName string, expected GDTypable, got GDTypable) GDRuntimeErr {
	return NewGDRuntimeErr(InvalidArgumentTypeCode, Sprintf("invalid argument type for `%@`: expected `%@` but got `%@`", argName, expected, got))
}

func UnsupportedOperationErr(operation string) GDRuntimeErr {
	return NewGDRuntimeErr(UnsupportedOperationCode, Sprintf("unsupported operation `%@`", operation))
}

func UnsupportedOperationBetweenTypesError(operation, a, b string) GDRuntimeErr {
	return NewGDRuntimeErr(UnsupportedOperationCode, Sprintf("unsupported operation `%@` between `%@` and `%@`", operation, a, b))
}

func UnsupportedTypeErr(typename string) GDRuntimeErr {
	return NewGDRuntimeErr(UnsupportedTypeCode, Sprintf("unsupported type `%@`", typename))
}

func AttributeNotFoundErr(name string) GDRuntimeErr {
	return NewGDRuntimeErr(AttrNotFoundCode, Sprintf("attribute `%@`, not found", name))
}

func SetConstObjectErr() GDRuntimeErr {
	return NewGDRuntimeErr(SetObjectWrongTypeErrCode, "can't set a constant object")
}

func ObjectNotFoundErr(key string) GDRuntimeErr {
	return NewGDRuntimeErr(ObjectNotFoundErrCode, Sprintf("object `%@` was not found", key))
}

func DuplicatedObjectCreationErr(key any) GDRuntimeErr {
	return NewGDRuntimeErr(DuplicatedObjectCreationCode, Sprintf("object `%@` was already created", key))
}

func AttributeAlreadyExistsErr(name string) GDRuntimeErr {
	return NewGDRuntimeErr(AttrNameAlreadyExistsCode, Sprintf("attribute `%@` already exists in the struct", name))
}

func MissingNumberOfArgumentsErr(expected, got uint) GDRuntimeErr {
	return NewGDRuntimeErr(FuncMissingArgsCode, Sprintf("missing number of arguments: expected `%@` but got `%@`", expected, got))
}
