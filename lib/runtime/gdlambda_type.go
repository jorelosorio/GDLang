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

type GDLambdaArgType GDKeyValue[GDIdent, GDTypable]

type GDLambdaArgTypes []GDLambdaArgType

type GDLambdaType struct {
	// Stores the ArgTypes and their types
	ArgTypes GDLambdaArgTypes
	// The type of the value that is expected to be returned
	ReturnType GDTypable
	// Meaning that the last argument is variadic
	// and can be repeated.
	IsVariadic bool
}

func (gd *GDLambdaType) GetCode() GDTypableCode {
	return GDLambdaTypeCode
}

func (gd *GDLambdaType) ToString() string {
	argsStr := JoinSlice(gd.ArgTypes, func(argType GDLambdaArgType, _ int) string {
		return Sprintf("%@: %@", argType.Key, argType.Value.ToString())
	}, ", ")

	var returnTypes string
	if gd.ReturnType != nil {
		returnTypes = Sprintf(" => %@", gd.ReturnType.ToString())
	}

	if gd.IsVariadic {
		return Sprintf("(%@, ...)%@", argsStr, returnTypes)
	}

	return Sprintf("(%@)%@", argsStr, returnTypes)
}

func (gd *GDLambdaType) CheckNumberOfArgs(argsLen uint) error {
	funcArgsLen := uint(len(gd.ArgTypes))
	if !gd.IsVariadic && argsLen != funcArgsLen {
		return MissingNumberOfArgumentsErr(funcArgsLen, argsLen)
	}

	if gd.IsVariadic && argsLen < funcArgsLen-1 {
		return MissingNumberOfArgumentsErr(funcArgsLen, argsLen)
	}

	return nil
}

func (gd *GDLambdaType) CheckArgAtIndex(index int, typ GDTypable, stack *GDSymbolStack) error {
	argTypesCount := len(gd.ArgTypes) - 1
	var lambdaArg GDLambdaArgType
	var lambdaErrArg GDTypable
	var inferredType = typ
	if index >= argTypesCount && gd.IsVariadic {
		lambdaArg = gd.ArgTypes[argTypesCount]

		switch ttype := typ.(type) {
		case GDSpreadableType:
			inferredType = ttype.GetIterableType()
		}

		lambdaErrArg = NewGDSpreadableType(NewGDArrayType(lambdaArg.Value))
	} else {
		lambdaArg = gd.ArgTypes[index]
		lambdaErrArg = lambdaArg.Value
	}

	// Check if the array elements are of the right type
	err := CanBeAssign(lambdaArg.Value, inferredType, stack)
	if err != nil {
		return InvalidArgumentTypeErr(lambdaArg.Key.ToString(), lambdaErrArg, typ)
	}

	return nil
}

func (gd *GDLambdaType) CheckReturn(retObject GDTypable, stack *GDSymbolStack) error {
	err := CanBeAssign(gd.ReturnType, retObject, stack)
	if err != nil {
		return err
	}

	return nil
}

func NewGDLambdaType(args GDLambdaArgTypes, returns GDTypable, isVariadic bool) *GDLambdaType {
	return &GDLambdaType{args, returns, isVariadic}
}
