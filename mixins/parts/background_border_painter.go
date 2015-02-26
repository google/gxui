// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parts

import (
	"gaze/gxui"
	"gaze/gxui/math"
	"gaze/gxui/mixins/outer"
)

type BackgroundBorderPainterOuter interface {
	outer.Redrawer
}

type BackgroundBorderPainter struct {
	outer BackgroundBorderPainterOuter
	brush gxui.Brush
	pen   gxui.Pen
}

func (b *BackgroundBorderPainter) Init(outer BackgroundBorderPainterOuter) {
	b.outer = outer
	b.brush = gxui.DefaultBrush
	b.pen = gxui.DefaultPen
}

func (b *BackgroundBorderPainter) PaintBackground(c gxui.Canvas, r math.Rect) {
	if b.brush.Color.A != 0 {
		borderF32 := b.pen.Width
		c.DrawRoundedRect(r, borderF32, borderF32, borderF32, borderF32, gxui.TransparentPen, b.brush)
	}
}

func (b *BackgroundBorderPainter) PaintBorder(c gxui.Canvas, r math.Rect) {
	if b.pen.Color.A != 0 {
		borderF32 := b.pen.Width
		c.DrawRoundedRect(r, borderF32, borderF32, borderF32, borderF32, b.pen, gxui.TransparentBrush)
	}
}

func (b *BackgroundBorderPainter) BackgroundBrush() gxui.Brush {
	return b.brush
}

func (b *BackgroundBorderPainter) SetBackgroundBrush(brush gxui.Brush) {
	if b.brush != brush {
		b.brush = brush
		b.outer.Redraw()
	}
}

func (b *BackgroundBorderPainter) BorderPen() gxui.Pen {
	return b.pen
}

func (b *BackgroundBorderPainter) SetBorderPen(pen gxui.Pen) {
	if b.pen != pen {
		b.pen = pen
		b.outer.Redraw()
	}
}
