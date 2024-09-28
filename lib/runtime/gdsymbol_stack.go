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

type StackContext byte

const (
	GlobalCtx StackContext = iota
	BlockCtx
	LambdaCtx
	ForCtx
	StructCtx
)

type GDSymbolStack struct {
	Parent  *GDSymbolStack
	Ctx     StackContext
	Symbols map[any]*GDSymbol
	Buffer  *GDBuffer
}

func (s *GDSymbolStack) Dispose() {
	s.Symbols = nil
	s.Buffer = nil
}

func (s *GDSymbolStack) NewSymbolStack(ctx StackContext) *GDSymbolStack {
	return &GDSymbolStack{
		s,
		ctx,
		make(map[any]*GDSymbol),
		nil,
	}
}

func (s *GDSymbolStack) AddSymbolStack(ident GDIdent, symbol *GDSymbol) error {
	_, ok := s.Symbols[ident.GetRawValue()]
	if ok {
		return DuplicatedObjectCreationErr(ident)
	}

	s.Symbols[ident.GetRawValue()] = symbol

	return nil
}

func (s *GDSymbolStack) AddSymbol(ident GDIdent, isPub, isConst bool, typ GDTypable, object GDObject) (*GDSymbol, error) {
	_, ok := s.Symbols[ident.GetRawValue()]
	if ok {
		return nil, DuplicatedObjectCreationErr(ident)
	}

	var symbol *GDSymbol
	if object != nil {
		inferredType, err := InferType(typ, GDNilType, s)
		if err != nil {
			return nil, err
		}

		symbol = NewGDSymbol(isPub, isConst, inferredType, object)
	} else {
		symbol = NewGDSymbol(isPub, isConst, typ, GDZNil)
	}

	s.Symbols[ident.GetRawValue()] = symbol

	return symbol, nil
}

func (s *GDSymbolStack) AddOrSetSymbol(ident GDIdent, object GDObject) error {
	symbol, ok := s.Symbols[ident.GetRawValue()]
	if !ok {
		symbol := NewGDSymbol(false, false, object.GetType(), object)
		s.Symbols[ident.GetRawValue()] = symbol

		return nil
	}

	symbol.Type = object.GetType()
	symbol.Object = object

	return nil
}

func (s *GDSymbolStack) SetSymbol(ident GDIdent, object GDObject, stack *GDSymbolStack) error {
	symbol, err := s.GetSymbol(ident)
	if err != nil {
		return err
	}

	err = symbol.SetObject(object, stack)
	if err != nil {
		return err
	}

	return nil
}

func (s *GDSymbolStack) GetSymbol(ident GDIdent) (*GDSymbol, error) {
	symbol, ok := s.Symbols[ident.GetRawValue()]
	if ok {
		return symbol, nil
	}

	if s.Parent != nil {
		return s.Parent.GetSymbol(ident)
	}

	return nil, ObjectNotFoundErr(ident.ToString())
}

// Push a new object to the buffer
func (s *GDSymbolStack) PushBuffer(obj GDObject) {
	// Do not add objects that are nil to avoid unnecessary memory usage
	// in the buffer.
	if obj == GDZNil {
		return
	}

	if s.Buffer == nil {
		s.Buffer = NewGDBuffer()
	}

	s.Buffer.Push(obj)
}

// Pop the last object from the buffer
func (s *GDSymbolStack) PopBuffer() GDObject {
	// If buffer is nil, return nil this compensates when a nil object is
	// trying to be pushed to the buffer.
	if s.Buffer == nil {
		return GDZNil
	}

	return s.Buffer.Pop()
}

// Pull the first object from the buffer
func (s *GDSymbolStack) PullBuffer() GDObject {
	if s.Buffer == nil {
		return GDZNil
	}

	return s.Buffer.Pull()
}

func NewGDSymbolStack() *GDSymbolStack {
	return &GDSymbolStack{
		nil,
		GlobalCtx,
		make(map[any]*GDSymbol),
		nil,
	}
}

func NewRootGDSymbolStack() *GDSymbolStack {
	return &GDSymbolStack{
		nil,
		GlobalCtx,
		make(map[any]*GDSymbol),
		nil,
	}
}
