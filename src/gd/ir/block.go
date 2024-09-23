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

package ir

import (
	"bytes"
	"fmt"
	"gdlang/lib/runtime"
	"gdlang/src/cpu"
)

type GDIRBlock struct {
	nodes []GDIRNode
	GDIRBaseNode
}

func (b *GDIRBlock) BuildAssembly(padding string) string {
	block := runtime.JoinSlice(b.nodes, func(node GDIRNode, _ int) string {
		return node.BuildAssembly(padding + spacing)
	}, "\n")

	return padding + fmt.Sprintf("begin\n%s\n%send", block, padding)
}

func (b *GDIRBlock) BuildBytecode(bytecode *bytes.Buffer, ctx *GDIRContext) error {
	err := Write(bytecode, cpu.BBegin)
	if err != nil {
		return err
	}

	// Block length bytes
	bLenMark := bytecode.Len()

	// Write space for the block length
	err = WriteUInt16(bytecode, 0)
	if err != nil {
		return err
	}

	// Where the block starts
	blockStart := bytecode.Len()

	for _, node := range b.nodes {
		err := node.BuildBytecode(bytecode, ctx)
		if err != nil {
			return err
		}
	}

	blockLen := bytecode.Len() - blockStart

	// Write length of block
	// +1 represents the BEnd opcode
	WriteUInt16At(bLenMark, bytecode, uint16(blockLen+1))

	err = Write(bytecode, cpu.BEnd)
	if err != nil {
		return err
	}

	return nil
}

func (b *GDIRBlock) AddNode(node ...GDIRNode) {
	b.nodes = append(b.nodes, node...)
}

func NewGDIRBlock(nodes ...GDIRNode) *GDIRBlock {
	return &GDIRBlock{nodes, GDIRBaseNode{}}
}
