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

var HttpRequestType = runtime.QuickGDStructType(
	"status", runtime.GDStringType,
	"statusCode", runtime.GDIntType,
	"body", runtime.GDStringType,
)

var HttpHeaderType = runtime.QuickGDStructType(
	"key", runtime.GDStringType,
	"value", runtime.GDStringType,
)

var HttpResponseType = runtime.QuickGDStructType(
	"status", runtime.GDIntType,
	"body", runtime.GDAnyType,
	"headers", runtime.NewGDArrayType(HttpHeaderType),
)

var HttpHandlerType = runtime.NewGDLambdaType(
	runtime.GDLambdaArgTypes{},
	HttpResponseType,
	false,
)

var HttpRouteType = runtime.QuickGDStructType(
	"method", runtime.GDStringType,
	"path", runtime.GDStringType,
	"handler", HttpHandlerType,
)

func HttpPackage() (*runtime.GDPackage[*runtime.GDSymbol], error) {
	pkg := runtime.NewGDPackage[*runtime.GDSymbol](runtime.NewGDStrIdent("http"), "http", runtime.PackageModeBuiltin)
	symbols := map[string]*runtime.GDSymbol{
		// Restful
		"fetch": fetch(),
		// Routing
		"get": get(),
		"ok":  ok(),
		// Server
		"host": host(),
		// Types
		"route":    runtime.NewGDSymbol(true, true, HttpRouteType, nil),
		"response": runtime.NewGDSymbol(true, true, HttpResponseType, nil),
	}

	for ident, symbol := range symbols {
		err := pkg.AddPublic(runtime.NewGDStrIdent(ident), symbol)
		if err != nil {
			return nil, err
		}
	}

	return pkg, nil
}

func fetch() *runtime.GDSymbol {
	url := runtime.NewGDStrRefType("url")

	typ := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{Key: url, Value: runtime.GDStringType},
		},
		HttpRequestType,
		false,
	)

	lambda := runtime.NewGDLambdaWithType(
		typ,
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
				HttpRequestType,
				runtime.GDString(response.Status),
				runtime.GDInt(response.StatusCode),
				runtime.GDString(body),
			)
		},
	)

	return runtime.NewGDSymbol(true, true, typ, lambda)
}

func get() *runtime.GDSymbol {
	path := runtime.NewGDStrRefType("path")
	handler := runtime.NewGDStrRefType("handler")

	typ := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{Key: path, Value: runtime.GDStringType},
			{Key: handler, Value: HttpHandlerType},
		},
		HttpRouteType,
		false,
	)

	lambda := runtime.NewGDLambdaWithType(
		typ,
		nil,
		func(stack *runtime.GDSymbolStack, args runtime.GDLambdaArgs) (runtime.GDObject, error) {
			path := args.Get(path).ToString()
			handler := args.Get(handler)

			return runtime.QuickGDStruct(
				stack,
				HttpRouteType,
				runtime.GDString("GET"),
				runtime.GDString(path),
				handler,
			)
		},
	)

	return runtime.NewGDSymbol(true, true, typ, lambda)
}

func ok() *runtime.GDSymbol {
	body := runtime.NewGDStrRefType("body")
	headers := runtime.NewGDStrRefType("headers")

	typ := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{Key: body, Value: runtime.GDAnyType},
			{Key: headers, Value: runtime.NewGDArrayType(HttpHeaderType)},
		},
		HttpResponseType,
		false,
	)

	lambda := runtime.NewGDLambdaWithType(
		typ,
		nil,
		func(stack *runtime.GDSymbolStack, args runtime.GDLambdaArgs) (runtime.GDObject, error) {
			body := args.Get(body)
			headers := args.Get(headers)

			return runtime.QuickGDStruct(
				stack,
				HttpResponseType,
				runtime.GDInt(200),
				body,
				headers,
			)
		},
	)

	return runtime.NewGDSymbol(true, true, typ, lambda)
}

func host() *runtime.GDSymbol {
	path := runtime.NewGDStrRefType("path")
	routes := runtime.NewGDStrRefType("routes")

	typ := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{Key: path, Value: runtime.GDStringType},
			{Key: routes, Value: runtime.NewGDArrayType(HttpRouteType)},
		},
		runtime.GDNilType,
		false,
	)

	lambda := runtime.NewGDLambdaWithType(
		typ,
		nil,
		func(stack *runtime.GDSymbolStack, args runtime.GDLambdaArgs) (runtime.GDObject, error) {
			path := args.Get(path).ToString()
			routes, isArray := args.Get(routes).(*runtime.GDArray)

			if !isArray {
				return nil, runtime.InvalidCastingWrongTypeErr(runtime.NewGDArrayType(HttpRouteType), routes.GetType())
			}

			mux := http.NewServeMux()

			for _, route := range routes.Objects {
				route, isStruct := route.(*runtime.GDStruct)
				if !isStruct {
					return nil, runtime.InvalidCastingWrongTypeErr(HttpRouteType, route.GetType())
				}

				method, err := route.GetAttr(runtime.NewGDStrIdent("method"))
				if err != nil {
					return nil, err
				}

				path, err := route.GetAttr(runtime.NewGDStrIdent("path"))
				if err != nil {
					return nil, err
				}

				handler, err := route.GetAttr(runtime.NewGDStrIdent("handler"))
				if err != nil {
					return nil, err
				}

				lambdaHandler, isLambda := handler.Object.(*runtime.GDLambda)
				if !isLambda {
					return nil, runtime.InvalidCastingWrongTypeErr(HttpRouteType, handler.Object.GetType())
				}

				routePath := method.Object.ToString() + " " + path.Object.ToString()
				mux.HandleFunc(routePath, func(w http.ResponseWriter, r *http.Request) {
					obj, _ := lambdaHandler.Call(runtime.NewGDArray())
					response, isStruct := obj.(*runtime.GDStruct)
					if !isStruct {
						return
					}

					status, err := response.GetAttr(runtime.NewGDStrIdent("status"))
					if err != nil {
						return
					}

					statusCodeInt, err := runtime.ToInt(status.Object)
					if err == nil {
						w.WriteHeader(int(statusCodeInt))
					}

					// Get headers
					headers, err := response.GetAttr(runtime.NewGDStrIdent("headers"))
					if err != nil {
						return
					}

					headersArray, isHeadersArray := headers.Object.(*runtime.GDArray)
					if !isHeadersArray {
						return
					}

					for _, header := range headersArray.Objects {
						header, isHeader := header.(*runtime.GDStruct)
						if !isHeader {
							continue
						}

						key, err := header.GetAttr(runtime.NewGDStrIdent("key"))
						if err != nil {
							continue
						}

						value, err := header.GetAttr(runtime.NewGDStrIdent("value"))
						if err != nil {
							continue
						}

						w.Header().Add(key.Object.ToString(), value.Object.ToString())
					}

					body, err := response.GetAttr(runtime.NewGDStrIdent("body"))
					if err != nil {
						return
					}

					_, err = w.Write([]byte(body.Object.ToString()))
					if err != nil {
						return
					}
				})
			}

			err := http.ListenAndServe(path, mux)
			if err != nil {
				return nil, err
			}

			return runtime.GDZNil, nil
		},
	)

	return runtime.NewGDSymbol(true, true, typ, lambda)
}
