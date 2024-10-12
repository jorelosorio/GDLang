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

func NewGDStringIdentType(ident string) runtime.GDIdent {
	return runtime.NewGDStrIdent(ident)
}

var (
	aParamIdent = NewGDStringIdentType("a")
	bParamIdent = NewGDStringIdentType("b")
	cParamIdent = NewGDStringIdentType("c")
	attr1Ident  = NewGDStringIdentType("attr1")

	// Structs
	structWithAttrAAsInt    = runtime.NewGDStructType(&runtime.GDStructAttrType{aParamIdent, runtime.GDIntTypeRef})
	structWithAttrAAsString = runtime.NewGDStructType(&runtime.GDStructAttrType{aParamIdent, runtime.GDStringTypeRef})
	structWithAttrsAStrBInt = runtime.NewGDStructType(&runtime.GDStructAttrType{aParamIdent, runtime.GDStringTypeRef}, &runtime.GDStructAttrType{bParamIdent, runtime.GDIntTypeRef})
	structWithAttrsBIntAStr = runtime.NewGDStructType(&runtime.GDStructAttrType{bParamIdent, runtime.GDIntTypeRef}, &runtime.GDStructAttrType{aParamIdent, runtime.GDStringTypeRef})
)
