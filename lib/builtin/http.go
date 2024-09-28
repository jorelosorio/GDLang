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

func HttpPackage() *runtime.GDPackage[*runtime.GDSymbol] {
	pkg := runtime.NewGDPackage[*runtime.GDSymbol](runtime.NewGDStringIdent("http"), "http", runtime.PackageModeBuiltin)

	err := pkg.AddPublic(runtime.NewGDStringIdent("get"), get())
	if err != nil {
		panic(err)
	}

	return pkg
}

func get() *runtime.GDSymbol {
	url := runtime.NewGDIdentRefType(runtime.NewGDStringIdent("url"))

	funcType := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{Key: url, Value: runtime.GDStringType},
		},
		runtime.NewGDTupleType(
			runtime.GDStringType,
			runtime.GDIntType,
			runtime.GDStringType,
		),
		false,
	)

	getFunc := runtime.NewGDLambdaWithType(
		funcType,
		nil,
		func(_ *runtime.GDSymbolStack, args runtime.GDLambdaArgs) (runtime.GDObject, error) {
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

			return runtime.NewGDTuple(
				runtime.GDString(response.Status),
				runtime.GDInt(response.StatusCode),
				runtime.GDString(body),
			), nil
		},
	)

	return runtime.NewGDSymbol(true, true, funcType, getFunc)
}
