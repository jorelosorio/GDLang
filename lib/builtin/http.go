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

var HttpRequestType = runtime.NewGDTypeAliasType(
	runtime.NewGDStrIdent("request"),
	runtime.QuickGDStructType(
		"status", runtime.GDStringTypeRef,
		"statusCode", runtime.GDIntTypeRef,
		"body", runtime.GDStringTypeRef,
	),
)

var HttpHeaderType = runtime.NewGDTypeAliasType(
	runtime.NewGDStrIdent("header"),
	runtime.QuickGDStructType(
		"key", runtime.GDStringTypeRef,
		"value", runtime.GDStringTypeRef,
	),
)

var HttpResponseType = runtime.NewGDTypeAliasType(
	runtime.NewGDStrIdent("response"),
	runtime.QuickGDStructType(
		"status", runtime.GDIntTypeRef,
		"body", runtime.GDAnyTypeRef,
		"headers", runtime.NewGDArrayType(HttpHeaderType),
	),
)

var HttpHandlerType = runtime.NewGDTypeAliasType(
	runtime.NewGDStrIdent("handler"),
	runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{},
		HttpResponseType,
		false,
	),
)

var HttpRouteType = runtime.NewGDTypeAliasType(
	runtime.NewGDStrIdent("route"),
	runtime.QuickGDStructType(
		"method", runtime.GDStringTypeRef,
		"path", runtime.GDStringTypeRef,
		"handler", HttpHandlerType,
	),
)

func HttpPackage() (*runtime.GDPackage[*runtime.GDSymbol], error) {
	pkg := runtime.NewGDPackage[*runtime.GDSymbol](
		runtime.NewGDStrIdent("http"),
		".",
		runtime.PackageModeBuiltin,
	)

	symbols := map[string]*runtime.GDSymbol{
		// Types
		"request":  runtime.NewGDSymbol(true, true, HttpRequestType, nil),
		"header":   runtime.NewGDSymbol(true, true, HttpHeaderType, nil),
		"response": runtime.NewGDSymbol(true, true, HttpResponseType, nil),
		"handler":  runtime.NewGDSymbol(true, true, HttpHandlerType, nil),
		"route":    runtime.NewGDSymbol(true, true, HttpRouteType, nil),
		// Restful
		"fetch": fetch(),
		// Routing
		"get": get(),
		"ok":  ok(),
		// Server
		"host": host(),
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
	url := runtime.NewGDStrIdent("url")

	typ := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{Key: url, Value: runtime.GDStringTypeRef},
		},
		HttpRequestType,
		false,
	)

	lambda := runtime.NewGDLambdaWithType(
		typ,
		nil,
		func(args runtime.GDLambdaArgs, stack *runtime.GDStack) (runtime.GDObject, error) {
			url, err := args.Get(url)
			if err != nil {
				return nil, err
			}

			response, err := http.Get(url.ToString())
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
				HttpRequestType.GDTypable.(runtime.GDStructType),
				runtime.GDString(response.Status),
				runtime.GDInt(response.StatusCode),
				runtime.GDString(body),
			)
		},
	)

	return runtime.NewGDSymbol(true, true, typ, lambda)
}

func get() *runtime.GDSymbol {
	path := runtime.NewGDStrIdent("path")
	handler := runtime.NewGDStrIdent("handler")

	typ := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{Key: path, Value: runtime.GDStringTypeRef},
			{Key: handler, Value: HttpHandlerType},
		},
		HttpRouteType,
		false,
	)

	lambda := runtime.NewGDLambdaWithType(
		typ,
		nil,
		func(args runtime.GDLambdaArgs, stack *runtime.GDStack) (runtime.GDObject, error) {
			path, err := args.Get(path)
			if err != nil {
				return nil, err
			}

			handler, err := args.Get(handler)
			if err != nil {
				return nil, err
			}

			return runtime.QuickGDStruct(
				stack,
				HttpRouteType.GDTypable.(runtime.GDStructType),
				runtime.GDString("GET"),
				runtime.GDString(path.ToString()),
				handler,
			)
		},
	)

	return runtime.NewGDSymbol(true, true, typ, lambda)
}

