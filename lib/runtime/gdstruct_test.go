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

package runtime_test

import (
	"gdlang/lib/runtime"
	"strings"
	"testing"
)

func TestStructWithNilInitialization(t *testing.T) {
	stack := runtime.NewRootGDSymbolStack()
	sType := runtime.NewGDStructType(runtime.GDStructAttrType{attr1Ident, runtime.GDStringType})
	structObj, err := runtime.NewGDStruct(sType, stack)
	if err != nil {
		t.Errorf("Error creating struct: %s", err.Error())
	}

	symbol, err := structObj.GetAttr(attr1Ident)
	if err != nil && err.Error() == runtime.AttributeNotFoundErr("attr1").Error() {
		t.Error("Attribute not found but it should be")
	}

	if symbol.Object != runtime.GDZNil {
		t.Errorf("Wrong structure attribute value for attr1, expected nil but got %q", symbol.Object.ToString())
	}
}

func TestChangeTheValueOfAnAttribute(t *testing.T) {
	stack := runtime.NewRootGDSymbolStack()
	sType := runtime.NewGDStructType(runtime.GDStructAttrType{attr1Ident, runtime.GDStringType})

	structObj, err := runtime.NewGDStruct(sType, stack)
	if err != nil {
		t.Errorf("Error creating struct: %s", err.Error())
	}

	_, err = structObj.SetAttr(attr1Ident, runtime.GDString("new value"))
	if err != nil && err.Error() == runtime.AttributeNotFoundErr("attr1").Error() {
		t.Error("Attribute not found but it should be")
	}

	symbol, err := structObj.GetAttr(attr1Ident)
	if err != nil {
		t.Errorf("Error getting attribute value: %s", err.Error())
	}

	if !runtime.EqualObjects(symbol.Object, runtime.GDString("new value")) {
		t.Errorf("Wrong structure attribute value for attr1, expected 'new value' but got %q", symbol.Object.ToString())
	}
}

func TestReturnedObjectFromStructAreCopies(t *testing.T) {
	stack := runtime.NewRootGDSymbolStack()

	sType := runtime.NewGDStructType(runtime.GDStructAttrType{attr1Ident, runtime.GDStringType})

	structObj, err := runtime.NewGDStruct(sType, stack)
	if err != nil {
		t.Errorf("Error creating struct: %s", err.Error())
	}

	_, err = structObj.SetAttr(attr1Ident, runtime.GDString("test"))
	if err != nil {
		t.Errorf("Error setting attribute value: %s", err.Error())
	}

	// Should store a copy of the object with the value "test"
	symbol1, err := structObj.GetAttr(attr1Ident)
	if err != nil {
		t.Errorf("Error getting attribute value: %s", err.Error())
	}
	eObj1 := symbol1.Object

	_, err = structObj.SetAttr(attr1Ident, runtime.GDString("new value"))
	if err != nil {
		t.Errorf("Error setting attribute value: %s", err.Error())
	}

	// Should store a copy of the object with the value "new value"
	symbol2, err := structObj.GetAttr(attr1Ident)
	if err != nil {
		t.Errorf("Error getting attribute value: %s", err.Error())
	}

	// Attr1Value1 and Attr1Value2 should be different objects
	if runtime.EqualObjects(eObj1, symbol2.Object) {
		t.Errorf("Returned object from struct is not a copy, was expecting %q but got %q", symbol1.Object.ToString(), symbol2.Object.ToString())
	}
}

func TestAddAttributeWithSameName(t *testing.T) {
	stack := runtime.NewRootGDSymbolStack()

	sType := runtime.NewGDStructType(runtime.GDStructAttrType{attr1Ident, runtime.GDStringType}, runtime.GDStructAttrType{attr1Ident, runtime.GDIntType})

	_, err := runtime.NewGDStruct(sType, stack)
	if err == nil {
		t.Error("Expected error adding attribute with the same name but got nil")
	}

	if err != nil && err.Error() != runtime.DuplicatedObjectCreationErr("attr1").Error() {
		t.Errorf("Expected error adding attribute with the same name but got %q", err.Error())
	}
}

func TestSetAttrWithDifferentType(t *testing.T) {
	stack := runtime.NewRootGDSymbolStack()

	sType := runtime.NewGDStructType(runtime.GDStructAttrType{attr1Ident, runtime.GDStringType})

	structObj, err := runtime.NewGDStruct(sType, stack)
	if err != nil {
		t.Errorf("Error creating struct: %s", err.Error())
	}

	_, err = structObj.SetAttr(attr1Ident, runtime.GDZInt)
	if err == nil {
		t.Error("Expected error setting attribute with different type but got nil")
	}

	errMsg := runtime.WrongTypesErr(runtime.GDStringType, runtime.GDIntType).Msg
	if err != nil && !strings.Contains(err.Error(), errMsg) {
		t.Errorf("Expected %q but got %q", errMsg, err.Error())
	}
}
