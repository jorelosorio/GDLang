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
	"gdlang/lib/builtin"
	"gdlang/lib/runtime"
	"gdlang/src/gd/ast"
	"gdlang/src/gd/scanner"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type SourceFile struct {
	// Scanner file reference for the source file, it contains the name and number of lines
	file *scanner.File
	// A reference to the package that the source file belongs to
	parentPackage *SourcePackage
	// A list all packages that are used in the source file
	packages []Package
	// Each source file has an internal package that stores all objects created in the file
	// as well references to other objects from other packages that are used in the file
	*runtime.GDPackage[any]
}

func NewSourceFile(file *scanner.File, parentPackage *SourcePackage) *SourceFile {
	return &SourceFile{
		file:          file,
		parentPackage: parentPackage,
		GDPackage:     runtime.NewGDPackage[any](runtime.NewGDStrIdent(file.Name()), file.Name(), runtime.PackageModeSource),
	}
}

type NodeWithSourceFile struct {
	ast.Node
	*SourceFile
}

type SourceFiles []*SourceFile

type Package interface {
	GetName() runtime.GDIdent
	GetPath() string
	GetMode() runtime.GDPackageMode
}

type SourcePackage struct {
	sourceFiles SourceFiles
	*runtime.GDPackage[*NodeWithSourceFile]
}

func (p *SourcePackage) GetName() runtime.GDIdent       { return p.GDPackage.Ident }
func (p *SourcePackage) GetPath() string                { return p.GDPackage.Path }
func (p *SourcePackage) GetMode() runtime.GDPackageMode { return p.GDPackage.Mode }

func NewSourcePackage(name runtime.GDIdent, path string) *SourcePackage {
	return &SourcePackage{
		GDPackage:   runtime.NewGDPackage[*NodeWithSourceFile](name, path, runtime.PackageModeSource),
		sourceFiles: make(SourceFiles, 0),
	}
}

type BuiltInPackage struct {
	*runtime.GDPackage[*runtime.GDSymbol]
}

func (p *BuiltInPackage) GetName() runtime.GDIdent       { return p.Ident }
func (p *BuiltInPackage) GetPath() string                { return p.Path }
func (p *BuiltInPackage) GetMode() runtime.GDPackageMode { return runtime.PackageModeBuiltin }

type PackageDependenciesAnalyzerOptions struct {
	// The package that is being analyzed
	ShouldLookUpFromMain bool
}

type PackageDependenciesAnalyzer struct {
	astBuilder      *ast.BuilderProc
	scanFileSet     *scanner.FileSet
	mainPackage     *SourcePackage
	trackedPackages map[string]Package
	trackedNodes    map[string]bool
	// Computed nodes based on the dependency hierarchy
	Nodes []ast.Node
	// A reference to the main entry point of the package
	MainEntry ast.Node
}

func (d *PackageDependenciesAnalyzer) Analyze(mainPackagePath string, options PackageDependenciesAnalyzerOptions) error {
	// Free memory from the scanner and ast processor from previous runs
	// and reset the state of the analyzer
	d.scanFileSet.Reset()
	d.astBuilder.Dispose()
	d.trackedPackages = make(map[string]Package)
	d.trackedNodes = make(map[string]bool)
	d.Nodes = make([]ast.Node, 0)

	mainIdent := runtime.NewGDStrIdent("main")
	d.mainPackage = NewSourcePackage(mainIdent, mainPackagePath)
	err := d.BuildPackage(d.mainPackage)
	if err != nil {
		return err
	}

	main, err := d.mainPackage.GetMember(mainIdent)
	if err != nil {
		return ErrorAt(scanner.ZeroPos).MainEntryWasNotFound()
	}

	d.MainEntry = main.Node

	if options.ShouldLookUpFromMain {
		d.trackedNodes = make(map[string]bool)
		return d.analyzeNode(main.Node, main.SourceFile)
	} else {
		var traversePkg func(pkg Package) error
		traversePkg = func(pkg Package) error {
			switch pkg := pkg.(type) {
			case *SourcePackage:
				for _, file := range pkg.sourceFiles {
					for _, filePkg := range file.packages {
						err := traversePkg(filePkg)
						if err != nil {
							return err
						}
					}

					for _, member := range file.Members {
						switch node := member.Value.(type) {
						case *NodeWithSourceFile:
							err := d.analyzeNode(node.Node, node.SourceFile)
							if err != nil {
								return err
							}
						case *runtime.GDSymbol:
							// Nothing to do
						}
					}
				}

				return nil
			case *BuiltInPackage:
				// Nothing to do
				return nil
			default:
				return nil
			}
		}

		return traversePkg(d.mainPackage)
	}
}

