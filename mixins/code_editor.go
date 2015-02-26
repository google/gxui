// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"fmt"
	"gxui"
	"gxui/interval"
	"gxui/math"
	"strings"
)

type CodeEditorOuter interface {
	TextBoxOuter
	CreateSuggestionList() gxui.List
}

type CodeEditor struct {
	TextBox
	outer              CodeEditorOuter
	layers             []gxui.CodeSyntaxLayer
	suggestionAdapter  *SuggestionAdapter
	suggestionList     gxui.List
	suggestionProvider gxui.CodeSuggestionProvider
	tabWidth           int
	theme              gxui.Theme
}

func (t *CodeEditor) updateSpans(edits []gxui.TextBoxEdit) {
	runeCount := len(t.controller.TextRunes())
	for _, l := range t.layers {
		l.UpdateSpans(runeCount, edits)
	}
}

func (t *CodeEditor) Init(outer CodeEditorOuter, driver gxui.Driver, theme gxui.Theme, font gxui.Font) {
	t.outer = outer
	t.tabWidth = 2
	t.theme = theme

	t.suggestionAdapter = &SuggestionAdapter{}
	t.suggestionList = t.outer.CreateSuggestionList()
	t.suggestionList.SetAdapter(t.suggestionAdapter)
	t.suggestionList.SetVisible(false)

	t.TextBox.Init(outer, driver, theme, font)
	t.controller.OnTextChanged(t.updateSpans)

	// Interface compliance test
	_ = gxui.CodeEditor(t)
}

func (t *CodeEditor) ItemSize(theme gxui.Theme) math.Size {
	return math.Size{W: math.MaxSize.W, H: t.font.GlyphMaxSize().H}
}

func (t *CodeEditor) CreateSuggestionList() gxui.List {
	l := t.theme.CreateList()
	l.SetBackgroundBrush(gxui.DefaultBrush)
	l.SetBorderPen(gxui.DefaultPen)
	return l
}

func (t *CodeEditor) SetSyntaxLayer(idx int, layer gxui.CodeSyntaxLayer) {
	if len(t.layers) <= idx {
		layers := make([]gxui.CodeSyntaxLayer, idx+1)
		copy(layers[:], t.layers[:])
		t.layers = layers
	}
	t.layers[idx] = layer
	t.onRedrawLines.Fire()
}

func (t *CodeEditor) ClearSyntaxLayers() {
	t.layers = t.layers[:0]
}

func (t *CodeEditor) TabWidth() int {
	return t.tabWidth
}

func (t *CodeEditor) SetTabWidth(tabWidth int) {
	t.tabWidth = tabWidth
}

func (t *CodeEditor) SpanAt(layerIdx, runeIdx int) *interval.IntData {
	return t.layers[layerIdx].SpanAt(runeIdx)
}

func (t *CodeEditor) SpansAt(at int) []interval.IntData {
	spans := []interval.IntData{}
	for _, layer := range t.layers {
		if s := layer.SpanAt(at); s != nil {
			spans = append(spans, *s)
		}
	}
	return spans
}

func (t *CodeEditor) SuggestionProvider() gxui.CodeSuggestionProvider {
	return t.suggestionProvider
}

func (t *CodeEditor) SetSuggestionProvider(provider gxui.CodeSuggestionProvider) {
	if t.suggestionProvider != provider {
		t.suggestionProvider = provider
		if t.IsSuggestionListShowing() {
			t.ShowSuggestionList() // Update list
		}
	}
}

func (t *CodeEditor) IsSuggestionListShowing() bool {
	return t.suggestionList.IsVisible()
}

func (t *CodeEditor) SortSuggestionList() {
	caret := t.controller.LastCaret()
	partial := t.controller.TextRange(t.controller.WordAt(caret))
	t.suggestionAdapter.Sort(partial)
}

