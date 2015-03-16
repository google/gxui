// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/google/gxui/math"
)

type TextBox interface {
	Focusable
	OnSelectionChanged(func()) EventSubscription
	OnTextChanged(func([]TextBoxEdit)) EventSubscription
	Padding() math.Spacing
	SetPadding(math.Spacing)
	Runes() []rune
	Text() string
	SetText(string)
	Font() Font
	SetFont(Font)
	Multiline() bool
	SetMultiline(bool)
	DesiredWidth() int
	SetDesiredWidth(desiredWidth int)
	TextColor() Color
	SetTextColor(Color)
	Select(TextSelectionList)
	SelectAll()
	Carets() []int
	RuneIndexAt(p math.Point) (idx int, found bool)
	TextAt(s, e int) string
	WordAt(runeIndex int) string
	ScrollToLine(int)
	ScrollToRune(int)
	LineIndex(runeIndex int) int
	LineStart(line int) int
	LineEnd(line int) int
}
