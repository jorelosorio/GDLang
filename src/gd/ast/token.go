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

package ast

import (
	"gdlang/src/gd/scanner"
)

// Map of scanner tokens to the grammar representation.
var ScannerTokens = map[scanner.Token]int{
	scanner.EOF:     -1,
	scanner.ILLEGAL: -2,

	scanner.IDENT:   LIDENT,
	scanner.INT:     LINT,
	scanner.FLOAT:   LFLOAT,
	scanner.STRING:  LSTRING,
	scanner.IMAG:    LIMAG,
	scanner.CHAR:    LCHAR,
	scanner.COMMENT: LCOMMENT,

	scanner.QMARK: LQMARK,
	scanner.NSAFE: LNSAFE,

	scanner.ADD: LADD,
	scanner.SUB: LSUB,
	scanner.MUL: LMUL,
	scanner.QUO: LQUO,
	scanner.REM: LREM,

	scanner.ADD_ASSIGN: LADD_ASSIGN,
	scanner.SUB_ASSIGN: LSUB_ASSIGN,
	scanner.MUL_ASSIGN: LMUL_ASSIGN,
	scanner.QUO_ASSIGN: LQUO_ASSIGN,
	scanner.REM_ASSIGN: LREM_ASSIGN,

	scanner.ARROW: LARROW,

	scanner.LAND: LLAND,
	scanner.LOR:  LLOR,
	scanner.NOT:  LNOT,

	scanner.OR:     LOR,
	scanner.LSHIFT: LLSHIFT,
	scanner.RSHIFT: LRSHIFT,

	scanner.EQL:    LEQL,
	scanner.LSS:    LLSS,
	scanner.GTR:    LGTR,
	scanner.ASSIGN: LASSIGN,

	scanner.NEQ:      LNEQ,
	scanner.LEQ:      LLEQ,
	scanner.GEQ:      LGEQ,
	scanner.ELLIPSIS: LELLIPSIS,

	scanner.LPAREN: LLPAREN,
	scanner.LBRACK: LLBRACK,
	scanner.LBRACE: LLBRACE,
	scanner.COMMA:  LCOMMA,
	scanner.PERIOD: LPERIOD,

	scanner.RPAREN:    LRPAREN,
	scanner.RBRACK:    LRBRACK,
	scanner.RBRACE:    LRBRACE,
	scanner.SEMICOLON: LSEMICOLON,
	scanner.COLON:     LCOLON,

	scanner.USE:       LUSE,
	scanner.SET:       LSET,
	scanner.PUB:       LPUB,
	scanner.CONST:     LCONST,
	scanner.ELSE:      LELSE,
	scanner.FOR:       LFOR,
	scanner.IN:        LIN,
	scanner.FUNC:      LFUNC,
	scanner.IF:        LIF,
	scanner.RETURN:    LRETURN,
	scanner.BREAK:     LBREAK,
	scanner.TYPEALIAS: LTYPEALIAS,
	scanner.AS:        LAS,

	scanner.TANY:     LTANY,
	scanner.TBOOL:    LTBOOL,
	scanner.TINT:     LTINT,
	scanner.TFLOAT:   LTFLOAT,
	scanner.TCOMPLEX: LTCOMPLEX,
	scanner.TSTRING:  LTSTRING,
	scanner.TCHAR:    LTCHAR,

	scanner.TRUE:  LTRUE,
	scanner.FALSE: LFALSE,
	scanner.NIL:   LNIL,
}

var LexTokens = map[string]scanner.Token{
	"LIDENT":   scanner.IDENT,
	"LINT":     scanner.INT,
	"LFLOAT":   scanner.FLOAT,
	"LSTRING":  scanner.STRING,
	"LIMAG":    scanner.IMAG,
	"LCHAR":    scanner.CHAR,
	"LCOMMENT": scanner.COMMENT,

	"LNSAFE": scanner.NSAFE,

	"LADD": scanner.ADD,
	"LSUB": scanner.SUB,
	"LMUL": scanner.MUL,
	"LQUO": scanner.QUO,
	"LREM": scanner.REM,

	"LADD_ASSIGN": scanner.ADD_ASSIGN,
	"LSUB_ASSIGN": scanner.SUB_ASSIGN,
	"LMUL_ASSIGN": scanner.MUL_ASSIGN,
	"LQUO_ASSIGN": scanner.QUO_ASSIGN,
	"LREM_ASSIGN": scanner.REM_ASSIGN,

	"LARROW": scanner.ARROW,

	"LLAND": scanner.LAND,
	"LLOR":  scanner.LOR,
	"LNOT":  scanner.NOT,

	"LEQL":    scanner.EQL,
	"LLSS":    scanner.LSS,
	"LGTR":    scanner.GTR,
	"LLSHIFT": scanner.LSHIFT,
	"LRSHIFT": scanner.RSHIFT,
	"LASSIGN": scanner.ASSIGN,

	"LNEQ":      scanner.NEQ,
	"LLEQ":      scanner.LEQ,
	"LGEQ":      scanner.GEQ,
	"LELLIPSIS": scanner.ELLIPSIS,

	"LLPAREN": scanner.LPAREN,
	"LLBRACK": scanner.LBRACK,
	"LLBRACE": scanner.LBRACE,
	"LCOMMA":  scanner.COMMA,
	"LPERIOD": scanner.PERIOD,

	"LRPAREN":    scanner.RPAREN,
	"LRBRACK":    scanner.RBRACK,
	"LRBRACE":    scanner.RBRACE,
	"LSEMICOLON": scanner.SEMICOLON,
	"LCOLON":     scanner.COLON,

	"LUSE":       scanner.USE,
	"LSET":       scanner.SET,
	"LPUB":       scanner.PUB,
	"LCONST":     scanner.CONST,
	"LELSE":      scanner.ELSE,
	"LFOR":       scanner.FOR,
	"LIN":        scanner.IN,
	"LFUNC":      scanner.FUNC,
	"LIF":        scanner.IF,
	"LRETURN":    scanner.RETURN,
	"LBREAK":     scanner.BREAK,
	"LTYPEALIAS": scanner.TYPEALIAS,
	"LAS":        scanner.AS,

	"LTANY":     scanner.TANY,
	"LTBOOL":    scanner.TBOOL,
	"LTINT":     scanner.TINT,
	"LTFLOAT":   scanner.TFLOAT,
	"LTCOMPLEX": scanner.TCOMPLEX,
	"LTSTRING":  scanner.TSTRING,
	"LTCHAR":    scanner.TCHAR,

	"LTRUE":  scanner.TRUE,
	"LFALSE": scanner.FALSE,
	"LNIL":   scanner.NIL,
}

func LookUpToken(lToken string) string {
	if tok, ok := LexTokens[lToken]; ok {
		return scanner.LookUpToken(tok)
	}

	return ""
}
