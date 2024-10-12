package staticcheck

import (
	"gdlang/lib/builtin"
	"gdlang/lib/runtime"
)

func ImportBuiltins(stack *runtime.GDStack) error {
	// Import core builtins
	coreBuiltins := builtin.GetCoreBuiltins()
	for strIdent, builtin := range coreBuiltins {
		symbol, err := builtin()
		if err != nil {
			return err
		}

		ident := runtime.NewGDStrIdent(strIdent)
		inference := NewInference(ident, symbol.Type)
		_, err = stack.AddNewSymbol(ident, symbol.IsPub, symbol.IsConst, symbol.Type, nil, inference)
		if err != nil {
			return err
		}
	}

	return nil
}
