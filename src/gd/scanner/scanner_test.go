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

package scanner

import (
	"fmt"
	"testing"
)

type tokenLitPos struct {
	tok Token
	lit string
	pos Position
}

var fset = NFileSet()

type errorCollector struct {
	cnt int      // number of errors encountered
	msg string   // last error message encountered
	pos Position // last error position encountered
}

var errors = []struct {
	src        string
	tok        Token
	posS, posE int
	lit        string
	err        string
}{
	{"\a", ILLEGAL, 1, 1, "", "illegal `char` U+0007"},
	{`#`, ILLEGAL, 1, 1, "", "illegal `char` U+0023 '#'"},
	{`…`, ILLEGAL, 1, 1, "", "illegal `char` U+2026 '…'"},
	{`' '`, CHAR, 0, 0, `' '`, ""},
	{`''`, CHAR, 1, 2, `''`, "illegal `char` literal"},
	{`'12'`, CHAR, 1, 4, `'12'`, "illegal `char` literal"},
	{`'123'`, CHAR, 1, 5, `'123'`, "illegal `char` literal"},
	{`'\0'`, CHAR, 4, 4, `'\0'`, "illegal `char` U+0027 ''' in escape sequence"},
	{`'\07'`, CHAR, 5, 5, `'\07'`, "illegal `char` U+0027 ''' in escape sequence"},
	{`'\8'`, CHAR, 3, 3, `'\8'`, "unknown escape sequence"},
	{`'\08'`, CHAR, 4, 4, `'\08'`, "illegal `char` U+0038 '8' in escape sequence"},
	{`'\x'`, CHAR, 4, 4, `'\x'`, "illegal `char` U+0027 ''' in escape sequence"},
	{`'\x0'`, CHAR, 5, 5, `'\x0'`, "illegal `char` U+0027 ''' in escape sequence"},
	{`'\x0g'`, CHAR, 5, 5, `'\x0g'`, "illegal `char` U+0067 'g' in escape sequence"},
	{`'\u'`, CHAR, 4, 4, `'\u'`, "illegal `char` U+0027 ''' in escape sequence"},
	{`'\u0'`, CHAR, 5, 5, `'\u0'`, "illegal `char` U+0027 ''' in escape sequence"},
	{`'\u00'`, CHAR, 6, 6, `'\u00'`, "illegal `char` U+0027 ''' in escape sequence"},
	{`'\u000'`, CHAR, 7, 7, `'\u000'`, "illegal `char` U+0027 ''' in escape sequence"},
	{`'\u000`, CHAR, 1, 6, `'\u000`, "escape sequence not terminated"},
	{`'\u0000'`, CHAR, 0, 0, `'\u0000'`, ""},
	{`'\U'`, CHAR, 4, 4, `'\U'`, "illegal `char` U+0027 ''' in escape sequence"},
	{`'\U0'`, CHAR, 5, 5, `'\U0'`, "illegal `char` U+0027 ''' in escape sequence"},
	{`'\U00'`, CHAR, 6, 6, `'\U00'`, "illegal `char` U+0027 ''' in escape sequence"},
	{`'\U000'`, CHAR, 7, 7, `'\U000'`, "illegal `char` U+0027 ''' in escape sequence"},
	{`'\U0000'`, CHAR, 8, 8, `'\U0000'`, "illegal `char` U+0027 ''' in escape sequence"},
	{`'\U00000'`, CHAR, 9, 9, `'\U00000'`, "illegal `char` U+0027 ''' in escape sequence"},
	{`'\U000000'`, CHAR, 10, 10, `'\U000000'`, "illegal `char` U+0027 ''' in escape sequence"},
	{`'\U0000000'`, CHAR, 11, 11, `'\U0000000'`, "illegal `char` U+0027 ''' in escape sequence"},
	{`'\U0000000`, CHAR, 1, 10, `'\U0000000`, "escape sequence not terminated"},
	{`'\U00000000'`, CHAR, 0, 0, `'\U00000000'`, ""},
	{`'\Uffffffff'`, CHAR, 3, 3, `'\Uffffffff'`, "escape sequence is invalid Unicode code point"},
	{`'`, CHAR, 1, 1, `'`, "`char` literal not terminated"},
	{`'9`, CHAR, 1, 2, `'9`, "`char` literal not terminated"},
	{`'\`, CHAR, 1, 2, `'\`, "escape sequence not terminated"},
	{"'\n", CHAR, 1, 2, "'", "`char` literal not terminated"},
	{"'\n   ", CHAR, 1, 2, "'", "`char` literal not terminated"},
	{`""`, STRING, 0, 0, `""`, ""},
	{`"abc`, STRING, 1, 4, `"abc`, "`string` literal not terminated"},
	{"\"abc\n", STRING, 1, 5, `"abc`, "`string` literal not terminated"},
	{"\"abc\n   ", STRING, 1, 5, `"abc`, "`string` literal not terminated"},
	{"``", STRING, 0, 0, "``", ""},
	{"`", STRING, 1, 1, "`", "raw `string` literal not terminated"},
	{"/**/", COMMENT, 0, 0, "/**/", ""},
	{"/*", COMMENT, 1, 2, "/*", "comment not terminated"},
	{"077", INT, 0, 0, "077", ""},
	{"078.", FLOAT, 0, 0, "078.", ""},
	{"0.i", IMAG, 0, 0, "0.i", ""},
	{"*", MUL, 0, 0, "", ""},
	{"0..i", FLOAT, 0, 0, "0.", ""},
	{"07801234567.", FLOAT, 0, 0, "07801234567.", ""},
	{"078e0", FLOAT, 0, 0, "078e0", ""},
	{"0E", FLOAT, 1, 2, "0E", "exponent has no digits"},
	{"078", INT, 3, 3, "078", "invalid digit '8' in octal literal"},
	{"07090000008", INT, 4, 4, "07090000008", "invalid digit '9' in octal literal"},
	{"0x", INT, 1, 2, "0x", "hexadecimal literal has no digits"},
	{"\"abc\x00def\"", STRING, 5, 5, "\"abc\x00def\"", "illegal `char` NUL"},
	{"\"abc\x80def\"", STRING, 5, 5, "\"abc\x80def\"", "illegal UTF-8 encoding"},
	{"\ufeff\ufeff", ILLEGAL, 4, 4, "\ufeff\ufeff", "illegal byte order mark"},                        // only first BOM is ignored
	{"//\ufeff", COMMENT, 3, 3, "//\ufeff", "illegal byte order mark"},                                // only first BOM is ignored
	{"'\ufeff" + `'`, CHAR, 2, 2, "'\ufeff" + `'`, "illegal byte order mark"},                         // only first BOM is ignored
	{`"` + "abc\ufeffdef" + `"`, STRING, 5, 5, `"` + "abc\ufeffdef" + `"`, "illegal byte order mark"}, // only first BOM is ignored
	{"abc\x00def", IDENT, 4, 4, "abc", "illegal `char` NUL"},
	{"abc\x00", IDENT, 4, 4, "abc", "illegal `char` NUL"},
	{"“abc”", ILLEGAL, 1, 1, "abc", "curly quotation mark `“` was found (use neutral `\"`) instead"},
	{"“abcdefgh”", ILLEGAL, 1, 1, "abcdefgh", "curly quotation mark `“` was found (use neutral `\"`) instead"},
}

func checkError(t *testing.T, src string, tok Token, posS, posE int, lit, err string) {
	var s Scanner
	var h errorCollector
	eh := func(msg string, pos Position) {
		h.cnt++
		h.msg = msg
		h.pos = pos
	}
	file, err1 := fset.AddFile("checkErrortest", fset.Base(), len(src))
	if err1 != nil {
		t.Fatalf("error while setting up file: %v", err)
	}

	err2 := s.Init(file, []byte(src), eh, ScanComments|dontInsertSemis)
	if err2 != nil {
		t.Fatalf("error while initializing scanner: %v", err)
	}

	_, _, tok0, lit0 := s.Scan()
	if tok0 != tok {
		t.Errorf("%q: got %d, expected %d", src, tok0, tok)
	}
	if tok0 != ILLEGAL && lit0 != lit {
		t.Errorf("%q: got literal %q, expected %q", src, lit0, lit)
	}
	cnt := 0
	if err != "" {
		cnt = 1
	}
	if h.cnt != cnt {
		t.Errorf("%q: got cnt %d, expected %d", src, h.cnt, cnt)
	}
	if h.msg != err {
		t.Errorf("%q: got msg %q, expected %q", src, h.msg, err)
	}
	if h.pos.ColStart != posS {
		t.Errorf("%q: got column start %d, expected %d", src, h.pos.ColStart, posS)
	}
	if h.pos.ColEnd != posE {
		t.Errorf("%q: got column end %d, expected %d", src, h.pos.ColEnd, posE)
	}
}

func TestScanErrors(t *testing.T) {
	for _, e := range errors {
		checkError(t, e.src, e.tok, e.posS, e.posE, e.lit, e.err)
	}
}

func TestTokenPositionStartEnd(t *testing.T) {
	for i, test := range []struct {
		src   string
		tlPos []tokenLitPos
	}{
		{
			"a?.b", []tokenLitPos{
				{IDENT, "a", Position{"test.gd", 1, 1, 1}},
				{NSAFE, "", Position{"test.gd", 1, 2, 3}},
				{IDENT, "b", Position{"test.gd", 1, 4, 4}},
			},
		},
		{
			"0.i", []tokenLitPos{
				{IMAG, "0.i", Position{"test.gd", 1, 1, 3}},
			},
		},
		{
			"0..", []tokenLitPos{
				{FLOAT, "0.", Position{"test.gd", 1, 1, 2}},
				{PERIOD, "", Position{"test.gd", 1, 3, 3}},
			},
		},
		{
			"0..i", []tokenLitPos{
				{FLOAT, "0.", Position{"test.gd", 1, 1, 2}},
				{PERIOD, "", Position{"test.gd", 1, 3, 3}},
				{IDENT, "i", Position{"test.gd", 1, 4, 4}},
			},
		},
		{
			"0i", []tokenLitPos{
				{IMAG, "0i", Position{"test.gd", 1, 1, 2}},
			},
		},
		{
			"set a = 1", []tokenLitPos{
				{SET, "", Position{"test.gd", 1, 1, 3}},
				{IDENT, "a", Position{"test.gd", 1, 5, 5}},
				{ASSIGN, "", Position{"test.gd", 1, 7, 7}},
				{INT, "1", Position{"test.gd", 1, 9, 9}},
			},
		},
		{
			"set aabcdefg 	= 	10000000000", []tokenLitPos{
				{SET, "", Position{"test.gd", 1, 1, 3}},
				{IDENT, "aabcdefg", Position{"test.gd", 1, 5, 12}},
				{ASSIGN, "", Position{"test.gd", 1, 15, 15}},
				{INT, "10000000000", Position{"test.gd", 1, 18, 28}},
			},
		},
		// Test with tabulations
		{
			"\t\t\tset aabcdefg\t=\t10000000000", []tokenLitPos{
				{SET, "", Position{"test.gd", 1, 4, 6}},
				{IDENT, "aabcdefg", Position{"test.gd", 1, 8, 15}},
				{ASSIGN, "", Position{"test.gd", 1, 17, 17}},
				{INT, "10000000000", Position{"test.gd", 1, 19, 29}},
			},
		},
		// Test use
		{
			"use a.b", []tokenLitPos{
				{USE, "", Position{"test.gd", 1, 1, 3}},
				{IDENT, "a", Position{"test.gd", 1, 5, 5}},
				{PERIOD, "", Position{"test.gd", 1, 6, 6}},
				{IDENT, "b", Position{"test.gd", 1, 7, 7}},
			},
		},
		// Test nested use
		{
			"use a.b.c", []tokenLitPos{
				{USE, "", Position{"test.gd", 1, 1, 3}},
				{IDENT, "a", Position{"test.gd", 1, 5, 5}},
				{PERIOD, "", Position{"test.gd", 1, 6, 6}},
				{IDENT, "b", Position{"test.gd", 1, 7, 7}},
				{PERIOD, "", Position{"test.gd", 1, 8, 8}},
				{IDENT, "c", Position{"test.gd", 1, 9, 9}},
			},
		},
	} {
		srcByte := []byte(test.src)
		file, err := fset.AddFile(fmt.Sprintf("test_%d", i), fset.Base(), len(srcByte))
		if err != nil {
			t.Fatalf("error while setting up file: %v", err)
		}

		var s Scanner
		s.Init(file, srcByte, nil, ScanComments|dontInsertSemis)

		index := 0
		for index < len(test.tlPos) {
			pos, epos, tok, lit := s.Scan()
			if tok == EOF {
				break
			}

			tokPos := file.PositionFrom(pos, epos)
			ttok := test.tlPos[index]

			if ttok.lit != lit || ttok.tok != tok {
				t.Errorf("bad token: got %s, expected %s when scanning %q", Tokens[tok], Tokens[ttok.tok], test.src)
			}

			checkPos(t, lit, tokPos, ttok.pos)

			index++
		}
	}
}

func checkPos(t *testing.T, lit string, pos Position, expected Position) {
	if pos.Line != expected.Line {
		t.Errorf("bad line for %q: got %d, expected %d", lit, pos.Line, expected.Line)
	}

	if pos.ColStart != expected.ColStart {
		t.Errorf("bad column start for %q: got %d, expected %d", lit, pos.ColStart, expected.ColStart)
	}

	if pos.ColEnd != expected.ColEnd {
		t.Errorf("bad column end for %q: got %d, expected %d", lit, pos.ColEnd, expected.ColEnd)
	}
}
