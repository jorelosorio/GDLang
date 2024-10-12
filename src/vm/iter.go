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

func (p *GDVMProc) evalIGet(stack *runtime.GDStack) (runtime.GDObject, error) {
	isNilSafe, err := p.ReadBool()
	if err != nil {
		return nil, err
	}

	idx, err := p.ReadIntObj(stack)
	if err != nil {
		return nil, err
	}

	iter, err := p.ReadIterObj(stack)
	if err != nil {
		return nil, err
	}

	intVal, err := runtime.ToInt(idx)
	if err != nil {
		return nil, err
	}

	obj, err := iter.GetObjectAt(int(intVal))
	if err != nil && !isNilSafe {
		return nil, err
	} else if err != nil && isNilSafe {
		stack.PushBuffer(runtime.GDZNil)
		return nil, nil
	}

	stack.PushBuffer(obj)

	return nil, nil
}

func (p *GDVMProc) evalILen(stack *runtime.GDStack) (runtime.GDObject, error) {
	iter, err := p.ReadIterObj(stack)
	if err != nil {
		return nil, err
	}

	lenVal := iter.Length()

	stack.PushBuffer(runtime.NewGDIntNumber(runtime.GDInt(lenVal)))

	return nil, nil
}
