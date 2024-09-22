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

type GDIdentMode byte

const (
	GDByteIdentMode GDIdentMode = iota
	GDUInt16IdentMode
	GDStringIdentMode
)

type GDIdent interface {
	GetMode() GDIdentMode
	GetRawValue() any
	GDPrintable
}

type GDByteIdent byte

func (r GDByteIdent) GetMode() GDIdentMode { return GDByteIdentMode }
func (r GDByteIdent) GetRawValue() any     { return byte(r) }
func (r GDByteIdent) ToString() string     { return Sprintf("%@", r) }

func NewGDByteIdent(value byte) GDIdent {
	return GDByteIdent(value)
}

type GDUInt16Ident uint16

func (r GDUInt16Ident) GetMode() GDIdentMode { return GDUInt16IdentMode }
func (r GDUInt16Ident) GetRawValue() any     { return uint16(r) }
func (r GDUInt16Ident) ToString() string     { return Sprintf("%@", r) }

func NewGDUInt16Ident(value uint16) GDIdent {
	return GDUInt16Ident(value)
}

type GDStringIdent string

func (r GDStringIdent) GetMode() GDIdentMode { return GDStringIdentMode }
func (r GDStringIdent) GetRawValue() any     { return string(r) }
func (r GDStringIdent) ToString() string     { return string(r) }

func NewGDStringIdent(value string) GDIdent {
	return GDStringIdent(value)
}
