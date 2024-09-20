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

type Pos int

const NoPos Pos = 0

var ZeroPos Position

func (p Pos) IsValid() bool {
	return p != NoPos
}

type Position struct {
	Filename string // filename, if any
	Line     int    // line number, starting at 1
	ColStart int    // column end, starting at 1 (byte count)
	ColEnd   int    // column start, starting at 1 (byte count)
}

func (pos *Position) IsValid() bool {
	return pos.Line > 0
}

func NZeroPostAt(fileName string) Position {
	return Position{
		Filename: fileName,
		Line:     0,
		ColStart: 0,
		ColEnd:   0,
	}
}
