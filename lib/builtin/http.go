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
	"io"
	"net/http"
)

var HTTPResponseType = runtime.QuickGDStructType(
	"status", runtime.GDStringType,
	"statusCode", runtime.GDIntType,
	"body", runtime.GDStringType,
)

func HttpPackage() (*runtime.GDPackage[*runtime.GDSymbol], error) {
	pkg := runtime.NewGDPackage[*runtime.GDSymbol](runtime.NewGDStringIdent("http"), "http", runtime.PackageModeBuiltin)
	symbols := map[string]*runtime.GDSymbol{
		"get": get(),
	}

	for ident, symbol := range symbols {
		err := pkg.AddPublic(runtime.NewGDStringIdent(ident), symbol)
		if err != nil {
			return nil, err
		}
	}

	return pkg, nil
}

func get() *runtime.GDSymbol {
	url := runtime.NewGDIdentRefType(runtime.NewGDStringIdent("url"))

	funcType := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{Key: url, Value: runtime.GDStringType},
		},
		HTTPResponseType,
		false,
	)

	getFunc := runtime.NewGDLambdaWithType(
		funcType,
		nil,
		func(stack *runtime.GDSymbolStack, args runtime.GDLambdaArgs) (runtime.GDObject, error) {
			url := args.Get(url).ToString()
			response, err := http.Get(url)
			if err != nil {
				return nil, err
			}
			defer response.Body.Close()

			body, err := io.ReadAll(response.Body)
			if err != nil {
				return nil, err
			}

			return runtime.QuickGDStruct(
				stack,
				HTTPResponseType,
				runtime.GDString(response.Status),
				runtime.GDInt(response.StatusCode),
				runtime.GDString(body),
			)
		},
	)

	return runtime.NewGDSymbol(true, true, funcType, getFunc)
}
