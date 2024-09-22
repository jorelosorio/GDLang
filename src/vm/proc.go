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
	"gdlang/lib/builtin"
	"gdlang/lib/runtime"
	"gdlang/src/cpu"
)

type GDVMProc struct {
	Stack       *runtime.GDSymbolStack
	CInstOffset uint
	CInst       cpu.GDInst
	*GDVMReader
}

type GDVMDisc struct {
	ident          runtime.GDIdent
	isPub, isConst bool
}

func (p *GDVMProc) Init(bytes []byte) error {
	p.GDVMReader = NewGDVMReader(bytes)
	p.Stack = runtime.NewRootGDSymbolStack()

	// Import builtins into the main stack
	err := builtin.Import(p.Stack)
	if err != nil {
		return err
	}

	return nil
}

func (p *GDVMProc) Dispose() {
	if p.Stack != nil {
		p.Stack.Dispose()
	}
	p.Stack = nil
}

func (p *GDVMProc) Run() error {
	for {
		if p.Off >= uint(len(p.Buff)) {
			break
		}

		_, err := p.evalInst(p.Stack)
		if err != nil {
			return RuntimeErr(err, p.CInst, p.CInstOffset)
		}
	}

	return nil
}

func (p *GDVMProc) evalInst(stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	instByte, err := p.ReadByte()
	if err != nil {
		return nil, err
	}

	p.CInst = cpu.GDInst(instByte)
	p.CInstOffset = p.Off - 1

	switch p.CInst {
	case cpu.BBegin:
		return p.evalBlock(stack)
	case cpu.BEnd:
		// It reached the end of the block
		// Nothing to do here!
		return nil, nil
	case cpu.Lambda:
		return p.evalLambda(stack)
	case cpu.Ret:
		return p.evalReturn(stack)
	case cpu.Call:
		return p.evalCall(stack)
	case cpu.Set:
		return p.evalSet(stack)
	case cpu.Mov:
		return p.evalMove(stack)
	case cpu.Jump:
		return p.evalJump(stack)
	case cpu.Operation:
		return p.evalOperation(stack)
	case cpu.TypeAlias:
		return p.evalTypeAlias(stack)
	case cpu.CastObj:
		return p.evalCastObj(stack)

	// Attributable
	case cpu.AGet:
		return p.evalAGet(stack)

	// Mutable collections
	case cpu.CSet:
		return p.evalCSet(stack)
	case cpu.CAdd:
		return p.evalCAdd(stack)
	case cpu.CRemove:
		return p.evalCRemove(stack)

	// Iterables
	case cpu.IGet:
		return p.evalIGet(stack)
	case cpu.ILen:
		return p.evalILen(stack)

	// Conditional instructions
	case cpu.Tif:
		return p.evalTif(stack)
	case cpu.CompareJump:
		return p.evalCompJump(stack)
	}

	panic("Unknown instruction: " + cpu.GetCPUInstName(cpu.GDInst(instByte)))
}

func (p *GDVMProc) evalDisc() (*GDVMDisc, error) {
	isPub, err := p.ReadBool()
	if err != nil {
		return nil, err
	}

	isConst, err := p.ReadBool()
	if err != nil {
		return nil, err
	}

	ident, err := p.ReadIdent()
	if err != nil {
		return nil, err
	}

	return &GDVMDisc{ident, isPub, isConst}, nil
}

func (p *GDVMProc) evalBlock(stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	blockStack := stack.NewSymbolStack(runtime.BlockCtx)
	defer blockStack.Dispose()

	bLen, err := p.ReadUInt16()
	if err != nil {
		return nil, err
	}

	// Read where the block starts
	startOfBlock := p.Off - 1

	// Calculate the end of the block
	// -1 to skip the block byte (cpu.BBegin)
	endOfBlock := startOfBlock + uint(bLen)
walk:
	// Walk the block
	var obj runtime.GDObject
	for {
		obj, err = p.evalInst(blockStack)
		if err != nil {
			return nil, err
		}

		// Stop walking the block if the block end is reached
		// or an object is returned.
		if p.Off >= uint(endOfBlock) || obj != nil {
			break
		}
	}

	if obj != nil {
		switch obj := obj.(type) {
		case VMJump:
			// Jump to the label
			jumpOff := uint(obj)
			if jumpOff >= startOfBlock && jumpOff <= uint(endOfBlock) {
				p.Off = jumpOff
				goto walk
			}

			// Returns sends the jump object to the next upper block
			// to be checked if is within the block range.
			return obj, nil
		}
	}

	// Jump to the end of the block
	p.Off = endOfBlock

	return obj, nil
}

