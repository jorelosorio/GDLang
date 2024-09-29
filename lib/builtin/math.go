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

package builtin

import (
	"gdlang/lib/runtime"
	"math"
)

var (
	SquareRootOfNegativeNumberErr = runtime.NewGDRuntimeErr(runtime.RuntimeErrorCode, "cannot calculate square root of a negative number")
	UnsupportedTypeOnMathOpErr    = func(typ runtime.GDTypable, op string) runtime.GDRuntimeErr {
		return runtime.NewGDRuntimeErr(runtime.RuntimeErrorCode, "unsupported `"+typ.ToString()+"` for `"+op+"` operation")
	}
)

func MathPackage() (*runtime.GDPackage[*runtime.GDSymbol], error) {
	pkg := runtime.NewGDPackage[*runtime.GDSymbol](runtime.NewGDStrIdent("math"), "math", runtime.PackageModeBuiltin)
	symbols := map[string]*runtime.GDSymbol{
		"abs":   abs(),
		"sqrt":  sqrt(),
		"pow":   pow(),
		"cos":   cos(),
		"sin":   sin(),
		"log":   log(),
		"tan":   tan(),
		"asin":  asin(),
		"acos":  acos(),
		"atan":  atan(),
		"atan2": atan2(),
		"exp":   exp(),
		"ceil":  ceil(),
		"floor": floor(),
		"round": round(),
		// Constants
		"pi": runtime.NewGDSymbol(true, false, runtime.GDFloatType, runtime.NewGDFloatNumber(runtime.GDFloat64(math.Pi))),
	}
	for ident, symbol := range symbols {
		err := pkg.AddPublic(runtime.NewGDStrIdent(ident), symbol)
		if err != nil {
			return nil, err
		}
	}

	return pkg, nil
}

func abs() *runtime.GDSymbol {
	return mathOp(func(num runtime.GDObject) (runtime.GDObject, error) {
		switch num := num.(type) {
		case runtime.GDFloat32, runtime.GDFloat64:
			fNum, err := runtime.ToFloat(num)
			if err != nil {
				return nil, err
			}

			if fNum < 0.0 {
				return -fNum, nil
			}
		case runtime.GDInt, runtime.GDInt8, runtime.GDInt16:
			intVal, err := runtime.ToInt(num)
			if err != nil {
				return nil, err
			}

			if intVal < 0 {
				return runtime.NewGDIntNumber(-intVal), nil
			}
		default:
			return nil, UnsupportedTypeOnMathOpErr(num.GetType(), "abs")
		}

		return num, nil
	})
}

func sqrt() *runtime.GDSymbol {
	return mathOp(func(num runtime.GDObject) (runtime.GDObject, error) {
		switch num := num.(type) {
		case runtime.GDFloat32, runtime.GDFloat64:
			fNum, err := runtime.ToFloat(num)
			if err != nil {
				return nil, err
			}

			if fNum < 0.0 {
				return nil, SquareRootOfNegativeNumberErr
			}

			return runtime.NewGDFloatNumber(runtime.GDFloat64(math.Sqrt(float64(fNum)))), nil
		case runtime.GDInt, runtime.GDInt8, runtime.GDInt16:
			intVal, err := runtime.ToInt(num)
			if err != nil {
				return nil, err
			}

			if intVal < 0 {
				return nil, SquareRootOfNegativeNumberErr
			}

			return runtime.NewGDFloatNumber(runtime.GDFloat64(math.Sqrt(float64(intVal)))), nil
		}

		return nil, UnsupportedTypeOnMathOpErr(num.GetType(), "sqrt")
	})
}

func pow() *runtime.GDSymbol {
	return mathOp(func(num runtime.GDObject) (runtime.GDObject, error) {
		switch num := num.(type) {
		case runtime.GDFloat32, runtime.GDFloat64:
			fNum, err := runtime.ToFloat(num)
			if err != nil {
				return nil, err
			}

			return runtime.NewGDFloatNumber(runtime.GDFloat64(math.Pow(float64(fNum), 2))), nil
		case runtime.GDInt, runtime.GDInt8, runtime.GDInt16:
			intVal, err := runtime.ToInt(num)
			if err != nil {
				return nil, err
			}

			return runtime.NewGDFloatNumber(runtime.GDFloat64(math.Pow(float64(intVal), 2))), nil
		}

		return nil, UnsupportedTypeOnMathOpErr(num.GetType(), "pow")
	})
}

func cos() *runtime.GDSymbol {
	return mathOp(func(num runtime.GDObject) (runtime.GDObject, error) {
		switch num := num.(type) {
		case runtime.GDFloat32, runtime.GDFloat64:
			fNum, err := runtime.ToFloat(num)
			if err != nil {
				return nil, err
			}

			return runtime.NewGDFloatNumber(runtime.GDFloat64(math.Cos(float64(fNum)))), nil
		case runtime.GDInt, runtime.GDInt8, runtime.GDInt16:
			intVal, err := runtime.ToInt(num)
			if err != nil {
				return nil, err
			}

			return runtime.NewGDFloatNumber(runtime.GDFloat64(math.Cos(float64(intVal)))), nil
		}

		return nil, UnsupportedTypeOnMathOpErr(num.GetType(), "cos")
	})
}

