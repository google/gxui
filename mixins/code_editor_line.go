// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/google/gxui"
	"github.com/google/gxui/interval"
	"github.com/google/gxui/math"
)

type CodeEditorLinePaintInfo struct {
	LineSpan     interval.IntData
	Runes        []rune
	GlyphOffsets []math.Point
	GlyphWidth   int
	LineHeight   int
	Font         gxui.Font
}

type CodeEditorLineOuter interface {
	DefaultTextBoxLineOuter
	PaintBackgroundSpans(c gxui.Canvas, info CodeEditorLinePaintInfo)
	PaintGlyphs(c gxui.Canvas, info CodeEditorLinePaintInfo)
	PaintBorders(c gxui.Canvas, info CodeEditorLinePaintInfo)
}

// CodeEditorLine
type CodeEditorLine struct {
	DefaultTextBoxLine
	outer CodeEditorLineOuter
	ce    *CodeEditor
}

func (l *CodeEditorLine) Init(outer CodeEditorLineOuter, theme gxui.Theme, ce *CodeEditor, lineIndex int) {
	l.DefaultTextBoxLine.Init(outer, theme, &ce.TextBox, lineIndex)
	l.outer = outer
	l.ce = ce
	// Interface compliance test
	_ = TextBoxLine(l)
}

func (t *CodeEditorLine) PaintBackgroundSpans(c gxui.Canvas, info CodeEditorLinePaintInfo) {
	start, _ := info.LineSpan.Span()
	offsets := info.GlyphOffsets
	remaining := interval.IntDataList{info.LineSpan}
	for _, l := range t.ce.layers {
		if l != nil && l.BackgroundColor() != nil {
			color := *l.BackgroundColor()
			for _, span := range l.Spans().Overlaps(info.LineSpan) {
				interval.Visit(&remaining, span, func(vs, ve uint64, _ int) {
					s, e := vs-start, ve-start
					r := math.CreateRect(offsets[s].X, 0, offsets[e-1].X+info.GlyphWidth, info.LineHeight)
					c.DrawRoundedRect(r, 3, 3, 3, 3, gxui.TransparentPen, gxui.Brush{Color: color})
				})
				interval.Remove(&remaining, span)
			}
		}
	}
}

func (t *CodeEditorLine) PaintGlyphs(c gxui.Canvas, info CodeEditorLinePaintInfo) {
	start, _ := info.LineSpan.Span()
	runes, offsets, font := info.Runes, info.GlyphOffsets, info.Font
	remaining := interval.IntDataList{info.LineSpan}
	for _, l := range t.ce.layers {
		if l != nil && l.Color() != nil {
			color := *l.Color()
			for _, span := range l.Spans().Overlaps(info.LineSpan) {
				interval.Visit(&remaining, span, func(vs, ve uint64, _ int) {
					s, e := vs-start, ve-start
					c.DrawRunes(font, runes[s:e], offsets[s:e], color)
				})
				interval.Remove(&remaining, span)
			}
		}
	}
	for _, span := range remaining {
		s, e := span.Span()
		s, e = s-start, e-start
		c.DrawRunes(font, runes[s:e], offsets[s:e], t.ce.textColor)
	}
}

func (t *CodeEditorLine) PaintBorders(c gxui.Canvas, info CodeEditorLinePaintInfo) {
	start, _ := info.LineSpan.Span()
	offsets := info.GlyphOffsets
	for _, l := range t.ce.layers {
		if l != nil && l.BorderColor() != nil {
			color := *l.BorderColor()
			interval.Visit(l.Spans(), info.LineSpan, func(vs, ve uint64, _ int) {
				s, e := vs-start, ve-start
				r := math.CreateRect(offsets[s].X, 0, offsets[e-1].X+info.GlyphWidth, info.LineHeight)
				c.DrawRoundedRect(r, 3, 3, 3, 3, gxui.CreatePen(0.5, color), gxui.TransparentBrush)
			})
		}
	}
}

// DefaultTextBoxLine overrides
func (t *CodeEditorLine) Paint(c gxui.Canvas) {
	font := t.ce.font
	rect := t.Size().Rect().OffsetX(t.caretWidth)
	controller := t.ce.controller
	runes := controller.LineRunes(t.lineIndex)
	start := controller.LineStart(t.lineIndex)
	end := controller.LineEnd(t.lineIndex)

	if start != end {
		lineSpan := interval.CreateIntData(start, end, nil)

		lineHeight := t.Size().H
		glyphWidth := font.GlyphMaxSize().W
		offsets := font.Layout(&gxui.TextBlock{
			Runes:     runes,
			AlignRect: rect,
			H:         gxui.AlignLeft,
			V:         gxui.AlignMiddle,
		})

		info := CodeEditorLinePaintInfo{
			LineSpan:     lineSpan,
			Runes:        runes, // TODO gxui.TextBlock?
			GlyphOffsets: offsets,
			GlyphWidth:   glyphWidth,
			LineHeight:   lineHeight,
			Font:         font,
		}

		// Background
		t.outer.PaintBackgroundSpans(c, info)

		// Selections
		if t.textbox.HasFocus() {
			t.outer.PaintSelections(c)
		}

		// Glyphs
		t.outer.PaintGlyphs(c, info)

		// Borders
		t.outer.PaintBorders(c, info)
	}

	// Carets
	if t.textbox.HasFocus() {
		t.outer.PaintCarets(c)
	}
}
