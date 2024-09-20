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
	"gdlang/lib/runtime"
	"testing"
)

func TestExprCases(t *testing.T) {
	tmpl := `
	pub func main(){
		set a=$SRC
		print(a)
	}`
	RunTestsWithTemplate(t, tmpl, []Test{
		{"1", runtime.NewGDIntNumber(1).ToString(), ""},
		{`(100)`, runtime.NewGDIntNumber(100).ToString(), ""},
		{`(100,)`, runtime.NewGDTuple(runtime.NewGDIntNumber(100)).ToString(), ""},
		{`((100,))`, runtime.NewGDTuple(runtime.NewGDIntNumber(100)).ToString(), ""},
		{`100;`, runtime.NewGDIntNumber(100).ToString(), ""},
		{`true`, runtime.GDBool(true).ToString(), ""},
		{`false`, runtime.GDBool(false).ToString(), ""},
		{`nil`, runtime.GDZNil.ToString(), ""},
		{`"hola"`, runtime.GDString("hola").ToString(), ""},
		{`'a'`, runtime.GDChar('a').ToString(), ""},
		{`1.0`, runtime.NewGDFloatNumber(1.0).ToString(), ""},
		{`1+1`, runtime.NewGDIntNumber(2).ToString(), ""},
		{`1-1`, runtime.NewGDIntNumber(0).ToString(), ""},
		{`10*10`, runtime.NewGDIntNumber(100).ToString(), ""},
		{`10/10`, runtime.NewGDIntNumber(1).ToString(), ""},
		{`10%10`, runtime.NewGDIntNumber(0).ToString(), ""},
		{`10*10`, runtime.NewGDIntNumber(100).ToString(), ""},
		{`10.0+1`, runtime.NewGDFloatNumber(11.0).ToString(), ""},
		{`0.5+0.5`, runtime.NewGDFloatNumber(1.0).ToString(), ""},
		{`0.5-0.5`, runtime.NewGDFloatNumber(0.0).ToString(), ""},
		{`0.5*0.5`, runtime.NewGDFloatNumber(0.25).ToString(), ""},
		{`0.5/0.5`, runtime.NewGDFloatNumber(1.0).ToString(), ""},
		{`1i+1`, runtime.NewGDComplexNumber(1i + 1).ToString(), ""},
		{`1i-1`, runtime.NewGDComplexNumber(1i - 1).ToString(), ""},
		{`1i*1`, runtime.NewGDComplexNumber(1i * 1).ToString(), ""},
		{`1i/1`, runtime.NewGDComplexNumber(1i / 1).ToString(), ""},
		{`5.0+1`, runtime.NewGDFloatNumber(5.0 + 1).ToString(), ""},
		{`"h"+"o"+"l"+"a"`, runtime.GDString("hola").ToString(), ""},
		{`"hola"+(1,)`, runtime.GDString("hola(1,)").ToString(), ""},
		{`"hola"+(3.1416)`, runtime.GDString("hola3.1416").ToString(), ""},
		{`"hola"+'s'`, runtime.GDString("holas").ToString(), ""},
		{`'h'+"ola"`, runtime.GDString("hola").ToString(), ""},
		{"'h'+`ola`", runtime.GDString("hola").ToString(), ""},
		{`1>2`, runtime.GDBool(false).ToString(), ""},
		{`1<2`, runtime.GDBool(true).ToString(), ""},
		{`1>=2`, runtime.GDBool(false).ToString(), ""},
		{`1<=2`, runtime.GDBool(true).ToString(), ""},
		{`1==2`, runtime.GDBool(false).ToString(), ""},
		{`1!=2`, runtime.GDBool(true).ToString(), ""},
		{`1==1`, runtime.GDBool(true).ToString(), ""},
		{`1!=1`, runtime.GDBool(false).ToString(), ""},
		{`true==true`, runtime.GDBool(true).ToString(), ""},
		{`true!=true`, runtime.GDBool(false).ToString(), ""},
		{`true==false`, runtime.GDBool(false).ToString(), ""},
		{`true!=false`, runtime.GDBool(true).ToString(), ""},
		{`false==false`, runtime.GDBool(true).ToString(), ""},
		{`false!=false`, runtime.GDBool(false).ToString(), ""},
		{`false==true`, runtime.GDBool(false).ToString(), ""},
		{`false!=true`, runtime.GDBool(true).ToString(), ""},
		{`true&&true`, runtime.GDBool(true).ToString(), ""},
		{`true&&false`, runtime.GDBool(false).ToString(), ""},
		{`false&&true`, runtime.GDBool(false).ToString(), ""},
		{`false&&false`, runtime.GDBool(false).ToString(), ""},
		{`true||true`, runtime.GDBool(true).ToString(), ""},
		{`true||false`, runtime.GDBool(true).ToString(), ""},
		{`false||true`, runtime.GDBool(true).ToString(), ""},
		{`false||false`, runtime.GDBool(false).ToString(), ""},
		{`!true`, runtime.GDBool(false).ToString(), ""},
		{`!false`, runtime.GDBool(true).ToString(), ""},
	})
}

