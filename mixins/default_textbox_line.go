// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/google/gxui"
	"github.com/google/gxui/interval"
	"github.com/google/gxui/math"
	"github.com/google/gxui/mixins/base"
)

type DefaultTextBoxLineOuter interface {
	base.ControlOuter
	MeasureRunes(s, e int) math.Size
	PaintText(c gxui.Canvas)
	PaintCarets(c gxui.Canvas)
	PaintCaret(c gxui.Canvas, top, bottom math.Point)
	PaintSelections(c gxui.Canvas)
	PaintSelection(c gxui.Canvas, top, bottom math.Point)
}

// DefaultTextBoxLine
type DefaultTextBoxLine struct {
	base.Control
	outer      DefaultTextBoxLineOuter
	textbox    *TextBox
	lineIndex  int
	caretWidth int
}

func (t *DefaultTextBoxLine) Init(outer DefaultTextBoxLineOuter, theme gxui.Theme, textbox *TextBox, lineIndex int) {
	t.Control.Init(outer, theme)
	t.outer = outer
	t.textbox = textbox
	t.lineIndex = lineIndex
	t.SetCaretWidth(2)
	t.OnAttach(func() {
		ev := t.textbox.OnRedrawLines(t.Redraw)
		t.OnDetach(ev.Unlisten)
	})

	// Interface compliance test
	_ = TextBoxLine(t)
}

func (t *DefaultTextBoxLine) SetCaretWidth(width int) {
	if t.caretWidth != width {
		t.caretWidth = width
	}
}

func (t *DefaultTextBoxLine) DesiredSize(min, max math.Size) math.Size {
	return max
}

func (t *DefaultTextBoxLine) Paint(c gxui.Canvas) {
	if t.textbox.HasFocus() {
		t.outer.PaintSelections(c)
	}

	t.outer.PaintText(c)

	if t.textbox.HasFocus() {
		t.outer.PaintCarets(c)
	}
}

func (t *DefaultTextBoxLine) MeasureRunes(s, e int) math.Size {
	controller := t.textbox.controller
	return t.textbox.font.Measure(&gxui.TextBlock{
		Runes: controller.TextRunes()[s:e],
	})
}

func (t *DefaultTextBoxLine) PaintText(c gxui.Canvas) {
	runes := []rune(t.textbox.controller.Line(t.lineIndex))
	f := t.textbox.font
	offsets := f.Layout(&gxui.TextBlock{
		Runes:     runes,
		AlignRect: t.Size().Rect().OffsetX(t.caretWidth),
		H:         gxui.AlignLeft,
		V:         gxui.AlignBottom,
	})
	c.DrawRunes(f, runes, offsets, t.textbox.textColor)
}

func (t *DefaultTextBoxLine) PaintCarets(c gxui.Canvas) {
	controller := t.textbox.controller
	for i, cnt := 0, controller.SelectionCount(); i < cnt; i++ {
		e := controller.Caret(i)
		l := controller.LineIndex(e)
		if l == t.lineIndex {
			s := controller.LineStart(l)
			m := t.outer.MeasureRunes(s, e)
			top := math.Point{X: t.caretWidth + m.W, Y: 0}
			bottom := top.Add(math.Point{X: 0, Y: t.Size().H})
			t.outer.PaintCaret(c, top, bottom)
		}
	}
}

func (t *DefaultTextBoxLine) PaintSelections(c gxui.Canvas) {
	controller := t.textbox.controller

	ls, le := controller.LineStart(t.lineIndex), controller.LineEnd(t.lineIndex)

	selections := controller.Selections()
	if t.textbox.selectionDragging {
		interval.Replace(&selections, t.textbox.selectionDrag)
	}
	interval.Visit(&selections, gxui.CreateTextSelection(ls, le, false), func(s, e uint64, _ int) {
		if s < e {
			x := t.outer.MeasureRunes(ls, int(s)).W
			m := t.outer.MeasureRunes(int(s), int(e))
			top := math.Point{X: t.caretWidth + x, Y: 0}
			bottom := top.Add(m.Point())
			t.outer.PaintSelection(c, top, bottom)
		}
	})
}

func (t *DefaultTextBoxLine) PaintCaret(c gxui.Canvas, top, bottom math.Point) {
	r := math.Rect{Min: top, Max: bottom}.ExpandI(t.caretWidth / 2)
	c.DrawRoundedRect(r, 1, 1, 1, 1, gxui.CreatePen(0.5, gxui.Gray70), gxui.WhiteBrush)
}

func (t *DefaultTextBoxLine) PaintSelection(c gxui.Canvas, top, bottom math.Point) {
	r := math.Rect{Min: top, Max: bottom}.ExpandI(t.caretWidth / 2)
	c.DrawRoundedRect(r, 1, 1, 1, 1, gxui.TransparentPen, gxui.Brush{Color: gxui.Gray40})
}

// TextBoxLine compliance
func (t *DefaultTextBoxLine) RuneIndexAt(p math.Point) int {
	font := t.textbox.font
	controller := t.textbox.controller

	x := p.X
	line := controller.Line(t.lineIndex)
	i := 0
	for ; i < len(line) && x > font.Measure(&gxui.TextBlock{Runes: []rune(line[:i+1])}).W; i++ {
	}

	return controller.LineStart(t.lineIndex) + i
}

func (t *DefaultTextBoxLine) PositionAt(runeIndex int) math.Point {
	font := t.textbox.font
	controller := t.textbox.controller

	x := runeIndex - controller.LineStart(t.lineIndex)
	line := controller.Line(t.lineIndex)
	return font.Measure(&gxui.TextBlock{Runes: []rune(line[:x])}).Point()
}
