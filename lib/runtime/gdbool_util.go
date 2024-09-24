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

func IsBool(value any) bool {
	switch value := value.(type) {
	case bool:
		return true
	case GDObject:
		return value.GetType() == GDBoolType
	}

	return false
}

func ToBool(value any) (GDBool, error) {
	switch value := value.(type) {
	case GDBool:
		return value, nil
	case GDObject:
		return false, InvalidCastingWrongTypeErr(GDBoolType, value.GetType())
	}

	return false, InvalidCastingExpectedTypeErr(GDBoolType)
}
