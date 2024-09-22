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

package ir_test

import (
	"bytes"
	"gdlang/lib/runtime"
	"gdlang/src/cpu"
	"gdlang/src/gd/ir"
	"testing"
)

func TestResolvePendingLabels(t *testing.T) {
	labelIdent := runtime.NewGDStringIdent("test")
	ctx := ir.NewGDIRContext()

	bytecode := &bytes.Buffer{}

	err := ir.Write(bytecode, cpu.Jump)
	if err != nil {
		t.Error(err)
	}

	// Add a mark to be fullfilled later, when the label is added
	_ = ctx.AddMark(bytecode, bytecode.Len(), labelIdent)

	// Write uint32 where the jump will go
	err = ir.WriteUInt16(bytecode, 0)
	if err != nil {
		t.Error(err)
	}

	// Writes another instruction
	err = ir.Write(bytecode, cpu.BBegin, cpu.BEnd)
	if err != nil {
		t.Error(err)
	}

	codeBytes := bytecode.Bytes()
	expectedBytes := []byte{byte(cpu.Jump), 0, 0, byte(cpu.BBegin), byte(cpu.BEnd)}
	if !bytes.Equal(expectedBytes, codeBytes) {
		t.Errorf("Expected %v, got %v", expectedBytes, codeBytes)
	}

	// A label was added to the context
	ctx.AddLabel(bytecode, 1, runtime.NewGDStringIdent("test"))

	// At this point, the label should be resolved
	codeBytes = bytecode.Bytes()
	expectedBytes = []byte{byte(cpu.Jump), 1, 0, byte(cpu.BBegin), byte(cpu.BEnd)}
	if !bytes.Equal(expectedBytes, codeBytes) {
		t.Errorf("Expected %v, got %v", expectedBytes, codeBytes)
	}

	// Ensure that marks now are empty
	if len(ctx.Marks) != 0 {
		t.Errorf("Expected 0, got %d", len(ctx.Marks))
	}
}