func ok() *runtime.GDSymbol {
	body := runtime.NewGDStrIdent("body")
	headers := runtime.NewGDStrIdent("headers")

	typ := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{Key: body, Value: runtime.GDAnyTypeRef},
			{Key: headers, Value: runtime.NewGDArrayType(HttpHeaderType)},
		},
		HttpResponseType,
		false,
	)

	lambda := runtime.NewGDLambdaWithType(
		typ,
		nil,
		func(args runtime.GDLambdaArgs, stack *runtime.GDStack) (runtime.GDObject, error) {
			body, err := args.Get(body)
			if err != nil {
				return nil, err
			}

			headers, err := args.Get(headers)
			if err != nil {
				return nil, err
			}

			return runtime.QuickGDStruct(
				stack,
				HttpResponseType.GDTypable.(runtime.GDStructType),
				runtime.GDInt(200),
				body,
				headers,
			)
		},
	)

	return runtime.NewGDSymbol(true, true, typ, lambda)
}

func host() *runtime.GDSymbol {
	path := runtime.NewGDStrIdent("path")
	routes := runtime.NewGDStrIdent("routes")

	typ := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{Key: path, Value: runtime.GDStringTypeRef},
			{Key: routes, Value: runtime.NewGDArrayType(runtime.NewGDStrTypeRefType("route"))},
		},
		runtime.GDNilTypeRef,
		false,
	)

	lambda := runtime.NewGDLambdaWithType(
		typ,
		nil,
		func(args runtime.GDLambdaArgs, stack *runtime.GDStack) (runtime.GDObject, error) {
			path, err := args.Get(path)
			if err != nil {
				return nil, err
			}

			routesArg, err := args.Get(routes)
			if err != nil {
				return nil, err
			}

			routes, isArray := routesArg.(*runtime.GDArray)
			if !isArray {
				return nil, runtime.InvalidCastingWrongTypeErr(runtime.NewGDArrayType(HttpRouteType), routes.GetType())
			}

			mux := http.NewServeMux()

			for _, route := range routes.GetObjects() {
				route, isStruct := route.(*runtime.GDStruct)
				if !isStruct {
					return nil, runtime.InvalidCastingWrongTypeErr(HttpRouteType, route.GetType())
				}

				methodSymbol, err := route.GetAttr(runtime.NewGDStrIdent("method"))
				if err != nil {
					return nil, err
				}

				pathSymbol, err := route.GetAttr(runtime.NewGDStrIdent("path"))
				if err != nil {
					return nil, err
				}

				handlerSymbol, err := route.GetAttr(runtime.NewGDStrIdent("handler"))
				if err != nil {
					return nil, err
				}

				methodObj, pathObj, handlerObj := methodSymbol.Value.(runtime.GDObject), pathSymbol.Value.(runtime.GDObject), handlerSymbol.Value.(runtime.GDObject)

				lambdaHandler, isLambda := handlerSymbol.Value.(*runtime.GDLambda)
				if !isLambda {
					return nil, runtime.InvalidCastingWrongTypeErr(HttpRouteType, handlerObj.GetType())
				}

				routePath := methodObj.ToString() + " " + pathObj.ToString()
				mux.HandleFunc(routePath, func(w http.ResponseWriter, r *http.Request) {
					obj, _ := lambdaHandler.Call(runtime.NewGDArray(stack))
					response, isStruct := obj.(*runtime.GDStruct)
					if !isStruct {
						return
					}

					status, err := response.GetAttr(runtime.NewGDStrIdent("status"))
					if err != nil {
						return
					}

					statusCodeInt, err := runtime.ToInt(status.Value)
					if err == nil {
						w.WriteHeader(int(statusCodeInt))
					}

					// Get headers
					headers, err := response.GetAttr(runtime.NewGDStrIdent("headers"))
					if err != nil {
						return
					}

					headersArray, isHeadersArray := headers.Value.(*runtime.GDArray)
					if !isHeadersArray {
						return
					}

					for _, header := range headersArray.GetObjects() {
						header, isHeader := header.(*runtime.GDStruct)
						if !isHeader {
							continue
						}

						keySymbol, err := header.GetAttr(runtime.NewGDStrIdent("key"))
						if err != nil {
							continue
						}

						valueSymbol, err := header.GetAttr(runtime.NewGDStrIdent("value"))
						if err != nil {
							continue
						}

						keyObj, valueObj := keySymbol.Value.(runtime.GDObject), valueSymbol.Value.(runtime.GDObject)

						w.Header().Add(keyObj.ToString(), valueObj.ToString())
					}

					bodySymbol, err := response.GetAttr(runtime.NewGDStrIdent("body"))
					if err != nil {
						return
					}

					bodyObj := bodySymbol.Value.(runtime.GDObject)

					_, err = w.Write([]byte(bodyObj.ToString()))
					if err != nil {
						return
					}
				})
			}

			err = http.ListenAndServe(path.ToString(), mux)
			if err != nil {
				return nil, err
			}

			return runtime.GDZNil, nil
		},
	)

	return runtime.NewGDSymbol(true, true, typ, lambda)
}
