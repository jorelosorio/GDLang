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

import (
	"fmt"
	"gdlang/lib/runtime"
	"gdlang/src/gd/scanner"
	"strings"
)

type SeverityType int8

var errorMessage = map[ErrCode]string{
	LabelErrCode:           "invalid label",
	DuplicatedLabelErrCode: "duplicated label `%@`",
}

func (e ErrCode) formatErrorMessage(params ...string) string {
	return fmt.Sprintf(errorMessage[e], params)
}

func (e ErrCode) Errorf(params ...string) ProcErr {
	return ProcErr{Code: e, Msg: e.formatErrorMessage(params...), Severity: FatalError}
}

const (
	WrongEndOfStatementErrMsg            = "wrong statement termination"
	NilAsATypeErrMsg                     = "hold up! `What's going on here?`, assigning `nil` as a type is not allowed"
	InvalidArraySpreadExpressionErrorMsg = "ellipsis expression can only be used in tuples or arrays"
	PackageNotFoundMsg                   = "package `%s` was not found"
	PackageSourceCodeFileParsingErrMsg   = "an error occurred while parsing a source code file from package `%s`"
	PublicObjectNotFoundErrMsg           = "public object `%s` was not found in package `%s`"
	UseOnlyAtHeaderErrMsg                = "`use` directive can only be used at the header of a file, not between statements"
	NoMainFunctionErrMsg                 = "no `main` function was found in the package"
	DuplicatedPublicObjectErrMsg         = "an object `%s` was already declared in the package `%s`"
	MisplacedBreakErrMsg                 = "`break` statement is not allowed here, it can only be used inside a control flow statement"
	NilAccessExceptionErrMsg             = "a `nil` was encountered while dereferencing an object"
)

const (
	FatalError SeverityType = iota
	Warning
	SyntaxError
)

var SeverityMap = map[SeverityType]string{
	FatalError:  "fatal error",
	Warning:     "warning",
	SyntaxError: "syntax error",
}

type ProcErr struct {
	Position scanner.Position
	Code     ErrCode
	Msg      string
	Severity SeverityType
	Hints    []string
}

func (e ProcErr) Error() string {
	hints := ""
	if len(e.Hints) > 0 {
		hints = "\n\t " + strings.Join(e.Hints, "\n\t-")
	}

	var positionInfo string
	if PrettyPrintErrors {
		var severityColor = ErrorHighlightColor
		if e.Severity == Warning {
			severityColor = WarningHighlightColor
		}

		severityStr := Colorizef(severityColor, "%s:", SeverityMap[e.Severity])
		fileName := Colorizef(TitleColor, "%s", e.Position.Filename)
		if e.Position.Line != 0 {
			positionInfo = Colorizef(PrimaryTextColor, " %d:%d-%d", e.Position.Line, e.Position.ColStart, e.Position.ColEnd)
		}
		errMsg := NewMarkdown(e.Msg).Stylize()
		return fmt.Sprintf("%s%s %s %s%s", fileName, positionInfo, severityStr, errMsg, hints)
	} else {
		if e.Position.Line != 0 {
			positionInfo = fmt.Sprintf(" %d:%d-%d", e.Position.Line, e.Position.ColStart, e.Position.ColEnd)
		}
		return fmt.Sprintf("%s%s %s: %s%s", e.Position.Filename, positionInfo, SeverityMap[e.Severity], e.Msg, hints)
	}
}

func NErr(code ErrCode, msg string, severity SeverityType, errPos scanner.Position, hints []string) ProcErr {
	return ProcErr{errPos, code, msg, severity, hints}
}

func WrapSyntaxErr(err error, errPos scanner.Position) ProcErr {
	return NErr(DefaultSyntaxErrCode, err.Error(), SyntaxError, errPos, nil)
}

func WrapFatalErr(err error, errPos scanner.Position) ProcErr {
	return NErr(DefaultFatalErrCode, err.Error(), FatalError, errPos, nil)
}

func CompilerErr(msg string, errPos scanner.Position) ProcErr {
	return NErr(DefaultCompilerErrCode, msg, FatalError, errPos, nil)
}

func AnalysisErr(msg string, errPos scanner.Position) ProcErr {
	return NErr(DefaultAnalysisErrCode, msg, FatalError, errPos, nil)
}

func DispatchCompilerWarning(msg string, errPos scanner.Position) {
	fmt.Println(NErr(DefaultCompilerWarningWarningCode, msg, Warning, errPos, nil))
}

func WrapCompilerErr(err error, errPos scanner.Position) ProcErr {
	return NErr(DefaultCompilerErrCode, err.Error(), FatalError, errPos, nil)
}

func CreateUnexpectedSyntaxError(unexpected string, expecting []string, errPos scanner.Position) ProcErr {
	switch len(expecting) {
	case 0:
		msg := fmt.Sprintf("an unexpected `%s` was encountered", unexpected)
		return NErr(DefaultSyntaxErrCode, msg, SyntaxError, errPos, nil)
	case 1:
		msg := fmt.Sprintf("an unexpected `%s` was encountered when expecting a `%s`", unexpected, expecting[0])
		return NErr(DefaultSyntaxErrCode, msg, SyntaxError, errPos, nil)
	case 2:
		msg := fmt.Sprintf("an unexpected `%s` was encountered when expecting a `%s` or `%s`", unexpected, expecting[0], expecting[1])
		return NErr(DefaultSyntaxErrCode, msg, SyntaxError, errPos, nil)
	default:
		msg := fmt.Sprintf("an unexpected `%s` was encountered when expecting one of the following: %s", unexpected, runtime.JoinSlice(expecting, func(token string, _ int) string {
			return fmt.Sprintf("`%s`", token)
		}, ", "))
		return NErr(DefaultSyntaxErrCode, msg, SyntaxError, errPos, nil)
	}
}

func CreateExpSymbolsErr(expecting []string, errPos scanner.Position) ProcErr {
	switch len(expecting) {
	case 1:
		msg := fmt.Sprintf("expected a `%s`", expecting[0])
		return NErr(DefaultSyntaxErrCode, msg, SyntaxError, errPos, nil)
	case 2:
		msg := fmt.Sprintf("expected a `%s` or `%s`", expecting[0], expecting[1])
		return NErr(DefaultSyntaxErrCode, msg, SyntaxError, errPos, nil)
	default:
		msg := fmt.Sprintf("expected one of the following symbols: %s", runtime.JoinSlice(expecting, func(token string, _ int) string {
			return fmt.Sprintf("`%s`", token)
		}, ", "))
		return NErr(DefaultSyntaxErrCode, msg, SyntaxError, errPos, nil)
	}
}

func GDFileParsingErr(pkgName string, errPos scanner.Position) ProcErr {
	msg := fmt.Sprintf(PackageSourceCodeFileParsingErrMsg, pkgName)
	return NErr(PackageSourceCodeFileParsingErrCode, msg, FatalError, errPos, nil)
}

func CreatePackageNotFoundErr(pkgName string, errPos scanner.Position) ProcErr {
	msg := fmt.Sprintf(PackageNotFoundMsg, pkgName)
	return NErr(PackageNotFoundErrCode, msg, FatalError, errPos, nil)
}
