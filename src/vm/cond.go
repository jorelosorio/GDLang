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

func (p *GDVMProc) evalTif(stack *runtime.GDStack) (runtime.GDObject, error) {
	// Read the else expression
	elseExpr, err := p.ReadObject(stack)
	if err != nil {
		return nil, err
	}

	// Read the then expression
	thenExpr, err := p.ReadObject(stack)
	if err != nil {
		return nil, err
	}

	// Read the expression
	expr, err := p.ReadObject(stack)
	if err != nil {
		return nil, err
	}

	condBoolVal, err := runtime.ToBool(expr)
	if err != nil {
		return nil, err
	}

	if condBoolVal {
		stack.PushBuffer(thenExpr)
	} else {
		stack.PushBuffer(elseExpr)
	}

	return nil, nil
}

func (p *GDVMProc) evalCompJump(stack *runtime.GDStack) (runtime.GDObject, error) {
	// Read the expression
	expr, err := p.ReadObject(stack)
	if err != nil {
		return nil, err
	}

	// Read the equalsTo expression
	equalsTo, err := p.ReadObject(stack)
	if err != nil {
		return nil, err
	}

	// Read the label
	labelOff, err := p.ReadUInt16()
	if err != nil {
		return nil, err
	}

	if runtime.EqualObjects(expr, equalsTo) {
		jumpOff := uint(labelOff)
		return VMJump(jumpOff), nil
	}

	return nil, nil
}
