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

func print(stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	printFunc := printFunc(stack, false)
	return printFunc, nil
}

func println(stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	printFunc := printFunc(stack, true)
	return printFunc, nil
}

func printFunc(stack *runtime.GDSymbolStack, newLine bool) runtime.GDObject {
	argIdent := runtime.GDStringIdentType("args")
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