func (p *GDVMProc) evalLambda(stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	// Read lambda type
	typ, err := p.ReadType(stack)
	if err != nil {
		return nil, err
	}

	lambdaType, ok := typ.(*runtime.GDLambdaType)
	if !ok {
		return nil, InvalidTypeErr("a `lambda` type", typ)
	}

	funcBlockStart := p.Off

	lambda := runtime.NewGDLambdaWithType(lambdaType, stack, func(stack *runtime.GDSymbolStack, args runtime.GDLambdaArgs) (runtime.GDObject, error) {
		lambdaStack := stack.NewSymbolStack(runtime.LambdaCtx)
		defer lambdaStack.Dispose()

		for _, arg := range args {
			// Arguments are not public and not constant
			symbol := runtime.NewGDSymbol(false, false, arg.Value.GetType(), arg.Value)
			err := lambdaStack.AddSymbolStack(arg.Key, symbol)
			if err != nil {
				return nil, err
			}
		}

		// Capture the return position, from where the function was called
		returnOff := p.Off

		// Jump to the function block start
		p.Off = funcBlockStart

		// Evaluate the function block
		obj, err := p.evalInst(lambdaStack)
		if err != nil {
			return nil, err
		}

		// Jump back to the return position
		p.Off = returnOff

		if obj != nil {
			return obj, nil
		}

		// Return nil if no object is returned
		return runtime.GDZNil, nil
	})

	// Read block byte
	_, err = p.ReadByte()
	if err != nil {
		return nil, err
	}

	// Read block length
	blockLen, err := p.ReadUInt16()
	if err != nil {
		return nil, err
	}

	// Jump to the end of the block
	p.Off += uint(blockLen)

	// Set the lambda to the buffer
	stack.PushBuffer(lambda)

	return nil, nil
}

func (p *GDVMProc) evalReturn(stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	obj, err := p.ReadObject(stack)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (p *GDVMProc) evalCall(stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	lambda, err := p.ReadLambdaObj(stack)
	if err != nil {
		return nil, err
	}

	args, err := p.ReadArrayObj(stack)
	if err != nil {
		return nil, err
	}

	obj, err := lambda.Call(args)
	if err != nil {
		return nil, err
	}

	// Push the return object to the buffer
	stack.PushBuffer(obj)

	return nil, nil
}

func (p *GDVMProc) evalSet(stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	disc, err := p.evalDisc()
	if err != nil {
		return nil, err
	}

	typ, err := p.ReadType(stack)
	if err != nil {
		return nil, err
	}

	expr, err := p.ReadObject(stack)
	if err != nil {
		return nil, err
	}

	symbol := runtime.NewGDSymbol(disc.isPub, disc.isConst, typ, expr)
	err = stack.AddSymbolStack(disc.ident, symbol)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (p *GDVMProc) evalMove(stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	// Read the target where the expression will be stored
	target, err := p.ReadType(stack)
	if err != nil {
		return nil, err
	}

	ident, isIdent := target.(runtime.GDIdent)
	if !isIdent {
		return nil, InvalidTypeErr("an `ident` object", target)
	}

	switch ident.GetMode() {
	case runtime.GDByteIdentMode:
		byte := ident.GetRawValue().(byte)
		switch cpu.GDReg(byte) {
		case cpu.RPop:
			obj := stack.PopBuffer()

			// Value must be captured after the pop from the stack
			// it is because value also might require to be popped
			value, err := p.ReadObject(stack)
			if err != nil {
				return nil, err
			}

			switch obj := obj.(type) {
			case *runtime.GDAttrIdObject:
				err := obj.SetAttr(obj.Ident, value)
				if err != nil {
					return nil, err
				}
			}
		default:
			value, err := p.ReadObject(stack)
			if err != nil {
				return nil, err
			}

			err = stack.AddOrSetSymbol(ident, value)
			if err != nil {
				return nil, err
			}
		}
	case runtime.GDStringIdentMode, runtime.GDUInt16IdentMode:
		value, err := p.ReadObject(stack)
		if err != nil {
			return nil, err
		}

		err = stack.SetSymbol(ident, value, stack)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (p *GDVMProc) evalJump(_ *runtime.GDSymbolStack) (runtime.GDObject, error) {
	labelOff, err := p.ReadUInt16()
	if err != nil {
		return nil, err
	}

	jumpOff := uint(labelOff)

	return VMJump(jumpOff), nil
}

func (p *GDVMProc) evalOperation(stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	opByte, err := p.ReadByte()
	if err != nil {
		return nil, err
	}

	op := runtime.ExprOperationType(opByte)

	right, err := p.ReadObject(stack)
	if err != nil {
		return nil, err
	}

	left, err := p.ReadObject(stack)
	if err != nil {
		return nil, err
	}

	left, right = runtime.Unwrap(left), runtime.Unwrap(right)
	obj, err := runtime.PerformExprOperation(op, left, right)
	if err != nil {
		return nil, err
	}

	stack.PushBuffer(obj)

	return nil, nil
}

func (p *GDVMProc) evalTypeAlias(stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	disc, err := p.evalDisc()
	if err != nil {
		return nil, err
	}

	typ, err := p.ReadType(stack)
	if err != nil {
		return nil, err
	}

	_, err = stack.AddSymbol(disc.ident, disc.isPub, disc.isConst, typ, nil)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (p *GDVMProc) evalCastObj(stack *runtime.GDSymbolStack) (runtime.GDObject, error) {
	typ, err := p.ReadType(stack)
	if err != nil {
		return nil, err
	}

	obj, err := p.ReadObject(stack)
	if err != nil {
		return nil, err
	}

	castObj, err := obj.CastToType(typ, stack)
	if err != nil {
		return nil, err
	}

	stack.PushBuffer(castObj)

	return nil, nil
}

func NewGDVMProc() *GDVMProc {
	return &GDVMProc{}
}
