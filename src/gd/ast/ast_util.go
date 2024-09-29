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

package ast

import "gdlang/lib/runtime"

func buildFuncType(args []Node, variadic bool, returnType runtime.GDTypable) *runtime.GDLambdaType {
	funcArgTypes := make(runtime.GDLambdaArgTypes, len(args))
	for index, arg := range args {
		typeIdent := arg.(*NodeIdentWithType)
		ident := runtime.NewGDStrIdent(typeIdent.Ident.Lit)
		funcArgTypes[index] = runtime.GDLambdaArgType{Key: ident, Value: typeIdent.Type}
	}

	return runtime.NewGDLambdaType(funcArgTypes, returnType, variadic)
}
