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

package runtime

type GDBuffer struct {
	Parent  *GDBuffer
	Objects []GDObject
}

func (b *GDBuffer) NewBuffer() *GDBuffer {
	buffer := &GDBuffer{b, nil}
	return buffer
}

func (b *GDBuffer) PopAll(gdType *GDArrayType) (*GDBuffer, *GDArray) {
	objects := b.Objects
	buffer := b.Parent

	b.Parent = nil
	b.Objects = nil

	buffArray := NewGDArrayWithTypeAndObjects(gdType, objects)

	return buffer, buffArray
}

func (b *GDBuffer) Push(obj GDObject) error {
	if b.Objects == nil {
		b.Objects = make([]GDObject, 0)
	}

	b.Objects = append(b.Objects, obj)

	return nil
}

func (b *GDBuffer) Pull() GDObject {
	if b.Objects == nil {
		return GDZNil
	}

	if len(b.Objects) > 0 {
		// Get and remove the first element
		obj := b.Objects[0]
		if len(b.Objects) == 1 {
			b.Objects = nil
		} else {
			b.Objects = b.Objects[1:]
		}

		return obj
	}

	return nil
}

func (b *GDBuffer) Pop() GDObject {
	if b.Objects == nil {
		return GDZNil
	}

	if len(b.Objects) > 0 {
		// Get and remove the last element
		obj := b.Objects[len(b.Objects)-1]
		if len(b.Objects) == 1 {
			b.Objects = nil
		} else {
			b.Objects = b.Objects[:len(b.Objects)-1]
		}

		return obj
	}

	return nil
}

func NewGDBuffer() *GDBuffer {
	return &GDBuffer{
		Objects: nil,
		Parent:  nil,
	}
}
