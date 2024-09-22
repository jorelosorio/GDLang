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

package ir

import (
	"crypto/rand"
	"fmt"
	"gdlang/lib/runtime"
	"time"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // 52 possibilities
	letterIdxBits = 6                                                      // 6 bits to represent 64 possibilities / indexes
	letterIdxMask = 1<<letterIdxBits - 1                                   // All 1-bits, as many as letterIdxBits
	genLength     = 4
)

func NewStrIdent() runtime.GDIdent {
	ident := randomAlphaString(genLength)
	return runtime.NewGDStringIdent(ident)
}

func randomAlphaString(length int) string {
	result := make([]byte, length)
	bufferSize := int(float64(length) * 1.3)
	for i, j, rBytes := 0, 0, []byte{}; i < length; j++ {
		if j%bufferSize == 0 {
			rBytes = randomBytes(bufferSize)
		}
		if idx := int(rBytes[j%length] & letterIdxMask); idx < len(letterBytes) {
			result[i] = letterBytes[idx]
			i++
		}
	}

	return string(result)
}

func randomBytes(length int) []byte {
	var randomBytes = make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		// A fallback to use the current time in nanoseconds
		return []byte(fmt.Sprintf("%d", time.Now().UnixNano()))
	}
	return randomBytes
}
