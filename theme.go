// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type Theme interface {
	Driver() Driver
	DefaultFont() Font
	SetDefaultFont(Font)
	DefaultMonospaceFont() Font
	SetDefaultMonospaceFont(Font)
	CreateBubbleOverlay() BubbleOverlay
	CreateButton() Button
	CreateCodeEditor() CodeEditor
	CreateDropDownList() DropDownList
	CreateImage() Image
	CreateLabel() Label
	CreateLinearLayout() LinearLayout
	CreateList() List
	CreatePanelHolder() PanelHolder
	CreateProgressBar() ProgressBar
	CreateScrollBar() ScrollBar
	CreateScrollLayout() ScrollLayout
	CreateSplitterLayout() SplitterLayout
	CreateTableLayout() TableLayout
	CreateTextBox() TextBox
	CreateTree() Tree
	CreateWindow(width, height int, title string) Window
}
