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

type GDSymbol struct {
	Ident   GDIdent
	IsPub   bool
	IsConst bool
	Type    GDTypable
	Object  GDObject
}

func (s *GDSymbol) SetType(typ GDTypable, stack *GDSymbolStack) error {
	if IsUntypedType(s.Type) {
		sType, err := InferType(typ, s.Type, stack)
		if err != nil {
			return err
		}

		s.Type = sType
	} else {
		err := CanBeAssign(s.Type, typ, stack)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *GDSymbol) SetObject(object GDObject, stack *GDSymbolStack) error {
	if s.IsConst {
		return SetConstObjectErr()
	}

	err := s.SetType(object.GetType(), stack)
	if err != nil {
		return err
	}

	s.Object = object

	return nil
}

func NewGDSymbol(isPub, isConst bool, typ GDTypable, object GDObject) *GDSymbol {
	return &GDSymbol{nil, isPub, isConst, typ, object}
}
