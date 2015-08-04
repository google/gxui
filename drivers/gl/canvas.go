// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"

	"github.com/google/gxui"
	"github.com/google/gxui/math"
	"github.com/goxjs/gl"
)

type drawStateStack []drawState

func (s *drawStateStack) head() *drawState {
	return &(*s)[len(*s)-1]
}
func (s *drawStateStack) push(ds drawState) {
	*s = append(*s, ds)
}
func (s *drawStateStack) pop() {
	*s = (*s)[:len(*s)-1]
}

type canvasOp func(ctx *context, dss *drawStateStack)

type drawState struct {
	// The below are all in window coordinates
	ClipPixels   math.Rect
	OriginPixels math.Point
}

type canvas struct {
	sizeDips          math.Size
	ops               []canvasOp
	built             bool
	buildingPushCount int
}

func newCanvas(sizeDips math.Size) *canvas {
	if sizeDips.W <= 0 || sizeDips.H < 0 {
		panic(fmt.Errorf("Canvas width and height must be positive. Size: %d", sizeDips))
	}
	c := &canvas{
		sizeDips: sizeDips,
	}
	return c
}

func (c *canvas) draw(ctx *context, dss *drawStateStack) {
	ds := dss.head()
	ctx.apply(ds)

	for _, op := range c.ops {
		op(ctx, dss)
	}
}

func (c *canvas) appendOp(name string, op canvasOp) {
	if c.built {
		panic(fmt.Errorf("%s() called after Complete()", name))
	}
	c.ops = append(c.ops, op)
}

// gxui.Canvas compliance
func (c *canvas) Size() math.Size {
	return c.sizeDips
}

func (c *canvas) IsComplete() bool {
	return c.built
}

func (c *canvas) Complete() {
	if c.built {
		panic("Complete() called twice")
	}
	if c.buildingPushCount != 0 {
		panic(fmt.Errorf("Push() count was %d when calling Complete", c.buildingPushCount))
	}
	c.built = true
}

func (c *canvas) Push() {
	c.buildingPushCount++
	c.appendOp("Push", func(ctx *context, dss *drawStateStack) {
		dss.push(*dss.head())
	})
}

func (c *canvas) Pop() {
	c.buildingPushCount--
	c.appendOp("Pop", func(ctx *context, dss *drawStateStack) {
		dss.pop()
		ctx.apply(dss.head())
	})
}

func (c *canvas) AddClip(r math.Rect) {
	c.appendOp("AddClip", func(ctx *context, dss *drawStateStack) {
		ds := dss.head()
		rectLocalPixels := ctx.resolution.rectDipsToPixels(r)
		rectWindowPixels := rectLocalPixels.Offset(ds.OriginPixels)
		ds.ClipPixels = ds.ClipPixels.Intersect(rectWindowPixels)
		ctx.apply(ds)
	})
}

func (c *canvas) Clear(color gxui.Color) {
	c.appendOp("Clear", func(ctx *context, dss *drawStateStack) {
		gl.ClearColor(
			color.R,
			color.G,
			color.B,
			color.A,
		)
		gl.Clear(gl.COLOR_BUFFER_BIT)
	})
}

func (c *canvas) DrawCanvas(cc gxui.Canvas, offsetDips math.Point) {
	if cc == nil {
		panic("Canvas cannot be nil")
	}
	childCanvas := cc.(*canvas)
	c.appendOp("DrawCanvas", func(ctx *context, dss *drawStateStack) {
		offsetPixels := ctx.resolution.pointDipsToPixels(offsetDips)
		dss.push(*dss.head())
		ds := dss.head()
		ds.OriginPixels = ds.OriginPixels.Add(offsetPixels)
		childCanvas.draw(ctx, dss)
		dss.pop()
		ctx.apply(dss.head())
	})
}

func (c *canvas) DrawRunes(f gxui.Font, r []rune, p []math.Point, col gxui.Color) {
	if f == nil {
		panic("Font cannot be nil")
	}
	runes := append([]rune{}, r...)
	points := append([]math.Point{}, p...)
	c.appendOp("DrawRunes", func(ctx *context, dss *drawStateStack) {
		f.(*font).DrawRunes(ctx, runes, points, col, dss.head())
	})
}

func (c *canvas) DrawLines(lines gxui.Polygon, pen gxui.Pen) {
	edge := openPolyToShape(lines, pen.Width)
	c.appendOp("DrawLines", func(ctx *context, dss *drawStateStack) {
		ds := dss.head()
		if edge != nil && pen.Color.A > 0 {
			ctx.blitter.blitShape(ctx, *edge, pen.Color, ds)
		}
	})
}

func (c *canvas) DrawPolygon(poly gxui.Polygon, pen gxui.Pen, brush gxui.Brush) {
	fill, edge := closedPolyToShape(poly, pen.Width)
	c.appendOp("DrawPolygon", func(ctx *context, dss *drawStateStack) {
		ds := dss.head()
		if fill != nil && brush.Color.A > 0 {
			ctx.blitter.blitShape(ctx, *fill, brush.Color, ds)
		}
		if edge != nil && pen.Color.A > 0 {
			ctx.blitter.blitShape(ctx, *edge, pen.Color, ds)
		}
	})
}

func (c *canvas) DrawRect(r math.Rect, brush gxui.Brush) {
	c.appendOp("DrawRect", func(ctx *context, dss *drawStateStack) {
		ctx.blitter.blitRect(ctx, ctx.resolution.rectDipsToPixels(r), brush.Color, dss.head())
	})
}

func (c *canvas) DrawRoundedRect(r math.Rect, tl, tr, bl, br float32, pen gxui.Pen, brush gxui.Brush) {
	if tl == 0 && tr == 0 && bl == 0 && br == 0 && pen.Color.A == 0 {
		c.DrawRect(r, brush)
		return
	}
	p := gxui.Polygon{
		gxui.PolygonVertex{Position: r.TL(), RoundedRadius: tl},
		gxui.PolygonVertex{Position: r.TR(), RoundedRadius: tr},
		gxui.PolygonVertex{Position: r.BR(), RoundedRadius: br},
		gxui.PolygonVertex{Position: r.BL(), RoundedRadius: bl},
	}
	c.DrawPolygon(p, pen, brush)
}

func (c *canvas) DrawTexture(t gxui.Texture, r math.Rect) {
	if t == nil {
		panic("Texture cannot be nil")
	}

	c.appendOp("DrawTexture", func(ctx *context, dss *drawStateStack) {
		tc := ctx.getOrCreateTextureContext(t.(*texture))
		ctx.blitter.blit(ctx, tc, tc.sizePixels.Rect(), ctx.resolution.rectDipsToPixels(r), dss.head())
	})
}
