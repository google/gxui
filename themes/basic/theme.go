// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package basic

import (
	"github.com/google/gxui"
)

type Theme struct {
	DriverInfo               gxui.Driver
	DefaultFontInfo          gxui.Font
	DefaultMonospaceFontInfo gxui.Font

	WindowBackground gxui.Color

	BubbleOverlayStyle        Style
	ButtonDefaultStyle        Style
	ButtonOverStyle           Style
	ButtonPressedStyle        Style
	CodeSuggestionListStyle   Style
	DropDownListDefaultStyle  Style
	DropDownListOverStyle     Style
	FocusedStyle              Style
	HighlightStyle            Style
	LabelStyle                Style
	PanelBackgroundStyle      Style
	ScrollBarBarDefaultStyle  Style
	ScrollBarBarOverStyle     Style
	ScrollBarRailDefaultStyle Style
	ScrollBarRailOverStyle    Style
	SplitterBarDefaultStyle   Style
	SplitterBarOverStyle      Style
	TabActiveHighlightStyle   Style
	TabDefaultStyle           Style
	TabOverStyle              Style
	TabPressedStyle           Style
	TextBoxDefaultStyle       Style
	TextBoxOverStyle          Style
}

// gxui.Theme compliance
func (t *Theme) Driver() gxui.Driver {
	return t.DriverInfo
}

func (t *Theme) DefaultFont() gxui.Font {
	return t.DefaultFontInfo
}

func (t *Theme) SetDefaultFont(f gxui.Font) {
	t.DefaultFontInfo = f
}

func (t *Theme) DefaultMonospaceFont() gxui.Font {
	return t.DefaultMonospaceFontInfo
}

func (t *Theme) SetDefaultMonospaceFont(f gxui.Font) {
	t.DefaultMonospaceFontInfo = f
}

func (t *Theme) CreateBubbleOverlay() gxui.BubbleOverlay {
	return CreateBubbleOverlay(t)
}

func (t *Theme) CreateButton() gxui.Button {
	return CreateButton(t)
}

func (t *Theme) CreateCodeEditor() gxui.CodeEditor {
	return CreateCodeEditor(t)
}

func (t *Theme) CreateDropDownList() gxui.DropDownList {
	return CreateDropDownList(t)
}

func (t *Theme) CreateImage() gxui.Image {
	return CreateImage(t)
}

func (t *Theme) CreateLabel() gxui.Label {
	return CreateLabel(t)
}

func (t *Theme) CreateLinearLayout() gxui.LinearLayout {
	return CreateLinearLayout(t)
}

func (t *Theme) CreateList() gxui.List {
	return CreateList(t)
}

func (t *Theme) CreatePanelHolder() gxui.PanelHolder {
	return CreatePanelHolder(t)
}

func (t *Theme) CreateProgressBar() gxui.ProgressBar {
	return CreateProgressBar(t)
}

func (t *Theme) CreateScrollBar() gxui.ScrollBar {
	return CreateScrollBar(t)
}

func (t *Theme) CreateScrollLayout() gxui.ScrollLayout {
	return CreateScrollLayout(t)
}

func (t *Theme) CreateSplitterLayout() gxui.SplitterLayout {
	return CreateSplitterLayout(t)
}

func (t *Theme) CreateTableLayout() gxui.TableLayout {
	return CreateTableLayout(t)
}

func (t *Theme) CreateTextBox() gxui.TextBox {
	return CreateTextBox(t)
}

func (t *Theme) CreateTree() gxui.Tree {
	return CreateTree(t)
}

func (t *Theme) CreateWindow(width, height int, title string) gxui.Window {
	return CreateWindow(t, width, height, title)
}
