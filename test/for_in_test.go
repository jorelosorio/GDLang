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

func TestForInCases(t *testing.T) {
	RunTests(t, []Test{
		{`pub func main() {
			for set a: int in [1, 2, 3] {
				print(a,"-")
			}
		}`, "1-2-3-", ""},
		{`pub func main() {
			for set a in [1, 2, 3] {
				print(a,"-")
			}
		}`, "1-2-3-", ""},
		{`pub func main() {
			for set a: char, b in "hola" {
				print(b,a)
			}
		}`, "h0o1l2a3", ""},
		{`pub func main() {
			for set t: int = 0, b = 1 in (1, 2, 3) {
				print(t+b)
			}
		}`, "135", ""},
		{`pub func main() {
			for set const b: int in [1, 2, 3] {
				print(b)
			}
		}`, "", "can't set a constant object"},
		{`pub func main() {
			set a = [1, 2, 3]
			for set b: int in a {
				print(b)
			}
		}`, "123", ""},
		{`pub func main() {
			func w() => [int] {
				return [1, 2, 3]
			}
			for set b: int in w() {
				print(b)
			}
		}`, "123", ""},
		// For with break
		{`pub func main() {
			for set a:int in [1, 2, 3, 4, 5, 6] {
				if a == 3 {
					break
				}
				print(a)
			}
		}`, "12", ""},
		{`pub func main() {
			for set a:int in [1, 2, 3, 4, 5, 6] {
				for set b: int in [1, 2, 3, 4, 5, 6] {
					if b == 3 && a == 3 {
						break
					}
					print(b)
				}
				print(a)
			}
		}`, "12345611234562123123456412345651234566", ""},
		{`pub func main() {
			for set a: int = 0 if a < 10 {
				if a == 5 {
					break
				}
				print(a)
				a += 1
			}
		}`, "01234", ""},
	})
}
