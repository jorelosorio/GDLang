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

	if b.Objects != nil {
		t.Errorf("Buffer should be empty")
	}
}

func TestPushBuffer(t *testing.T) {
	b := runtime.NewGDBuffer()
	buffer := b.NewBuffer()

	if buffer.Parent == nil {
		t.Errorf("Buffer parent should not be nil")
	}

	if buffer.Objects != nil {
		t.Errorf("Buffer should be empty")
	}
}

func TestPopBuffer(t *testing.T) {
	b := runtime.NewGDBuffer()
	buffer := b.NewBuffer()

	nBuffer, objects := buffer.PopAll(runtime.NewGDArrayType(runtime.GDUntypedType))
	if len(objects.Objects) != 0 {
		t.Errorf("Buffer should be empty")
	}

	if nBuffer == nil {
		t.Errorf("Buffer should not be nil")
	}
}

func TestPushBufferObject(t *testing.T) {
	b := runtime.NewGDBuffer()

	if b.Objects != nil {
		t.Errorf("Buffer should be empty")
	}

	b.Push(runtime.GDInt(0))

	if len(b.Objects) != 1 {
		t.Errorf("Buffer should have 1 object")
	}
}

func TestSubBufferPushObject(t *testing.T) {
	b := runtime.NewGDBuffer()
	buffer := b.NewBuffer()

	if buffer.Objects != nil {
		t.Errorf("Buffer should be empty")
	}

	buffer.Push(runtime.GDInt(0))

	if len(buffer.Objects) != 1 {
		t.Errorf("Buffer should have 1 object")
	}

	if buffer.Parent == nil {
		t.Errorf("Buffer parent should not be nil")
	}
}

func TestLastPopObj(t *testing.T) {
	b := runtime.NewGDBuffer()

	if b.Objects != nil {
		t.Errorf("Buffer should be empty")
	}

	b.Push(runtime.GDInt(0))

	if len(b.Objects) != 1 {
		t.Errorf("Buffer should have 1 object")
	}

	obj := b.Pop()
	if obj == nil {
		t.Errorf("Object should not be nil")
	}

	if b.Objects != nil {
		t.Errorf("Buffer should be empty")
	}
}
