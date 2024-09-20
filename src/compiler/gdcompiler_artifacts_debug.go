//go:build debug
// +build debug

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

package compiler

import (
	"os"
	"path"
)

func (c *GDCompiler) GenerateArtifacts(pkgName, outputPath string) error {
	// IR, Bytecode and SourceMap writing process

	// IR assembly
	asmPath := path.Join(outputPath, pkgName+".gdasm")
	err := c.writeAsm(asmPath)
	if err != nil {
		return err
	}

	// Bytecode
	bytecodePath := path.Join(outputPath, pkgName+".gdbin")
	err = c.writeBytecode(bytecodePath)
	if err != nil {
		return err
	}

	// Source map
	srcMapPath := path.Join(outputPath, pkgName+".gdmap")
	err = c.writeSourceMap(srcMapPath)
	if err != nil {
		return err
	}

	return nil
}

func (c *GDCompiler) writeAsm(outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(c.Root.BuildAssembly("")))
	if err != nil {
		return err
	}

	return nil
}
