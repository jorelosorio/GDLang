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

const (
	GDNilTypeCode GDTypableCode = iota
	GDAnyTypeCode
	GDBoolTypeCode
	GDCharTypeCode
	GDStringTypeCode

	GDTupleTypeCode
	GDLambdaTypeCode
	GDArrayTypeCode
	GDStructTypeCode

	// Internal Types
	GDUnionTypeCode
	GDSpreadableTypeCode
	GDUntypedTypeCode
	GDTypeRefTypeCode
	GDObjRefTypeCode

	// Number types ordered
	// from lowest to highest precision
	GDInt8TypeCode
	GDInt16TypeCode
	GDIntTypeCode
	GDFloat32TypeCode
	GDFloat64TypeCode
	GDFloatTypeCode
	GDComplex64TypeCode
	GDComplex128TypeCode
	GDComplexTypeCode
)

var GDTypeCodeMap = [...]string{
	GDNilTypeCode:    "nil",
	GDAnyTypeCode:    "any",
	GDBoolTypeCode:   "bool",
	GDCharTypeCode:   "char",
	GDStringTypeCode: "string",

	GDTupleTypeCode:  "tuple",
	GDLambdaTypeCode: "func",
	GDArrayTypeCode:  "array",
	GDStructTypeCode: "struct",

	// Internal Types
	GDUnionTypeCode:      "unionType",
	GDSpreadableTypeCode: "spreadable",
	GDUntypedTypeCode:    "untyped",

	GDTypeRefTypeCode: "type_ref",
	GDObjRefTypeCode:  "obj_ref",

	// Number types
	GDInt8TypeCode:       "int8",
	GDInt16TypeCode:      "int16",
	GDIntTypeCode:        "int",
	GDFloat32TypeCode:    "float32",
	GDFloat64TypeCode:    "float64",
	GDFloatTypeCode:      "float",
	GDComplex64TypeCode:  "complex64",
	GDComplex128TypeCode: "complex128",
	GDComplexTypeCode:    "complex",
}

type GDType GDTypableCode

var (
	// Primitive Types
	GDNilType     = GDType(GDNilTypeCode)
	GDAnyType     = GDType(GDAnyTypeCode)
	GDBoolType    = GDType(GDBoolTypeCode)
	GDCharType    = GDType(GDCharTypeCode)
	GDIntType     = GDType(GDIntTypeCode)
	GDFloatType   = GDType(GDFloatTypeCode)
	GDComplexType = GDType(GDComplexTypeCode)
	GDStringType  = GDType(GDStringTypeCode)

	// Internal Types
	GDUntypedType = GDType(GDUntypedTypeCode)

	// Sub-Types
	GDInt8Type       = GDType(GDInt8TypeCode)
	GDInt16Type      = GDType(GDInt16TypeCode)
	GDFloat32Type    = GDType(GDFloat32TypeCode)
	GDFloat64Type    = GDType(GDFloat64TypeCode)
	GDComplex64Type  = GDType(GDComplex64TypeCode)
	GDComplex128Type = GDType(GDComplex128TypeCode)
)

func (t GDType) GetCode() GDTypableCode { return GDTypableCode(t) }
func (t GDType) ToString() string       { return GDTypeCodeMap[t] }
