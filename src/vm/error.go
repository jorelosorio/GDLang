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

package vm

import (
	"gdlang/lib/runtime"
	"gdlang/src/comn"
	"gdlang/src/cpu"
)

type VmErr struct{ msg string }

func (c VmErr) Error() string { return c.msg }

var (
	EOFErr         = VmErr{"EOF"}
	InvalidTypeErr = func(expected string, got runtime.GDTypable) VmErr {
		return VmErr{"Invalid type. Expected " + expected + "`, got " + got.ToString()}
	}
	InvalidObjErr = func(expected string, got runtime.GDObject) VmErr {
		return VmErr{"Invalid object. Expected " + expected + ", got " + got.ToString()}
	}
	InvalidTypeCodeReadingObjectErr = func(code byte) VmErr {
		return VmErr{"Invalid type code: " + runtime.GDTypeCodeMap[code] + " reading object"}
	}
	RuntimeErr = func(err error, inst cpu.GDInst, instOff uint) VmErr {
		errMsg := formatRuntimeError(err.Error(), cpu.GetCPUInstName(inst), uint(inst), instOff)
		mdMsg := comn.NewMarkdown(errMsg)
		return VmErr{mdMsg.Stylize()}
	}
)

// Error template for different runtime exceptions
var errorTemplate = "`[ERROR]` Runtime Exception: `%@`\n" +
	"Instruction:	`%@` (code: `%@`)\n" +
	"Byte Offset:	`%@`\n"

// Function to generate error message
func formatRuntimeError(errorType, instruction string, opcode uint, byteOffset uint) string {
	return runtime.Sprintf(
		errorTemplate,
		errorType,
		instruction,
		opcode,
		byteOffset,
	)
}
