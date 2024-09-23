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

func TestBlockLen(t *testing.T) {
	tif, _ := ir.NewGDIRTIf(
		ir.NewGDIRObject(runtime.GDBool(true), nil),
		ir.NewGDIRObject(runtime.GDString("ok"), nil),
		ir.NewGDIRObject(runtime.GDString("no"), nil),
		nil,
	)
	b := ir.NewGDIRBlock()
	b.AddNode(tif)

	bytecode := &bytes.Buffer{}
	err := b.BuildBytecode(bytecode, ir.NewGDIRContext())
	if err != nil {
		t.Error(err)
	}

	// At this point, the label should be resolved
	codeBytes := bytecode.Bytes()
	expectedBytes := []byte{
		byte(cpu.BBegin), 14, 0, // Block and length
		byte(cpu.Tif),                                                                        // Tif
		byte(runtime.GDStringTypeCode), byte(runtime.GDInt8TypeCode) /* int8 */, 2, 'n', 'o', // String
		byte(runtime.GDStringTypeCode), byte(runtime.GDInt8TypeCode) /* int8 */, 2, 'o', 'k', // String
		byte(runtime.GDBoolTypeCode), 1, // Bool true
		byte(cpu.BEnd), // End block
	}
	if !bytes.Equal(expectedBytes, codeBytes) {
		t.Errorf("Expected %v, got %v", expectedBytes, codeBytes)
	}
}
