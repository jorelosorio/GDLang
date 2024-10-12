package ir

import (
	"bytes"
	"fmt"
	"gdlang/lib/runtime"
	"gdlang/src/cpu"
	"gdlang/src/gd/ast"
)

type IRASet struct {
	isNilSafe bool
	ident     runtime.GDIdent
	expr      GDIRNode
	obj       GDIRNode
	GDIRBaseNode
}

func (i *IRASet) BuildAssembly(padding string) string {
	return padding + fmt.Sprintf("%s%s %s %s %s", cpu.GetCPUInstName(cpu.ASet), tif(i.isNilSafe, " nilsafe", ""), i.expr.BuildAssembly(""), i.ident.ToString(), i.obj.BuildAssembly(""))
}

func (i *IRASet) BuildBytecode(bytecode *bytes.Buffer, ctx *GDIRContext) error {
	ctx.AddMapping(bytecode, i.GetPosition())

	err := Write(bytecode, cpu.ASet, i.isNilSafe)
	if err != nil {
		return err
	}

	err = Write(bytecode, i.ident)
	if err != nil {
		return err
	}

	err = i.expr.BuildBytecode(bytecode, ctx)
	if err != nil {
		return err
	}

	err = i.obj.BuildBytecode(bytecode, ctx)
	if err != nil {
		return err
	}

	return nil
}

func NewIRASet(isNilSafe bool, ident runtime.GDIdent, expr GDIRNode, obj GDIRNode, node ast.Node) *IRASet {
	return &IRASet{isNilSafe, ident, expr, obj, GDIRBaseNode{node}}
}
