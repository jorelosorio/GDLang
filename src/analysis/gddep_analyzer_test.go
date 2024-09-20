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

package analysis_test

import (
	"gdlang/src/analysis"
	"gdlang/src/comn"
	"gdlang/src/test_helper"
	"strings"
	"testing"
)

func init() {
	comn.PrettyPrintErrors = false
}

func TestDependencyAnalysisCases(t *testing.T) {
	depAnalyzer := analysis.NDepAnalyzerProc()

	tests := []struct {
		ErrMsg string
		Pkgs   test_helper.FNode
	}{
		{
			"public object `pkg1A` was not found in package `pkg1`",
			test_helper.NDir("",
				test_helper.NDir("pkg1", test_helper.NFile("pk1.gd", `set pkg1A = 1`)),
				test_helper.NMFile(`use pkg1.pkg1A
				pub func main() {
					print(pkg1A)
				}`),
			),
		},
		{
			"",
			test_helper.NDir("",
				test_helper.NDir("pkg1", test_helper.NDir("subpkg1", test_helper.NFile("pk1.gd", `pub set a = 1;`))),
				test_helper.NMFile(`use pkg1.subpkg1.a

				pub func main() {
					print(a)
				}`),
			),
		},
	}

	for _, test := range tests {
		err := test_helper.BuildPackageTree(test.Pkgs, func(tmpDir string) error {
			defer depAnalyzer.Dispose()

			return depAnalyzer.Build(tmpDir, analysis.DepAnalyzerOpt{MainAsEntryPoint: false})
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
}