func (t *CodeEditor) ShowSuggestionList() {
	if t.suggestionProvider == nil {
		return
	}

	caret := t.controller.LastCaret()
	s, _ := t.controller.WordAt(caret)

	suggestions := t.suggestionProvider.SuggestionsAt(s)
	if len(suggestions) == 0 {
		t.HideSuggestionList()
		return
	}

	t.suggestionAdapter.SetSuggestions(suggestions)
	t.SortSuggestionList()

	// Position the suggestion list below the last caret
	lineIdx := t.controller.LineIndex(caret)
	// TODO: What if the last caret is not visible?
	bounds := t.Bounds().Size().Rect().Contract(t.Padding())
	line := t.Line(lineIdx)
	lineOffset := gxui.ChildToParent(math.ZeroPoint, line, t.outer)
	target := line.PositionAt(caret).Add(lineOffset)
	cs := t.suggestionList.DesiredSize(math.ZeroSize, bounds.Size())
	t.suggestionList.Layout(cs.Rect().Offset(target).Intersect(bounds))
	t.suggestionList.Select(t.suggestionList.Adapter().ItemId(0))
	t.suggestionList.SetVisible(true)
}

func (t *CodeEditor) HideSuggestionList() {
	t.suggestionList.SetVisible(false)
}

func (t *CodeEditor) Line(idx int) TextBoxLine {
	id := gxui.AdapterItemId(idx)
	return gxui.FindControl(t.Item(id), func(c gxui.Control) bool {
		_, b := c.(TextBoxLine)
		return b
	}).(TextBoxLine)
}

// mixins.List overrides
func (t *CodeEditor) LayoutChildren() {
	t.List.LayoutChildren()
	if t.suggestionList.Parent() != t.outer {
		t.AddChild(t.suggestionList)
	}
}

func (t *CodeEditor) Click(ev gxui.MouseEvent) (consume bool) {
	t.HideSuggestionList()
	return t.TextBox.Click(ev)
}

func (t *CodeEditor) KeyPress(ev gxui.KeyboardEvent) (consume bool) {
	switch ev.Key {
	case gxui.KeyTab:
		replace := true
		for _, sel := range t.controller.Selections() {
			s, e := sel.Range()
			if t.controller.LineIndex(s) != t.controller.LineIndex(e) {
				replace = false
				break
			}
		}
		switch {
		case replace:
			t.controller.ReplaceAll(strings.Repeat(" ", t.tabWidth))
			t.controller.Deselect(false)
		case ev.Modifier.Shift():
			t.controller.UnindentSelection(t.tabWidth)
		default:
			t.controller.IndentSelection(t.tabWidth)
		}
		return true
	case gxui.KeySpace:
		if ev.Modifier.Control() {
			t.ShowSuggestionList()
			return
		}
	case gxui.KeyUp:
		fallthrough
	case gxui.KeyDown:
		if t.IsSuggestionListShowing() {
			return t.suggestionList.KeyPress(ev)
		}
	case gxui.KeyLeft:
		t.HideSuggestionList()
	case gxui.KeyRight:
		t.HideSuggestionList()
	case gxui.KeyEnter:
		controller := t.controller
		if t.IsSuggestionListShowing() {
			text := t.suggestionAdapter.Suggestion(t.suggestionList.Selected()).Code()
			s, e := controller.WordAt(t.controller.LastCaret())
			controller.SetSelection(gxui.CreateTextSelection(s, e, false))
			controller.ReplaceAll(text)
			controller.Deselect(false)
			t.HideSuggestionList()
		} else {
			t.controller.ReplaceWithNewlineKeepIndent()
		}
		return true
	case gxui.KeyEscape:
		if t.IsSuggestionListShowing() {
			t.HideSuggestionList()
			return true
		}
	}
	return t.TextBox.KeyPress(ev)
}

func (t *CodeEditor) KeyStroke(ev gxui.KeyStrokeEvent) (consume bool) {
	consume = t.TextBox.KeyStroke(ev)
	if t.IsSuggestionListShowing() {
		t.SortSuggestionList()
	}
	return
}

// mixins.TextBox overrides
func (t *CodeEditor) CreateLine(theme gxui.Theme, index int) (TextBoxLine, gxui.Control) {
	lineNumber := theme.CreateLabel()
	lineNumber.SetText(fmt.Sprintf("%.4d", index+1)) // Displayed lines start at 1

	line := &CodeEditorLine{}
	line.Init(line, theme, t, index)

	layout := theme.CreateLinearLayout()
	layout.SetOrientation(gxui.Horizontal)
	layout.AddChild(lineNumber)
	layout.AddChild(line)

	return line, layout
}
