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

func TestStructCases(t *testing.T) {
	RunTests(t, []Test{
		{`pub func main() {
			set s = {
				a: 1,
				b: 2,
			}
			print(s)
		}`, "{a: 1, b: 2}", ""},
		// Test struct with type prefix
		{`pub func main() {
			set s: {
				a: int,
				b: int,
			} = {a: 1, b: 1}
			print(s)
		}`, "{a: 1, b: 1}", ""},
		// Test struct with type prefix and type mismatch
		{`pub func main() {
			set s: {
				a: int,
				b: int,
			} = {a: 1, b: "1"}
			print(s)
		}`, "", "expected `{a: int, b: int}` but got `{a: int, b: string}`"},
		// Structs as part of a function return
		{`pub func main() {
			func f() => {
				a: int,
				b: int,
			} {
				return {a: 1, b: 2}
			}
			print(f())
		}`, "{a: 1, b: 2}", ""},
		// Structs as a result of a function call
		{`pub func main() {
			func f() => {
				a: int,
				b: int,
			} {
				return {a: 1, b: 2}
			}
			print(f(), f().a, f().b, f().a + f().b)
		}`, "{a: 1, b: 2}123", ""},
		// Nested structs
		{`pub func main() {
			set s = {
				a: 1,
				b: {
					c: 2,
					d: 3,
				},
			}
			print(s)
		}`, "{a: 1, b: {c: 2, d: 3}}", ""},
		// Accessing nested structs
		{`pub func main() {
			set s = {
				a: 1,
				b: {
					c: 2,
					d: 3,
				},
			}
			print(s.b.c)
		}`, "2", ""},
		// Trickier nested structs
		{`pub func main() {
			typealias t = {
				a: int,
				b: {
					c: int,
				},
			}
			func c(user: t) => int {
				return user?.b.c
			}
			func f() => t {
				return {a: 1, b: {c: c(nil)}}
			}
			set nU = f()
			set s = nU?.b.c
			print(s)
		}`, "nil", ""},
		// Access struct fields
		{`pub func main() {
			set s = {
				a: 1,
				b: 2,
			}
			print(s.a, s.b)
		}`, "12", ""},
		// Access struct lambda fields
		{`pub func main() {
			set s: {
				handler: func() => int,
			} = {
				handler: func() => int {
					return 1
				},
			}
			print(s.handler())
		}`, "1", ""},
		// Update a struct attribute
		{`pub func main() {
			set s = {
				a: 1,
				b: 2,
			}
			s.a = 3
			print(s)
		}`, "{a: 3, b: 2}", ""},
		// Update a struct attribute with a lambda
		{`pub func main() {
			set s: {
				handler: func() => int,
			} = {
				handler: func() => int {
					return 1
				},
			}
			s.handler = func() => int {
				return 2
			}
			print(s.handler())
		}`, "2", ""},
	})
}
