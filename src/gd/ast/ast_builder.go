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
	"errors"
	"gdlang/src/gd/scanner"
)

// A processor that will process the AST
type AstBuilderProc struct {
	Ast *Ast
	// Compiles all the grammar and parser Errors into an array
	Errors []error
}

func (p *AstBuilderProc) Init(file *scanner.File, src []byte) error {
	// Reset errors
	p.Errors = make([]error, 0)
	err := p.Ast.Init(file, src, p.astErrorListener)
	if err != nil {
		return err
	}

	return nil
}

func (p *AstBuilderProc) Dispose() {
	p.Errors = nil
	p.Ast.Dispose()
}

func (p *AstBuilderProc) Build() (Node, error) {
	rootNode := p.Ast.Build()
	if len(p.Errors) > 0 {
		return nil, errors.Join(p.Errors...)
	}

	return rootNode, nil
}

// All the reported errors from the ast will be appended to the errors array
// Ast errors are the errors that are within the grammar and parser
func (p *AstBuilderProc) astErrorListener(err error) {
	p.AppendError(err)
}

func (p *AstBuilderProc) AppendError(err error) {
	if err != nil {
		p.Errors = append(p.Errors, err)
	}
}

func NAstBuilderProc() *AstBuilderProc {
	return &AstBuilderProc{NAst(), make([]error, 0)}
}
