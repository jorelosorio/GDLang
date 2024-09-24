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

package comn

type ErrCode uint16

const (
	// A general fatal error code
	DefaultFatalErrCode ErrCode = iota

	// A general syntax error code
	DefaultSyntaxErrCode

	// A general compiler error code
	DefaultCompilerErrCode

	// Analysis error code
	DefaultAnalysisErrCode

	// A general warning error code
	DefaultCompilerWarningWarningCode

	// Label
	LabelErrCode

	// A duplicated label
	DuplicatedLabelErrCode

	// Trying to use `nil` as a type is not allowed
	NilAsATypeErrCode

	// Not valid spreadable type
	InvalidSpreadableTypeErrCode

	// A package was not found
	PackageNotFoundErrCode

	// Error was found while parsing a source code from a path
	PackageSourceCodeFileParsingErrCode

	// A public object was not found in the package
	PublicObjectNotFoundErrCode

	// Use directive can only be used at the header of a file example:
	// use io::read{*}
	UseOnlyAtHeaderErrCode
)
