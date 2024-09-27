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

package staticcheck_test

import (
	"gdlang/lib/builtin"
	"gdlang/lib/runtime"
	"gdlang/src/comn"
	"gdlang/src/gd/analysis"
	"gdlang/src/gd/staticcheck"
	"gdlang/src/test_helper"
	"strings"
	"testing"
)

func init() {
	comn.PrettyPrintErrors = false
}

func TestCases(t *testing.T) {
	stack := runtime.NewGDSymbolStack()
	err := builtin.Import(stack)
	if err != nil {
		t.Fatal(err)
	}

	dependencyAnalyzer := analysis.NewPackageDependenciesAnalyzer()
	staticCheck := staticcheck.NewStaticCheck(dependencyAnalyzer)

	type Test struct {
		SrcNode test_helper.FNode
		ErrMsg  string // Expected error
	}

	evalTests := func(tests []Test) error {
		for _, test := range tests {
			err := test_helper.BuildPackageTree(test.SrcNode, func(tmpDir string) error {
				mainStack := stack.NewSymbolStack(runtime.GlobalCtx)
				defer (func() {
					mainStack.Dispose()
					dependencyAnalyzer.Dispose()
				})()

				err := dependencyAnalyzer.Analyze(tmpDir, analysis.PackageDependenciesAnalyzerOptions{ShouldLookUpFromMain: false})
				if err != nil {
					return err
				}

				err = staticCheck.Check(mainStack)
				if err != nil {
					return err
				}

				return nil
			})

			if err != nil {
				if test.ErrMsg == "" {
					t.Errorf("Expected no error but got %q", err.Error())
				} else if !strings.Contains(err.Error(), test.ErrMsg) {
					t.Errorf("Expected an error to contain %q but got %q", test.ErrMsg, err.Error())
				}

				continue
			} else if test.ErrMsg != "" {
				t.Errorf("Expected an error to contain %q but got no error", test.ErrMsg)
			}
		}

		return nil
	}

	t.Run("Test typealiases", func(t *testing.T) {
		// Test typealiases
		typealiasCases := []Test{
			{test_helper.NMFile(`typealias MyInt = int
		pub func main() {
			set a: MyInt = 1
		}`), ""},
		}

		err := evalTests(typealiasCases)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Test returns", func(t *testing.T) {
		returnCases := []Test{
			{test_helper.NMFile(`pub func main() { return 1; }`), "expected `nil` but got `int`"},
			{test_helper.NMFile(`pub func main() { return "string"; }`), "expected `nil` but got `string`"},
			{test_helper.NMFile(`pub func main() { return 1.0; }`), "expected `nil` but got `float`"},
		}

		err := evalTests(returnCases)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Test packages", func(t *testing.T) {
		// General packageCases
		packageCases := []Test{
			// Check use of packages
			{test_helper.NRDir(
				test_helper.NDir(
					"pkg1",
					test_helper.NFile(
						"pkg1.gd",
						`pub func pkg1() => int {
							return 1
						}
						pub func pkg1_2() => string {
							return 1
						}`,
					),
				),
				test_helper.NMFile(
					`use pkg1{pkg1}
					pub func main() {
						set p1 = pkg1()
					}`,
				),
			), "expected `string` but got `int`"},
			// Pkg1 is never used, so it should not throw an error even though it returns an int while expected a string
			{test_helper.NRDir(
				test_helper.NDir("pkg1",
					test_helper.NFile(
						"pkg1.gd",
						`pub func pkg1() => string {
							return 1
						}`,
					),
				),
				test_helper.NMFile(`pub func main() {}`),
			), ""},
			// Sub-package non used function should throw an error
			{test_helper.NRDir(
				test_helper.NDir("pkg1",
					test_helper.NDir("pkg2",
						test_helper.NFile(
							"pkg2.gd",
							`pub func pkg2() => string {
								return "pkg2"
							}
							pub func pkg2_2() => int {
								return "wrong"
							}`,
						),
					),
					test_helper.NFile(
						"pkg1.gd",
						`use pkg1.pkg2 {pkg2}
						pub func pkg1() => string {
							return pkg2()
						}`,
					),
				),
				test_helper.NMFile(
					`use pkg1 {pkg1}
					pub func main() {
						set p1 = pkg1()
					}`,
				),
			), "expected `int` but got `string`"},
			// Try to use outer package function in inner package
			{test_helper.NRDir(
				test_helper.NDir(
					"pkg1",
					test_helper.NDir(
						"pkg2",
						test_helper.NFile(
							"pkg2.gd",
							`use pkg1 {pkg1}
							pub func pkg2() => string {
								return pkg1()
							}`,
						),
					),
					test_helper.NFile(
						"pkg1.gd",
						`pub func pkg1() => string {
							return "hola, soy pkg1"
						}`,
					),
				),
				test_helper.NMFile(
					`use pkg1.pkg2 {pkg2}
					pub func main() {
						set p1 = pkg2()
					}`,
				),
			), ""},
		}

		err := evalTests(packageCases)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Test sets", func(t *testing.T) {
		// Test sets
		SetCases := []Test{
			{test_helper.NMFile(`pub func main() {
			set a = 1, b = 2
		}`), ""},
			{test_helper.NMFile(`pub func main() {
				set (a, b) = (1, 2)
				a = 2
				b = 1
			}`), ""},
			{test_helper.NMFile(`pub func main() {
			set a, b = (1, 2)
			a = 2
			b = 1
		}`), "main.gd 4:8-8 fatal error: expected `(int, int)` but got `int`"},
			// Try to destructure an array
			{test_helper.NMFile(`pub func main() {
			set (a, b) = [1, 2]
			a = 2
			b = 1
		}`), ""},
			{test_helper.NMFile(`pub func main() {
			set (a, b) = [1, 2]
			a = 2
			b = 1
		}`), ""},
			{test_helper.NMFile(`pub func main() {
			set (a, b) = [1, 2, 3], c = 1
			a = 1
			b = 2
			c = 3
		}`), ""},
		}

		err := evalTests(SetCases)
		if err != nil {
			t.Fatal(err)
		}
	})
}
