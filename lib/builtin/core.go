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

package builtin

import "gdlang/lib/runtime"

type BuiltinMap = map[string]func(stack *runtime.GDSymbolStack) (runtime.GDObject, error)

var coreBuiltins = BuiltinMap{
	"print":   print,
	"println": println,
	"typeof":  typeof,
}

func ImportCoreBuiltins(stack *runtime.GDSymbolStack) error {
	for ident, fn := range coreBuiltins {
		obj, err := fn(stack)
		if err != nil {
			return err
		}

		ident := runtime.NewGDStringIdent(ident)
		_, err = stack.AddSymbol(ident, true, true, obj.GetType(), obj)
		if err != nil {
			return err
		}
	}

	return nil
}

// Type functions

func typeof(stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	objParam := runtime.NewStrRefType("obj")
	funcType := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{Key: objParam, Value: runtime.GDAnyType},
		},
		runtime.GDStringType,
		false,
	)
	typeOfFunc := runtime.NewGDLambdaWithType(
		funcType,
		stack,
		func(_ *runtime.GDSymbolStack, args runtime.GDLambdaArgs) (runtime.GDObject, error) {
			typeName := args.Get(objParam).GetType().ToString()
			return runtime.GDString(typeName), nil
		},
	)

	return typeOfFunc, nil
}

// Print functions

func print(stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	printFunc := printFunc(stack, false)
	return printFunc, nil
}

func println(stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	printFunc := printFunc(stack, true)
	return printFunc, nil
}

func printFunc(stack *runtime.GDSymbolStack, newLine bool) runtime.GDObject {
	argIdent := runtime.NewStrRefType("args")
	argType := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{Key: argIdent, Value: runtime.GDAnyType},
		},
		runtime.GDNilType,
		true,
	)
	printFunc := runtime.NewGDLambdaWithType(
		argType,
		stack,
		func(_ *runtime.GDSymbolStack, args runtime.GDLambdaArgs) (runtime.GDObject, error) {
			t := args.Get(argIdent).(*runtime.GDArray).Objects
			for _, arg := range t {
				argVal := runtime.GDString(arg.ToString())
				runtime.Print(argVal.Escape())
			}

			if newLine {
				runtime.Print("\n")
			}

			return runtime.GDZNil, nil
		},
	)

	return printFunc
}
