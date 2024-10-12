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
	GDTypeAliasTypeCode
	GDTypeRefTypeCode
	GDObjectRefTypeCode

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

	GDTypeAliasTypeCode: "typealias",
	GDTypeRefTypeCode:   "typeref",
	GDObjectRefTypeCode: "objref",

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

var (
	// Primitive Types
	GDNilTypeRef     = GDType(GDNilTypeCode)
	GDAnyTypeRef     = GDType(GDAnyTypeCode)
	GDBoolTypeRef    = GDType(GDBoolTypeCode)
	GDCharTypeRef    = GDType(GDCharTypeCode)
	GDIntTypeRef     = GDType(GDIntTypeCode)
	GDFloatTypeRef   = GDType(GDFloatTypeCode)
	GDComplexTypeRef = GDType(GDComplexTypeCode)
	GDStringTypeRef  = GDStringType(GDStringTypeCode)

	// Internal Types
	GDUntypedTypeRef = GDType(GDUntypedTypeCode)

	// Sub-Types
	GDInt8TypeRef       = GDType(GDInt8TypeCode)
	GDInt16TypeRef      = GDType(GDInt16TypeCode)
	GDFloat32TypeRef    = GDType(GDFloat32TypeCode)
	GDFloat64TypeRef    = GDType(GDFloat64TypeCode)
	GDComplex64TypeRef  = GDType(GDComplex64TypeCode)
	GDComplex128TypeRef = GDType(GDComplex128TypeCode)
)

type GDType GDTypableCode

func (t GDType) GetCode() GDTypableCode { return GDTypableCode(t) }
func (t GDType) ToString() string       { return GDTypeCodeMap[t] }

type GDStringType GDTypableCode

func (t GDStringType) GetCode() GDTypableCode { return GDTypableCode(t) }
func (t GDStringType) ToString() string       { return GDTypeCodeMap[t] }

func (t GDStringType) GetTypeAt(index int) GDTypable { return GDCharTypeRef }
func (t GDStringType) GetIterableType() GDTypable    { return GDCharTypeRef }