func TestExpOpCases(t *testing.T) {
	tmpl := `
	pub func main(){
		set a=$SRC
		print(a)
	}`
	RunTestsWithTemplate(t, tmpl, []Test{
		{`1+1`, runtime.NewGDIntNumber(2).ToString(), ""},
		{`1-1`, runtime.NewGDIntNumber(0).ToString(), ""},
		{`1*1`, runtime.NewGDIntNumber(1).ToString(), ""},
		{`1/1`, runtime.NewGDIntNumber(1).ToString(), ""},
		{`1%1`, runtime.NewGDIntNumber(0).ToString(), ""},
		{`1+1.0`, runtime.NewGDFloatNumber(2.0).ToString(), ""},
		{`1-1.0`, runtime.NewGDFloatNumber(0.0).ToString(), ""},
		{`1*1.0`, runtime.NewGDFloatNumber(1.0).ToString(), ""},
		{`1/1.0`, runtime.NewGDFloatNumber(1.0).ToString(), ""},
		{`1%1.0`, runtime.NewGDFloatNumber(0.0).ToString(), ""},
		{`1+1i`, runtime.NewGDComplexNumber(1 + 1i).ToString(), ""},
		{`1-1i`, runtime.NewGDComplexNumber(1 - 1i).ToString(), ""},
		{`1*1i`, runtime.NewGDComplexNumber(1 * 1i).ToString(), ""},
		{`1/1i`, runtime.NewGDComplexNumber(1 / 1i).ToString(), ""},
		{`1+1.0`, runtime.NewGDFloatNumber(2.0).ToString(), ""},
		{`1+1i`, runtime.NewGDComplexNumber(1 + 1i).ToString(), ""},
		{`1+1.0`, runtime.NewGDFloatNumber(2.0).ToString(), ""},
		{`1+1i`, runtime.NewGDComplexNumber(1 + 1i).ToString(), ""},
		{`1+1.0`, runtime.NewGDFloatNumber(2.0).ToString(), ""},
		{`1+1i`, runtime.NewGDComplexNumber(1 + 1i).ToString(), ""},
		{`1+2+3+4+5+6+7+8+9+10`, runtime.NewGDIntNumber(55).ToString(), ""},
		{`10-9-8-7-6-5-4-3-2-1`, runtime.NewGDIntNumber(-35).ToString(), ""},
		{`2*3*4*5*6*7*8*9*10`, runtime.NewGDIntNumber(3628800).ToString(), ""},
		{`100/10/5`, runtime.NewGDIntNumber(2).ToString(), ""},
		{`10%3%2`, runtime.NewGDIntNumber(1).ToString(), ""},
		{`2+3*4-5/2`, runtime.NewGDFloatNumber(2 + 3*4 - 5/2).ToString(), ""},
		{`(2+3)*(4-5)/2.0`, runtime.NewGDFloatNumber((2 + 3) * (4 - 5) / 2.0).ToString(), ""},
		{`2+3*(4-5)/2.0`, runtime.NewGDFloatNumber(2 + 3*(4-5)/2.0).ToString(), ""},
		{`(2+3)*4-5/2.0`, runtime.NewGDFloatNumber((2+3)*4 - 5/2.0).ToString(), ""},
		{`2+3*4-(5/2.0)`, runtime.NewGDFloatNumber(2 + 3*4 - (5 / 2.0)).ToString(), ""},
		{`(2+3)*(4-(5/2.0))`, runtime.NewGDFloatNumber((2 + 3) * (4 - (5 / 2.0))).ToString(), ""},
		{`2+3*(4-(5/2.0))`, runtime.NewGDFloatNumber(2 + 3*(4-(5/2.0))).ToString(), ""},
		{`(2+3)*4-(5/2.0)`, runtime.NewGDFloatNumber((2+3)*4 - (5 / 2.0)).ToString(), ""},
		{`2+3*4-5/2+6-7*8/9`, runtime.NewGDIntNumber(2 + 3*4 - 5/2 + 6 - 7*8/9).ToString(), ""},
		{`(2+3)*(4-5)/2.0+6-(7*8)/9.0`, runtime.NewGDFloatNumber((2+3)*(4-5)/2.0 + 6 - (7*8)/9.0).ToString(), ""},
		{`(2+3)*4-5/2.0+6-(7*8)/9.0`, runtime.NewGDFloatNumber((2+3)*4 - 5/2.0 + 6 - (7*8)/9.0).ToString(), ""},
		{`2+3*4-(5/2.0)+6-7*8/9.0`, runtime.NewGDFloatNumber(2 + 3*4 - (5 / 2.0) + 6 - 7*8/9.0).ToString(), ""},
		{`(2+3)*(4-(5/2.0))+6-(7*8)/9.0`, runtime.NewGDFloatNumber((2+3)*(4-(5/2.0)) + 6 - (7*8)/9.0).ToString(), ""},
		{`2+3*(4-(5/2.0))+6-7*8/9.0`, runtime.NewGDFloatNumber(2 + 3*(4-(5/2.0)) + 6 - 7*8/9.0).ToString(), ""},
		{`(2+3)*4-(5/2.0)+6-(7*8)/9.0`, runtime.NewGDFloatNumber((2+3)*4 - (5 / 2.0) + 6 - (7*8)/9.0).ToString(), ""},
		{`2+3*4-5/2.0+6-7*8/9.0+10-11*12.0/13.0`, runtime.GDFloat32(2 + 3*4 - 5/runtime.GDFloat32(2.0) + 6 - 7*8/runtime.GDFloat32(9.0) + 10 - 11*12/runtime.GDFloat32(13.0)).ToString(), ""},
		{`(2+3)*(4-5)/2.0+6-(7*8)/9.0+10-(11*12)/13.0`, runtime.GDFloat32((2+3)*(4-5)/runtime.GDFloat32(2.0) + 6 - (7*8)/runtime.GDFloat32(9.0) + 10 - (11*12)/runtime.GDFloat32(13.0)).ToString(), ""},
		{`2+3*(4-5)/2.0+6-7*8/9.0+10-11*12/13.0`, runtime.GDFloat32(2 + 3*(4-5)/runtime.GDFloat32(2.0) + 6 - 7*8/runtime.GDFloat32(9.0) + 10 - 11*12/runtime.GDFloat32(13.0)).ToString(), ""},
		{`(2+3)*4-5/2.0+6-(7*8)/9.0+10-(11*12)/13.0`, runtime.GDFloat32((2+3)*4 - 5/float64(2.0) + 6 - (7*8)/float64(9.0) + 10 - (11*12)/float64(13.0)).ToString(), ""},
		{`true && false`, runtime.GDBool(false).ToString(), ""},
		{`true || false`, runtime.GDBool(true).ToString(), ""},
		{`true && true || false`, runtime.GDBool(true).ToString(), ""},
		{`true || false && true`, runtime.GDBool(true).ToString(), ""},
		{`true && (false || true)`, runtime.GDBool(true).ToString(), ""},
		{`(true && false) || true`, runtime.GDBool(true).ToString(), ""},
		{`true && false || true`, runtime.GDBool(true).ToString(), ""},
		{`true || (false && true)`, runtime.GDBool(true).ToString(), ""},
		{`true && false && true || false`, runtime.GDBool(false).ToString(), ""},
		{`true || false || true && true`, runtime.GDBool(true).ToString(), ""},
		{`true && (false || true) && true`, runtime.GDBool(true).ToString(), ""},
		{`(true && false) || true || false`, runtime.GDBool(true).ToString(), ""},
		{`true && false || true || false`, runtime.GDBool(true).ToString(), ""},
		{`true || (false && true) || false`, runtime.GDBool(true).ToString(), ""},
		{`true && false && true || false || true`, runtime.GDBool(true).ToString(), ""},
		{`true || false || true && true && false`, runtime.GDBool(true).ToString(), ""},
		{`true && (false || true) && true && false`, runtime.GDBool(false).ToString(), ""},
		{`(true && false) || true || false || true`, runtime.GDBool(true).ToString(), ""},
		{`true && false || true || false || true`, runtime.GDBool(true).ToString(), ""},
		{`true || (false && true) || false || true`, runtime.GDBool(true).ToString(), ""},
		{`true && false && true || false || true || false`, runtime.GDBool(true).ToString(), ""},
		{`true || false || true && true && false || true`, runtime.GDBool(true).ToString(), ""},
		{`true && (false || true) && true && false || true`, runtime.GDBool(true).ToString(), ""},
		{`(true && false) || true || false || true || false`, runtime.GDBool(true).ToString(), ""},
		{`true && false || true || false || true || false`, runtime.GDBool(true).ToString(), ""},
		{`true || (false && true) || false || true || false`, runtime.GDBool(true).ToString(), ""},
		{`2+3*4-5/2+6-(7*8)/9+10-(11*12)/13+14*15-16/17+18-(19*20)/21+22-(23*24)/25+26*27-28/29+30-(31*32)/33+34-(35*36)/37+38*39-40/41+42-(43*44)/45+46-(47*48)/49+50*51-52/53+54-(55*56)/57+58-(59*60)/61+62*63-64/65+66-(67*68)/69+70-(71*72)/73+74*75-76/77+78-(79*80)/81+82-(83*84)/85+86*87-88/89+90-(91*92)/93+94-(95*96)/97+98*99-100/101`, runtime.NewGDIntNumber(31596).ToString(), ""},
	})
}
