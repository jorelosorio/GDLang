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

func TestTupleCases(t *testing.T) {
	RunTests(t, []Test{
		{`pub func main(){set a = (1,);print(a);}`, "(1,)", ""},
		{`pub func main(){set a = (1,2,3);print(a);}`, "(1, 2, 3)", ""},
		{`set t: (int,) = (1,);pub func main(){print(t);}`, "(1,)", ""},
		{`set t: ((int,),) = ((1,),);pub func main(){print(t);}`, "((1,),)", ""},
		{`pub func main() {
			set (a,b,c) = (1,2,3)
			print(a,b,c);
		}`, "123", ""},
		{`pub func main() {
			set (a,b) = (1,2,3), c=0
			print(a,b,c);
		}`, "120", ""},
		{`pub func main() {
			set t:(int,int,int) = (nil,nil,nil)
			set (t0) = t
			print(t0);
		}`, "nil", ""},
		{`pub func main() {
			func createTuple() => (int,int,int) {
				return (1,nil,nil)
			}
			set t = createTuple()
			set (t0) = t
			print(t0);
		}`, "1", ""},
	})
}
