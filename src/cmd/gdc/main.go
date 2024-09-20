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

package main

import (
	"flag"
	"gdlang/src/comn"
	"gdlang/src/compiler"
	"os"
	"path"
)

var (
	version     = "dev"
	buildNumber = "0"
	arch        = "host"
)

var (
	pkgPath     = flag.String("pkg", "", "path to the main package")
	outputPath  = flag.String("o", "", "output path for the compiled files")
	showVersion = flag.Bool("version", false, "prints the GDLang Compiler version")
)

func main() {
	flag.Parse()

	if *showVersion {
		versionMsg := comn.NewMarkdown("GDLang Compiler `version`: " + version + " `build`: " + buildNumber + " `arch`: " + arch)
		println(versionMsg.Stylize())
		os.Exit(0)
	}

	c := compiler.NewGDCompiler()
	defer c.Dispose()

	err := c.Compile(*pkgPath)
	if err != nil {
		print(err.Error())
		os.Exit(1)
	}

	if outputPath == nil || *outputPath == "" {
		*outputPath = *pkgPath
	}

	pkgName := path.Base(*pkgPath)

	err = c.GenerateArtifacts(pkgName, *outputPath)
	if err != nil {
		print(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
