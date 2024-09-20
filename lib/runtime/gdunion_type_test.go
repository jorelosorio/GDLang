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

func TestCanBeAssignToTheLeft(t *testing.T) {
	l := runtime.NewGDUnionType(runtime.GDIntType, runtime.GDFloatType)
	r := runtime.GDFloatType
	if err := runtime.CanBeAssign(l, r, nil); err != nil {
		t.Errorf("Expected %q to be assignable to %q", r.ToString(), l.ToString())
	}

	r = runtime.GDIntType
	if err := runtime.CanBeAssign(l, r, nil); err != nil {
		t.Errorf("Expected %q to be assignable to %q", l.ToString(), r.ToString())
	}
}

func TestWrongTypeAssign(t *testing.T) {
	l := runtime.NewGDUnionType(runtime.GDIntType, runtime.GDFloatType)
	r := runtime.GDStringType
	if err := runtime.CanBeAssign(l, r, nil); err == nil {
		t.Errorf("Expected %q to be not assignable to %q", r.ToString(), l.ToString())
	}

	r = runtime.GDBoolType
	if err := runtime.CanBeAssign(l, r, nil); err == nil {
		t.Errorf("Expected %q to be not assignable to %q", r.ToString(), l.ToString())
	}
}

func TestArrayTypesWithUnionTypes(t *testing.T) {
	l := runtime.NewGDUnionType(runtime.NewGDArrayType(runtime.GDIntType), runtime.NewGDArrayType(runtime.GDFloatType))
	var r runtime.GDTypable = runtime.GDIntType

	if err := runtime.CanBeAssign(l, r, nil); err == nil {
		t.Errorf("Expected %q to be not assignable to %q but got %v", r.ToString(), l.ToString(), err)
	}

	r = runtime.NewGDArrayType(runtime.GDIntType)
	if err := runtime.CanBeAssign(l, r, nil); err != nil {
		t.Errorf("Expected %q to be assignable to %q but got %v", r.ToString(), l.ToString(), err)
	}

	r = runtime.NewGDArrayType(runtime.GDFloatType)
	if err := runtime.CanBeAssign(l, r, nil); err != nil {
		t.Errorf("Expected %q to be assignable to %q but got %v", r.ToString(), l.ToString(), err)
	}
}

func TestComplexUnionTypes(t *testing.T) {
	l := runtime.NewGDArrayType(runtime.NewGDUnionType(runtime.GDIntType, runtime.GDFloatType))
	var r runtime.GDTypable = runtime.GDIntType

	if err := runtime.CanBeAssign(l, r, nil); err == nil {
		t.Errorf("Expected %q to be not assignable to %q", r.ToString(), l.ToString())
	}

	// Push to a type to an array of union types
	if err := runtime.CanBeAssign(l.SubType, r, nil); err != nil {
		t.Errorf("Expected %q to be assignable to %q", r.ToString(), l.ToString())
	}
}

func TestTwoUnion(t *testing.T) {
	l := runtime.NewGDUnionType(runtime.GDIntType, runtime.GDFloatType)
	r := runtime.NewGDUnionType(runtime.GDIntType, runtime.GDFloatType)

	if err := runtime.CanBeAssign(l, r, nil); err != nil {
		t.Errorf("Expected %q to be assignable to %q", r.ToString(), l.ToString())
	}

	r = runtime.NewGDUnionType(runtime.GDFloatType, runtime.GDIntType)
	if err := runtime.CanBeAssign(l, r, nil); err != nil {
		t.Errorf("Expected %q to be assignable to %q", r.ToString(), l.ToString())
	}
}

func TestNoUnionWithUnion(t *testing.T) {
	l := runtime.GDIntType
	r := runtime.NewGDUnionType(runtime.GDIntType, runtime.GDStringType)

	if err := runtime.CanBeAssign(l, r, nil); err == nil {
		t.Errorf("Expected %q to be not assignable to %q", r.ToString(), l.ToString())
	}
}
