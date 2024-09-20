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

func TestFunctionWithParameters(t *testing.T) {
	RunTest(t, Test{`
	func add() => int {
		return 0
	}
	pub func main() {
		add(1, 0)
	}`, "", "expected `0` but got `2`"})
}

func TestUnknownTypesWhileCreatingFunction(t *testing.T) {
	RunTest(t, Test{`
	func add(a: test) {
		return 0
	}
	pub func main() {
		add(1)
	}`, "", "object `test` was not found"})
}

func TestSingleReturnType(t *testing.T) {
	RunTest(t, Test{`
		func add() => int {
			return 0
		}
		pub func main() {
			print(add())
		}
	`, "0", ""})
}

func TestReturnTupleWithOneElement(t *testing.T) {
	RunTest(t, Test{`
		func add() => (int,) {
			return (0,)
		}
		pub func main() {
			print(add())
		}
	`, "(0,)", ""})
}

func TestMultipleReturnTypes(t *testing.T) {
	RunTest(t, Test{`
		func add() => (int, float) {
			return (0, 2.1)
		}
		pub func main() {
			print(add())
		}
	`, "(0, 2.1)", ""})
}

func TestUnknownTypesOnReturnWhileCreatingFunction(t *testing.T) {
	RunTest(t, Test{`
	func add() => test {
		return 0
	}
	pub func main() {
		add()
	}
	`, "", "object `test` was not found"})
}

func TestSendBlockStmtAsParameter(t *testing.T) {
	RunTest(t, Test{`
		func add(a: int) => int {
			return a+1
		}
		pub func main() {
			print(add(func()=>int{return 1;}()))
		}`, "2", ""})
}

func TestFuncCases(t *testing.T) {
	RunTests(t, []Test{
		{`func f(){
			func f1(){
				print(true);
			};
			f1();
		};
		pub func main(){
			f();
		}`, "true", ""}, // Function inside function
		{`func f(){
			func f1(){
				func f2(){
					print('*', true, '*');
				};
				f2();
			};
			f1();
		};
		pub func main(){
			f();
		}`, "*true*", ""},
		{`pub func w(){
			print(1);
		};
		pub func main(){
			w();
		}`, "1", ""},
		{`func a1(n:int)=>int{
			return n+1;
		};
		pub func main(){
			print(a1(a1(a1(1))));
		}`, "4", ""},
	})
}
