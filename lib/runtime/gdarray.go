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

type GDArray struct {
	Objects []GDObject
	*GDArrayType
}

func (gd *GDArray) GetType() GDTypable    { return gd.GDArrayType }
func (gd *GDArray) GetSubType() GDTypable { return nil }
func (gd *GDArray) ToString() string {
	vals := JoinSlice(gd.Objects, func(object GDObject, _ int) string {
		return ObjectToStringForInternalData(object)
	}, ", ")

	return Sprintf("[%@]", vals)
}
func (gd *GDArray) CastToType(typ GDTypable, stack *GDSymbolStack) (GDObject, error) {
	switch typ := typ.(type) {
	case GDType:
		switch typ {
		case GDStringType:
			return GDString(gd.ToString()), nil
		}
	case *GDArrayType:
		for i, obj := range gd.Objects {
			castObj, err := obj.CastToType(typ.SubType, stack)
			if err != nil {
				return nil, TypeCastingWrongTypeWithHierarchyError(typ, obj.GetType(), err)
			}

			gd.Objects[i] = castObj
		}

		gd.GDArrayType = typ

		return gd, nil
	}

	return nil, InvalidCastingWrongTypeErr(typ, gd.GetType())
}

func (gd *GDArray) Length() int   { return len(gd.Objects) }
func (gd *GDArray) IsEmpty() bool { return len(gd.Objects) == 0 }
func (gd *GDArray) Get(index int) (GDObject, error) {
	if err := gd.checkIndex(index); err != nil {
		return nil, err
	}

	return gd.Objects[index], nil
}
func (gd *GDArray) GetObjects() []GDObject { return gd.Objects }

func (gd *GDArray) Dispose() { gd.Objects = nil }
func (gd *GDArray) AddObject(object GDObject, stack *GDSymbolStack) error {
	err := CanBeAssign(gd.SubType, object.GetType(), stack)
	if err != nil {
		return err
	}

	gd.Objects = append(gd.Objects, object)

	return nil
}
func (gd *GDArray) AddObjects(objects []GDObject, stack *GDSymbolStack) error {
	for _, object := range objects {
		if err := gd.AddObject(object, stack); err != nil {
			return err
		}
	}

	return nil
}
func (gd *GDArray) Remove(index int) (GDObject, error) {
	if err := gd.checkIndex(index); err != nil {
		return nil, err
	}

	obj := gd.Objects[index]
	gd.Objects = append(gd.Objects[:index], gd.Objects[index+1:]...)

	return obj, nil
}
func (gd *GDArray) Set(index int, object GDObject, stack *GDSymbolStack) error {
	if err := gd.checkIndex(index); err != nil {
		return err
	}

	err := CanBeAssign(gd.Objects[index].GetType(), object.GetType(), stack)
	if err != nil {
		return err
	}

	gd.Objects[index] = object

	return nil
}

func (gd *GDArray) checkIndex(index int) error {
	if index < 0 || index >= len(gd.Objects) {
		return IndexOutOfBoundsErr
	}
	return nil
}

func NewGDArray(values ...GDObject) *GDArray {
	return &GDArray{
		Objects: values,
		// Determine array type by computing the object types.
		GDArrayType: NewGDArrayType(ComputeTypeFromObjects(values)),
	}
}

func NewGDEmptyArray() *GDArray {
	return &GDArray{
		Objects:     []GDObject{},
		GDArrayType: NewGDEmptyArrayType(),
	}
}

func NewGDArrayWithType(subType GDTypable) *GDArray {
	return &GDArray{
		GDArrayType: NewGDArrayType(subType),
		Objects:     make([]GDObject, 0),
	}
}

func NewGDArrayWithSubTypeAndObjects(subType GDTypable, values []GDObject, stack *GDSymbolStack) (*GDArray, error) {
	arrayType := NewGDArrayType(subType)
	array := GDArray{
		GDArrayType: arrayType,
		Objects:     make([]GDObject, len(values)),
	}

	for i, value := range values {
		err := CanBeAssign(arrayType.SubType, value.GetType(), stack)
		if err != nil {
			return nil, WrongTypesErr(arrayType.SubType, value.GetType())
		}

		array.Objects[i] = value
	}

	return &array, nil
}

// Only used when type and array elements are pre-computed.
func NewGDArrayWithTypeAndObjects(arrayType *GDArrayType, objects []GDObject) *GDArray {
	return &GDArray{
		GDArrayType: arrayType,
		Objects:     objects,
	}
}

func NewGDArrayWithObjects(objects []GDObject, stack *GDSymbolStack) (*GDArray, error) {
	return NewGDArrayWithSubTypeAndObjects(ComputeTypeFromObjects(objects), objects, stack)
}
