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

package comn

import (
	"fmt"
)

type Color string

var (
	PrimaryHighlightColor Color = "\x1b[32;1m"
	WarningHighlightColor Color = "\x1b[33;1m"
	ErrorHighlightColor   Color = "\x1b[32;1m"
	TitleColor            Color = "\x1b[42;30;2m"
	PrimaryTextColor      Color = "\x1b[32;2m"
	Reset                 Color = "\x1b[0m"
)

func Colorize(color Color, s string) string {
	return string(color) + s + string(Reset)
}

func Colorizef(color Color, format string, a ...any) string {
	return string(color) + fmt.Sprintf(format, a...) + string(Reset)
}
