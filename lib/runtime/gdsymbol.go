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
	IsPub   bool
	IsConst bool
	Type    GDTypable
	Value   GDRawValue
}

func (s *GDSymbol) GetType() GDTypable {
	return s.Type
}

func (s *GDSymbol) SetType(typ GDTypable, stack *GDStack) error {
	if s.IsConst {
		return SetConstObjectErr()
	}

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

func (s *GDSymbol) SetObject(typ GDTypable, value GDRawValue, stack *GDStack) error {
	err := s.SetType(typ, stack)
	if err != nil {
		return err
	}

	s.Value = value

	return nil
}

func NewGDSymbol(isPub, isConst bool, typ GDTypable, value GDRawValue) *GDSymbol {
	return &GDSymbol{isPub, isConst, typ, value}
}
