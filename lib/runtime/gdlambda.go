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

type GDLambdaArg GDKeyValue[GDIdent, GDObject]

// Args are what the function receives after is called, they are identified by a name
// func (a: int, b: int) => int
type GDLambdaArgs []*GDLambdaArg

func (gd GDLambdaArgs) Get(ident GDIdent) (GDObject, error) {
	for _, arg := range gd {
		if arg.Key == ident {
			return arg.Value, nil
		}
	}

	return nil, AttributeNotFoundErr(ident.ToString())
}

type GDLambdaCallback func(args GDLambdaArgs, stack *GDStack) (GDObject, error)

type GDLambda struct {
	Type     *GDLambdaType
	stack    *GDStack         // A reference to the stack where the function was created
	callback GDLambdaCallback // The Callback function that represents the function
}

func (gd *GDLambda) GetType() GDTypable                         { return gd.Type }
func (gd *GDLambda) GetSubType() GDTypable                      { return nil }
func (gd *GDLambda) ToString() string                           { return gd.Type.ToString() }
func (gd *GDLambda) CastToType(typ GDTypable) (GDObject, error) { return nil, nil }

func (gd *GDLambda) Call(args *GDArray) (GDObject, error) {
	if gd.callback == nil {
		return nil, NoFunctionCallbackErr
	}

	mappedArgValues, err := gd.CheckArgValues(args)
	if err != nil {
		return nil, err
	}

	returnObj, err := gd.callback(mappedArgValues, gd.stack)
	if err != nil {
		return nil, err
	}

	err = gd.Type.CheckLambdaReturnType(returnObj.GetType(), any(gd.stack).(*GDStack))
	if err != nil {
		return nil, err
	}

	return returnObj, nil
}

func (gd *GDLambda) CheckArgValues(args *GDArray) (GDLambdaArgs, error) {
	funcArgsLen := len(gd.Type.ArgTypes)
	err := gd.Type.CheckNumberOfArgs(uint(len(args.GetObjects())))
	if err != nil {
		return nil, err
	}

	gdFuncArgObjects := make(GDLambdaArgs, funcArgsLen)
	if gd.Type.IsVariadic {
		funcArgsLen--
		var vargs = make([]GDObject, 0)
		for i := funcArgsLen; i < len(args.GetObjects()); i++ {
			obj, err := args.GetObjectAt(i)
			if err != nil {
				return nil, err
			}

			switch argObj := obj.(type) {
			case *GDSpreadable:
				vargs = append(vargs, argObj.GetObjects()...)
			case GDObject:
				vargs = append(vargs, argObj)
			}
		}
		// TODO: Optimize the creation of the GDFuncArgObjects
		// by no traversing the args.Objects twice
		arrayOfArgs, err := NewGDArrayWithSubTypeAndObjects(gd.Type.ArgTypes[funcArgsLen].Value, vargs, gd.stack)
		if err != nil {
			return nil, err
		}
		gdFuncArgObjects[funcArgsLen] = &GDLambdaArg{gd.Type.ArgTypes[funcArgsLen].Key, arrayOfArgs}
	}
	for i := 0; i < funcArgsLen; i++ {
		funcArg := gd.Type.ArgTypes[i]
		argObj, err := args.GetObjectAt(i)
		if err != nil {
			return nil, err
		}

		err = CanBeAssign(funcArg.Value, argObj.GetType(), gd.stack)
		if err != nil {
			return nil, InvalidArgumentTypeErr(funcArg.Key.ToString(), funcArg.Value, argObj.GetType())
		}

		gdFuncArgObjects[i] = &GDLambdaArg{funcArg.Key, argObj}
	}

	return gdFuncArgObjects, nil
}

func NewGDLambda(args GDLambdaArgTypes, returns GDTypable, isVariadic bool, stack *GDStack, funcCb GDLambdaCallback) *GDLambda {
	return &GDLambda{NewGDLambdaType(args, returns, isVariadic), stack, funcCb}
}

func NewGDLambdaWithType(typ *GDLambdaType, stack *GDStack, funcCb GDLambdaCallback) *GDLambda {
	return &GDLambda{typ, stack, funcCb}
}
