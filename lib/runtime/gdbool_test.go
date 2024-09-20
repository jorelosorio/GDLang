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

func TestBoolTrueValue(t *testing.T) {
	b := runtime.GDBool(true)

	if b != true {
		t.Error("Wrong value")
	}
}

func TestBoolFalseValue(t *testing.T) {
	b := runtime.GDBool(false)

	if b != false {
		t.Error("Wrong value")
	}
}

func TestBoolToStringTrue(t *testing.T) {
	b := runtime.GDBool(true)

	if b.ToString() != "true" {
		t.Error("Wrong string")
	}
}

func TestBoolToStringFalse(t *testing.T) {
	b := runtime.GDBool(false)

	if b.ToString() != "false" {
		t.Error("Wrong string")
	}
}

func TestBoolFromStringTrue(t *testing.T) {
	b, err := runtime.NewGDBoolFromString("true")

	if err != nil {
		t.Errorf("Error while parsing bool: %v", err)
	}

	if b != true {
		t.Error("Wrong bool value")
	}
}

func TestBoolFromStringFalse(t *testing.T) {
	b, err := runtime.NewGDBoolFromString("false")

	if err != nil {
		t.Error("Error")
	}

	if b != false {
		t.Error("Wrong bool value")
	}
}

func TestBoolFromStringInvalid(t *testing.T) {
	_, err := runtime.NewGDBoolFromString("invalid")

	if err == nil {
		t.Error("Expect error but got nil")
	}
}

func TestIsBoolWithDifferentTypes(t *testing.T) {
	testCases := []struct {
		value any
		want  bool
	}{
		{true, true},
		{false, true},
		{runtime.GDBool(true), true},
		{runtime.GDBool(false), true},
		{1, false},
		{1.0, false},
		{"true", false},
	}

	for _, tc := range testCases {
		if runtime.IsBool(tc.value) != tc.want {
			t.Errorf("IsBool(%v) = %v, want %v", tc.value, !tc.want, tc.want)
		}
	}
}
