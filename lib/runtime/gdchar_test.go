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

func TestCreateChar(t *testing.T) {
	char := runtime.GDChar('a')
	if char != 'a' {
		t.Error("Expected 'a' but got ", char)
	}
}

func TestCreateCharFromString(t *testing.T) {
	char, err := runtime.GDCharFromString("a")
	if err != nil {
		t.Error("Expected no error but got ", err)
	}

	if char != 'a' {
		t.Error("Expected 'a' but got ", char)
	}
}

func TestCreateCharFromInvalidString(t *testing.T) {
	_, err := runtime.GDCharFromString("ab")
	if err == nil {
		t.Error("Expected error but got nil")
	}
}

func TestCharType(t *testing.T) {
	char := runtime.GDChar('a')
	if char.GetType() != runtime.GDCharType {
		t.Error("Expected CharType but got ", char.GetType())
	}
}

func TestCharToString(t *testing.T) {
	char := runtime.GDChar('a')
	if char.ToString() != "a" {
		t.Error("Expected 'a' but got ", char.ToString())
	}
}

func TestCharEqual(t *testing.T) {
	char1 := runtime.GDChar('a')
	char2 := runtime.GDChar('a')
	if !runtime.EqualObjects(char1, char2) {
		t.Error("Expected char1 to be equal to char2")
	}
}

func TestCharNotEqual(t *testing.T) {
	char1 := runtime.GDChar('a')
	char2 := runtime.GDChar('b')
	if runtime.EqualObjects(char1, char2) {
		t.Error("Expected char1 to not be equal to char2")
	}
}

func TestCharNotEqualType(t *testing.T) {
	char := runtime.GDChar('a')
	number := runtime.NewGDIntNumber(1)
	if runtime.EqualObjects(char, number) {
		t.Error("Expected char to not be equal to number")
	}
}

func TestCharNotEqualValue(t *testing.T) {
	char1 := runtime.GDChar('a')
	char2 := runtime.GDChar('b')
	if runtime.EqualObjects(char1, char2) {
		t.Error("Expected char1 to not be equal to char2")
	}
}
