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

package vm

import (
	"gdlang/lib/runtime"
	"gdlang/src/cpu"
	"unicode/utf8"
	"unsafe"
)

type GDNumberConstraints interface {
	~uint8 | ~uint16 | ~uint32 | ~int | ~int8 | ~int16 | ~float32 | ~float64 | ~complex64 | ~complex128
}

type GDVMReader struct {
	Off  uint
	Buff []byte
}

// Type

func (p *GDVMReader) ReadType(stack *runtime.GDStack) (runtime.GDTypable, error) {
	byteValue, err := p.ReadByte()
	if err != nil {
		return nil, err
	}

	typeCode := runtime.GDTypableCode(byteValue)
	switch typeCode {
	case runtime.GDStructTypeCode:
		attrsLen, err := p.ReadByte()
		if err != nil {
			return nil, err
		}

		attrs := make([]*runtime.GDStructAttrType, attrsLen)
		for i := range attrsLen {
			attrIdent, err := p.ReadIdent()
			if err != nil {
				return nil, err
			}

			attrType, err := p.ReadType(stack)
			if err != nil {
				return nil, err
			}

			attrs[i] = &runtime.GDStructAttrType{Ident: attrIdent, Type: attrType}
		}

		return runtime.NewGDStructType(attrs...), nil
	case runtime.GDObjectRefTypeCode:
		ident, err := p.ReadIdent()
		if err != nil {
			return nil, err
		}

		oIdent := runtime.NewGDObjectRefType(ident)

		return oIdent, nil
	case runtime.GDTypeAliasTypeCode:
		ident, err := p.ReadIdent()
		if err != nil {
			return nil, err
		}

		symbol, err := stack.GetSymbol(ident)
		if err != nil {
			return nil, err
		}

		return symbol.Type, nil
	case runtime.GDArrayTypeCode:
		subType, err := p.ReadType(stack)
		if err != nil {
			return nil, err
		}

		return runtime.NewGDArrayType(subType), nil
	case runtime.GDUnionTypeCode:
		uLen, err := p.ReadByte()
		if err != nil {
			return nil, err
		}

		uTypes := make([]runtime.GDTypable, uLen)
		for i := range uLen {
			uType, err := p.ReadType(stack)
			if err != nil {
				return nil, err
			}

			uTypes[i] = uType
		}

		return runtime.NewGDUnionType(uTypes...), nil
	case runtime.GDTupleTypeCode:
		tLen, err := p.ReadByte()
		if err != nil {
			return nil, err
		}

		tTypes := make([]runtime.GDTypable, tLen)
		for i := range tLen {
			tType, err := p.ReadType(stack)
			if err != nil {
				return nil, err
			}
			tTypes[i] = tType
		}

		return runtime.NewGDTupleType(tTypes...), nil
	case runtime.GDLambdaTypeCode:
		rType, err := p.ReadType(stack)
		if err != nil {
			return nil, err
		}

		argLen, err := p.ReadByte()
		if err != nil {
			return nil, err
		}

		argTypes := make(runtime.GDLambdaArgTypes, argLen)
		for i := range argLen {
			argIdent, err := p.ReadIdent()
			if err != nil {
				return nil, err
			}

			argType, err := p.ReadType(stack)
			if err != nil {
				return nil, err
			}

			argTypes[i] = &runtime.GDLambdaArgType{Key: argIdent, Value: argType}
		}

		isVariadic, err := p.ReadBool()
		if err != nil {
			return nil, err
		}

		return runtime.NewGDLambdaType(argTypes, rType, isVariadic), nil
	default:
		return runtime.GDType(typeCode), nil
	}
}

