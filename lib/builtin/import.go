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

type GDBuiltinMap = map[string]func(stack *runtime.GDSymbolStack) (runtime.GDObject, error)

var builtins = GDBuiltinMap{
	// fmt
	"print":   print,
	"println": println,
	// types
	"typeof": typeof,
	// math
	"abs":   abs,
	"sqrt":  sqrt,
	"pow":   pow,
	"sin":   sin,
	"cos":   cos,
	"tan":   tan,
	"log":   log,
	"asin":  asin,
	"acos":  acos,
	"atan":  atan,
	"atan2": atan2,
	"exp":   exp,
	"ceil":  ceil,
	"floor": floor,
	"round": round,
}

func Import(stack *runtime.GDSymbolStack) error {
	for ident, fn := range builtins {
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
