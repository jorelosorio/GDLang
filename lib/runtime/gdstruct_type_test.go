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

package runtime_test

import (
	"gdlang/lib/runtime"
	"strings"
	"testing"
)

func TestSimpleStructType(t *testing.T) {
	stack := runtime.NewGDStack()

	subStructTypeWithInt := runtime.NewGDStructType(&runtime.GDStructAttrType{bParamIdent, runtime.GDIntTypeRef})
	subStructTypeWithArray := runtime.NewGDStructType(&runtime.GDStructAttrType{bParamIdent, runtime.NewGDArrayType(runtime.GDIntTypeRef)})

	for _, test := range []struct {
		attrs    []*runtime.GDStructAttrType
		expected string
		errMsg   string
	}{
		{[]*runtime.GDStructAttrType{}, "{}", ""},
		{[]*runtime.GDStructAttrType{{aParamIdent, runtime.GDIntTypeRef}}, "{a: int}", ""},
		{[]*runtime.GDStructAttrType{{aParamIdent, runtime.GDIntTypeRef}, {bParamIdent, runtime.GDStringTypeRef}}, "{a: int, b: string}", ""},
		{[]*runtime.GDStructAttrType{{aParamIdent, runtime.GDIntTypeRef}, {bParamIdent, runtime.GDStringTypeRef}, {cParamIdent, runtime.GDBoolTypeRef}}, "{a: int, b: string, c: bool}", ""},
		// Struct with nested struct
		{[]*runtime.GDStructAttrType{{aParamIdent, subStructTypeWithInt}}, "{a: {b: int}}", ""},
		// Struct with nested struct and array
		{[]*runtime.GDStructAttrType{{aParamIdent, subStructTypeWithArray}}, "{a: {b: [int]}}", ""},
	} {
		structType := runtime.NewGDStructType(test.attrs...)
		_, err := runtime.InferType(structType, runtime.GDUntypedTypeRef, stack)
		if err != nil {
			if test.errMsg != "" {
				if !strings.Contains(err.Error(), test.errMsg) {
					t.Errorf("Expected error to contain %q but got %q", test.errMsg, err.Error())
				}
				return
			} else {
				t.Errorf("Error creating struct: %s", err.Error())
				return
			}
		}

		if len(structType) != len(test.attrs) {
			t.Errorf("Expected %d attributes but got %d", len(test.attrs), len(structType))
			return
		}

		if structType.ToString() != "" && structType.ToString() != test.expected {
			t.Errorf("Expected %s but got %s", test.expected, structType.ToString())
		}
	}
}
