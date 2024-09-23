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

import (
	"os"
	"strconv"
	"strings"
)

func Printf(format string, args ...any) {
	Print(Sprintf(format, args...))
}

func Print(value string) {
	os.Stdout.WriteString(value)
}

func Sprintf(format string, args ...any) string {
	var buffer strings.Builder
	defer buffer.Reset()

	i := 0
	argIndex := 0
	fmtLen := len(format)
	argsLen := len(args)
	for ; i < fmtLen; i++ { // Always move forward
		ch := rune(format[i])
		switch ch {
		case '%':
			l := i + 1 // Look ahead
			if l < fmtLen && format[l] == '@' && argIndex < argsLen {
				buffer.WriteString(convertAny(args[argIndex]))
				argIndex++
				i += 1 // skip %
				continue
			}
		}
		buffer.WriteRune(ch)
	}

	return buffer.String()
}

func convertAny(value any) string {
	switch v := value.(type) {
	case GDTypable:
		return v.ToString()
	case GDObject:
		return v.ToString()
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(int64(v.(int)), 10)
	case GDByteIdent:
		return strconv.FormatUint(uint64(byte(v)), 10)
	case GDUInt16Ident:
		return strconv.FormatUint(uint64(uint16(v)), 10)
	case GDIdent:
		return v.ToString()
	case uint, uint16, uint32, uint64:
		return strconv.FormatUint(uint64(v.(uint)), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case float32, float64:
		return strconv.FormatFloat(float64(v.(float64)), 'f', -1, 64)
	case complex64, complex128:
		return strconv.FormatComplex(v.(complex128), 'f', -1, 128)
	case bool:
		return strconv.FormatBool(v)
	case error:
		return v.Error()
	case string:
		return v
	}

	panic(NewGDRuntimeErr(UnsupportedTypeCode, "Unsupported type when converting to string"))
}
