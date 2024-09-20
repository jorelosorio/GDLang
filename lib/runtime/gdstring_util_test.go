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

func TestIsStringWithDifferentTypes(t *testing.T) {
	for i, test := range []struct {
		value any
		want  bool
	}{
		{"", true},
		{runtime.GDString(""), true},
		{1, false},
		{1.0, false},
		{runtime.GDInt(1), false},
		{runtime.GDInt8(1), false},
		{runtime.GDInt16(1), false},
		{runtime.GDFloat32(1.0), false},
		{runtime.GDFloat64(1.0), false},
		{runtime.GDComplex64(1.0), false},
		{runtime.GDComplex128(1.0), false},
		{'a', false},
	} {
		resultValue := runtime.IsString(test.value)
		if resultValue != test.want {
			t.Errorf("IsString(%v) = %v, want %v for test case %d", test.value, resultValue, test.want, i+1)
		}
	}
}
