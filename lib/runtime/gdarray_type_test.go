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

func TestSimpleArrayType(t *testing.T) {
	arrayType := runtime.NewGDArrayType(runtime.GDStringType)

	strRepresentation := "[string]"
	if arrayType.ToString() != strRepresentation {
		t.Errorf("Expected %s but got %s", strRepresentation, arrayType.ToString())
	}
}

func TestSubArrayType(t *testing.T) {
	subArray := runtime.NewGDArrayType(runtime.GDStringType)
	arrayType := runtime.NewGDArrayType(subArray)

	strRepresentation := "[[string]]"
	if arrayType.ToString() != strRepresentation {
		t.Errorf("Expected %s but got %s", strRepresentation, arrayType.ToString())
	}
}

func TestComplexArrayType(t *testing.T) {
	arrayType := runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.NewGDArrayType(runtime.GDStringType))))

	strRepresentation := "[[[[string]]]]"
	if arrayType.ToString() != strRepresentation {
		t.Errorf("Expected %s but got %s", strRepresentation, arrayType.ToString())
	}
}
