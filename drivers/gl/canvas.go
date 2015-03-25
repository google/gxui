// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/google/gxui"
	"github.com/google/gxui/math"
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
	if sizeDips.W <= 0 || sizeDips.H < 0 {
		panic(fmt.Errorf("Canvas width and height must be positive. Size: %d", sizeDips))
	}
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
	if c.built {
		panic(fmt.Errorf("%s() called after Complete()", name))
	}
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
	if c.built {
		panic("Complete() called twice")
	}
	if c.buildingPushCount != 0 {
		panic(fmt.Errorf("Push() count was %d when calling Complete", c.buildingPushCount))
	}
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
			color.R,
			color.G,
			color.B,
			color.A,
		)
		gl.Clear(gl.COLOR_BUFFER_BIT)
	})
}

func (c *Canvas) DrawCanvas(canvas gxui.Canvas, offsetDips math.Point) {
	if canvas == nil {
		panic("Canvas cannot be nil")
	}
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

func (c *Canvas) DrawRunes(f gxui.Font, r []rune, p []math.Point, col gxui.Color) {
	if f == nil {
		panic("Font cannot be nil")
	}
	runes := append([]rune{}, r...)
	points := append([]math.Point{}, p...)
	c.appendOp("DrawRunes", func(ctx *Context, dss *DrawStateStack) {
		f.(*Font).DrawRunes(ctx, runes, points, col, dss.Head())
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
	if t == nil {
		panic("Texture cannot be nil")
	}

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
