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

func TestLambdaCases(t *testing.T) {
	RunTestsWithMainTemplate(t, []Test{
		{`func (){}()`, "", ""},
		{`func ()=>int{return 1;}()`, "", ""},
		{`func (){print(1);}()`, "1", ""},
		{`set x = func (){
			return
		}()
		x = 1;
		print(x);`, "1", ""},
		{`set x = func ()=>bool{return true;}();print(x);`, "true", ""},
		{`func (){return;print(1);}()`, "", ""},
		{`set a=1;func(){set a=2;print(a);}()`, "2", ""},
		{`set a=1;func(){set a=2;print(a);}();print(a);`, "21", ""},
		{`set a=1
			func(){
				set a=2
				print(a)
				func(){
					set a=3
					print(a)
				}()
			}()
			print(a)
		`, "231", ""},
		{`set a=func(){print(1);}();`, "1", ""},
		{`set a=func()=>int{return 1;return 2;}();print(a)`, "1", ""},
		{`print(func()=>char{return 'a';}());`, "a", ""},
		{`print(func(){}());`, "nil", ""},
		{`set a=func()=>int{return func()=>int{return 1;}();}();print(a);`, "1", ""},
		{`set a=true;a=func()=>bool{return false;}();print(a)`, "false", ""},
		{`set a=true;a=func()=>int{return 1;}();print(a);`, "", "expected `bool` but got `int`"},
		{`set a=true;a=func()=>char{return 'a';}();print(a);`, "", "expected `bool` but got `char`"},
		{`set a:bool=func()=>string{return "hello";}();print(a);`, "", "expected `bool` but got `string`"},
	})
}
