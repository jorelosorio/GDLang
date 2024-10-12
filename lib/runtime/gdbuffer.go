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
	objects []GDObject
	isEmpty bool
}

func (b *GDBuffer) Push(obj GDObject) {
	b.objects = append(b.objects, obj)
	if b.isEmpty {
		b.isEmpty = false
	}
}

func (b *GDBuffer) Pop() GDObject {
	if !b.isEmpty {
		// Get and remove the last element
		obj := b.objects[len(b.objects)-1]
		if len(b.objects) == 1 {
			b.objects = make([]GDObject, 0)
			b.isEmpty = true
		} else {
			b.objects = b.objects[:len(b.objects)-1]
		}

		return obj
	}

	return GDZNil
}

func (b *GDBuffer) IsEmpty() bool { return b.isEmpty }

func NewGDBuffer() *GDBuffer { return &GDBuffer{make([]GDObject, 0), true} }