func (p *GDVMReader) ReadIdent() (runtime.GDIdent, error) {
	modeValue, err := p.ReadByte()
	if err != nil {
		return nil, err
	}

	mode := runtime.GDIdentMode(modeValue)
	switch mode {
	case runtime.GDByteIdentMode:
		b, err := p.ReadByte()
		if err != nil {
			return nil, err
		}

		return runtime.NewGDByteIdent(b), nil
	case runtime.GDUInt16IdentMode:
		uInt16, err := p.ReadUInt16()
		if err != nil {
			return nil, err
		}

		return runtime.NewGDUInt16Ident(uInt16), nil
	case runtime.GDStringIdentMode:
		strLen, err := p.ReadByte()
		if err != nil {
			return nil, err
		}

		str, err := p.readString(uint(strLen))
		if err != nil {
			return nil, err
		}

		return runtime.NewGDStrIdent(str), nil
	}

	return nil, WrongTypeErr(runtime.Sprintf("an `ident` type was expected, but got an invalid mode: %@", byte(mode)))
}

// Objects

func (p *GDVMReader) ReadIntObj(stack *runtime.GDStack) (runtime.GDObject, error) {
	obj, err := p.ReadObject(stack)
	if err != nil {
		return nil, err
	}

	isInt := runtime.IsInt(obj)
	if !isInt {
		return nil, InvalidObjErr("an `int` object", obj)
	}

	return obj, nil
}

func (p *GDVMReader) ReadLambdaObj(stack *runtime.GDStack) (*runtime.GDLambda, error) {
	obj, err := p.ReadObject(stack)
	if err != nil {
		return nil, err
	}

	lambda, ok := obj.(*runtime.GDLambda)
	if !ok {
		return nil, InvalidObjErr("a `lambda` object", obj)
	}

	return lambda, nil
}

func (p *GDVMReader) ReadArrayObj(stack *runtime.GDStack) (*runtime.GDArray, error) {
	obj, err := p.ReadObject(stack)
	if err != nil {
		return nil, err
	}

	arr, ok := obj.(*runtime.GDArray)
	if !ok {
		return nil, InvalidObjErr("an `array` object", obj)
	}

	return arr, nil
}

func (p *GDVMReader) ReadIterObj(stack *runtime.GDStack) (runtime.GDIterableCollection, error) {
	obj, err := p.ReadObject(stack)
	if err != nil {
		return nil, err
	}

	iter, ok := obj.(runtime.GDIterableCollection)
	if !ok {
		return nil, InvalidObjErr("an `iterable` object", obj)
	}

	return iter, nil
}

func (p *GDVMReader) ReadMutCollectionObj(stack *runtime.GDStack) (runtime.GDMutableCollection, error) {
	obj, err := p.ReadObject(stack)
	if err != nil {
		return nil, err
	}

	iter, ok := obj.(runtime.GDMutableCollection)
	if !ok {
		return nil, InvalidObjErr("a `mutable collection` object", obj)
	}

	return iter, nil
}

func (p *GDVMReader) ReadAttributable(stack *runtime.GDStack) (runtime.GDAttributable, error) {
	obj, err := p.ReadObject(stack)
	if err != nil {
		return nil, err
	}

	attr, ok := obj.(runtime.GDAttributable)
	if !ok {
		return nil, InvalidObjErr("an `attributable` object", obj)
	}

	return attr, nil
}

