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
	"strings"
)

func MapSlice[S ~[]E, E any, T any](slice S, fn func(E, int) T) []T {
	vals := make([]T, len(slice))
	for i, v := range slice {
		vals[i] = fn(v, i)
	}
	return vals
}

func JoinSlice[S ~[]E, E any](slice S, fn func(E, int) string, separator string) string {
	return strings.Join(MapSlice(slice, fn), separator)
}

func ContainsAny(target string, list ...string) bool {
	if len(list) == 0 {
		return false
	}

	result := false
	for _, str := range list {
		result = result || strings.Contains(target, str)
	}

	return result
}

func ContainsAll(target string, list ...string) bool {
	if len(list) == 0 {
		return false
	}

	for _, str := range list {
		if !strings.Contains(target, str) {
			return false
		}
	}

	return true
}
