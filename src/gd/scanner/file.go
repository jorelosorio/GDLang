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
	"sort"
	"sync"
)

type File struct {
	name string // file name as provided to AddFile
	base int    // Pos value range for this file is [base...base+size]
	size int    // file size as provided to AddFile

	// lines and infos are protected by mutex
	mutex sync.Mutex
	lines []int // lines contains the offset of the first character for each line (the first entry is always 0)
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Size() int {
	return f.size
}

func (f *File) Pos(offset int) Pos {
	return Pos(f.base + f.fixOffset(offset))
}

func (f *File) Offset(p Pos) int {
	return f.fixOffset(int(p) - f.base)
}

func (f *File) unpack(offset int) (filename string, line, column int) {
	f.mutex.Lock()
	filename = f.name
	if i := searchInts(f.lines, offset); i >= 0 {
		line, column = i+1, offset-f.lines[i]+1
	}
	// Follow up https://go.dev/issue/38471 to use `defer f.mutex.Unlock()` here.
	f.mutex.Unlock()
	return
}

func (f *File) Position(p Pos) (pos Position) {
	return f.PositionFor(p)
}

func (f *File) PositionFrom(start Pos, end Pos) (pos Position) {
	spos := f.PositionFor(start)
	// A -1 is required to get the correct column end,
	// Due to the way the scanner works, the end position is always one byte ahead.
	if end <= start {
		end = start
	} else {
		end--
	}
	spos.ColEnd = spos.ColStart + (int(end) - int(start))
	return spos
}

func (f *File) PositionFor(p Pos) (pos Position) {
	if p != NoPos {
		pos = f.position(p)
	}
	return
}

// AddLine adds the line offset for a new line.
// The line offset must be larger than the offset for the previous line
// and smaller than the file size; otherwise the line offset is ignored.
func (f *File) AddLine(offset int) {
	f.mutex.Lock()
	if i := len(f.lines); (i == 0 || f.lines[i-1] < offset) && offset < f.size {
		f.lines = append(f.lines, offset)
	}
	f.mutex.Unlock()
}

func (f *File) position(p Pos) (pos Position) {
	offset := f.fixOffset(int(p) - f.base)
	pos.Filename, pos.Line, pos.ColStart = f.unpack(offset)
	return
}

// fixOffset fixes an out-of-bounds offset such that 0 <= offset <= f.size.
func (f *File) fixOffset(offset int) int {
	switch {
	case offset < 0:
		return 0
	case offset > f.size:
		return f.size
	default:
		return offset
	}
}

func searchInts(a []int, x int) int {
	// Follow up changes src/go/token/position.go
	return sort.Search(len(a), func(i int) bool { return a[i] > x }) - 1
}
