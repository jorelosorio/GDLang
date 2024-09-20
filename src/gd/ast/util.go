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

func GetStartEndPosition[T Node](nodes []T) scanner.Position {
	if len(nodes) == 0 {
		return scanner.ZeroPos
	}

	if len(nodes) == 1 {
		return nodes[0].GetPosition()
	}

	return buildStartEndPos(nodes[0], nodes[len(nodes)-1])
}

func buildStartEndPos(firstItem, lastItem Node) scanner.Position {
	var fp, lp scanner.Position = scanner.ZeroPos, scanner.ZeroPos

	if firstItem != nil {
		fp = firstItem.GetPosition()
	}

	if lastItem != nil {
		lp = lastItem.GetPosition()
	}

	return scanner.Position{
		Filename: fp.Filename,
		Line:     fp.Line,
		ColStart: fp.ColStart,
		ColEnd:   lp.ColEnd,
	}
}
