// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/google/gxui"
	"github.com/google/gxui/math"
	"github.com/google/gxui/mixins/base"
	"github.com/google/gxui/mixins/parts"
)

type ProgressBarOuter interface {
	base.ControlOuter
	PaintProgress(gxui.Canvas, math.Rect, float32)
}

type ProgressBar struct {
	base.Control
	parts.BackgroundBorderPainter

	outer            ProgressBarOuter
	desiredSize      math.Size
	progress, target int
}

func (b *ProgressBar) Init(outer ProgressBarOuter, theme gxui.Theme) {
	b.outer = outer
	b.Control.Init(outer, theme)
	b.BackgroundBorderPainter.Init(outer)
	b.desiredSize = math.MaxSize
	b.target = 100

	// Interface compliance test
	_ = gxui.ProgressBar(b)
}

func (b *ProgressBar) Paint(c gxui.Canvas) {
	frac := math.Saturate(float32(b.progress) / float32(b.target))
	r := b.outer.Size().Rect()
	b.PaintBackground(c, r)
	b.outer.PaintProgress(c, r, frac)
	b.PaintBorder(c, r)
}

func (b *ProgressBar) PaintProgress(c gxui.Canvas, r math.Rect, frac float32) {
	r.Max.X = math.Lerp(r.Min.X, r.Max.X, frac)
	c.DrawRect(r, gxui.CreateBrush(gxui.Gray50))
}

func (b *ProgressBar) DesiredSize(min, max math.Size) math.Size {
	return b.desiredSize.Clamp(min, max)
}

// gxui.ProgressBar compliance
func (b *ProgressBar) SetDesiredSize(size math.Size) {
	b.desiredSize = size
	b.Relayout()
}

func (b *ProgressBar) SetProgress(progress int) {
	if b.progress != progress {
		b.progress = progress
		b.Redraw()
	}
}

func (b *ProgressBar) Progress() int {
	return b.progress
}

func (b *ProgressBar) SetTarget(target int) {
	if b.target != target {
		b.target = target
		b.Redraw()
	}
}

func (b *ProgressBar) Target() int {
	return b.target
}
