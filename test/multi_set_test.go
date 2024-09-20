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

func TestMultiSetSets(t *testing.T) {
	RunTests(t, []Test{
		{`pub func main() {
			set x = 1, y = 2
			print(x, y)
		}`, "12", ""},
		{`pub func main() {
			set x = 1, const y = 2
			print(x, y)
		}`, "12", ""},
		{`pub func main() {
			set x = 1, const y = 2
			y = 2
			print(y)
		}`, "", "can't set a constant object"},
		{`pub func main() {
			set x = 1, y = x
			print(x, y)
		}`, "11", ""},
		{`pub func main() {
			set x = 0, y = x + 1
			print(x, y)
		}`, "01", ""},
		{`pub func main() {
			set b = func()=>int{
				set x: int = 1, y = x
				return y+1
			}(), c = b
			print(b, c+1)
		}`, "23", ""},
	})
}
