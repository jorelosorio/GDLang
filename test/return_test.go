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

func TestReturnTypes(t *testing.T) {
	RunTests(t, []Test{
		{"func f(){return 1;};pub func main(){print(f());}", "", "expected `nil` but got `int`"},
		{"func f()=>int{return 1;};pub func main(){print(f());}", "1", ""},
		{`func f()=>int{return "OK";};pub func main(){print(f());}`, "", "expected `int` but got `string`"},
		{`func f()=>int{return 1.0;};pub func main(){print(f());}`, "", "expected `int` but got `float`"},
		{`func f()=>(int,){return (1,);};pub func main(){print(f());};`, "(1,)", ""},
		{`func f()=>(int,){return nil;};pub func main(){print(f());};`, "nil", ""},
		{`func f()=>(int,int){return (1,2);};pub func main(){print(f());};`, "(1, 2)", ""},
		{`func f()=>(int,){return (1,2,3);};pub func main(){print(f());};`, "", "expected `(int,)` but got `(int, int, int)`"},
		{`func f()=>[int]{return [];};pub func main(){print(f());};`, "[]", ""},
		{`
		set c = 1
		func f() {
		    func c() => int {
		        return 1
		    }
		    return c
		}
		pub func main(){print(f());}
		`, "", "expected `nil` but got `func() => int`"},
		{`
		set a = 1
		set b = a
		pub func main(){
			a = b
			print(b)
			func c()=>([string],(int,)){
				print(a)
				return (["hello", "world"], (1,))
			}
			func b() {
				func c() => int {
					return 1
				}
				return c
			}
			print(b())
		}`, "", "expected `nil` but got `func() => int`"},
		{`func f()=>[int] {
			return []
		}
		pub func main() {
			print(typeof(f()));
		}`, "[int]", ""},
	})
}
