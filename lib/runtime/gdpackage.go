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

package runtime

type GDPackageMode byte

const (
	PackageModeSource GDPackageMode = iota
	PackageModeBuiltin
	PackageModeLib
)

// Map of package types
var PackageModeMap = map[GDPackageMode]string{
	PackageModeSource:  "sourced",
	PackageModeBuiltin: "builtin",
	PackageModeLib:     "lib",
}

type GDMemberType byte

const (
	GDMemberPublic GDMemberType = iota
	GDMemberLocal
)

type GDMember[T any] struct {
	Type  GDMemberType
	Value T
}

type GDMembers[T any] map[any]GDMember[T]

type GDPackage[T any] struct {
	Ident   GDIdent
	Path    string        // The name and the path of the package
	Mode    GDPackageMode // The type of the package
	Members GDMembers[T]  // A list of public references in the package
}

func (p *GDPackage[T]) AddPublic(ident GDIdent, member T) error {
	return p.add(ident, GDMemberPublic, member)
}

func (p *GDPackage[T]) AddLocal(ident GDIdent, member T) error {
	return p.add(ident, GDMemberLocal, member)
}

func (p *GDPackage[T]) GetMember(ident GDIdent) (T, error) {
	if member, ok := p.Members[ident.GetRawValue()]; ok {
		return member.Value, nil
	}

	var zT T
	return zT, ObjectNotFoundErr(ident.ToString())
}

func (p *GDPackage[T]) add(ident GDIdent, memberType GDMemberType, value T) error {
	if _, ok := p.Members[ident.GetRawValue()]; ok {
		return DuplicatedObjectCreationErr(ident)
	}

	p.Members[ident.GetRawValue()] = GDMember[T]{memberType, value}

	return nil
}

func NewGDPackage[T any](name GDIdent, path string, typ GDPackageMode) *GDPackage[T] {
	return &GDPackage[T]{name, path, typ, make(GDMembers[T])}
}
