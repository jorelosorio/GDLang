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

func TestArrays(t *testing.T) {
	RunTestsWithMainTemplate(t, []Test{
		{`print([]);`, "[]", ""},
		{`print([1]);`, "[1]", ""},
		{`print([1, 2]);`, "[1, 2]", ""},
		{`print([1, 2, 3]);`, "[1, 2, 3]", ""},
		{`print(['1', "2", 3.1]);`, "['1', \"2\", 3.1]", ""},
		{`set a = [
			func ()=>int{
				return 1
			}(),
		]
		print(a);`, "[1]", ""},
		{`print([
			func () => (any, any, any) {
				return ([1, 2, 3],1,1)
			}(),
		]...)`, "([1, 2, 3], 1, 1)", ""},
		{`set table = [0]
		print(table[0])`, "0", ""},
		{`set ra = 1
			set table = [ra]
			print(table[ra-1])`, "1", ""},
		{`set j = [1][0]
			print(j)`, "1", ""},
		{`print([1, 2, 3][2])`, "3", ""},
		{`func a() => [int] {
			return [1, 2, 3]
		}
		set j = a()
		print(j[2])`, "3", ""},
		{`func a() => [int] {
			return [1, 2, 3]
		}
		print(a()[1])`, "2", ""},
		{`func a(a: [[int]]) => [int] {
			return a[0]
		}
		print(a([[1, 2, 3]])[0])`, "1", ""},
		// Set array value
		{`set a = [1, 2, 3]
			a[0] = 4
			print(a[0])`, "4", ""},
		{`set a = [1, 2, 3]
			a[1] = 4
			print(a[1])`, "4", ""},
		{`set a = [1, 2, 3]
			a[2] = 4
			print(a[2])`, "4", ""},
		{`func a(array: [int]) => [int] {
				array[0] = 4
				return array
			}
			print(a([1, 2, 3])[0])`, "4", ""},
		{`set pos = 0
			func a(array: [int], val: int) => [int] {
				array[pos] = val
				return array
			}
			set res = a([1, 2, 3], -1)
			print(res[pos])
			pos = pos + 1
			print(a(res, -2)[pos])
			print(res)`, "-1-2[-1, -2, 3]", ""},
		// Array add
		{`set a = [1, 2, 3]
			set b = a << 4
			print(a, b)`, "[1, 2, 3, 4]4", ""},
		{`set a = [], b = 2, c = [0] << (a << b)
			print(a,b,c) `, "[2]22", ""},
		{`set a = [1, 2, 3] << 4
			print(a)`, "4", ""},
		// Remove array value
		{`set a = [1, 2, 3]
			print(a >> 0, a)`, "1[2, 3]", ""},
	})
}
