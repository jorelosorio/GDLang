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

func TestContainsAny(t *testing.T) {
	testCases := []struct {
		target string
		list   []string
		want   bool
	}{
		{"hola", []string{"hola"}, true},
		{"hola", []string{"hola", "adios"}, true},
		{"hola", []string{"adios"}, false},
		{"hola", []string{}, false},
		{"hola", []string{"adios", "hola"}, true},
		{"hola", []string{"adios", "adios"}, false},
		{"hola", []string{"hola", "adios", "hello"}, true},
		{"hola", []string{"adios", "hello"}, false},
		{"hola", []string{"hello", "hola", "adios"}, true},
		{"hola", []string{"hello", "adios", "hello"}, false},
		{"hola", []string{"hello", "adios", "hello", "hola"}, true},
		{"hola", []string{"hello", "adios", "hello", "adios"}, false},
	}

	for _, tc := range testCases {
		if runtime.ContainsAny(tc.target, tc.list...) != tc.want {
			t.Errorf("runtime.ContainsAny(%q, %q) = %v, want %v", tc.target, tc.list, !tc.want, tc.want)
		}
	}
}

func TestContainsAll(t *testing.T) {
	testCases := []struct {
		target string
		list   []string
		want   bool
	}{
		{"hola", []string{"hola"}, true},
		{"hola", []string{"hola", "adios"}, false},
		{"hola", []string{"adios"}, false},
		{"hola", []string{}, false},
		{"hola", []string{"adios", "hola"}, false},
		{"hola", []string{"adios", "adios"}, false},
		{"hola", []string{"hola", "adios", "hello"}, false},
		{"hola", []string{"adios", "hello"}, false},
		{"hola", []string{"hello", "hola", "adios"}, false},
		{"hola", []string{"hello", "adios", "hello"}, false},
		{"hola", []string{"hello", "adios", "hello", "hola"}, false},
		{"hola", []string{"hello", "adios", "hello", "adios"}, false},
		{"lorem ipsum dolor sit amet", []string{"lorem", "ipsum", "dolor", "sit", "amet"}, true},
		{"lorem ipsum dolor sit amet", []string{"ipsum", "dolor", "sit", "lorem", "amet"}, true},
	}

	for _, tc := range testCases {
		if runtime.ContainsAll(tc.target, tc.list...) != tc.want {
			t.Errorf("ContainsAll(%q, %q) = %v, want %v", tc.target, tc.list, !tc.want, tc.want)
		}
	}
}
