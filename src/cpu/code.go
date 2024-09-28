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

type GDReg byte
type GDInst byte

const (
	TypeAlias   GDInst = iota // Define a type alias
	CastObj                   // Cast an object
	Use                       // Use a package
	Set                       // Instance a new object
	Mov                       // Move a value to a register
	CSet                      // Set a value to an iterable
	CAdd                      // Add a value to a collectable
	CRemove                   // Remove a value from an collectable
	IGet                      // Get a value from an iterable
	ILen                      // Get the length of an iterable
	AGet                      // Get a value from a struct
	Lambda                    // Define a lambda
	BBegin                    // Define a block
	BEnd                      // End a block
	Ret                       // Return
	Call                      // Call a function
	Operation                 // Perform an operation
	Tif                       // If expression
	CompareJump               // Compare two values and jump if equals
	Jump                      // Jump to a label
	Label                     // Define a label
)

const (
	RPop GDReg = iota // Pop register
	Rx                // Return value register
	Ra
	Rb
	Rc
	Rd
	Re
	Rf
	Rg
	Rh
	Ri
)
