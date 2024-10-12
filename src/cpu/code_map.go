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

package cpu

var cpuOpMap = map[GDInst]string{
	TypeAlias:   "typealias",
	CastObj:     "cast",
	Set:         "set",
	Use:         "use",
	Lambda:      "lambda",
	BBegin:      "block",
	BEnd:        "end",
	Ret:         "ret",
	Call:        "call",
	Mov:         "mov",
	CSet:        "cset",
	IGet:        "iget",
	ILen:        "ilen",
	AGet:        "aget",
	ASet:        "aset",
	CAdd:        "cadd",
	Tif:         "tif",
	CRemove:     "cremove",
	Operation:   "op",
	CompareJump: "cmpjump",
	Jump:        "jump",
	Label:       "label",
}

var cpuRegMap = map[GDReg]string{
	RPop: "rpop",
	Rx:   "rx",
	Ra:   "ra",
	Rb:   "rb",
	Rc:   "rc",
	Rd:   "rd",
	Re:   "re",
	Rf:   "rf",
	Rg:   "rg",
	Rh:   "rh",
	Ri:   "ri",
}

func GetCPUInstName(op GDInst) string { return cpuOpMap[op] }
func GetCPURegName(reg GDReg) string  { return cpuRegMap[reg] }
