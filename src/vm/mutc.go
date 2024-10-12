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

func (p *GDVMProc) evalCSet(stack *runtime.GDStack) (runtime.GDObject, error) {
	idx, err := p.ReadIntObj(stack)
	if err != nil {
		return nil, err
	}

	mutCObj, err := p.ReadMutCollectionObj(stack)
	if err != nil {
		return nil, err
	}

	obj, err := p.ReadObject(stack)
	if err != nil {
		return nil, err
	}

	intVal, err := runtime.ToInt(idx)
	if err != nil {
		return nil, err
	}

	err = mutCObj.Set(int(intVal), obj, stack)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (p *GDVMProc) evalCAdd(stack *runtime.GDStack) (runtime.GDObject, error) {
	left, err := p.ReadMutCollectionObj(stack)
	if err != nil {
		return nil, err
	}

	right, err := p.ReadObject(stack)
	if err != nil {
		return nil, err
	}

	err = left.AddObject(right, stack)
	if err != nil {
		return nil, err
	}

	// Push the added object to the stack
	stack.PushBuffer(right)

	return nil, nil
}

func (p *GDVMProc) evalCRemove(stack *runtime.GDStack) (runtime.GDObject, error) {
	left, err := p.ReadMutCollectionObj(stack)
	if err != nil {
		return nil, err
	}

	idx, err := p.ReadIntObj(stack)
	if err != nil {
		return nil, err
	}

	idxVal, err := runtime.ToInt(idx)
	if err != nil {
		return nil, err
	}

	obj, err := left.Remove(int(idxVal))
	if err != nil {
		return nil, err
	}

	// Push the removed object to the stack
	stack.PushBuffer(obj)

	return nil, nil
}
