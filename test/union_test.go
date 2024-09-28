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

import "testing"

func TestCompositeInBuiltinFunction(t *testing.T) {
	RunTests(t, []Test{
		{`use math {abs}
		pub func main() {
		set a = abs(1)
		print(typeof(a), a)
		}`, "int1", ""},
		{`use math {abs}
		pub func main() {
			set a: int = abs(-1.0)
			print(typeof(a), a)
		}`, "", "expected `int` but got `(int | float | complex)`"},
		{`use math {abs}
		pub func main() {
			set a: any = abs(-1.0)
			print(typeof(a), a)
		}`, "float1", ""},
	})
}

func TestCompositeOnSetObjects(t *testing.T) {
	RunTests(t, []Test{
		// Composite syntax
		{`pub func main() {
			set a: (int | float) = 1
			print(typeof(a), a)
		}`, "int1", ""},
		{`pub func main() {
			set a: (int | float) = 1.0
			print(typeof(a), a)
		}`, "float1", ""},
	})
}

func TestCompositeOnFunc(t *testing.T) {
	RunTests(t, []Test{
		{`pub func main() {
			func add(a: (int | float), b: (int | float)) => (int | float) {
				return a + b
			}
			set a = add(1, 2)
			print(typeof(a), a)
		}`, "int3", ""},
		{`pub func main() {
			func add(a: float, b: (int | float)) => (int | float) {
				return a + b
			}
			set a = add(1.0, 2)
			print(typeof(a), a)
		}`, "float3", ""},
		{`pub func main() {
			func add(a: (int | float), b: (int | float)) => int {
				return a + b
			}
			set a = add(1, 2)
			print(typeof(a), a)
		}`, "", "expected `int` but got `(int | float)`"},
		{`pub func main() {
			func add(a: (int | float), b: (int | float)) => int {
				return a + b
			}
			set a = add(1.0, 2.0)
			print(typeof(a), a)
		}`, "", "expected `int` but got `(int | float)`"},
	})
}

func TestCompositeOnTypes(t *testing.T) {
	RunTests(t, []Test{
		{`pub func main() {
			set a: [(int | float)] = [1, 1]
			print(typeof(a), a)
		}`, "[(int | float)][1, 1]", ""},
		{`pub func main() {
			set a: ((int | float),) = (1,)
			print(typeof(a), a)
		}`, "(int,)(1,)", ""},
		{`pub func main() {
			set a: ((int | float),) = (1.0,)
			print(typeof(a), a)
		}`, "(float,)(1,)", ""},
		{`pub func main() {
			set a: ((int | float), (string | char)) = (1, 'a')
			print(typeof(a), a)
		}`, "(int, char)(1, 'a')", ""},
		{`pub func main() {
			set a: ((int | float), (string | char)) = (1, "hola")
			print(typeof(a), a)
		}`, "(int, string)(1, \"hola\")", ""},
		{`pub func main() {
			set a: [(int | float)] = [1, 1.0]
			print(typeof(a), a)
		}`, "[(int | float)][1, 1]", ""},
		{`pub func main() {
        	set a: [(int | float)] = [1, 1, "t"]
            print(typeof(a), a)
        }`, "", "expected `[(int | float)]` but got `[(int | string)]`"},
	})
}
