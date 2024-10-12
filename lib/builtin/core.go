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

import (
	"gdlang/lib/runtime"
	"sync"
)

// Declare the variables
var (
	coreBuiltins     map[string]func() (*runtime.GDSymbol, error)
	coreBuiltinsOnce sync.Once
)

// Lazy initializer function
func GetCoreBuiltins() map[string]func() (*runtime.GDSymbol, error) {
	coreBuiltinsOnce.Do(func() {
		coreBuiltins = map[string]func() (*runtime.GDSymbol, error){
			"print":   print,
			"println": println,
			"typeof":  typeof,
		}
	})

	return coreBuiltins
}

// Type functions

func typeof() (*runtime.GDSymbol, error) {
	objIdent := runtime.NewGDStrIdent("obj")
	fnType := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{Key: objIdent, Value: runtime.GDAnyTypeRef},
		},
		runtime.GDStringTypeRef,
		false,
	)

	fn := runtime.NewGDLambdaWithType(
		fnType,
		nil,
		func(args runtime.GDLambdaArgs, stack *runtime.GDStack) (runtime.GDObject, error) {
			objType, err := args.Get(objIdent)
			if err != nil {
				return nil, err
			}

			return runtime.GDString(objType.GetType().ToString()), nil
		},
	)

	return runtime.NewGDSymbol(true, true, fnType, fn), nil
}

// Print functions

func print() (*runtime.GDSymbol, error)   { return printFunc(false), nil }
func println() (*runtime.GDSymbol, error) { return printFunc(true), nil }

func printFunc(newLine bool) *runtime.GDSymbol {
	argsIdent := runtime.NewGDStrIdent("args")
	fnType := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{Key: argsIdent, Value: runtime.GDAnyTypeRef},
		},
		runtime.GDNilTypeRef,
		true,
	)

	fn := runtime.NewGDLambdaWithType(
		fnType,
		nil,
		func(args runtime.GDLambdaArgs, stack *runtime.GDStack) (runtime.GDObject, error) {
			argsArg, err := args.Get(argsIdent)
			if err != nil {
				return nil, err
			}

			argsArray, isArray := argsArg.(*runtime.GDArray)
			if !isArray {
				return nil, runtime.InvalidCastingWrongTypeErr(runtime.NewGDArrayType(runtime.GDAnyTypeRef), argsArg.GetType())
			}

			for _, arg := range argsArray.GetObjects() {
				argVal := runtime.GDString(arg.ToString())
				runtime.Print(argVal.Escape())
			}

			if newLine {
				runtime.Print("\n")
			}

			return runtime.GDZNil, nil
		},
	)

	return runtime.NewGDSymbol(true, true, fnType, fn)
}
