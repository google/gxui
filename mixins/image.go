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

type ImageOuter interface {
	base.ControlOuter
}

type Image struct {
	base.Control
	parts.BackgroundBorderPainter

	outer        ImageOuter
	texture      gxui.Texture
	canvas       gxui.Canvas
	scalingMode  gxui.ScalingMode
	aspectMode   gxui.AspectMode
	explicitSize math.Size
}

func (i *Image) calculateDrawRect() math.Rect {
	r := i.outer.Size().Rect()
	texW, texH := i.texture.Size().WH()
	aspectSrc := float32(texH) / float32(texW)
	aspectDst := float32(r.H()) / float32(r.W())
	switch i.aspectMode {
	case gxui.AspectCorrectLetterbox, gxui.AspectCorrectCrop:
		if (aspectDst < aspectSrc) != (i.aspectMode == gxui.AspectCorrectLetterbox) {
			contract := r.H() - int(float32(r.W())*aspectSrc)
			r = r.Contract(math.Spacing{T: contract / 2, B: contract / 2})
		} else {
			contract := r.W() - int(float32(r.H())/aspectSrc)
			r = r.Contract(math.Spacing{L: contract / 2, R: contract / 2})
		}
	}
	return r
}

func (i *Image) Init(outer ImageOuter, theme gxui.Theme) {
	i.outer = outer
	i.Control.Init(outer, theme)
	i.BackgroundBorderPainter.Init(outer)
	i.SetBorderPen(gxui.TransparentPen)
	i.SetBackgroundBrush(gxui.TransparentBrush)

	// Interface compliance test
	_ = gxui.Image(i)
}

func (i *Image) Texture() gxui.Texture {
	return i.texture
}

func (i *Image) SetTexture(tex gxui.Texture) {
	if i.texture != tex {
		i.texture = tex
		i.canvas = nil
		i.outer.Relayout()
	}
}

func (i *Image) Canvas() gxui.Canvas {
	return i.canvas
}

func (i *Image) SetCanvas(canvas gxui.Canvas) {
	if !canvas.IsComplete() {
		panic("SetCanvas() called with an incomplete canvas")
	}

	if i.canvas != canvas {
		i.canvas = canvas
		i.texture = nil
		i.outer.Relayout()
	}
}

func (i *Image) ScalingMode() gxui.ScalingMode {
	return i.scalingMode
}

func (i *Image) SetScalingMode(mode gxui.ScalingMode) {
	if i.scalingMode != mode {
		i.scalingMode = mode
		i.outer.Relayout()
	}
}

func (i *Image) AspectMode() gxui.AspectMode {
	return i.aspectMode
}

func (i *Image) SetAspectMode(mode gxui.AspectMode) {
	if i.aspectMode != mode {
		i.aspectMode = mode
		i.outer.Redraw()
	}
}

func (i *Image) SetExplicitSize(explicitSize math.Size) {
	if i.explicitSize != explicitSize {
		i.explicitSize = explicitSize
		i.outer.Relayout()
	}
	i.SetScalingMode(gxui.ScalingExplicitSize)
}

func (i *Image) PixelAt(p math.Point) (math.Point, bool) {
	ir := i.calculateDrawRect()
	if tex := i.Texture(); tex != nil {
		s := tex.SizePixels()
		p = p.Sub(ir.Min).
			ScaleX(float32(s.W) / float32(ir.W())).
			ScaleY(float32(s.H) / float32(ir.H()))
		if s.Rect().Contains(p) {
			return p, true
		}
	}
	return math.Point{X: -1, Y: -1}, false
}

func (i *Image) DesiredSize(min, max math.Size) math.Size {
	s := max
	switch i.scalingMode {
	case gxui.ScalingExplicitSize:
		s = i.explicitSize
	case gxui.Scaling1to1:
		switch {
		case i.texture != nil:
			s = i.texture.Size()
		case i.canvas != nil:
			s = i.canvas.Size()
		}
	}
	return s.Expand(math.CreateSpacing(int(i.BorderPen().Width))).Clamp(min, max)
}

func (i *Image) Paint(c gxui.Canvas) {
	r := i.outer.Size().Rect()
	i.PaintBackground(c, r)
	switch {
	case i.texture != nil:
		c.DrawTexture(i.texture, i.calculateDrawRect())
	case i.canvas != nil:
		c.DrawCanvas(i.canvas, math.ZeroPoint)
	}
	i.PaintBorder(c, r)
}