func (d *PackageDependenciesAnalyzer) BuildPackage(pkg Package) error {
	switch pkg := pkg.(type) {
	case *SourcePackage:
		// Walk through all the files in the package
		err := filepath.WalkDir(pkg.Path, func(sourceFilePath string, dir fs.DirEntry, err error) error {
			if err != nil {
				return ErrorAt(scanner.ZeroPos).PackageNotFound(pkg.Ident.ToString())
			}

			if dir.IsDir() {
				if sourceFilePath != pkg.Path {
					return fs.SkipDir
				}

				return nil
			}

			if strings.HasSuffix(sourceFilePath, ".gd") {
				sourceFile, err := d.BuildSourceFile(sourceFilePath, pkg)
				if err != nil {
					return err
				}

				pkg.sourceFiles = append(pkg.sourceFiles, sourceFile)
			}

			return nil
		})

		return err
	case *BuiltInPackage:
		// Nothing to do
		return nil
	default:
		return nil
	}
}

func (d *PackageDependenciesAnalyzer) BuildSourceFile(sourceFilePath string, pkg *SourcePackage) (*SourceFile, error) {
	sourceBytes, err := os.ReadFile(sourceFilePath)
	if err != nil {
		return nil, ErrorAt(scanner.ZeroPos).ReadingSourceFile(sourceFilePath)
	}

	file, err := d.scanFileSet.AddFile(sourceFilePath, d.scanFileSet.Base(), len(sourceBytes))
	if err != nil {
		return nil, ErrorAt(scanner.ZeroPos).ReadingSourceFile(sourceFilePath)
	}

	// Initialize the AST Builder to be used, it resets
	// some state variables that are used in the AST
	err = d.astBuilder.Init(file, sourceBytes)
	if err != nil {
		return nil, err
	}

	root, err := d.astBuilder.Build()
	if err != nil {
		return nil, err
	}

	// Free memory (Ast is no longer needed)
	defer (func() {
		d.astBuilder.Dispose()
		d.scanFileSet.Reset()
		root = nil
	})()

	sourceFile := NewSourceFile(file, pkg)

	// Traverse nodes and register members
	for _, node := range root.Nodes {
		var err error
		switch node := node.(type) {
		case *ast.NodeFunc:
			nodeIdent := runtime.NewGDStrIdent(node.Ident.Lit)
			err = d.addMemberToSourceFile(nodeIdent, node.IsPub, node, sourceFile)
		case *ast.NodeTypeAlias:
			nodeIdent := runtime.NewGDStrIdent(node.Ident.Lit)
			err = d.addMemberToSourceFile(nodeIdent, node.IsPub, node, sourceFile)
		case *ast.NodeSets:
			for _, node := range node.Nodes {
				set, isNodeSet := node.(*ast.NodeSet)
				if !isNodeSet {
					panic("Invalid node type: expected *ast.NodeSet")
				}

				nodeIdent := runtime.NewGDStrIdent(set.IdentWithType.Ident.Lit)
				err := d.addMemberToSourceFile(nodeIdent, set.IsPub, set, sourceFile)
				if err != nil {
					return nil, err
				}
			}
		}

		if err != nil {
			return nil, err
		}
	}

	// Check for `use` imports in the source file
	for _, node := range root.Packages {
		nodePackage, isNodePackage := node.(*ast.NodePackage)
		if !isNodePackage {
			panic("Invalid node type: expected *ast.NodePackage")
		}

		var pkg Package
		var err error

		// Check if the package exists
		packagePath := nodePackage.GetPath()
		packageAbsolutePath := path.Join(d.mainPackage.Path, nodePackage.GetPath())
		pkg, isTracked := d.trackedPackages[packageAbsolutePath]

		// If the package was already tracked,
		// then add it to the source file
		if isTracked {
			goto next
		}

		// Try to find the package in the filesystem
		_, err = os.Stat(packageAbsolutePath)
		// If no error, then the package exists,
		// and it is possible to build it
		if err == nil {
			ident := runtime.NewGDStrIdent(nodePackage.GetName())
			pkg = NewSourcePackage(ident, packageAbsolutePath)

			err := d.BuildPackage(pkg)
			if err != nil {
				return nil, err
			}

			goto next
		}

		// Check if the package was not in the filesystem
		// try to find it in the builtin packages
		if builtInPackage, ok := builtin.Packages[nodePackage.GetPath()]; ok {
			pkg = &BuiltInPackage{builtInPackage}
			goto next
		}

		return nil, ErrorAt(nodePackage.GetPosition()).PackageNotFound(nodePackage.GetName())

	next:
		for _, ident := range nodePackage.Imports {
			identNode, isIdentNode := ident.(*ast.NodeIdent)
			if !isIdentNode {
				panic("Invalid node type: expected *ast.NodeIdent")
			}

			// Check if `use` imports exist, and add them as a reference in the source file
			// as an indication of all usages of objects for the source file
			ident := runtime.NewGDStrIdent(identNode.Lit)
			switch pkg := pkg.(type) {
			case *SourcePackage:
				nod, err := pkg.GetMember(ident)
				if err != nil {
					return nil, ErrorAt(identNode.GetPosition()).PackageObjectWasNotFound(ident.ToString(), nodePackage.GetName())
				}

				// Add the node reference to the source file
				err = sourceFile.AddPublic(ident, nod)
				if err != nil {
					return nil, ErrorAt(nod.Node.GetPosition()).DuplicatedObject(ident.ToString())
				}
			case *BuiltInPackage:
				obj, err := pkg.GetMember(ident)
				if err != nil {
					return nil, ErrorAt(identNode.GetPosition()).PackageObjectWasNotFound(ident.ToString(), nodePackage.GetName())
				}

				// Add the object reference to the source file
				err = sourceFile.AddPublic(ident, obj)
				if err != nil {
					return nil, ErrorAt(identNode.GetPosition()).DuplicatedObject(ident.ToString())
				}
			default:
				panic("Invalid package type, expected *SourcePackage or *BuiltInPackage")
			}
		}

		if !isTracked {
			d.trackedPackages[packageAbsolutePath] = pkg

			// Set inferred values for the package
			nodePackage.InferredPath = packagePath
			nodePackage.InferredAbsolutePath = packageAbsolutePath
			nodePackage.InferredMode = pkg.GetMode()

			// Append the package to the hierarchy
			ident := "package: " + nodePackage.GetName() + "@" + sourceFile.file.Name()
			if !d.trackIdent(ident) {
				d.Nodes = append(d.Nodes, nodePackage)
			}
		}

		// Register the package into the source file
		sourceFile.packages = append(sourceFile.packages, pkg)
	}

	return sourceFile, nil
}

