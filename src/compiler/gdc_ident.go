//go:build !debug
// +build !debug

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

package compiler

import (
	"gdlang/lib/runtime"
	"gdlang/src/gd/ast"
)

func (c *GDCompiler) DeriveIdent(node ast.Node) runtime.GDIdent {
	if ident := node.RuntimeIdent(); ident != nil {
		return ident
	}

	return node.InferredIdent()
}

func (c *GDCompiler) DeriveType(node ast.Node) runtime.GDTypable {
	if typ := node.RuntimeType(); typ != nil {
		return typ
	}

	return node.InferredType()
}
