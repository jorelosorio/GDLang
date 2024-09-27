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
	"gdlang/src/comn"
	"gdlang/src/gd/scanner"
)

const (
	ErrorPackageNotFound comn.ErrCode = 100 + iota
	ErrorReadingSourceFile
	ErrorPackageObjectNotFound
)

type ErrorAt scanner.Position

func (e ErrorAt) PackageNotFound(pkgName string) comn.Error {
	return comn.NewErrorf(ErrorPackageNotFound, comn.FatalError, scanner.Position(e), "package `%s` was not found", pkgName)
}

func (e ErrorAt) ReadingSourceFile(filePath string) comn.Error {
	return comn.NewErrorf(ErrorReadingSourceFile, comn.FatalError, scanner.Position(e), "an error occurred reading `%s`", filePath)
}

func (e ErrorAt) PackageObjectWasNotFound(objName, pkgName string) comn.Error {
	return comn.NewErrorf(ErrorPackageObjectNotFound, comn.FatalError, scanner.Position(e), "public object `%s` was not found in package `%s`", objName, pkgName)
}

func (e ErrorAt) DuplicatedObject(objName string) comn.Error {
	return comn.NewErrorf(ErrorPackageObjectNotFound, comn.FatalError, scanner.Position(e), "duplicated object `%s`", objName)
}

func (e ErrorAt) MainEntryWasNotFound() comn.Error {
	return comn.NewErrorf(ErrorPackageObjectNotFound, comn.FatalError, scanner.Position(e), "no `main` function was found in the package")
}