func (d *PackageDependenciesAnalyzer) Dispose() {
	d.astBuilder.Dispose()
	d.scanFileSet.Reset()
	d.mainPackage = nil
	d.trackedPackages = nil
	d.trackedNodes = nil
	d.Nodes = nil
	d.MainEntry = nil
}

func (d *PackageDependenciesAnalyzer) analyzeType(typ runtime.GDTypable, astNode ast.Node, sourceFile *SourceFile) error {
	switch typ := typ.(type) {
	case runtime.GDIdent:
		typeIdent := typ.ToString()
		nodeIdent := runtime.NewGDStrIdent(typeIdent)
		if nodeRef := d.getNodeReference(nodeIdent, sourceFile); nodeRef != nil {
			switch nodeRef := nodeRef.(type) {
			case *NodeWithSourceFile:
				return d.analyzeNode(nodeRef.Node, nodeRef.SourceFile)
			case *runtime.GDSymbol:
				// Nothing to do
			default:
				panic("Invalid member type, expected *NodeWithSourceFile or *runtime.GDSymbol")
			}
		}

		return nil
	case runtime.GDUnionType:
		for _, typ := range typ {
			err := d.analyzeType(typ, astNode, sourceFile)
			if err != nil {
				return err
			}
		}

		return nil
	case *runtime.GDArrayType:
		return d.analyzeType(typ.SubType, astNode, sourceFile)
	case runtime.GDTupleType:
		for _, typ := range typ {
			err := d.analyzeType(typ, astNode, sourceFile)
			if err != nil {
				return err
			}
		}

		return nil
	case runtime.GDStructType:
		for _, attr := range typ {
			err := d.analyzeType(attr.Type, astNode, sourceFile)
			if err != nil {
				return err
			}
		}

		return nil
	case *runtime.GDLambdaType:
		for _, arg := range typ.ArgTypes {
			err := d.analyzeType(arg.Value, astNode, sourceFile)
			if err != nil {
				return err
			}
		}

		return d.analyzeType(typ.ReturnType, astNode, sourceFile)
	case runtime.GDType:
		// Nothing to do
		return nil
	default:
		panic("Type not supported")
	}
}

