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

func TestNilComparison(t *testing.T) {
	for _, test := range []struct {
		compare runtime.GDObject
	}{
		{runtime.NewGDIntNumber(0)},
		{runtime.GDInt8(0)},
		{runtime.GDString("")},
		{runtime.GDBool(false)},
		{runtime.NewGDComplexNumber(0)},
		{runtime.NewGDFloatNumber(0)},
	} {
		if runtime.GDZNil == test.compare {
			t.Errorf("Expected nil != %v", test.compare)
		}
	}
}

func TestNilToString(t *testing.T) {
	if runtime.GDZNil.ToString() != "nil" {
		t.Errorf("Expected nil, got %v", runtime.GDZNil.ToString())
	}
}

func TestNilType(t *testing.T) {
	if runtime.GDZNil.GetType() != runtime.GDNilType {
		t.Errorf("Expected nil, got %v", runtime.GDZNil.GetType())
	}
}

func TestNilOnlyEqualsToNil(t *testing.T) {
	objNil := runtime.GDZNil
	if objNil != runtime.GDZNil {
		t.Errorf("Expected nil == nil")
	}
}
