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

import "gdlang/lib/runtime"

func Istr(ident string) runtime.GDIdentType {
	return runtime.GDStringIdentType(ident)
}

var (
	aParamIdent = Istr("a")
	bParamIdent = Istr("b")
	cParamIdent = Istr("c")
	attr1Ident  = Istr("attr1")

	// Structs
	structWithAttrAAsInt    = runtime.NewGDStructType(runtime.GDStructAttrType{aParamIdent, runtime.GDIntType})
	structWithAttrAAsString = runtime.NewGDStructType(runtime.GDStructAttrType{aParamIdent, runtime.GDStringType})
	structWithAttrsAStrBInt = runtime.NewGDStructType(runtime.GDStructAttrType{aParamIdent, runtime.GDStringType}, runtime.GDStructAttrType{bParamIdent, runtime.GDIntType})
	structWithAttrsBIntAStr = runtime.NewGDStructType(runtime.GDStructAttrType{bParamIdent, runtime.GDIntType}, runtime.GDStructAttrType{aParamIdent, runtime.GDStringType})
)
