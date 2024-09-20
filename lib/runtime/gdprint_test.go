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

package runtime_test

import (
	"gdlang/lib/runtime"
	"testing"
)

func TestPrintf(t *testing.T) {
	tests := []struct {
		format   string
		args     []any
		expected string
	}{
		{"", []any{}, ""},
		{"hola", []any{}, "hola"},
		{"%@", []any{"hi"}, "hi"},
		{"hola %@", []any{"mundo"}, "hola mundo"},
		{"%@ %@", []any{"hola", "mundo"}, "hola mundo"},
		{"%@ %@ %@", []any{"hola"}, "hola %@ %@"},
		{"%@ %@ %@", []any{"hola", "mundo"}, "hola mundo %@"},
		{"%", []any{"hola"}, "%"},
		{"% @", []any{}, "% @"},
		{"@%@", []any{"hola"}, "@hola"},
		{"@%@@", []any{"hola"}, "@hola@"},
	}

	for _, test := range tests {
		result := runtime.Sprintf(test.format, test.args...)
		if result != test.expected {
			t.Errorf("Expected %q but got %q", test.expected, result)
		}
	}
}