func (d *PackageDependenciesAnalyzer) analyzeNode(astNode ast.Node, sourceFile *SourceFile) error {
	switch astNode := astNode.(type) {
	case *ast.NodeLiteral:
		// Nothing to do
		return nil
	case *ast.NodeIdent:
		nodeIdent := runtime.NewGDStrIdent(astNode.Lit)
		if nodeRef := d.getNodeReference(nodeIdent, sourceFile); nodeRef != nil {
			switch nodeRef := nodeRef.(type) {
			case *NodeWithSourceFile:
				return d.analyzeNode(nodeRef.Node, nodeRef.SourceFile)
			case *runtime.GDSymbol:
				// Nothing to do
			default:
				panic("Invalid member type, expected *NodeWithSourceFile or *runtime.GDSymbol")
			}
		}

		return nil
	case *ast.NodeFunc:
		ident := astNode.Ident.Lit + "@" + sourceFile.file.Name()
		if d.trackIdent(ident) {
			return nil
		}

		err := d.analyzeNode(astNode.NodeLambda, sourceFile)
		if err != nil {
			return err
		}

		// Only append and analyse the node if it is part of the first citizen objects in the source file
		nodeIdent := runtime.NewGDStrIdent(astNode.Ident.Lit)
		if d.getNodeReference(nodeIdent, sourceFile) != nil {
			d.Nodes = append(d.Nodes, astNode)
		}

		return nil
	case *ast.NodeLambda:
		return d.analyzeNode(astNode.Block, sourceFile)
	case *ast.NodeBlock:
		for _, node := range astNode.Nodes {
			err := d.analyzeNode(node, sourceFile)
			if err != nil {
				return err
			}
		}

		return nil
	case *ast.NodeExprOperation:
		if astNode.R != nil {
			return d.analyzeNode(astNode.R, sourceFile)
		}

		if astNode.L != nil {
			err := d.analyzeNode(astNode.L, sourceFile)
			if err != nil {
				return err
			}
		}

		return nil
	case *ast.NodeEllipsisExpr:
		return d.analyzeNode(astNode.Expr, sourceFile)
	case *ast.NodeTuple:
		for i := len(astNode.Nodes) - 1; i >= 0; i-- {
			node := astNode.Nodes[i]
			err := d.analyzeNode(node, sourceFile)
			if err != nil {
				return err
			}
		}

		return nil
	case *ast.NodeStruct:
		for i := len(astNode.Nodes) - 1; i >= 0; i-- {
			node := astNode.Nodes[i]
			err := d.analyzeNode(node, sourceFile)
			if err != nil {
				return err
			}
		}

		return nil
	case *ast.NodeStructAttr:
		return d.analyzeNode(astNode.Expr, sourceFile)
	case *ast.NodeArray:
		for i := len(astNode.Nodes) - 1; i >= 0; i-- {
			node := astNode.Nodes[i]
			err := d.analyzeNode(node, sourceFile)
			if err != nil {
				return err
			}
		}

		return nil
	case *ast.NodeReturn:
		if astNode.Expr != nil {
			return d.analyzeNode(astNode.Expr, sourceFile)
		}

		return nil
	case *ast.NodeBreak:
		// Nothing to do
		return nil
	case *ast.NodeIterIdxExpr:
		err := d.analyzeNode(astNode.Expr, sourceFile)
		if err != nil {
			return err
		}

		return d.analyzeNode(astNode.IdxExpr, sourceFile)
	case *ast.NodeCallExpr:
		for i := len(astNode.Args) - 1; i >= 0; i-- {
			arg := astNode.Args[i]
			err := d.analyzeNode(arg, sourceFile)
			if err != nil {
				return err
			}
		}

		return d.analyzeNode(astNode.Expr, sourceFile)
	case *ast.NodeSafeDotExpr:
		err := d.analyzeNode(astNode.Expr, sourceFile)
		if err != nil {
			return err
		}

		// TODO: Check if this ident needs to be analyzed
		return d.analyzeNode(astNode.Ident, sourceFile)
	case *ast.NodeSets:
		for _, node := range astNode.Nodes {
			err := d.analyzeNode(node, sourceFile)
			if err != nil {
				return err
			}
		}

		return nil
	case *ast.NodeSet:
		ident := astNode.IdentWithType.Ident.Lit + "@" + sourceFile.file.Name()
		if d.trackIdent(ident) {
			return nil
		}

		if astNode.IdentWithType.Type != nil {
			err := d.analyzeType(astNode.IdentWithType.Type, astNode, sourceFile)
			if err != nil {
				return err
			}
		}

		if astNode.Expr != nil {
			err := d.analyzeNode(astNode.Expr, sourceFile)
			if err != nil {
				return err
			}
		}

		nodeIdent := runtime.NewGDStrIdent(astNode.IdentWithType.Ident.Lit)
		if d.getNodeReference(nodeIdent, sourceFile) != nil {
			d.Nodes = append(d.Nodes, astNode)
		}

		return nil
	case *ast.NodeSharedExpr:
		return d.analyzeNode(astNode.Expr, sourceFile)
	case *ast.NodeUpdateSet:
		err := d.analyzeNode(astNode.Expr, sourceFile)
		if err != nil {
			return err
		}

		return d.analyzeNode(astNode.IdentExpr, sourceFile)
	case *ast.NodeLabel:
		// Nothing to do for now!
		return nil
	case *ast.NodeIf:
		for _, node := range astNode.Conditions {
			err := d.analyzeNode(node, sourceFile)
			if err != nil {
				return err
			}
		}

		return d.analyzeNode(astNode.Block, sourceFile)
	case *ast.NodeIfElse:
		err := d.analyzeNode(astNode.If, sourceFile)
		if err != nil {
			return err
		}

		for _, node := range astNode.ElseIf {
			err := d.analyzeNode(node, sourceFile)
			if err != nil {
				return err
			}
		}

		if astNode.Else != nil {
			return d.analyzeNode(astNode.Else, sourceFile)
		}

		return nil
	case *ast.NodeTernaryIf:
		err := d.analyzeNode(astNode.Expr, sourceFile)
		if err != nil {
			return err
		}

		err = d.analyzeNode(astNode.Then, sourceFile)
		if err != nil {
			return err
		}

		return d.analyzeNode(astNode.Else, sourceFile)
	case *ast.NodeForIn:
		err := d.analyzeNode(astNode.Sets, sourceFile)
		if err != nil {
			return err
		}

		err = d.analyzeNode(astNode.Expr, sourceFile)
		if err != nil {
			return err
		}

		return d.analyzeNode(astNode.Block, sourceFile)
	case *ast.NodeForIf:
		if astNode.Sets != nil {
			err := d.analyzeNode(astNode.Sets, sourceFile)
			if err != nil {
				return err
			}
		}

		if astNode.Conditions != nil {
			for _, node := range astNode.Conditions {
				err := d.analyzeNode(node, sourceFile)
				if err != nil {
					return err
				}
			}
		}

		return d.analyzeNode(astNode.Block, sourceFile)
	case *ast.NodeMutCollectionOp:
		err := d.analyzeNode(astNode.L, sourceFile)
		if err != nil {
			return err
		}

		return d.analyzeNode(astNode.R, sourceFile)
	case *ast.NodeTypeAlias:
		ident := astNode.Ident.Lit + "@" + sourceFile.file.Name()
		if d.trackIdent(ident) {
			return nil
		}

		err := d.analyzeType(astNode.Type, astNode, sourceFile)
		if err != nil {
			return err
		}

		nodeIdent := runtime.NewGDStrIdent(astNode.Ident.Lit)
		if d.getNodeReference(nodeIdent, sourceFile) != nil {
			d.Nodes = append(d.Nodes, astNode)
		}

		return nil
	case *ast.NodeCastExpr:
		err := d.analyzeType(astNode.Type, astNode, sourceFile)
		if err != nil {
			return err
		}

		err = d.analyzeNode(astNode.Expr, sourceFile)
		if err != nil {
			return err
		}

		return nil
	default:
		panic("Node type not supported")
	}
}

