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

package analysis

import (
	"fmt"
	"gdlang/lib/runtime"
	"gdlang/src/comn"
	"gdlang/src/gd/ast"
	"gdlang/src/gd/scanner"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// Dependency Analysis and Compilation Ordering

type GDFileNode struct {
	ast.Node   // Node reference
	*GDSrcFile // In which file the node is stored
}

type FileNodeMap map[string]*GDFileNode

type GDSrcFile struct {
	File  *scanner.File
	Pkg   *GDPackage   // A reference to the package that the file belongs to
	Pkgs  []*GDPackage // All the packages that the file uses
	Uses  FileNodeMap  // Objects that file have access to from different packages
	Pub   FileNodeMap  // Public object references
	Local FileNodeMap  // Local object references
	Nodes []ast.Node   // Ast ordered nodes
}

type GDPackage struct {
	Name  string       // Package name
	Path  string       // Location where the package is stored
	Pub   FileNodeMap  // All public objects that the package has
	Files []*GDSrcFile // Source files that conforms the package
}

// Contains the reusable objects that are used by all the processes
// to evaluate the source files, as well a reference to the main package
type GDDepAnalyzer struct {
	fileSet     *scanner.FileSet
	astBuilder  *ast.AstBuilderProc
	basePath    string                // Main path where the source files are stored
	mainPkg     *GDPackage            // The main root package of the package tree
	visitedPkgs map[string]*GDPackage // All visited packages
	FileNodes   []*GDFileNode
	MainRef     *GDFileNode
}

// Dependency analysis options
type DepAnalyzerOpt struct {
	// It takes the main package as the main entry point
	// and start the dependency analysis from there.
	// this means it collects all the dependencies that are actually used
	// by the main package.
	MainAsEntryPoint bool
}

// Base path is where the source files are stored
func (p *GDDepAnalyzer) Build(basePath string, op DepAnalyzerOpt) error {
	p.basePath = basePath
	p.mainPkg = &GDPackage{
		Name: "main",
		Path: basePath,
		Pub:  make(FileNodeMap),
	}

	// Follows all the packages and files that are referenced in the package files
	// and builds the package tree.
	err := p.buildPkg(p.mainPkg, ast.NewNodePackage([]ast.Node{}, []ast.Node{}))
	if err != nil {
		return err
	}

	// Free memory from the scanner and ast processor
	p.fileSet.Reset()
	p.astBuilder.Dispose()

	// Check for the main entry
	mainRef := p.mainPkg.Pub["main"]
	if mainRef == nil {
		return comn.NErr(comn.DefaultFatalErrCode, comn.NoMainFunctionErrMsg, comn.FatalError, scanner.NZeroPostAt(basePath), nil)
	}

	p.MainRef = mainRef
	visitedObjects := make(map[string]bool)

	defer (func() {
		visitedObjects = nil
	})()

	getNodeRef := func(id string, file *GDSrcFile) *GDFileNode {
		if pubUses, ok := file.Uses[id]; ok {
			return pubUses
		} else if pkgPubUses, ok := file.Pkg.Pub[id]; ok {
			return pkgPubUses
		} else if local, ok := file.Local[id]; ok {
			return local
		}
		return nil
	}

	// Determine the order of top-level objects and their dependencies
	var walkType func(file *GDSrcFile, typ runtime.GDTypable) error
	walkType = func(file *GDSrcFile, typ runtime.GDTypable) error {
		switch typ := typ.(type) {
		case runtime.GDUnionType:
			for _, field := range typ {
				err := walkType(file, field)
				if err != nil {
					return err
				}
			}
		case runtime.GDIdent:
			identStr := typ.ToString()
			path := identStr + "@" + file.File.Name()
			if _, visited := visitedObjects[path]; visited {
				return nil
			}
			visitedObjects[path] = true

			objNodRef := getNodeRef(identStr, file)
			if objNodRef != nil {
				p.FileNodes = append(p.FileNodes, objNodRef)
			}
		case *runtime.GDArrayType:
			return walkType(file, typ.SubType)
		case runtime.GDTupleType:
			for _, elem := range typ {
				err := walkType(file, elem)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}

	var traverseAST func(node ast.Node, file *GDSrcFile) error
	traverseAST = func(node ast.Node, file *GDSrcFile) error {
		switch node := node.(type) {
		case *ast.NodeIdent:
			switch node.Token {
			case scanner.IDENT:
				// Check if the identifier is a public object
				// but it is not found, then it does not throw an error, it is because
				// Other analysis processes will check for identities using the stack.
				if objNodRef := getNodeRef(node.Lit, file); objNodRef != nil {
					err := traverseAST(objNodRef.Node, objNodRef.GDSrcFile)
					if err != nil {
						return err
					}
				}
			}
		case *ast.NodeFunc:
			path := node.Ident.Lit + "@" + file.File.Name()
			if _, visited := visitedObjects[path]; visited {
				return nil
			}
			visitedObjects[path] = true

			err := traverseAST(node.NodeLambda, file)
			if err != nil {
				return err
			}

			objNodRef := getNodeRef(node.Ident.Lit, file)
			if objNodRef != nil {
				p.FileNodes = append(p.FileNodes, objNodRef)
			}
		case *ast.NodeLambda:
			return traverseAST(node.Block, file)
		case *ast.NodeMutCollectionOp:
			err := traverseAST(node.L, file)
			if err != nil {
				return err
			}

			return traverseAST(node.R, file)
		case *ast.NodeIf:
			for _, node := range node.Conds {
				err := traverseAST(node, file)
				if err != nil {
					return err
				}
			}

			return traverseAST(node.Block, file)
		case *ast.NodeTernaryIf:
			err := traverseAST(node.Expr, file)
			if err != nil {
				return err
			}

			err = traverseAST(node.Then, file)
			if err != nil {
				return err
			}

			return traverseAST(node.Else, file)
		case *ast.NodeForIn:
			err := traverseAST(node.Sets, file)
			if err != nil {
				return err
			}

			err = traverseAST(node.Expr, file)
			if err != nil {
				return err
			}

			return traverseAST(node.Block, file)
		case *ast.NodeForIf:
			if node.Sets != nil {
				err := traverseAST(node.Sets, file)
				if err != nil {
					return err
				}
			}

			if node.Conds != nil {
				for _, node := range node.Conds {
					err := traverseAST(node, file)
					if err != nil {
						return err
					}
				}
			}

			return traverseAST(node.Block, file)
		case *ast.NodeTypeAlias:
			return walkType(file, node.Type)
		case *ast.NodeCastExpr:
			err := traverseAST(node.Expr, file)
			if err != nil {
				return err
			}

			return walkType(file, node.Type)
		case *ast.NodeIfElse:
			err := traverseAST(node.If, file)
			if err != nil {
				return err
			}

			for _, elseIf := range node.ElseIf {
				err := traverseAST(elseIf, file)
				if err != nil {
					return err
				}
			}

			if node.Else != nil {
				return traverseAST(node.Else, file)
			}

			return nil
		case *ast.NodeExprOperation:
			err := traverseAST(node.L, file)
			if err != nil {
				return err
			}

			return traverseAST(node.R, file)
		case *ast.NodeUpdateSet:
			err := traverseAST(node.IdentExpr, file)
			if err != nil {
				return err
			}

			return traverseAST(node.Expr, file)
		case *ast.NodeSets:
			for _, node := range node.Nodes {
				err = traverseAST(node, file)
				if err != nil {
					return err
				}
			}
		case *ast.NodeSharedExpr:
			return traverseAST(node.Expr, file)
		case *ast.NodeSet:
			if node.IdentWithType.Type != nil {
				err := walkType(file, node.IdentWithType.Type)
				if err != nil {
					return err
				}
			}

			if node.Expr != nil {
				err = traverseAST(node.Expr, file)
				if err != nil {
					return err
				}
			}

			path := node.IdentWithType.Ident.Lit + "@" + file.File.Name()
			if _, visited := visitedObjects[path]; visited {
				return nil
			}
			visitedObjects[path] = true

			objNodRef := getNodeRef(node.IdentWithType.Ident.Lit, file)
			if objNodRef != nil {
				p.FileNodes = append(p.FileNodes, objNodRef)
			}
		case *ast.NodeEllipsisExpr:
			return traverseAST(node.Expr, file)
		case *ast.NodeArray:
			for _, node := range node.Nodes {
				err := traverseAST(node, file)
				if err != nil {
					return err
				}
			}
		case *ast.NodeTuple:
			for _, node := range node.Nodes {
				err := traverseAST(node, file)
				if err != nil {
					return err
				}
			}
		case *ast.NodeStruct:
			for _, attr := range node.Nodes {
				err := traverseAST(attr, file)
				if err != nil {
					return err
				}
			}
		case *ast.NodeStructAttr:
			return traverseAST(node.Expr, file)
		case *ast.NodeReturn:
			return traverseAST(node.Expr, file)
		case *ast.NodeSafeDotExpr:
			err := traverseAST(node.Expr, file)
			if err != nil {
				return err
			}

			return traverseAST(node.Ident, file)
		case *ast.NodeIterIdxExpr:
			err := traverseAST(node.Expr, file)
			if err != nil {
				return err
			}

			return traverseAST(node.IdxExpr, file)
		case *ast.NodeCallExpr:
			// If the expression is an identifier it means that it is a function call at
			// the top level.
			err := traverseAST(node.Expr, file)
			if err != nil {
				return err
			}

			for _, nArg := range node.Args {
				err := traverseAST(nArg, file)
				if err != nil {
					return err
				}
			}
		case *ast.NodeBlock:
			for _, node := range node.Nodes {
				err := traverseAST(node, file)
				if err != nil {
					return err
				}
			}

			return nil
		}

		return nil
	}

	// Start from the main function and traverse the ast tree
	if op.MainAsEntryPoint {
		err = traverseAST(mainRef.Node, mainRef.GDSrcFile)
		if err != nil {
			return err
		}
	} else {
		var traversePkg func(p *GDPackage) error
		traversePkg = func(p *GDPackage) error {
			for _, file := range p.Files {
				for _, node := range file.Pkgs {
					err := traversePkg(node)
					if err != nil {
						return err
					}
				}

				for _, node := range file.Nodes {
					err := traverseAST(node, file)
					if err != nil {
						return err
					}
				}
			}

			return nil
		}

		err := traversePkg(p.mainPkg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *GDDepAnalyzer) Dispose() {
	p.astBuilder.Dispose()
	p.fileSet.Reset()
	p.basePath = ""
	p.mainPkg = nil
	p.FileNodes = make([]*GDFileNode, 0)
}

func (p *GDDepAnalyzer) buildGDFileFromPath(filePath string, parentPkg *GDPackage, nodPkg *ast.NodePackage) (*GDSrcFile, error) {
	srcBytes, err := readSrcFile(filePath)
	if err != nil {
		return nil, comn.GDFileParsingErr(nodPkg.GetName(), nodPkg.GetPosition())
	}

	// Builds file ast tree and evaluate other packages
	fsFile, err := p.fileSet.AddFile(filePath, p.fileSet.Base(), len(srcBytes))
	if err != nil {
		return nil, comn.GDFileParsingErr(nodPkg.GetName(), nodPkg.GetPosition())
	}

	err = p.astBuilder.Init(fsFile, srcBytes)
	if err != nil {
		return nil, comn.GDFileParsingErr(nodPkg.GetName(), nodPkg.GetPosition())
	}

	// Errors with scanner or parser
	rootNode, err := p.astBuilder.Build()
	if err != nil {
		return nil, err
	}

	// Free memory (Ast is not longer needed)
	defer (func() {
		p.astBuilder.Dispose()
		rootNode = nil
	})()

	// Load other packages referenced in the file
	nodGDFile, isNodGDFile := rootNode.(*ast.NodeFile)
	if !isNodGDFile {
		return nil, comn.GDFileParsingErr(nodPkg.GetName(), nodPkg.GetPosition())
	}

	// Order nodes
	sort.SliceStable(nodGDFile.Nodes, func(i, j int) bool {
		return nodGDFile.Nodes[i].Order() < nodGDFile.Nodes[j].Order()
	})

	gdFile := &GDSrcFile{
		File:  fsFile,
		Pkg:   parentPkg,
		Pkgs:  make([]*GDPackage, 0),
		Nodes: nodGDFile.Nodes,
		Uses:  make(FileNodeMap),
		Pub:   make(FileNodeMap),
		Local: make(FileNodeMap),
	}

	// Read all packages that are used / referenced in the file
	// as: `use package { pub_functions }`.
	for _, nPkgs := range nodGDFile.Packages {
		nPkg, isNodPackage := nPkgs.(*ast.NodePackage)
		if !isNodPackage {
			return nil, comn.NErr(comn.DefaultFatalErrCode, "invalid package", comn.FatalError, nPkgs.GetPosition(), nil)
		}

		// Try find the package using the relative path
		pkgPath := path.Join(parentPkg.Path, nPkg.GetPath())
		// Check if package directory exists
		if _, err := os.Stat(pkgPath); err != nil {
			// Try to find the package using the absolute path
			pkgPath = path.Join(p.basePath, nPkg.GetPath())
		}

		pkg := p.visitedPkgs[pkgPath]

		if pkg == nil {
			// Add the package to the list of packages that the file uses
			filePkg := &GDPackage{
				Name: nPkg.GetName(),
				// Previous package path + new package path
				// TODO: Look for "std" package if relative path does not exist
				Path: pkgPath,
				Pub:  make(FileNodeMap),
			}

			// Load a package
			err := p.buildPkg(filePkg, nPkg)
			if err != nil {
				return nil, err
			}

			// Set the package reference
			pkg = filePkg

			gdFile.Pkgs = append(gdFile.Pkgs, filePkg)
			p.visitedPkgs[pkgPath] = filePkg
		}

		// Read all the public references the file is referencing
		// And store from which file they come from.
		for _, pub := range nPkg.Pub {
			if ident, ok := pub.(*ast.NodeIdent); ok {
				fileNodeRef := pkg.Pub[ident.Lit]
				// Check that the public object that is being used by the file exists
				if fileNodeRef == nil {
					msg := fmt.Sprintf(comn.PublicObjectNotFoundErrMsg, ident.Lit, pkg.Name)
					return nil, comn.NErr(comn.PublicObjectNotFoundErrCode, msg, comn.FatalError, ident.GetPosition(), nil)
				}

				// Register the public object that the file uses and his node / file reference
				gdFile.Uses[ident.Lit] = fileNodeRef
			} else {
				panic("invalid package")
			}
		}
	}

	return gdFile, nil
}

func (p *GDDepAnalyzer) buildPkg(pkg *GDPackage, nodPkg *ast.NodePackage) error {
	err := filepath.WalkDir(pkg.Path, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return comn.CreatePackageNotFoundErr(pkg.Name, nodPkg.GetPosition())
		}

		if d.IsDir() {
			if filePath != pkg.Path {
				return filepath.SkipDir
			}

			return nil
		}

		if strings.HasSuffix(filePath, ".gd") {
			gdFile, err := p.buildGDFileFromPath(filePath, pkg, nodPkg)
			if err != nil {
				return err
			}

			// Check for top-node level objects
			// and register them in the package as public or private objects
			// also check for duplicated objects
			for _, node := range gdFile.Nodes {
				switch node := node.(type) {
				case *ast.NodeFunc:
					funcName := node.Ident.Lit
					fileNodeRef := &GDFileNode{node, gdFile}

					if node.IsPub {
						// Do not allow duplicated pub objects
						if _, ok := pkg.Pub[funcName]; ok {
							msg := fmt.Sprintf(comn.DuplicatedPublicObjectErrMsg, node.Ident.Lit, pkg.Name)
							return comn.NErr(comn.DefaultFatalErrCode, msg, comn.FatalError, node.GetPosition(), nil)
						}

						// Public object references
						gdFile.Pub[funcName] = fileNodeRef

						// Package reference for pub objects in all files
						pkg.Pub[funcName] = fileNodeRef
					} else {
						if _, ok := gdFile.Local[funcName]; ok {
							msg := fmt.Sprintf(comn.DuplicatedPublicObjectErrMsg, node.Ident.Lit, pkg.Name)
							return comn.NErr(comn.DefaultFatalErrCode, msg, comn.FatalError, node.GetPosition(), nil)
						}

						// Private object references
						gdFile.Local[funcName] = fileNodeRef
					}
				case *ast.NodeTypeAlias:
					if _, ok := gdFile.Local[node.Ident.Lit]; ok {
						msg := fmt.Sprintf(comn.DuplicatedPublicObjectErrMsg, node.Ident.Lit, pkg.Name)
						return comn.NErr(comn.DefaultFatalErrCode, msg, comn.FatalError, node.GetPosition(), nil)
					}

					fileNodeRef := &GDFileNode{node, gdFile}
					if node.IsPub {
						gdFile.Pub[node.Ident.Lit] = fileNodeRef

						// Package reference for pub objects in all files
						pkg.Pub[node.Ident.Lit] = fileNodeRef
					} else {
						gdFile.Local[node.Ident.Lit] = fileNodeRef
					}
				case *ast.NodeSets:
					for _, set := range node.Nodes {
						if set, ok := set.(*ast.NodeSet); ok {
							objName := set.IdentWithType.Ident.Lit
							// TODO: Check for duplicated objects in public and local
							if _, ok := gdFile.Local[objName]; ok {
								msg := fmt.Sprintf(comn.DuplicatedPublicObjectErrMsg, objName, pkg.Name)
								return comn.NErr(comn.DefaultFatalErrCode, msg, comn.FatalError, set.GetPosition(), nil)
							}

							fileNodeRef := &GDFileNode{set, gdFile}
							if set.IsPub {
								gdFile.Pub[objName] = fileNodeRef
								// Package reference for pub objects in all files
								pkg.Pub[objName] = fileNodeRef
							} else {
								gdFile.Local[objName] = fileNodeRef
							}
						} else {
							panic("invalid NodeSet")
						}
					}
				}
			}

			// Append new file to the package
			pkg.Files = append(pkg.Files, gdFile)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func readSrcFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	srcBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return srcBytes, nil
}

func NDepAnalyzerProc() *GDDepAnalyzer {
	return &GDDepAnalyzer{
		astBuilder:  ast.NAstBuilderProc(),
		fileSet:     scanner.NFileSet(),
		visitedPkgs: make(map[string]*GDPackage, 0),
		FileNodes:   make([]*GDFileNode, 0),
	}
}
