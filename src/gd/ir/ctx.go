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

type GDIRLabelMark struct {
	Offset int
	Label  runtime.GDIdent
}

type GDIRContext struct {
	SrcMap       *GDSourceMap
	LabelOffsets map[any]uint16
	Marks        []GDIRLabelMark
}

func (c *GDIRContext) ResolveMarks(bytecode *bytes.Buffer) {
	pendingMarks := make([]GDIRLabelMark, 0)
	for _, mark := range c.Marks {
		ok := c.ResolveMark(mark, bytecode)
		if !ok {
			pendingMarks = append(pendingMarks, mark)
		}
	}

	c.Marks = pendingMarks
}

func (c *GDIRContext) ResolveMark(mark GDIRLabelMark, bytecode *bytes.Buffer) bool {
	offsetValue, ok := c.LabelOffsets[mark.Label.GetRawValue()]
	if !ok {
		return false
	}

	WriteUInt16At(mark.Offset, bytecode, offsetValue)

	return true
}

func (c *GDIRContext) AddMark(bytecode *bytes.Buffer, offset int, ident runtime.GDIdent) bool {
	mark := GDIRLabelMark{offset, ident}
	ok := c.ResolveMark(mark, bytecode)
	// It was not possible to find a label yet, so add it to the pending list
	// to be resolved later, when a label is defined.
	if !ok {
		c.Marks = append(c.Marks, mark)
	}

	return ok
}

func (c *GDIRContext) AddLabel(bytecode *bytes.Buffer, offset uint16, label runtime.GDIdent) {
	if _, ok := c.LabelOffsets[label.GetRawValue()]; ok {
		panic(fmt.Sprintf("Label %s already defined", label))
	}

	c.LabelOffsets[label.GetRawValue()] = offset

	// Resolve pending marks
	c.ResolveMarks(bytecode)
}

func (c *GDIRContext) AddMapping(bytecode *bytes.Buffer, pos scanner.Position) {
	c.SrcMap.AddMapping(bytecode.Len(), pos)
}

func NewGDIRContext() *GDIRContext {
	return &GDIRContext{NewGDIRSourceMap(), make(map[any]uint16), make([]GDIRLabelMark, 0)}
}
