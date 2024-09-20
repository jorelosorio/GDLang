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

package test

import (
	"bytes"
	"gdlang/src/comn"
	"gdlang/src/compiler"
	"gdlang/src/test_helper"
	"gdlang/src/vm"
	"os"
	"strings"
	"testing"
)

func init() {
	comn.PrettyPrintErrors = false
}

type Test struct {
	Src    string
	Output string
	ErrMsg string
}

func RunTests(t *testing.T, tests []Test) {
	RunTestsWithTemplate(t, "", tests)
}

func RunTestsWithTemplate(t *testing.T, tmpl string, tests []Test) {
	for _, test := range tests {
		var src string
		if tmpl != "" {
			src = strings.ReplaceAll(tmpl, "$SRC", test.Src)
		} else {
			src = test.Src
		}
		t.Run(src, func(t *testing.T) {
			test_helper.BuildPackageTree(test_helper.NMFile(src), func(tmpDir string) error {
				output := CaptureStdout(func() {
					_, _, err := RunFileTest(tmpDir)
					// TODO: Dispose is failing
					// defer proc.Dispose()
					if err != nil {
						if test.ErrMsg == "" {
							t.Errorf("Expected no errors but got %s when running %s", err.Error(), src)
						} else if !strings.Contains(err.Error(), test.ErrMsg) {
							t.Errorf("Expected error message to contain %q but got: %q", test.ErrMsg, err.Error())
						}

						return
					} else if test.ErrMsg != "" {
						t.Errorf("Expected error message to contain %q but got no error", test.ErrMsg)
					}
				})

				if strings.Contains(output, "warning") && test.ErrMsg != "" {
					if !strings.Contains(output, test.ErrMsg) {
						t.Errorf("Expected %q but got %q when running %s", test.Output, output, src)
					}
					return nil
				}

				if output != test.Output {
					t.Errorf("Expected %q but got %q when running %s", test.Output, output, src)
				}

				return nil
			})
		})
	}
}

func RunTestsWithMainTemplate(t *testing.T, tests []Test) {
	tmpl := `pub func main() {
		$SRC
	}`
	RunTestsWithTemplate(t, tmpl, tests)
}

func RunTest(t *testing.T, test Test) {
	RunTests(t, []Test{test})
}

func RunFileTest(pkgPath string) (*vm.GDVMProc, *compiler.GDCompiler, error) {
	compiler := compiler.NewGDCompiler()
	defer compiler.Dispose()

	err := compiler.Compile(pkgPath)
	if err != nil {
		return nil, nil, err
	}

	buffer := &bytes.Buffer{}
	err = compiler.Root.BuildBytecode(buffer, compiler.Ctx)
	if err != nil {
		print(err.Error())
		os.Exit(1)
	}

	vmProc := vm.NewGDVMProc()
	err = vmProc.Init(buffer.Bytes())
	if err != nil {
		return nil, nil, err
	}

	err = vmProc.Run()
	if err != nil {
		return nil, nil, err
	}

	return vmProc, compiler, nil
}
