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

type Token rune

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF

	// Literals

	IDENT
	INT     // 12345
	FLOAT   // 123.45
	BOOL    // true, false
	STRING  // "abc"
	IMAG    // 123.45i
	CHAR    // 'a'
	COMMENT // // or /* */

	// Operators and delimiters

	QMARK // ?
	NSAFE // ?.

	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=

	ARROW // =>

	LAND // &&
	LOR  // ||

	NOT    // !
	OR     // |
	RSHIFT // >>
	LSHIFT // <<

	ASSIGN // =

	EQL // ==
	LSS // <
	GTR // >

	NEQ      // !=
	LEQ      // <=
	GEQ      // >=
	ELLIPSIS // ...

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .

	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	SEMICOLON // ;
	COLON     // :

	// Keywords

	keyword_beg

	USE
	SET
	PUB
	CONST
	ELSE
	FOR
	IN
	FUNC
	IF
	RETURN
	BREAK
	TYPEALIAS
	AS

	TANY     // any
	TBOOL    // bool
	TINT     // int
	TFLOAT   // float
	TCOMPLEX // complex
	TSTRING  // string
	TCHAR    // char

	TRUE  // true
	FALSE // false
	NIL   // nil

	keyword_end
)

var Tokens = [...]string{
	ILLEGAL: "ILLEGAL",

	EOF: "eof",

	IDENT:   "IDENT",
	INT:     "INT",
	FLOAT:   "FLOAT",
	STRING:  "STRING",
	IMAG:    "COMPLEX",
	CHAR:    "CHAR",
	COMMENT: "COMMENT",

	QMARK: "?",
	NSAFE: "?.",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	QUO: "/",
	REM: "%",

	ADD_ASSIGN: "+=",
	SUB_ASSIGN: "-=",
	MUL_ASSIGN: "*=",
	QUO_ASSIGN: "/=",
	REM_ASSIGN: "%=",

	ARROW: "=>",

	LAND: "&&",
	LOR:  "||",

	NOT:    "!",
	OR:     "|",
	LSHIFT: "<<",
	RSHIFT: ">>",

	EQL: "==",
	LSS: "<",
	GTR: ">",

	ASSIGN: "=",

	NEQ:      "!=",
	LEQ:      "<=",
	GEQ:      ">=",
	ELLIPSIS: "...",

	LPAREN: "(",
	LBRACK: "[",
	LBRACE: "{",
	COMMA:  ",",
	PERIOD: ".",

	RPAREN:    ")",
	RBRACK:    "]",
	RBRACE:    "}",
	SEMICOLON: ";",
	COLON:     ":",

	USE:       "use",
	SET:       "set",
	CONST:     "const",
	PUB:       "pub",
	ELSE:      "else",
	FOR:       "for",
	IN:        "in",
	FUNC:      "func",
	IF:        "if",
	RETURN:    "return",
	BREAK:     "break",
	TYPEALIAS: "typealias",
	AS:        "as",

	TANY:     "any",
	TBOOL:    "bool",
	TINT:     "int",
	TFLOAT:   "float",
	TCOMPLEX: "complex",
	TSTRING:  "string",
	TCHAR:    "char",

	TRUE:  "true",
	FALSE: "false",
	NIL:   "nil",
}

// String returns the string corresponding to the token.
func (tok Token) Fmt() string {
	s := ""
	if 0 <= tok && tok < Token(len(Tokens)) {
		s = Tokens[tok]
	}

	return s
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token, keyword_end-(keyword_beg+1))
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[Tokens[i]] = i
	}
}

// Lookup maps an identifier to its keyword token or [IDENT] (if not a keyword).
func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}

	return IDENT
}

func LookUpToken(token Token) string {
	return Tokens[token]
}

func (tok Token) IsKeyword() bool {
	return keyword_beg < tok && tok < keyword_end
}