func (p *GDVMReader) ReadObject(stack *runtime.GDStack) (runtime.GDObject, error) {
	typ, err := p.ReadType(stack)
	if err != nil {
		return nil, err
	}

	switch typ.GetCode() {
	case runtime.GDNilTypeCode:
		return runtime.GDZNil, nil
	case runtime.GDObjectRefTypeCode:
		objRef, ok := typ.(runtime.GDObjectRefType)
		if !ok {
			return nil, InvalidTypeErr("an `object reference` type", typ)
		}

		ident := objRef.GDIdent
		switch ident.GetMode() {
		case runtime.GDByteIdentMode:
			b := ident.GetRawValue().(byte)
			switch cpu.GDReg(b) {
			case cpu.RPop:
				return stack.PopBuffer(), nil
			default:
				symbol, err := stack.GetSymbol(ident)
				if err != nil {
					return nil, err
				}

				return symbol.Value.(runtime.GDObject), nil
			}
		case runtime.GDStringIdentMode, runtime.GDUInt16IdentMode:
			symbol, err := stack.GetSymbol(ident)
			if err != nil {
				return nil, err
			}

			return symbol.Value.(runtime.GDObject), nil
		default:
			return nil, WrongTypeErr("an `ident` type was expected, but got an invalid mode")
		}
	case runtime.GDIntTypeCode:
		intVal, err := p.ReadInt()
		if err != nil {
			return nil, err
		}

		return runtime.GDInt(intVal), nil
	case runtime.GDInt8TypeCode:
		int8Val, err := p.ReadInt8()
		if err != nil {
			return nil, err
		}

		return runtime.GDInt8(int8Val), nil
	case runtime.GDInt16TypeCode:
		int16Val, err := p.ReadInt16()
		if err != nil {
			return nil, err
		}

		return runtime.GDInt16(int16Val), nil
	case runtime.GDFloat32TypeCode:
		float32Val, err := p.ReadFloat32()
		if err != nil {
			return nil, err
		}

		return runtime.GDFloat32(float32Val), nil
	case runtime.GDFloat64TypeCode:
		float64Val, err := p.ReadFloat64()
		if err != nil {
			return nil, err
		}

		return runtime.GDFloat64(float64Val), nil
	case runtime.GDComplex64TypeCode:
		complex64Val, err := p.ReadComplex64()
		if err != nil {
			return nil, err
		}

		return runtime.GDComplex64(complex64Val), nil
	case runtime.GDComplex128TypeCode:
		complex128Val, err := p.ReadComplex128()
		if err != nil {
			return nil, err
		}

		return runtime.GDComplex128(complex128Val), nil
	case runtime.GDBoolTypeCode:
		boolVal, err := p.ReadBool()
		if err != nil {
			return nil, err
		}

		return runtime.GDBool(boolVal), nil
	case runtime.GDStringTypeCode:
		sLen, err := p.ReadObject(stack)
		if err != nil {
			return nil, err
		}

		sLenIntVal, err := runtime.ToInt(sLen)
		if err != nil {
			return nil, err
		}

		strVal, err := p.readString(uint(sLenIntVal))
		if err != nil {
			return nil, err
		}

		return runtime.GDString(strVal), nil
	case runtime.GDCharTypeCode:
		runeVal, _, err := p.ReadRune()
		if err != nil {
			return nil, err
		}

		return runtime.GDChar(runeVal), nil
	case runtime.GDTupleTypeCode:
		tLen, err := p.ReadByte()
		if err != nil {
			return nil, err
		}

		tVals := make([]runtime.GDObject, tLen)
		for i := range tLen {
			tVal, err := p.ReadObject(stack)
			if err != nil {
				return nil, err
			}

			tVals[tLen-i-1] = tVal
		}

		return runtime.NewGDTuple(tVals...), nil
	case runtime.GDStructTypeCode:
		sType := typ.(runtime.GDStructType)

		sLen, err := p.ReadByte()
		if err != nil {
			return nil, err
		}

		sObj, err := runtime.NewGDStruct(typ.(runtime.GDStructType), stack)
		if err != nil {
			return nil, err
		}

		for i := range sLen {
			attrObj, err := p.ReadObject(stack)
			if err != nil {
				return nil, err
			}

			_, err = sObj.SetAttr(sType[sLen-i-1].Ident, attrObj.GetType(), attrObj)
			if err != nil {
				return nil, err
			}
		}

		return sObj, nil
	case runtime.GDSpreadableTypeCode:
		exprObj, err := p.ReadObject(stack)
		if err != nil {
			return nil, err
		}

		iter, isIter := exprObj.(runtime.GDIterableCollection)
		if !isIter {
			return nil, InvalidTypeErr("an iterable collection", exprObj.GetType())
		}

		return runtime.NewGDSpreadable(iter), nil
	case runtime.GDArrayTypeCode:
		aLen, err := p.ReadByte()
		if err != nil {
			return nil, err
		}

		aObjs := make([]runtime.GDObject, aLen)
		for i := range aLen {
			aObj, err := p.ReadObject(stack)
			if err != nil {
				return nil, err
			}

			aObjs[aLen-i-1] = aObj
		}

		return runtime.NewGDArrayWithTypeAndObjects(typ.(*runtime.GDArrayType), aObjs), nil
	default:
		return nil, InvalidTypeCodeReadingObjectErr(byte(typ.GetCode()))
	}
}

