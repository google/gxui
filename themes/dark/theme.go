// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dark

import (
	"fmt"

	"github.com/google/gxui"
	"github.com/google/gxui/gxfont"
)

type Theme struct {
	driver      gxui.Driver
	defaultFont gxui.Font

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

func CreateTheme(driver gxui.Driver) gxui.Theme {
	defaultFont, err := driver.CreateFont(gxfont.Default, 12)
	if err == nil {
		defaultFont.LoadGlyphs(32, 126)
	} else {
		fmt.Printf("Warning: Failed to load default font - %v\n", err)
	}

	scrollBarRailDefaultBg := gxui.Black
	scrollBarRailDefaultBg.A = 0.7

	scrollBarRailOverBg := gxui.Gray20
	scrollBarRailOverBg.A = 0.7

	neonBlue := gxui.ColorFromHex(0xFF5C8CFF)
	focus := gxui.ColorFromHex(0xA0C4D6FF)

	return &Theme{
		driver:      driver,
		defaultFont: defaultFont,

		WindowBackground: gxui.Black,

		//                                   fontColor    brushColor   penColor
		BubbleOverlayStyle:        CreateStyle(gxui.Gray80, gxui.Gray20, gxui.Gray40, 1.0),
		ButtonDefaultStyle:        CreateStyle(gxui.Gray80, gxui.Gray10, gxui.Gray20, 1.0),
		ButtonOverStyle:           CreateStyle(gxui.Gray90, gxui.Gray15, gxui.Gray50, 1.0),
		ButtonPressedStyle:        CreateStyle(gxui.Gray20, gxui.Gray70, gxui.Gray30, 1.0),
		CodeSuggestionListStyle:   CreateStyle(gxui.Gray80, gxui.Gray20, gxui.Gray10, 1.0),
		DropDownListDefaultStyle:  CreateStyle(gxui.Gray80, gxui.Gray10, gxui.Gray20, 1.0),
		DropDownListOverStyle:     CreateStyle(gxui.Gray80, gxui.Gray15, gxui.Gray50, 1.0),
		FocusedStyle:              CreateStyle(gxui.Gray80, gxui.Transparent, focus, 1.0),
		HighlightStyle:            CreateStyle(gxui.Gray80, gxui.Transparent, neonBlue, 2.0),
		LabelStyle:                CreateStyle(gxui.Gray80, gxui.Transparent, gxui.Transparent, 0.0),
		PanelBackgroundStyle:      CreateStyle(gxui.Gray80, gxui.Gray10, gxui.Gray15, 1.0),
		ScrollBarBarDefaultStyle:  CreateStyle(gxui.Gray80, gxui.Gray30, gxui.Gray40, 1.0),
		ScrollBarBarOverStyle:     CreateStyle(gxui.Gray80, gxui.Gray50, gxui.Gray60, 1.0),
		ScrollBarRailDefaultStyle: CreateStyle(gxui.Gray80, scrollBarRailDefaultBg, gxui.Transparent, 1.0),
		ScrollBarRailOverStyle:    CreateStyle(gxui.Gray80, scrollBarRailOverBg, gxui.Gray20, 1.0),
		SplitterBarDefaultStyle:   CreateStyle(gxui.Gray80, gxui.Gray10, gxui.Gray10, 1.0),
		SplitterBarOverStyle:      CreateStyle(gxui.Gray80, gxui.Gray10, gxui.Gray50, 1.0),
		TabActiveHighlightStyle:   CreateStyle(gxui.Gray90, neonBlue, neonBlue, 0.0),
		TabDefaultStyle:           CreateStyle(gxui.Gray80, gxui.Gray30, gxui.Gray40, 1.0),
		TabOverStyle:              CreateStyle(gxui.Gray90, gxui.Gray30, gxui.Gray50, 1.0),
		TabPressedStyle:           CreateStyle(gxui.Gray20, gxui.Gray70, gxui.Gray30, 1.0),
		TextBoxDefaultStyle:       CreateStyle(gxui.Gray80, gxui.Gray10, gxui.Gray20, 1.0),
		TextBoxOverStyle:          CreateStyle(gxui.Gray80, gxui.Gray10, gxui.Gray50, 1.0),
	}
}

// gxui.Theme compliance
func (t *Theme) Driver() gxui.Driver {
	return t.driver
}

func (t *Theme) DefaultFont() gxui.Font {
	return t.defaultFont
}

func (t *Theme) SetDefaultFont(f gxui.Font) {
	t.defaultFont = f
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

func (t *Theme) CreateTextBox() gxui.TextBox {
	return CreateTextBox(t)
}

func (t *Theme) CreateTree() gxui.Tree {
	return CreateTree(t)
}

func (t *Theme) CreateWindow(width, height int, title string) gxui.Window {
	return CreateWindow(t, width, height, title)
}
