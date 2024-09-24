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

func TestSingleIfCondition(t *testing.T) {
	RunTests(t, []Test{
		{`
		pub func main() {
			if true {
				print("true")
			}
		}`, "true", ""},
		{`
		pub func main() {
			if false {
				print("true")
			}
		}`, "", ""},
		{`
		pub func main() {
			if true {
				print("true")
			} else {
				print("false")
			}
		}`, "true", ""},
		{`
		pub func main() {
			if false {
				print("true")
			} else {
				print("false")
			}
		}`, "false", ""},
		{`
		pub func main() {
			if true {
				if true {
					print("true")
				}
			} 
		}`, "true", ""},
		{`
		pub func main() {
			if true {
				if false {
					print("true")
				}
			} 
		}`, "", ""},
		{`
		pub func main() {
			if false {
				if true {
					print("true")
				}
			} 
		}`, "", ""},
		{`
		pub func main() {
			if false {
				if false {
					print("true")
				}
			} 
		}`, "", ""},
		// Double conditions
		{`
		pub func main() {
			if true, true {
				print("true")
			}
		}`, "true", ""},
		{`
		pub func main() {
			if true, false {
				print("true")
			}
		}`, "true", ""},
		{`
		pub func main() {
			if false, true {
				print("true")
			}
		}`, "true", ""},
		{`
		pub func main() {
			if false, false {
				print("true")
			} else {
				print("false")
			}
		}`, "false", ""},
		// Complex chain of conditions
		{`
		pub func main() {
			if true && true {
				if true || false {
					if true, true {
						print("true")
					}
				}
			}
		}`, "true", ""},
		{`
		pub func main() {
			if true && true {
				if true || false {
					if true, false {
						print("true")
					}
				}
			}
		}`, "true", ""},
		{`
		pub func main() {
			if true && true {
				if true || false {
					if false, true {
						print("true")
					}
				}
			}
		}`, "true", ""},
		{`
		pub func main() {
			if false, false {
				print("true")
			}
		}`, "", ""},
		{`
		pub func main() {
			set x = 10
			set y = 30
			if (x == 5 && y == 20) || y == 30 {
				print("another complex condition met")
			}
		}`, "another complex condition met", ""},
		{`
		pub func main() {
			set a = 2*2*2
			if a == 2*2*2, a == 0 {
				print("true in the first condition")
			}
		}`, "true in the first condition", ""},
		{`
		pub func main() {
			set a = 1i
			if a == 1i {
				print("true")
			}
		}`, "true", ""},
		{`
		pub func main() {
			set a = "hello"
			if a == "hello" {
				print("true")
			}
		}`, "true", ""},
		{`
		pub func main() {
			set a = 1.0
			if a == 1.0 {
				print("true")
			}
		}`, "true", ""},
		// If with return
		{`
		pub func main() {
			func test() => int {
				if true {
					return 10
				}
				return 20
			}
			print(test())
		}`, "10", ""},
		{`
		pub func main() {
			func test() => int {
				set x = func()=>bool{
					return true;
				}()
				if x {
					return 10
				}
				return 20
			}
			print(test())
		}`, "10", ""},
		// Short if expressions
		{`pub func main() {
			set x = 10 == 10 ? true : false
			print(x)
		}`, "true", ""},
		{`pub func main() {
			func a() => int {
				return true ? 10 : 0
			}
			print(a())
		}`, "10", ""},
		{`pub func main() {
			print(true ? 10 : 20)
		}`, "10", ""},
		{`pub func main() {
			set x = 20 * 20 == 20 * 20 ? 10-11 : 10+11
			print(x)
		}`, "-1", ""},
		{`pub func main() {
			if 1 < 3 ? true : false {
				print("true")
			}
		}`, "true", ""},
	})
}