// Check if the identifier is a public object,
// but it is not found, then it should not throw an error
// because it could be a local object that must be evaluated later
// in another stage with a stack
func (d *PackageDependenciesAnalyzer) getNodeReference(ident runtime.GDIdent, sourceFile *SourceFile) any {
	// References are first looked up in the source file,
	// with local references having higher priority
	fileRef, _ := sourceFile.GetMember(ident)
	if fileRef != nil {
		return fileRef
	}

	// If the reference was not found in the source file,
	// then look up in the parent package, where all public references are stored
	packageRef, _ := sourceFile.parentPackage.GetMember(ident)
	if packageRef != nil {
		return packageRef
	}

	return nil
}

func (d *PackageDependenciesAnalyzer) trackIdent(ident string) bool {
	// Check if the node was already tracked, if so, return true
	if _, isTracked := d.trackedNodes[ident]; isTracked {
		return true
	}

	// Register the node as tracked
	d.trackedNodes[ident] = true

	// If false, means that the has not been tracked yet
	return false
}

func (d *PackageDependenciesAnalyzer) addMemberToSourceFile(ident runtime.GDIdent, isPub bool, node ast.Node, sourceFile *SourceFile) error {
	sourceNode := &NodeWithSourceFile{node, sourceFile}
	if isPub {
		err := sourceFile.parentPackage.AddPublic(ident, sourceNode)
		if err != nil {
			return err
		}

		return sourceFile.AddPublic(ident, sourceNode)
	}

	err := sourceFile.parentPackage.AddLocal(ident, sourceNode)
	if err != nil {
		return err
	}

	return sourceFile.AddLocal(ident, sourceNode)
}

func NewPackageDependenciesAnalyzer() *PackageDependenciesAnalyzer {
	return &PackageDependenciesAnalyzer{
		astBuilder:      ast.NAstBuilderProc(),
		scanFileSet:     scanner.NewFileSet(),
		trackedPackages: make(map[string]Package),
		trackedNodes:    make(map[string]bool),
		Nodes:           make([]ast.Node, 0),
	}
}
