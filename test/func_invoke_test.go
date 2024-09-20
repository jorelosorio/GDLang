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

package test

import (
	"gdlang/lib/runtime"
	"testing"
)

func TestFuncInvoke(t *testing.T) {
	tmpl := `
		$SRC
		pub func main(){print(r);}
	`
	RunTestsWithTemplate(t, tmpl, []Test{
		{"func f()=>(int,){return (1,);};set r=f()", runtime.NewGDTuple(runtime.NewGDIntNumber(1)).ToString(), ""},
		{"func f()=>(int,int){return (1,1);};set r=f()", runtime.NewGDTuple(runtime.NewGDIntNumber(1), runtime.NewGDIntNumber(1)).ToString(), ""},
		{"func f()=>(float,int){\nreturn (1.0,1)\n}\nset r=f()", runtime.NewGDTuple(runtime.NewGDFloatNumber(1), runtime.NewGDIntNumber(1)).ToString(), ""},
		{"func f()=>(,){\nreturn\n}\nset r=f()", runtime.GDZNil.ToString(), ""},
		{"func f()=>(,){\nreturn nil\n}\nset r=f()", runtime.GDZNil.ToString(), ""},
		{"func f()=>(,){\nreturn (,)\n}\nset r=f()", runtime.NewGDTuple().ToString(), ""},
		{"func f()=>(int,int){\nreturn (nil,nil)\n}\nset r=f()", runtime.NewGDTuple(runtime.GDZNil, runtime.GDZNil).ToString(), ""},
		{"func f(){\nreturn nil\n}\nset r=f()", runtime.GDZNil.ToString(), ""},
		{"func f(){\nreturn\n}\nset r=f()", runtime.GDZNil.ToString(), ""},
		{"func f(){\nreturn\nreturn\n}\nset r=f()", runtime.GDZNil.ToString(), ""},
		{"func f()=>int{\nreturn 500\nreturn 1000\n}\nset r=f()", runtime.NewGDIntNumber(500).ToString(), ""},
		{"func f()=>string{\nreturn \"hola\"\n}\nset r=f()", runtime.GDString("hola").ToString(), ""},
		{"func f()=>(string,){\nreturn (\"hola\",)\n}\nset r=f()", runtime.NewGDTuple(runtime.GDString("hola")).ToString(), ""},
		{"func f()=>(string,){\nreturn (nil,)\n}\nset r=f()", runtime.NewGDTuple(runtime.GDZNil).ToString(), ""},
		{"func f()=>(char,){\nreturn ('a',)\n}\nset r=f()", runtime.NewGDTuple(runtime.GDChar('a')).ToString(), ""},
		{"func f()=>(char,){\nreturn (nil,)\n}\nset r=f()", runtime.NewGDTuple(runtime.GDZNil).ToString(), ""},
	})
}
