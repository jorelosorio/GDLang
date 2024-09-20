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

package ir

import (
	"bytes"
	"fmt"
	"gdlang/lib/runtime"
	"gdlang/src/gd/scanner"
)

type GDIRDiscoverable struct {
	isPub, isConst bool
	ident          runtime.GDIdentType
	GDIRBaseNode
}

func (d *GDIRDiscoverable) IR(padding string) string {
	return fmt.Sprintf("%s%s%s", tif(d.isPub, "pub ", ""), tif(d.isConst, "const ", ""), IRTypeToString(d.ident))
}

func (d *GDIRDiscoverable) Bytecode(bytecode *bytes.Buffer, ctx *GDIRContext) error {
	err := Write(bytecode, d.isPub, d.isConst, d.ident)
	if err != nil {
		return err
	}

	return nil
}

func NewGDIRDiscoverable(isPub, isConst bool, ident runtime.GDIdentType) *GDIRDiscoverable {
	// TODO: ZeroPos is not defined, it should use the position where the node is created
	return &GDIRDiscoverable{isPub, isConst, ident, GDIRBaseNode{scanner.ZeroPos}}
}
