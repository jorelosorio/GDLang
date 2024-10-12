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

type GDTuple struct {
	GDTupleType
	Objects []GDObject
}

// Object interface

func (gd *GDTuple) GetType() GDTypable    { return gd.GDTupleType }
func (gd *GDTuple) GetSubType() GDTypable { return nil }
func (gd *GDTuple) ToString() string {
	if len(gd.Objects) == 0 {
		return "(,)"
	}

	if len(gd.Objects) == 1 {
		return "(" + gd.Objects[0].ToString() + ",)"
	}

	return "(" + JoinSlice(gd.Objects, func(object GDObject, _ int) string {
		switch object.GetType() {
		case GDStringTypeRef:
			return "\"" + object.ToString() + "\""
		case GDCharTypeRef:
			return "'" + object.ToString() + "'"
		}

		return object.ToString()
	}, ", ") + ")"
}
func (gd *GDTuple) CastToType(typ GDTypable) (GDObject, error) {
	switch typ := typ.(type) {
	case GDStringType:
		return GDString(gd.ToString()), nil
	case GDTupleType:
		if len(gd.Objects) != len(typ) {
			return nil, InvalidCastingWrongTypeErr(typ, gd.GetType())
		}

		objects := make([]GDObject, len(gd.Objects))
		for i, obj := range gd.Objects {
			castObj, err := obj.CastToType(typ[i])
			if err != nil {
				return nil, TypeCastingWrongTypeWithHierarchyError(typ, obj.GetType(), err)
			}

			objects[i] = castObj
		}

		return NewGDTupleWithType(typ, objects...), nil
	}

	return nil, InvalidCastingWrongTypeErr(typ, gd.GetType())
}

// Iterable interface

func (gd *GDTuple) Length() int   { return len(gd.Objects) }
func (gd *GDTuple) IsEmpty() bool { return len(gd.Objects) == 0 }
func (gd *GDTuple) GetObjectAt(index int) (GDObject, error) {
	if err := gd.checkIndex(index); err != nil {
		return nil, err
	}

	return gd.Objects[index], nil
}
func (gd *GDTuple) GetObjects() []GDObject { return gd.Objects }

func (gd *GDTuple) checkIndex(index int) error {
	if index < 0 || index >= len(gd.Objects) {
		return IndexOutOfBoundsErr
	}
	return nil
}

func NewGDTuple(elements ...GDObject) *GDTuple {
	tupleTypes := MapSlice(elements, func(gdObject GDObject, _ int) GDTypable {
		return gdObject.GetType()
	})
	return &GDTuple{NewGDTupleType(tupleTypes...), elements}
}

func NewGDTupleWithType(typ GDTupleType, elements ...GDObject) *GDTuple {
	return &GDTuple{typ, elements}
}
