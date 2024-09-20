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
	// Used internally to spot mistakes
	PenaltyDebugErrCode ErrCode = iota

	// Execution error might be caused by a runtime error
	RuntimeExceptionErrCode

	// A general fatal error code
	DefaultFatalErrCode

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

	// A statement was not terminated correctly, missing braces, semicolons, etc.
	WrongEndOfStatementErrCode

	// Invoker is nil, this happends when trying either to call a method on a nil object
	// or accessing a property of a nil object
	InvokerIsNilErrCode

	// Not a function when invoking a function
	InvokerNotAFunctionErrCode

	// Not valid spreadable type
	InvalidSpreadableTypeErrCode

	// A package was not found
	PackageNotFoundErrCode

	// Error was found while parsing a source code from a path
	PackageSourceCodeFileParsingErrCode

	// A public object was not found in the package
	PublicObjectNotFoundErrCode

	// Duplicated pub object
	DuplicatedPublicObjectErrCode

	// Use directive can only be used at the header of a file example:
	// use io::read{*}
	UseOnlyAtHeaderErrCode
)
