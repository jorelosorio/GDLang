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
	objects []GDObject
	*GDArrayType
}

func (gd *GDArray) GetType() GDTypable    { return gd.GDArrayType }
func (gd *GDArray) GetSubType() GDTypable { return nil }
func (gd *GDArray) ToString() string {
	vals := JoinSlice(gd.objects, func(object GDObject, _ int) string {
		return ConvertObjectToString(object)
	}, ", ")

	return Sprintf("[%@]", vals)
}
func (gd *GDArray) CastToType(typ GDTypable) (GDObject, error) {
	switch typ := typ.(type) {
	case GDStringType:
		return GDString(gd.ToString()), nil
	case *GDArrayType:
		objects := make([]GDObject, len(gd.objects))
		for i, obj := range gd.objects {
			castObj, err := obj.CastToType(typ.SubType)
			if err != nil {
				return nil, TypeCastingWrongTypeWithHierarchyError(typ, obj.GetType(), err)
			}

			objects[i] = castObj
		}

		return NewGDArrayWithTypeAndObjects(typ, objects), nil
	}

	return nil, InvalidCastingWrongTypeErr(typ, gd.GetType())
}

func (gd *GDArray) Length() int   { return len(gd.objects) }
func (gd *GDArray) IsEmpty() bool { return len(gd.objects) == 0 }
func (gd *GDArray) GetObjectAt(index int) (GDObject, error) {
	if err := gd.checkIndex(index); err != nil {
		return nil, err
	}

	return gd.objects[index], nil
}
func (gd *GDArray) GetObjects() []GDObject { return gd.objects }

func (gd *GDArray) Dispose() { gd.objects = nil }
func (gd *GDArray) AddObject(object GDObject, stack *GDStack) error {
	err := CanBeAssign(gd.SubType, object.GetType(), stack)
	if err != nil {
		return err
	}

	gd.objects = append(gd.objects, object)

	return nil
}
func (gd *GDArray) AddObjects(stack *GDStack, objects ...GDObject) error {
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

	obj := gd.objects[index]
	gd.objects = append(gd.objects[:index], gd.objects[index+1:]...)

	return obj, nil
}
func (gd *GDArray) Set(index int, object GDObject, stack *GDStack) error {
	if err := gd.checkIndex(index); err != nil {
		return err
	}

	err := CanBeAssign(gd.objects[index].GetType(), object.GetType(), stack)
	if err != nil {
		return err
	}

	gd.objects[index] = object

	return nil
}

func (gd *GDArray) checkIndex(index int) error {
	if index < 0 || index >= len(gd.objects) {
		return IndexOutOfBoundsErr
	}
	return nil
}

func NewGDArray(stack *GDStack, objects ...GDObject) *GDArray {
	typ := NewGDArrayType(ComputeTypeFromObjects(objects, stack))
	return &GDArray{
		objects: objects,
		// Determine array type by computing the object types.
		GDArrayType: typ,
	}
}

func NewGDEmptyArray() *GDArray {
	typ := NewGDEmptyArrayType()
	return &GDArray{
		GDArrayType: typ,
		objects:     []GDObject{},
	}
}

func NewGDArrayWithType(subType GDTypable) *GDArray {
	typ := NewGDArrayType(subType)
	return &GDArray{
		GDArrayType: typ,
		objects:     make([]GDObject, 0),
	}
}

func NewGDArrayWithSubTypeAndObjects(subType GDTypable, values []GDObject, stack *GDStack) (*GDArray, error) {
	arrayType := NewGDArrayType(subType)
	array := GDArray{
		GDArrayType: arrayType,
		objects:     make([]GDObject, 0),
	}

	for _, value := range values {
		err := array.AddObject(value, stack)
		if err != nil {
			return nil, err
		}
	}

	return &array, nil
}

// Only used when type and array elements are pre-computed.
func NewGDArrayWithTypeAndObjects(arrayType *GDArrayType, objects []GDObject) *GDArray {
	return &GDArray{
		GDArrayType: arrayType,
		objects:     objects,
	}
}

func NewGDArrayWithObjects(objects []GDObject, stack *GDStack) (*GDArray, error) {
	return NewGDArrayWithSubTypeAndObjects(ComputeTypeFromObjects(objects, stack), objects, stack)
}
