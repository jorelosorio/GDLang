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

package ir

import (
	"bytes"
	"gdlang/lib/runtime"
	"gdlang/src/cpu"
	"unsafe"
)

func Write(bytecode *bytes.Buffer, data ...any) error {
	for _, value := range data {
		var err error
		switch value := value.(type) {
		case bool:
			err = WriteBool(bytecode, value)
		case cpu.GDInst:
			err = bytecode.WriteByte(byte(value))
		case runtime.GDObject:
			err = writeObjectWithType(bytecode, value)
		case runtime.GDIdent:
			err = WriteIdent(bytecode, value)
		case runtime.GDTypable:
			err = WriteType(bytecode, value)
		default:
			panic("Unknown type")
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func WriteType(bytecode *bytes.Buffer, typ runtime.GDTypable) error {
	err := bytecode.WriteByte(byte(typ.GetCode()))
	if err != nil {
		return err
	}
	switch t := typ.(type) {
	case runtime.GDStructType:
		err := WriteInt8(bytecode, int8(len(t)))
		if err != nil {
			return err
		}

		for _, attr := range t {
			err := WriteIdent(bytecode, attr.Ident)
			if err != nil {
				return err
			}

			err = WriteType(bytecode, attr.Type)
			if err != nil {
				return err
			}
		}
	case runtime.GDIdent:
		// Write the ident mode and raw value
		return WriteIdent(bytecode, t)
	case *runtime.GDArrayType:
		return WriteType(bytecode, t.SubType)
	case runtime.GDUnionType:
		err := WriteInt8(bytecode, int8(len(t)))
		if err != nil {
			return err
		}

		for _, field := range t {
			err := WriteType(bytecode, field)
			if err != nil {
				return err
			}
		}
	case runtime.GDTupleType:
		err := WriteInt8(bytecode, int8(len(t)))
		if err != nil {
			return err
		}
		for _, tItemType := range t {
			err := WriteType(bytecode, tItemType)
			if err != nil {
				return err
			}
		}
	case *runtime.GDLambdaType:
		err = WriteType(bytecode, t.ReturnType)
		if err != nil {
			return err
		}

		err = WriteInt8(bytecode, int8(len(t.ArgTypes)))
		if err != nil {
			return err
		}

		for _, argType := range t.ArgTypes {
			err = WriteIdent(bytecode, argType.Key)
			if err != nil {
				return err
			}
			err = WriteType(bytecode, argType.Value)
			if err != nil {
				return err
			}
		}

		return WriteBool(bytecode, t.IsVariadic)
	}

	return nil
}

func writeObject(bytecode *bytes.Buffer, obj runtime.GDObject) error {
	var err error
	switch obj := obj.(type) {
	case runtime.GDInt:
		err = WriteInt(bytecode, int(obj))
	case runtime.GDInt8:
		err = WriteInt8(bytecode, int8(obj))
	case runtime.GDInt16:
		err = WriteInt16(bytecode, int16(obj))
	case runtime.GDFloat32:
		err = WriteFloat32(bytecode, float32(obj))
	case runtime.GDFloat64:
		err = WriteFloat64(bytecode, float64(obj))
	case runtime.GDComplex64:
		err = WriteComplex64(bytecode, complex64(obj))
	case runtime.GDComplex128:
		err = WriteComplex128(bytecode, complex128(obj))
	case runtime.GDBool:
		err = WriteBool(bytecode, bool(obj))
	case runtime.GDString:
		err = WriteString(bytecode, string(obj))
	case runtime.GDChar:
		err = WriteChar(bytecode, rune(obj))
	case *runtime.GDTuple:
		err = WriteByte(bytecode, byte(len(obj.Objects)))
		if err != nil {
			return err
		}

		for _, item := range obj.Objects {
			err := writeObjectWithType(bytecode, item)
			if err != nil {
				return err
			}
		}
	case *runtime.GDStruct:
		err = WriteByte(bytecode, byte(len(obj.Type)))
		if err != nil {
			return err
		}

		for _, item := range obj.Type {
			obj, err := obj.GetAttr(item.Ident)
			if err != nil {
				return err
			}

			err = writeObjectWithType(bytecode, obj.Object)
			if err != nil {
				return err
			}
		}
	case *runtime.GDSpreadable:
		err = writeObjectWithType(bytecode, obj.Iterable)
	case *runtime.GDArray:
		err = WriteByte(bytecode, byte(len(obj.Objects)))
		if err != nil {
			return err
		}

		for _, item := range obj.Objects {
			err := writeObjectWithType(bytecode, item)
			if err != nil {
				return err
			}
		}
	}

	return err
}

func writeObjectType(bytecode *bytes.Buffer, obj runtime.GDObject) error {
	var err error
	if typ := obj.GetSubType(); typ != nil {
		err = WriteType(bytecode, typ)
	} else {
		err = WriteType(bytecode, obj.GetType())
	}

	if err != nil {
		return err
	}

	return nil
}

func writeObjectWithType(bytecode *bytes.Buffer, obj runtime.GDObject) error {
	err := writeObjectType(bytecode, obj)
	if err != nil {
		return err
	}

	err = writeObject(bytecode, obj)
	if err != nil {
		return err
	}

	return nil
}

// Specific types

func WriteUInt16At(off int, bytecode *bytes.Buffer, value uint16) {
	ui32bytes := UI16toBytes(value)
	bytes := bytecode.Bytes()
	for _, intByte := range ui32bytes {
		bytes[off] = intByte
		off++
	}
}

func WriteChar(bytecode *bytes.Buffer, val rune) error {
	_, err := bytecode.WriteRune(val)
	return err
}

func WriteByte(bytecode *bytes.Buffer, val byte) error {
	return bytecode.WriteByte(val)
}

// Ident are 256 characters long, to fit in a byte
func WriteIdent(bytecode *bytes.Buffer, ident runtime.GDIdent) error {
	// Write the ident mode
	err := WriteByte(bytecode, byte(ident.GetMode()))
	if err != nil {
		return err
	}

	// Write the raw value
	switch mode := ident.GetMode(); mode {
	case runtime.GDByteIdentMode:
		return WriteByte(bytecode, byte(ident.GetRawValue().(byte)))
	case runtime.GDUInt16IdentMode:
		return WriteUInt16(bytecode, ident.GetRawValue().(uint16))
	case runtime.GDStringIdentMode:
		str := ident.GetRawValue().(string)

		err := bytecode.WriteByte(byte(len(str)))
		if err != nil {
			return err
		}

		_, err = bytecode.WriteString(str)
		if err != nil {
			return err
		}
	}

	return nil
}

func WriteString(bytecode *bytes.Buffer, str string) error {
	err := writeObjectWithType(bytecode, runtime.NewGDIntNumber(runtime.GDInt(len(str))))
	if err != nil {
		return err
	}
	_, err = bytecode.WriteString(str)
	return err
}

func WriteBool(bytecode *bytes.Buffer, val bool) error {
	if val {
		return bytecode.WriteByte(1)
	}
	return bytecode.WriteByte(0)
}

func WriteInt(bytecode *bytes.Buffer, val int) error {
	_, err := bytecode.Write(ItoBytes(val))
	return err
}

func WriteInt8(bytecode *bytes.Buffer, val int8) error {
	_, err := bytecode.Write(I8toBytes(val))
	return err
}

func WriteInt16(bytecode *bytes.Buffer, val int16) error {
	_, err := bytecode.Write(I16toBytes(val))
	return err
}

func WriteUInt16(bytecode *bytes.Buffer, val uint16) error {
	_, err := bytecode.Write(UI16toBytes(val))
	return err
}

func WriteUInt32(bytecode *bytes.Buffer, val uint32) error {
	_, err := bytecode.Write(UI32toBytes(val))
	return err
}

func WriteInt32(bytecode *bytes.Buffer, val int32) error {
	_, err := bytecode.Write(ItoBytes(int(val)))
	return err
}

func WriteUInt(bytecode *bytes.Buffer, val uint) error {
	_, err := bytecode.Write(ItoBytes(int(val)))
	return err
}

func WriteFloat32(bytecode *bytes.Buffer, val float32) error {
	_, err := bytecode.Write(F32toBytes(val))
	return err
}

func WriteFloat64(bytecode *bytes.Buffer, val float64) error {
	_, err := bytecode.Write(F64toBytes(val))
	return err
}

func WriteComplex64(bytecode *bytes.Buffer, val complex64) error {
	_, err := bytecode.Write(C64toBytes(val))
	return err
}

func WriteComplex128(bytecode *bytes.Buffer, val complex128) error {
	_, err := bytecode.Write(C128toBytes(val))
	return err
}

// Converters

// Int

func ItoBytes(i int) []byte {
	data := *(*[unsafe.Sizeof(i)]byte)(unsafe.Pointer(&i))
	return data[:]
}

// Int16

func I16toBytes(i int16) []byte {
	data := *(*[unsafe.Sizeof(i)]byte)(unsafe.Pointer(&i))
	return data[:]
}

func UI16toBytes(i uint16) []byte {
	data := *(*[unsafe.Sizeof(i)]byte)(unsafe.Pointer(&i))
	return data[:]
}

// Int32

func I32toBytes(i int32) []byte {
	data := *(*[unsafe.Sizeof(i)]byte)(unsafe.Pointer(&i))
	return data[:]
}

func UI32toBytes(i uint32) []byte {
	data := *(*[unsafe.Sizeof(i)]byte)(unsafe.Pointer(&i))
	return data[:]
}

// Int8

func I8toBytes(i int8) []byte {
	data := *(*[unsafe.Sizeof(i)]byte)(unsafe.Pointer(&i))
	return data[:]
}

// Float32

func F32toBytes(f float32) []byte {
	data := *(*[unsafe.Sizeof(f)]byte)(unsafe.Pointer(&f))
	return data[:]
}

// Float64

func F64toBytes(f float64) []byte {
	data := *(*[unsafe.Sizeof(f)]byte)(unsafe.Pointer(&f))
	return data[:]
}

// Complex64

func C64toBytes(c complex64) []byte {
	data := *(*[unsafe.Sizeof(c)]byte)(unsafe.Pointer(&c))
	return data[:]
}

// Complex128

func C128toBytes(c complex128) []byte {
	data := *(*[unsafe.Sizeof(c)]byte)(unsafe.Pointer(&c))
	return data[:]
}
