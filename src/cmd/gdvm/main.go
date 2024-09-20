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
	"gdlang/src/vm"
	"os"
)

var (
	version     = "dev"
	buildNumber = "0"
	arch        = "host"
)

var (
	gdbin       = flag.String("gdbin", "", "path to the GDLang binary file")
	showVersion = flag.Bool("version", false, "prints the GDLang VM version")
)

func main() {
	flag.Parse()

	if *showVersion {
		versionMsg := comn.NewMarkdown("GDLang VM `version`: " + version + " `build`: " + buildNumber + " `arch`: " + arch)
		println(versionMsg.Stylize())
		os.Exit(0)
	}

	bytes, err := os.ReadFile(*gdbin)
	if err != nil {
		print(err.Error())
		os.Exit(1)
	}

	vmProc := vm.NewGDVMProc()
	err = vmProc.Init(bytes)
	if err != nil {
		print(err.Error())
		os.Exit(1)
	}

	err = vmProc.Run()
	if err != nil {
		print(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
