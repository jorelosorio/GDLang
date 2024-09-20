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

import "testing"

func TestSimpleBacltick(t *testing.T) {
	s := "`a`"
	result := NewMarkdown(s).Stylize()

	if result != Colorize(PrimaryHighlightColor, "a") {
		t.Errorf("Expected %s but got %s", Colorize(PrimaryHighlightColor, "a"), result)
	}
}

func TestMultipleBackticks(t *testing.T) {
	s := "`a` `b`"
	result := NewMarkdown(s).Stylize()

	expected := Colorize(PrimaryHighlightColor, "a") + Colorize(PrimaryTextColor, " ") + Colorize(PrimaryHighlightColor, "b")
	if result != expected {
		t.Errorf("Expected %q but got %q", expected, result)
	}
}

func TestNoBackticks(t *testing.T) {
	s := "a"
	result := NewMarkdown(s).Stylize()

	if result != Colorize(PrimaryTextColor, "a") {
		t.Errorf("Expected a but got %s", result)
	}
}

func TestEmptyString(t *testing.T) {
	s := ""
	result := NewMarkdown(s).Stylize()
	if result != "" {
		t.Error("Expected empty string but got", result)
	}
}

func TestBacktickAtEnd(t *testing.T) {
	s := "`a` `b` `"
	result := NewMarkdown(s).Stylize()

	if result != Colorize(PrimaryHighlightColor, "a")+Colorize(PrimaryTextColor, " ")+Colorize(PrimaryHighlightColor, "b")+Colorize(PrimaryTextColor, " ") {
		t.Errorf("Expected %s but got %s", Colorize(PrimaryHighlightColor, "a")+" "+Colorize(PrimaryHighlightColor, "b")+" ", result)
	}
}
