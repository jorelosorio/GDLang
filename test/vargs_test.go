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
	"testing"
)

func TestVariadicArgs(t *testing.T) {
	RunTests(t, []Test{
		{`func f(a:any,...){print(a);};pub func main(){f(1,2,3);}`, "[1, 2, 3]", ""},
		{`func f(a:any){print(a);};pub func main(){f(1,2,3);}`, "", "missing number of arguments: expected `1` but got `3`"},
		{`func f(a:any,...){print(a...);};pub func main(){f(1,2,3,"ok");}`, "123ok", ""},
		{`func f(a:any,...){print(a);};func ff(b:any,...){f(b);};pub func main(){ff(1,2,3);}`, "[[1, 2, 3]]", ""},
		{`func f(a:any,...){print(a...);};func ff(b:any,...){f(b);};pub func main(){ff(1,2,3);}`, "[1, 2, 3]", ""},
		{`func f(a:any,...){print(a...);};func ff(b:any,...){f(b...);};pub func main(){ff(1,2,3);}`, "123", ""},
		{`func f(a:any,...){print(a...);};func ff(b:any,...){f(b...);};pub func main(){ff((1,2),3);}`, "(1, 2)3", ""},
		{`func f(a:any,...){print(a);};func ff(b:any,...){f(b...);};pub func main(){ff((1,2),3);}`, "[(1, 2), 3]", ""},
		{`func f(a:any,...){print(a);};func ff(b:any,...){f(b);};pub func main(){ff((1,2),3);}`, "[[(1, 2), 3]]", ""},
		{`func f(a:int,...){print(a...);};func ff(b:any,...){f(b...);};pub func main(){ff(1);}`, "", "invalid argument type for `a`: expected `int...` but got `any...`"},
		{`func f(a:any,...){print(a...);};func ff(b:int,...){f(b...);};pub func main(){ff(1);}`, "1", ""},
		{`func f(a:int,...){print(a...);};func ff(b:any,...){f(b...);};pub func main(){ff("string");}`, "", "invalid argument type for `a`: expected `int...` but got `any...`"},
		{`func f(a:int,...){print(a...);};func ff(b:any,...){f(b);};pub func main(){ff("string");}`, "", "expected `int...` but got `[any]`"},
		{`func f(a:any,...){print(a...);};func ff(b:string,...){f(b...);};pub func main(){ff("string");}`, `string`, ""},
		{`func f(a:string,b:any,...){print(a,b...);};func ff(b:string,...){f(b...);};pub func main(){ff();}`, ``, "invalid argument type for `a`: expected `string` but got `string...`"},
		{`func f(a:string,b:any,...){print(a,b...);};func ff(b:string,...){f(b...);};pub func main(){ff("hi","hello","hola");}`, "", "invalid argument type for `a`: expected `string` but got `string...`"},
		{`func f(c:string,...){print(c...);};func ff(a:float,b:string,...){f(b...);};pub func main(){ff(20.0);}`, ``, ""},
		{`func f(a:string,...){print(a...);};func ff(a:float,b:string,...){print(a);f(b...);};pub func main(){ff(20.1);}`, `20.1`, ""},
		{`func f(a:string,...){print(a...);};func ff(a:float,b:string,...){print(a);f(b...);};pub func main(){ff(20.1,"hola");}`, `20.1hola`, ""},
		{`func f(a:string,...)=>any{return a;};pub func main(){print(f("hello","world")...);}`, ``, "ellipsis expression can only be used in tuples or arrays"},
		{`func f(a:string,...)=>[any]{return a;};pub func main(){print(f("hello","world")...);}`, `helloworld`, ""},
	})
}
