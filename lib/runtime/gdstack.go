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

type GDStack struct {
	Parent  *GDStack
	Ctx     StackContext
	Symbols map[any]*GDSymbol
	Buffer  *GDBuffer
}

func (s *GDStack) Dispose() {
	s.Symbols = nil
	s.Buffer = nil
	s.Parent = nil
}

func (s *GDStack) NewStack(ctx StackContext) *GDStack {
	return &GDStack{
		s,
		ctx,
		make(map[any]*GDSymbol),
		NewGDBuffer(),
	}
}

func (s *GDStack) AddSymbol(ident GDIdent, symbol *GDSymbol) error {
	_, ok := s.Symbols[ident.GetRawValue()]
	if ok {
		return DuplicatedSymbolErr(ident)
	}

	s.Symbols[ident.GetRawValue()] = symbol

	return nil
}

// If the assign type is nil, there is no type inference
func (s *GDStack) AddNewSymbol(ident GDIdent, isPub, isConst bool, symbolType, valueType GDTypable, value GDRawValue) (*GDSymbol, error) {
	_, ok := s.Symbols[ident.GetRawValue()]
	if ok {
		return nil, DuplicatedSymbolErr(ident)
	}

	if valueType == nil {
		valueType = GDUntypedTypeRef
	}

	inferredType, err := InferType(symbolType, valueType, s)
	if err != nil {
		return nil, err
	}

	symbol := NewGDSymbol(isPub, isConst, inferredType, value)

	s.Symbols[ident.GetRawValue()] = symbol

	return symbol, nil
}

func (s *GDStack) AddOrSetSymbol(ident GDIdent, valueType GDTypable, value GDRawValue) error {
	symbol, ok := s.Symbols[ident.GetRawValue()]
	if !ok {
		symbol := NewGDSymbol(false, false, valueType, value)
		s.Symbols[ident.GetRawValue()] = symbol

		return nil
	}

	symbol.Type = valueType
	symbol.Value = value

	return nil
}

func (s *GDStack) SetSymbol(ident GDIdent, valueType GDTypable, value GDRawValue) error {
	symbol, err := s.GetSymbol(ident)
	if err != nil {
		return err
	}

	err = symbol.SetObject(valueType, value, s)
	if err != nil {
		return err
	}

	return nil
}

func (s *GDStack) GetSymbol(ident GDIdent) (*GDSymbol, error) {
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
func (s *GDStack) PushBuffer(value GDObject) {
	// Do not add objects that are nil to avoid unnecessary memory usage
	// in the buffer.
	if value == GDZNil {
		return
	}

	s.Buffer.Push(value)
}

// Pop the last object from the buffer
func (s *GDStack) PopBuffer() GDObject { return s.Buffer.Pop() }

func NewGDStack() *GDStack {
	return &GDStack{
		nil,
		GlobalCtx,
		make(map[any]*GDSymbol),
		NewGDBuffer(),
	}
}