func sin() *runtime.GDSymbol {
	return mathOp(func(num runtime.GDObject) (runtime.GDObject, error) {
		switch num := num.(type) {
		case runtime.GDFloat32, runtime.GDFloat64:
			fNum, err := runtime.ToFloat(num)
			if err != nil {
				return nil, err
			}

			return runtime.NewGDFloatNumber(runtime.GDFloat64(math.Sin(float64(fNum)))), nil
		case runtime.GDInt, runtime.GDInt8, runtime.GDInt16:
			intVal, err := runtime.ToInt(num)
			if err != nil {
				return nil, err
			}

			return runtime.NewGDFloatNumber(runtime.GDFloat64(math.Sin(float64(intVal)))), nil
		}

		return nil, UnsupportedTypeOnMathOpErr(num.GetType(), "sin")
	})
}

func log() *runtime.GDSymbol {
	return mathOp(func(num runtime.GDObject) (runtime.GDObject, error) {
		switch num := num.(type) {
		case runtime.GDFloat32, runtime.GDFloat64:
			fNum, err := runtime.ToFloat(num)
			if err != nil {
				return nil, err
			}

			return runtime.NewGDFloatNumber(runtime.GDFloat64(math.Log(float64(fNum)))), nil
		case runtime.GDInt, runtime.GDInt8, runtime.GDInt16:
			intVal, err := runtime.ToInt(num)
			if err != nil {
				return nil, err
			}

			return runtime.NewGDFloatNumber(runtime.GDFloat64(math.Log(float64(intVal)))), nil
		}

		return nil, UnsupportedTypeOnMathOpErr(num.GetType(), "log")
	})
}

func tan() *runtime.GDSymbol {
	return mathOp(func(num runtime.GDObject) (runtime.GDObject, error) {
		switch num := num.(type) {
		case runtime.GDFloat32, runtime.GDFloat64:
			fNum, err := runtime.ToFloat(num)
			if err != nil {
				return nil, err
			}

			return runtime.NewGDFloatNumber(runtime.GDFloat64(math.Tan(float64(fNum)))), nil
		case runtime.GDInt, runtime.GDInt8, runtime.GDInt16:
			intVal, err := runtime.ToInt(num)
			if err != nil {
				return nil, err
			}

			return runtime.NewGDFloatNumber(runtime.GDFloat64(math.Tan(float64(intVal)))), nil
		}

		return nil, UnsupportedTypeOnMathOpErr(num.GetType(), "tan")
	})
}

func asin() *runtime.GDSymbol {
	return mathOp(func(num runtime.GDObject) (runtime.GDObject, error) {
		switch num := num.(type) {
		case runtime.GDFloat32, runtime.GDFloat64:
			fNum, err := runtime.ToFloat(num)
			if err != nil {
				return nil, err
			}

			return runtime.NewGDFloatNumber(runtime.GDFloat64(math.Asin(float64(fNum)))), nil
		case runtime.GDInt, runtime.GDInt8, runtime.GDInt16:
			intVal, err := runtime.ToInt(num)
			if err != nil {
				return nil, err
			}

			return runtime.NewGDFloatNumber(runtime.GDFloat64(math.Asin(float64(intVal)))), nil
		}

		return nil, UnsupportedTypeOnMathOpErr(num.GetType(), "asin")
	})
}

func acos() *runtime.GDSymbol {
	return mathOp(func(num runtime.GDObject) (runtime.GDObject, error) {
		switch num := num.(type) {
		case runtime.GDFloat32, runtime.GDFloat64:
			fNum, err := runtime.ToFloat(num)
			if err != nil {
				return nil, err
			}

			return runtime.NewGDFloatNumber(runtime.GDFloat64(math.Acos(float64(fNum)))), nil
		case runtime.GDInt, runtime.GDInt8, runtime.GDInt16:
			intVal, err := runtime.ToInt(num)
			if err != nil {
				return nil, err
			}

			return runtime.NewGDFloatNumber(runtime.GDFloat64(math.Acos(float64(intVal)))), nil
		}

		return nil, UnsupportedTypeOnMathOpErr(num.GetType(), "acos")
	})
}

func atan() *runtime.GDSymbol {
	return mathOp(func(num runtime.GDObject) (runtime.GDObject, error) {
		switch num := num.(type) {
		case runtime.GDFloat32, runtime.GDFloat64:
			fNum, err := runtime.ToFloat(num)
			if err != nil {
				return nil, err
			}

			return runtime.NewGDFloatNumber(runtime.GDFloat64(math.Atan(float64(fNum)))), nil
		case runtime.GDInt, runtime.GDInt8, runtime.GDInt16:
			intVal, err := runtime.ToInt(num)
			if err != nil {
				return nil, err
			}

			return runtime.NewGDFloatNumber(runtime.GDFloat64(math.Atan(float64(intVal)))), nil
		}

		return nil, UnsupportedTypeOnMathOpErr(num.GetType(), "atan")
	})
}

