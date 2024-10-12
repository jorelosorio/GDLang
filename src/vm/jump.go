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

import "gdlang/lib/runtime"

type VMJump uint

func (gd VMJump) GetType() runtime.GDTypable    { return runtime.GDIntTypeRef }
func (gd VMJump) GetSubType() runtime.GDTypable { return nil }
func (gd VMJump) ToString() string              { return "" }
func (gd VMJump) CastToType(typ runtime.GDTypable) (runtime.GDObject, error) {
	return nil, nil
}
