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
	"fmt"
	"gdlang/lib/runtime"
	"gdlang/src/comn"
	"gdlang/src/gd/scanner"
)

// Errors that can be thrown by the AST and the scanner
type AstErrorListener func(err error)

// Ast is the abstract syntax tree of the source code
type Ast struct {
	scanner  scanner.Scanner
	CTok     *NodeTokenInfo
	Root     *NodeFile
	File     *scanner.File
	Src      []byte
	listener AstErrorListener
}

func init() {
	yyErrorVerbose = true
}

func (a *Ast) Init(file *scanner.File, src []byte, errHandler AstErrorListener) error {
	// Reset it because it can be reused
	a.Root = nil
	a.CTok = nil
	a.listener = errHandler
	a.File = file
	a.Src = src

	err := a.scanner.Init(file, src, a.scannerErrorHandler, 0)
	if err != nil {
		return err
	}

	return nil
}

func (a *Ast) Dispose() {
	a.Root = nil
	a.CTok = nil
	a.listener = nil
	a.File = nil
	a.Src = nil
	a.scanner.Dispose()
}

// Build calls the yyParse and builds the AST
// and returns the root node of the AST.
func (a *Ast) Build() *NodeFile {
	// Returned token is ignored since the exception is already caught
	_ = yyParse(a)
	return a.Root
}

// Errors reported by the scanner
func (a *Ast) scannerErrorHandler(msg string, pos scanner.Position) {
	a.dispatchError(comn.NewError(comn.DefaultSyntaxErrCode, msg, comn.SyntaxError, pos, nil))
}

func (a *Ast) dispatchError(err error) {
	if a.listener != nil {
		a.listener(err)
	}
}

func (a *Ast) Lex(union *yySymType) int {
	offsS, offsE, tok, lit := a.scanner.Scan()
	position := a.File.PositionFrom(offsS, offsE)
	tokId := int(scanner.ILLEGAL)

	if offsS != scanner.NoPos && offsE != scanner.NoPos {
		if tok, ok := ScannerTokens[tok]; ok {
			tokId = tok
		} else {
			panic(fmt.Sprintf("token (%q %q %d-%d) not implemented yet: ", tok, lit, offsS, offsE))
		}

		// Special cases
		switch tok {
		case scanner.STRING, scanner.CHAR:
			if len(lit) > 1 {
				lit = lit[1 : len(lit)-1]
			}
		}

		nToken := NewNodeTokenInfo(tok, position, lit)
		a.CTok = nToken

		union.token = nToken

		return tokId
	} else {
		panic(fmt.Sprintf("token (%q %q %d-%d) not implemented yet: ", tok, lit, offsS, offsE))
	}
}

// Errors reported by the `yyParse`
func (a *Ast) Error(msg string) {
	switch {
	case runtime.ContainsAny(msg, "unexpected", "expecting"):
		expectedTokens := ScanTokens(msg)
		if len(expectedTokens) > 0 {
			switch expectedTokens[0] {
			case ";":
				a.dispatchError(comn.CreateExpSymbolsErr(expectedTokens, a.CTok.Position))
				return
			default:
				a.dispatchError(comn.CreateUnexpectedSyntaxError(expectedTokens[0], expectedTokens[1:], a.CTok.Position))
				return
			}
		}
	case runtime.ContainsAny(msg, "NIL_AS_A_TYPE_ERR"):
		a.dispatchError(comn.NewError(comn.NilAsATypeErrCode, comn.NilAsATypeErrMsg, comn.SyntaxError, a.CTok.Position, nil))
		return
	case runtime.ContainsAny(msg, "USE_ONLY_AT_HEADER_ERR"):
		a.dispatchError(comn.NewError(comn.UseOnlyAtHeaderErrCode, comn.UseOnlyAtHeaderErrMsg, comn.SyntaxError, a.CTok.Position, nil))
		return
	}

	panic("missing error handling for: " + msg)
}

func NAst() *Ast {
	return &Ast{}
}
