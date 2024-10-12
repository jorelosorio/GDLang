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

func TestBufferShouldBeEmpty(t *testing.T) {
	b := runtime.NewGDBuffer()
	if !b.IsEmpty() {
		t.Errorf("Buffer should be empty")
	}
}

func TestPopBuffer(t *testing.T) {
	b := runtime.NewGDBuffer()
	b.Push(runtime.GDInt(0))

	if b.IsEmpty() {
		t.Errorf("Buffer should have 1 object")
	}

	obj := b.Pop()
	if obj == runtime.GDZNil {
		t.Errorf("Object should not be nil")
	}

	if !b.IsEmpty() {
		t.Errorf("Buffer should be empty")
	}

	if obj != runtime.GDInt(0) {
		t.Errorf("Object should be an integer")
	}
}

func TestPushBufferObject(t *testing.T) {
	b := runtime.NewGDBuffer()
	b.Push(runtime.GDInt(0))

	if b.IsEmpty() {
		t.Errorf("Buffer should have 1 object")
	}
}

func TestPopEmptyBuffer(t *testing.T) {
	b := runtime.NewGDBuffer()
	obj := b.Pop()

	if obj != runtime.GDZNil {
		t.Errorf("Object should be nil")
	}
}
