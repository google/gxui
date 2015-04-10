// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parts

import (
	"github.com/google/gxui"
	"github.com/google/gxui/math"
	"github.com/google/gxui/mixins/outer"
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
		w := b.pen.Width
		c.DrawRoundedRect(r, w, w, w, w, gxui.TransparentPen, b.brush)
	}
}

func (b *BackgroundBorderPainter) PaintBorder(c gxui.Canvas, r math.Rect) {
	if b.pen.Color.A != 0 && b.pen.Width != 0 {
		w := b.pen.Width
		c.DrawRoundedRect(r, w, w, w, w, b.pen, gxui.TransparentBrush)
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