func atan2() *runtime.GDSymbol {
	yParam := runtime.NewGDStrRefType("y")
	xParam := runtime.NewGDStrRefType("x")

	typ := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{Key: yParam, Value: runtime.NewGDUnionType(runtime.GDIntType, runtime.GDFloatType)},
			{Key: xParam, Value: runtime.NewGDUnionType(runtime.GDIntType, runtime.GDFloatType)},
		},
		runtime.GDFloatType,
		false,
	)

	lambda := runtime.NewGDLambdaWithType(
		typ,
		nil,
		func(_ *runtime.GDSymbolStack, args runtime.GDLambdaArgs) (runtime.GDObject, error) {
			x := args.Get(xParam)
			y := args.Get(yParam)

			xFloat, err := runtime.ToFloat(x)
			if err != nil {
				return nil, err
			}

			yFloat, err := runtime.ToFloat(y)
			if err != nil {
				return nil, err
			}

			return runtime.NewGDFloatNumber(runtime.GDFloat64(math.Atan2(float64(xFloat), float64(yFloat)))), nil
		},
	)

	return runtime.NewGDSymbol(true, true, typ, lambda)
}

func exp() *runtime.GDSymbol {
	return mathOp(func(num runtime.GDObject) (runtime.GDObject, error) {
		switch num := num.(type) {
		case runtime.GDFloat32, runtime.GDFloat64:
			fNum, err := runtime.ToFloat(num)
			if err != nil {
				return nil, err
			}

			return runtime.NewGDFloatNumber(runtime.GDFloat64(math.Exp(float64(fNum)))), nil
		case runtime.GDInt, runtime.GDInt8, runtime.GDInt16:
			intVal, err := runtime.ToInt(num)
			if err != nil {
				return nil, err
			}

			return runtime.NewGDFloatNumber(runtime.GDFloat64(math.Exp(float64(intVal)))), nil
		}

		return nil, UnsupportedTypeOnMathOpErr(num.GetType(), "exp")
	})
}

func ceil() *runtime.GDSymbol {
	return mathOp(func(num runtime.GDObject) (runtime.GDObject, error) {
		switch num := num.(type) {
		case runtime.GDFloat32, runtime.GDFloat64:
			fNum, err := runtime.ToFloat(num)
			if err != nil {
				return nil, err
			}

			return runtime.NewGDFloatNumber(runtime.GDFloat64(math.Ceil(float64(fNum)))), nil
		case runtime.GDInt, runtime.GDInt8, runtime.GDInt16:
			return num, nil
		}

		return nil, UnsupportedTypeOnMathOpErr(num.GetType(), "ceil")
	})
}

func floor() *runtime.GDSymbol {
	return mathOp(func(num runtime.GDObject) (runtime.GDObject, error) {
		switch num := num.(type) {
		case runtime.GDFloat32, runtime.GDFloat64:
			fNum, err := runtime.ToFloat(num)
			if err != nil {
				return nil, err
			}

			return runtime.NewGDFloatNumber(runtime.GDFloat64(math.Floor(float64(fNum)))), nil
		case runtime.GDInt, runtime.GDInt8, runtime.GDInt16:
			return num, nil
		}

		return nil, UnsupportedTypeOnMathOpErr(num.GetType(), "floor")
	})
}

func round() *runtime.GDSymbol {
	return mathOp(func(num runtime.GDObject) (runtime.GDObject, error) {
		switch num := num.(type) {
		case runtime.GDFloat32, runtime.GDFloat64:
			fNum, err := runtime.ToFloat(num)
			if err != nil {
				return nil, err
			}

			return runtime.NewGDFloatNumber(runtime.GDFloat64(math.Round(float64(fNum)))), nil
		case runtime.GDInt, runtime.GDInt8, runtime.GDInt16:
			return num, nil
		}

		return nil, UnsupportedTypeOnMathOpErr(num.GetType(), "round")
	})
}

func mathOp(opFunc func(num runtime.GDObject) (runtime.GDObject, error)) *runtime.GDSymbol {
	num := runtime.NewGDStrRefType("num")

	typ := runtime.NewGDLambdaType(
		runtime.GDLambdaArgTypes{
			{Key: num, Value: runtime.NewGDUnionType(runtime.GDIntType, runtime.GDFloatType, runtime.GDComplexType)},
		},
		runtime.NewGDUnionType(runtime.GDIntType, runtime.GDFloatType, runtime.GDComplexType),
		false,
	)

	lambda := runtime.NewGDLambdaWithType(
		typ,
		nil,
		func(_ *runtime.GDSymbolStack, args runtime.GDLambdaArgs) (runtime.GDObject, error) {
			num := args.Get(num)
			return opFunc(num)
		},
	)

	return runtime.NewGDSymbol(true, true, typ, lambda)
}
