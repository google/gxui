// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"gxui"
	"gxui/assert"
	"gxui/math"

	"github.com/go-gl/gl"
)

type DrawStateStack []DrawState

func (s *DrawStateStack) Head() *DrawState {
	return &(*s)[len(*s)-1]
}
func (s *DrawStateStack) Push(ds DrawState) {
	*s = append(*s, ds)
}
func (s *DrawStateStack) Pop() {
	*s = (*s)[:len(*s)-1]
}

type CanvasOp func(ctx *Context, dss *DrawStateStack)

type DrawState struct {
	// The below are all in window coordinates
	ClipPixels   math.Rect
	OriginPixels math.Point
}

type Canvas struct {
	refCounted
	sizeDips          math.Size
	resources         []RefCounted
	ops               []CanvasOp
	built             bool
	buildingPushCount int
}

func CreateCanvas(sizeDips math.Size) *Canvas {
	assert.True(sizeDips.W > 0, "Canvas width must be positive. Got: %d", sizeDips.W)
	assert.True(sizeDips.H > 0, "Canvas height must be positive. Got: %d", sizeDips.H)
	c := &Canvas{
		sizeDips: sizeDips,
	}
	c.init()
	globalStats.CanvasCount++
	return c
}

func (c *Canvas) draw(ctx *Context, dss *DrawStateStack) {
	c.AssertAlive("draw")
	ds := dss.Head()
	ctx.Apply(ds)

	for _, op := range c.ops {
		op(ctx, dss)
	}
}

func (c *Canvas) appendOp(name string, op CanvasOp) {
	c.AssertAlive(name)
	assert.False(c.built, "%s() called after Complete()", name)
	c.ops = append(c.ops, op)
}

func (c *Canvas) appendResource(r RefCounted) {
	r.AddRef()
	c.resources = append(c.resources, r)
}

// gxui.Canvas compliance
func (c *Canvas) Size() math.Size {
	return c.sizeDips
}

func (c *Canvas) Complete() {
	assert.False(c.built, "Complete() called twice")
	assert.Equals(0, c.buildingPushCount, "Push count")
	c.built = true
}

func (c *Canvas) Push() {
	c.buildingPushCount++
	c.appendOp("Push", func(ctx *Context, dss *DrawStateStack) {
		dss.Push(*dss.Head())
	})
}

func (c *Canvas) Pop() {
	c.buildingPushCount--
	c.appendOp("Pop", func(ctx *Context, dss *DrawStateStack) {
		dss.Pop()
		ctx.Apply(dss.Head())
	})
}

func (c *Canvas) AddClip(r math.Rect) {
	c.appendOp("AddClip", func(ctx *Context, dss *DrawStateStack) {
		ds := dss.Head()
		rectLocalPixels := ctx.RectDipsToPixels(r)
		rectWindowPixels := rectLocalPixels.Offset(ds.OriginPixels)
		ds.ClipPixels = ds.ClipPixels.Intersect(rectWindowPixels)
		ctx.Apply(ds)
	})
}

func (c *Canvas) Clear(color gxui.Color) {
	c.appendOp("Clear", func(ctx *Context, dss *DrawStateStack) {
		gl.ClearColor(
			gl.GLclampf(color.R),
			gl.GLclampf(color.G),
			gl.GLclampf(color.B),
			gl.GLclampf(color.A),
		)
		gl.Clear(gl.COLOR_BUFFER_BIT)
	})
}

func (c *Canvas) DrawCanvas(canvas gxui.Canvas, offsetDips math.Point) {
	assert.NotNil(canvas, "canvas")
	childCanvas := canvas.(*Canvas)
	c.appendOp("DrawCanvas", func(ctx *Context, dss *DrawStateStack) {
		offsetPixels := ctx.PointDipsToPixels(offsetDips)
		dss.Push(*dss.Head())
		ds := dss.Head()
		ds.OriginPixels = ds.OriginPixels.Add(offsetPixels)
		childCanvas.draw(ctx, dss)
		dss.Pop()
		ctx.Apply(dss.Head())
	})
	c.appendResource(childCanvas)
}

func (c *Canvas) DrawText(f gxui.Font, s string, col gxui.Color, r math.Rect, h gxui.HorizontalAlignment, v gxui.VerticalAlignment) {
	assert.NotNil(f, "font")
	c.appendOp("DrawText", func(ctx *Context, dss *DrawStateStack) {
		f.(*Font).Draw(ctx, s, col, r, h, v, dss.Head())
	})
}

func (c *Canvas) DrawRunes(f gxui.Font, r []rune, col gxui.Color, p []math.Point, o math.Point) {
	assert.NotNil(f, "font")
	runes := append([]rune{}, r...)
	points := append([]math.Point{}, p...)
	c.appendOp("DrawRunes", func(ctx *Context, dss *DrawStateStack) {
		f.(*Font).DrawRunes(ctx, runes, col, points, o, dss.Head())
	})
}

func (c *Canvas) DrawLines(lines gxui.Polygon, pen gxui.Pen) {
	edge := openPolyToShape(lines, pen.Width)
	c.appendOp("DrawLines", func(ctx *Context, dss *DrawStateStack) {
		ds := dss.Head()
		if edge != nil && pen.Color.A > 0 {
			ctx.Blitter.BlitShape(ctx, *edge, pen.Color, ds)
		}
	})
	if edge != nil {
		c.appendResource(edge)
		edge.Release()
	}
}

func (c *Canvas) DrawPolygon(poly gxui.Polygon, pen gxui.Pen, brush gxui.Brush) {
	fill, edge := closedPolyToShape(poly, pen.Width)
	c.appendOp("DrawPolygon", func(ctx *Context, dss *DrawStateStack) {
		ds := dss.Head()
		if fill != nil && brush.Color.A > 0 {
			ctx.Blitter.BlitShape(ctx, *fill, brush.Color, ds)
		}
		if edge != nil && pen.Color.A > 0 {
			ctx.Blitter.BlitShape(ctx, *edge, pen.Color, ds)
		}
	})
	if fill != nil {
		c.appendResource(fill)
		fill.Release()
	}
	if edge != nil {
		c.appendResource(edge)
		edge.Release()
	}
}

func (c *Canvas) DrawRect(r math.Rect, brush gxui.Brush) {
	c.appendOp("DrawRect", func(ctx *Context, dss *DrawStateStack) {
		ctx.Blitter.BlitRect(ctx, ctx.RectDipsToPixels(r), brush.Color, dss.Head())
	})
}

func (c *Canvas) DrawRoundedRect(r math.Rect, tl, tr, bl, br float32, pen gxui.Pen, brush gxui.Brush) {
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

func (c *Canvas) DrawTexture(t gxui.Texture, r math.Rect) {
	assert.NotNil(t, "texture")
	c.appendOp("DrawTexture", func(ctx *Context, dss *DrawStateStack) {
		tc := ctx.GetOrCreateTextureContext(t.(*Texture))
		ctx.Blitter.Blit(ctx, tc, tc.SizePixels().Rect(), ctx.RectDipsToPixels(r), dss.Head())
	})
	c.appendResource(t)
}

func (c *Canvas) Release() {
	if c.release() {
		for _, r := range c.resources {
			r.Release()
		}
		c.ops = nil
		c.resources = nil
		globalStats.CanvasCount--
	}
}
