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

package ast_test

import (
	"gdlang/src/comn"
	"gdlang/src/gd/ast"
	"gdlang/src/gd/scanner"
)

var (
	fileSet = scanner.NFileSet()
	astObj  = ast.NAst()
)

func init() {
	comn.PrettyPrintErrors = false
}

func buildAst(src string, errHandler ast.AstErrorListener) (ast.Node, error) {
	fileSet.Reset()
	srcFile, err := newTestFile(len(src))
	if err != nil {
		return nil, err
	}

	err = astObj.Init(srcFile, []byte(src), errHandler)
	if err != nil {
		return nil, err
	}

	return astObj.Build(), nil
}

func newTestFile(srcLen int) (*scanner.File, error) {
	return fileSet.AddFile("test.gd", fileSet.Base(), srcLen)
}