func (p *GDVMReader) ReadBool() (bool, error) {
	b, err := p.ReadByte()
	if err != nil {
		return false, err
	}
	return b == 1, nil
}

func (p *GDVMReader) ReadRune() (rune, int, error) {
	if p.Off >= uint(len(p.Buff)) {
		return 0, 0, EOFErr
	}
	if c := p.Buff[p.Off]; c < utf8.RuneSelf {
		p.Off++
		return rune(c), 1, nil
	}
	ch, size := utf8.DecodeRune(p.Buff[p.Off:])
	p.Off += uint(size)
	return ch, size, nil
}

func (p *GDVMReader) ReadByte() (byte, error) {
	if p.Off >= uint(len(p.Buff)) {
		return 0, EOFErr
	}
	b := p.Buff[p.Off]
	p.Off++
	return b, nil
}

// Numbers

func (p *GDVMReader) ReadInt() (int, error) {
	d, n := readNumber[int](p.Off, p.Buff)
	p.Off += d
	return n, nil
}

func (p *GDVMReader) ReadInt8() (int8, error) {
	d, n := readNumber[int8](p.Off, p.Buff)
	p.Off += d
	return n, nil
}

func (p *GDVMReader) ReadInt16() (int16, error) {
	d, n := readNumber[int16](p.Off, p.Buff)
	p.Off += d
	return n, nil
}

func (p *GDVMReader) ReadFloat32() (float32, error) {
	d, n := readNumber[float32](p.Off, p.Buff)
	p.Off += d
	return n, nil
}

func (p *GDVMReader) ReadFloat64() (float64, error) {
	d, n := readNumber[float64](p.Off, p.Buff)
	p.Off += d
	return n, nil
}

func (p *GDVMReader) ReadComplex64() (complex64, error) {
	d, n := readNumber[complex64](p.Off, p.Buff)
	p.Off += d
	return n, nil
}

func (p *GDVMReader) ReadComplex128() (complex128, error) {
	d, n := readNumber[complex128](p.Off, p.Buff)
	p.Off += d
	return n, nil
}

// Internal

func (p *GDVMReader) ReadUInt16() (uint16, error) {
	d, n := readNumber[uint16](p.Off, p.Buff)
	p.Off += d
	return n, nil
}

// Strings

func (p *GDVMReader) readString(strLen uint) (string, error) {
	if strLen == 0 {
		return "", nil
	}

	if p.Off >= uint(len(p.Buff)) {
		return "", EOFErr
	}

	b := make([]byte, int(strLen))
	n := copy(b, p.Buff[p.Off:])
	p.Off += uint(n)

	return string(b), nil
}

func readNumber[T GDNumberConstraints](off uint, bytes []byte) (uint, T) {
	var zV T
	nL := unsafe.Sizeof(zV)
	b := make([]byte, nL)
	n := copy(b, bytes[off:])
	return uint(n), bytesToNumber[T](b)
}

func bytesToNumber[T GDNumberConstraints](b []byte) T {
	return *(*T)(unsafe.Pointer(&b[0]))
}

func NewGDVMReader(bytes []byte) *GDVMReader {
	return &GDVMReader{0, bytes}
}
