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

// This wrapper for GDIdent represents a type, commonly for type alias.
type GDTypeRefType struct{ GDIdent }

func (t GDTypeRefType) GetCode() GDTypableCode { return GDTypeRefTypeCode }
func (t GDTypeRefType) ToString() string       { return t.GDIdent.ToString() }

func NewGDTypeRefType(ident GDIdent) GDTypeRefType { return GDTypeRefType{ident} }
func NewGDStrTypeRefType(ident string) GDTypeRefType {
	return NewGDTypeRefType(NewGDStrIdent(ident))
}

// This wrapper for GDIdent represents a type, commonly for object reference.

type GDObjectRefType struct{ GDIdent }

func (t GDObjectRefType) GetCode() GDTypableCode { return GDObjectRefTypeCode }
func (t GDObjectRefType) ToString() string       { return t.GDIdent.ToString() }

func NewGDObjectRefType(ident GDIdent) GDObjectRefType { return GDObjectRefType{ident} }
func NewGDStrObjectRefType(ident string) GDObjectRefType {
	return NewGDObjectRefType(NewGDStrIdent(ident))
}
