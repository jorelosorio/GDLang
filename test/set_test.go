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

func TestSetCases(t *testing.T) {
	RunTests(t, []Test{
		{`set a = 1
		pub func main() {
			print(a)
		}`, "1", ""},
		{`set const a = 1;pub func main(){print(a);}`, "1", ""},
		{`set const a = 1;pub func main(){a=2;print(a);}`, "", "can't set a constant object"},
		{`set a = [];pub func main(){print(a,typeof(a));}`, "[][untyped]", ""},
		{`set a:[int] = [];pub func main(){print(a,typeof(a));}`, "[][untyped]", ""},
		{`set a:[int] = [0];pub func main(){print(a,typeof(a));}`, "[0][int]", ""},
		{`set a = []
		pub func main(){
			a = [1, 2, 3]
			print(a,typeof(a))
		}`, "[1, 2, 3][int]", ""},
		{`set a = [nil];pub func main(){print(a);}`, "[nil]", ""},
		{`set a = nil;pub func main(){print(a);}`, "nil", ""},
		{`set a
		pub func main() {
			print(a)
		}`, "nil", ""},
		{`set a = nil
		pub func main(){
			a = 1
			print(typeof(a))
		}`, "int", ""},
		{`set a = 1
		set b = 2
		pub func main() {
			print(a)
			print(b)
		}`, "12", ""},
		{`set a: int = 1;pub func main(){print(a);}`, "1", ""},
		{`set a: int = 1
		set b: int = 2
		pub func main() {
			print(a)
			print(b)
		}
		`, "12", ""},
		{`set a: float = "test";pub func main(){print(a);}`, "", "expected `float` but got `string`"},
		{`set a = 1
		pub func main() {
			a = "test"
		}`, "", "expected `int` but got `string`"},
		{`set a = 1
		set a = 2;pub func main(){}`, "", "object `a` was already created"},
		{`
		set a: float = 1.9
		set b: float = 2.9
		pub func main() {
			a = b
			b = 3.0
			print(typeof(a), typeof(b))
		}`, "floatfloat", ""},
		{`set a = nil
		set b = a
		pub func main(){
			print(a)
			print(b)
		}
		`, "nilnil", ""},
		{`set a: any = 1
		set b: any = "test"
		pub func main(){
			print(a)
			print(b)
			print(typeof(a), typeof(b))
		}`, "1testintstring", ""},
		{`pub func main(){abcde = 1;}`, "", "object `abcde` was not found"},
		{`set abcde = 1;pub func main(){abcde = true;}`, "", "expected `int` but got `bool`"},
		{`set abcde:string = 1;pub func main(){print(abcde);}`, "", "expected `string` but got `int`"},
		{`set abcde:nil = 1;pub func main(){print(abcde);}`, "", "assigning `nil` as a type is not allowed"},
		{`set a = {};pub func main(){print(a);}`, "{}", ""},
		{`set a = func(){return 1;}();pub func main(){print(a);}`, "", "expected `nil` but got `int`"},
		{`set a = func(){return nil;}();pub func main(){print(a);}`, "nil", ""},
	})
}

func TestDifferentTypeAssigments(t *testing.T) {
	tmpl := `set a = $SRC
	pub func main() {
		print(typeof(a))
	}`
	RunTestsWithTemplate(t, tmpl, []Test{
		{`1`, runtime.GDIntTypeRef.ToString(), ""},
		{`"test"`, runtime.GDStringTypeRef.ToString(), ""},
		{`true`, runtime.GDBoolTypeRef.ToString(), ""},
		{`false`, runtime.GDBoolTypeRef.ToString(), ""},
		{`nil`, runtime.GDNilTypeRef.ToString(), ""},
		{`1.1`, runtime.GDFloatTypeRef.ToString(), ""},
		{`0b100`, runtime.GDIntTypeRef.ToString(), ""},
		{"`abc`", runtime.GDStringTypeRef.ToString(), ""},
		{`'2'`, runtime.GDCharTypeRef.ToString(), ""},
		{`1_000`, runtime.GDIntTypeRef.ToString(), ""},
		{"0", runtime.GDIntTypeRef.ToString(), ""},
		{"1", runtime.GDIntTypeRef.ToString(), ""},
		{"01234567", runtime.GDIntTypeRef.ToString(), ""},
		{"0xcafebabe", runtime.GDIntTypeRef.ToString(), ""},
		{"0.0", runtime.GDFloatTypeRef.ToString(), ""},
		{"3.14159265", runtime.GDFloatTypeRef.ToString(), ""},
		{"1e0", runtime.GDFloatTypeRef.ToString(), ""},
		{"1e+100", runtime.GDFloatTypeRef.ToString(), ""},
		{"1e-100", runtime.GDFloatTypeRef.ToString(), ""},
		{"2.71828e-1000", runtime.GDFloatTypeRef.ToString(), ""},
		{"0i", runtime.GDComplexTypeRef.ToString(), ""},
		{"1i", runtime.GDComplexTypeRef.ToString(), ""},
		{"012345678901234567889i", runtime.GDComplexTypeRef.ToString(), ""},
		{"123456789012345678890i", runtime.GDComplexTypeRef.ToString(), ""},
		{"0.0i", runtime.GDComplexTypeRef.ToString(), ""},
		{"3.14159265i", runtime.GDComplexTypeRef.ToString(), ""},
		{"1e0i", runtime.GDComplexTypeRef.ToString(), ""},
		{"1e+100i", runtime.GDComplexTypeRef.ToString(), ""},
		{"1e-100i", runtime.GDComplexTypeRef.ToString(), ""},
		{"2.71828e-1000i", runtime.GDComplexTypeRef.ToString(), ""},
		{"'a'", runtime.GDCharTypeRef.ToString(), ""},
		{"`foobar`", runtime.GDStringTypeRef.ToString(), ""},
		{"`" + `foo
		                    bar` +
			"`",
			runtime.GDStringTypeRef.ToString(), ""},
		{"`\r`", runtime.GDStringTypeRef.ToString(), ""},
		{"`foo\r\nbar`", runtime.GDStringTypeRef.ToString(), ""},
		{"(1, 2, 3)", runtime.NewGDTupleType(runtime.GDIntTypeRef, runtime.GDIntTypeRef, runtime.GDIntTypeRef).ToString(), ""},
		{"((1,), 2, 3)", runtime.NewGDTupleType(runtime.NewGDTupleType(runtime.GDIntTypeRef), runtime.GDIntTypeRef, runtime.GDIntTypeRef).ToString(), ""},
		{`func()=>int{
			return 1;
		}();`, runtime.GDIntTypeRef.ToString(), ""},
		{"[0]", runtime.NewGDArrayType(runtime.GDIntTypeRef).ToString(), ""},
		{"[0,'c']", runtime.NewGDArrayType(runtime.NewGDUnionType(runtime.GDIntTypeRef, runtime.GDCharTypeRef)).ToString(), ""},
		{"[[1, 2, 3]]", runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDIntTypeRef)).ToString(), ""},
		{"[[1, 2, 3], 1]", runtime.NewGDArrayType(runtime.NewGDUnionType(runtime.NewGDArrayType(runtime.GDIntTypeRef), runtime.GDIntTypeRef)).ToString(), ""},
	})
}
