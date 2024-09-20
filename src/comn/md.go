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

import "strings"

type MdScanner struct {
	src string
}

func NewMarkdown(src string) *MdScanner {
	return &MdScanner{src}
}

func (m *MdScanner) scanBacktick(i int) (string, int) {
	offs := i
	for ; i < len(m.src); i++ {
		ch := m.src[i]
		if ch == '`' {
			return m.src[offs:i], i
		}
	}

	// No backtick found
	lstr := m.src[offs:]
	return lstr, i
}

func (m *MdScanner) Stylize() string {
	if m.src == "" {
		return ""
	}

	var result strings.Builder
	defer result.Reset()

	var buffer strings.Builder
	defer buffer.Reset()

	i := 0
	for ; i < len(m.src); i++ {
		ch := rune(m.src[i])
		switch ch {
		case '`':
			if buffer.Len() > 0 {
				result.WriteString(Colorize(PrimaryTextColor, buffer.String()))
				buffer.Reset()
			}

			i++ // skip the backtick
			str, w := m.scanBacktick(i)
			if w > i {
				result.WriteString(Colorize(PrimaryHighlightColor, str))
			}
			i = w
		default:
			buffer.WriteRune(ch)
		}
	}

	if buffer.Len() > 0 {
		result.WriteString(Colorize(PrimaryTextColor, buffer.String()))
	}

	return result.String()
}
