package ir

import (
	"bytes"
	"fmt"
	"gdlang/lib/runtime"
	"gdlang/src/cpu"
)

type GDIRUse struct {
	ident   runtime.GDIdent
	mode    runtime.GDPackageMode
	imports []runtime.GDIdent
	GDIRBaseNode
}

func (u *GDIRUse) BuildAssembly(padding string) string {
	imports := runtime.JoinSlice(u.imports, func(importIdent runtime.GDIdent, _ int) string {
		return importIdent.ToString()
	}, ", ")
	return fmt.Sprintf("%s%s %s %s %s", padding, cpu.GetCPUInstName(cpu.Use), runtime.PackageModeMap[u.mode], u.ident.ToString(), imports)
}

func (u *GDIRUse) BuildBytecode(bytecode *bytes.Buffer, ctx *GDIRContext) error {
	ctx.AddMapping(bytecode, u.GetPosition())

	err := Write(bytecode, cpu.Use)
	if err != nil {
		return err
	}

	err = WriteByte(bytecode, byte(u.mode))
	if err != nil {
		return err
	}

	err = WriteIdent(bytecode, u.ident)
	if err != nil {
		return err
	}

	err = WriteByte(bytecode, byte(len(u.imports)))
	if err != nil {
		return err
	}

	for _, ident := range u.imports {
		err = WriteIdent(bytecode, ident)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewGDIRUse(mode runtime.GDPackageMode, ident runtime.GDIdent, imports []runtime.GDIdent) *GDIRUse {
	return &GDIRUse{ident, mode, imports, GDIRBaseNode{}}
}
