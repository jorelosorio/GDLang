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

func TestStrings(t *testing.T) {
	RunTests(t, []Test{
		{`pub func main() {
			print("Hello, World!")
		}`, "Hello, World!", ""},
		{`pub func main() {
			set x = "hello"
			set y = "world"
			print(x + ", " + y + "!")
		}`, "hello, world!", ""},
		{`pub func main() {
			set x = "hello"[0]
			print(x)
		}`, "h", ""},
		{`pub func main() {
			set x = "hello"
			x[0] = "H"
			print(x)
		}`, "", "invalid collectable type: `string`"},
	})
}
