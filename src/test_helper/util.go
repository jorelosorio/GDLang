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

package test_helper

import (
	"os"
	"path/filepath"
)

type FNode interface{ Name() string }

type FNodeBase struct{ name string }

func (n FNodeBase) Name() string { return n.name }

func NewFNodeBase(name string) FNodeBase { return FNodeBase{name} }

type FFile struct {
	Src string
	FNode
}

func NFile(name, src string) FFile { return FFile{src, FNodeBase{name}} }

func NMFile(src string) FFile { return FFile{src, FNodeBase{"main.gd"}} }

type FDir struct {
	Items []FNode
	FNode
}

func NDir(name string, items ...FNode) FDir { return FDir{items, FNodeBase{name}} }

func NRDir(items ...FNode) FDir { return FDir{items, FNodeBase{""}} }

func BuildPackageTree(node FNode, run func(tmpDir string) error) error {
	tmpDir, err := os.MkdirTemp("", "gdl-*")
	if err != nil {
		return err
	}

	defer os.RemoveAll(tmpDir)

	var traverseTree func(node FNode, tmpDir string) error
	traverseTree = func(node FNode, tmpDir string) error {
		switch n := node.(type) {
		case FFile:
			tmpFile := filepath.Join(tmpDir, n.Name())
			err = os.WriteFile(tmpFile, []byte(n.Src), 0644)
			if err != nil {
				return err
			}
		case FDir:
			if n.Name() != "" {
				tmpDir = filepath.Join(tmpDir, n.Name())
				err := os.Mkdir(tmpDir, 0755)
				if err != nil {
					return err
				}
			}

			for _, item := range n.Items {
				err := traverseTree(item, tmpDir)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}

	err = traverseTree(node, tmpDir)
	if err != nil {
		return err
	}

	return run(tmpDir)
}

func BuildSrc(src string, run func(tmpDir string) error) error {
	return BuildPackageTree(FFile{src, FNodeBase{"main.gd"}}, run)
}
