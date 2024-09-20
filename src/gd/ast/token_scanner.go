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
	"gdlang/src/comn"
	"strings"
)

func ScanTokens(msg string) []string {
	tokens := []string{}
	var buffer strings.Builder

	for _, ch := range msg {
		switch {
		case (ch == ' ' || ch == ',') && buffer.Len() > 0:
			bufferValue := mapToken(buffer.String())
			tokens = append(tokens, bufferValue)
			buffer.Reset()
		case ch == 'L' || ch == '$':
			buffer.WriteRune(ch)
		default:
			if buffer.Len() > 0 {
				buffer.WriteRune(ch)
			}
		}
	}

	if buffer.Len() > 0 {
		bufferValue := mapToken(buffer.String())
		tokens = append(tokens, bufferValue)
		buffer.Reset()
	}

	return tokens
}

func mapToken(token string) string {
	switch token {
	case "$end":
		return comn.WrongEndOfStatementErrMsg
	case "LIDENT":
		return "identifier"
	default:
		return LookUpToken(token)
	}
}
