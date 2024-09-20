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
	"gdlang/src/gd/scanner"
)

type GDSourceMap struct {
	Version  byte          `json:"version"`
	Sources  []string      `json:"sources"`
	Mappings map[int][]int `json:"mappings"`

	fileCount int
	files     map[string]int
}

func (mf *GDSourceMap) AddMapping(bytecodeOffset int, pos scanner.Position) {
	if pos == scanner.ZeroPos {
		return
	}

	if filePos, ok := mf.files[pos.Filename]; ok {
		mf.Mappings[bytecodeOffset] = []int{filePos, pos.Line, pos.ColStart, pos.ColEnd}
	} else {
		mf.files[pos.Filename] = mf.fileCount
		mf.Sources = append(mf.Sources, pos.Filename)
		mf.fileCount++
	}
}

func NewGDIRSourceMap() *GDSourceMap {
	return &GDSourceMap{Version: 1, Mappings: make(map[int][]int), Sources: make([]string, 0), files: make(map[string]int)}
}
