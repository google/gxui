// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/google/gxui"
	"github.com/google/gxui/math"
	"github.com/google/gxui/mixins/base"
)

type BubbleOverlayOuter interface {
	base.ContainerOuter
}

type BubbleOverlay struct {
	base.Container
	outer       BubbleOverlayOuter
	targetPoint math.Point
	arrowLength int
	arrowWidth  int
	brush       gxui.Brush
	pen         gxui.Pen
}

func (o *BubbleOverlay) Init(outer BubbleOverlayOuter, theme gxui.Theme) {
	o.Container.Init(outer, theme)
	o.outer = outer
	o.arrowLength = 20
	o.arrowWidth = 15

	// Interface compliance test
	_ = gxui.BubbleOverlay(o)
}

func (o *BubbleOverlay) LayoutChildren() {
	for _, child := range o.outer.Children() {
		bounds := o.outer.Size().Rect().Contract(o.outer.Padding())
		arrowPadding := math.CreateSpacing(o.arrowLength)
		cm := child.Control.Margin()
		cs := child.Control.DesiredSize(math.ZeroSize, bounds.Size().Contract(cm).Max(math.ZeroSize))
		cr := cs.Expand(arrowPadding).EdgeAlignedFit(bounds, o.targetPoint).Contract(arrowPadding)
		child.Layout(cr)
	}
}

func (o *BubbleOverlay) DesiredSize(min, max math.Size) math.Size {
	return max
}

func (o *BubbleOverlay) Show(control gxui.Control, target math.Point) {
	o.Hide()
	o.outer.AddChild(control)
	o.targetPoint = target
}

func (o *BubbleOverlay) Hide() {
	o.outer.RemoveAll()
}

func (o *BubbleOverlay) Brush() gxui.Brush {
	return o.brush
}

func (o *BubbleOverlay) SetBrush(brush gxui.Brush) {
	if o.brush != brush {
		o.brush = brush
		o.Redraw()
	}
}

func (o *BubbleOverlay) Pen() gxui.Pen {
	return o.pen
}

func (o *BubbleOverlay) SetPen(pen gxui.Pen) {
	if o.pen != pen {
		o.pen = pen
		o.Redraw()
	}
}

func (o *BubbleOverlay) Paint(c gxui.Canvas) {
	if !o.IsVisible() {
		return
	}
	for _, child := range o.outer.Children() {
		b := child.Bounds().Expand(o.outer.Padding())
		t := o.targetPoint
		a := o.arrowWidth / 2
		var p gxui.Polygon

		switch {
		case t.X < b.Min.X:
			/*
			    A-----------------B
			    G                 |
			 F                    |
			    E                 |
			    D-----------------C
			*/
			p = gxui.Polygon{
				/*A*/ {Position: b.TL(), RoundedRadius: 5},
				/*B*/ {Position: b.TR(), RoundedRadius: 5},
				/*C*/ {Position: b.BR(), RoundedRadius: 5},
				/*D*/ {Position: b.BL(), RoundedRadius: 5},
				/*E*/ {Position: math.Point{X: b.Min.X, Y: math.Clamp(t.Y+a, b.Min.Y+a, b.Max.Y)}, RoundedRadius: 0},
				/*F*/ {Position: t, RoundedRadius: 0},
				/*G*/ {Position: math.Point{X: b.Min.X, Y: math.Clamp(t.Y-a, b.Min.Y, b.Max.Y-a)}, RoundedRadius: 0},
			}
			// fmt.Printf("A: %+v\n", p)
		case t.X > b.Max.X:
			/*
			   A-----------------B
			   |                 C
			   |                    D
			   |                 E
			   G-----------------F
			*/
			p = gxui.Polygon{
				/*A*/ {Position: b.TL(), RoundedRadius: 5},
				/*B*/ {Position: b.TR(), RoundedRadius: 5},
				/*C*/ {Position: math.Point{X: b.Max.X, Y: math.Clamp(t.Y-a, b.Min.Y, b.Max.Y-a)}, RoundedRadius: 0},
				/*D*/ {Position: t, RoundedRadius: 0},
				/*E*/ {Position: math.Point{X: b.Max.X, Y: math.Clamp(t.Y+a, b.Min.Y+a, b.Max.Y)}, RoundedRadius: 0},
				/*F*/ {Position: b.BR(), RoundedRadius: 5},
				/*G*/ {Position: b.BL(), RoundedRadius: 5},
			}
			// fmt.Printf("B: %+v\n", p)
		case t.Y < b.Min.Y:
			/*
			                 C
			                / \
			   A-----------B   D-E
			   |                 |
			   |                 |
			   G-----------------F
			*/
			p = gxui.Polygon{
				/*A*/ {Position: b.TL(), RoundedRadius: 5},
				/*B*/ {Position: math.Point{X: math.Clamp(t.X-a, b.Min.X, b.Max.X-a), Y: b.Min.Y}, RoundedRadius: 0},
				/*C*/ {Position: t, RoundedRadius: 0},
				/*D*/ {Position: math.Point{X: math.Clamp(t.X+a, b.Min.X+a, b.Max.X), Y: b.Min.Y}, RoundedRadius: 0},
				/*E*/ {Position: b.TR(), RoundedRadius: 5},
				/*F*/ {Position: b.BR(), RoundedRadius: 5},
				/*G*/ {Position: b.BL(), RoundedRadius: 5},
			}
			// fmt.Printf("C: %+v\n", p)
		default:
			/*
			   A-----------------B
			   |                 |
			   |                 |
			   G-----------F   D-C
			                \ /
			                 E
			*/
			p = gxui.Polygon{
				/*A*/ {Position: b.TL(), RoundedRadius: 5},
				/*B*/ {Position: b.TR(), RoundedRadius: 5},
				/*C*/ {Position: b.BR(), RoundedRadius: 5},
				/*D*/ {Position: math.Point{X: math.Clamp(t.X+a, b.Min.X+a, b.Max.X), Y: b.Max.Y}, RoundedRadius: 0},
				/*E*/ {Position: t, RoundedRadius: 0},
				/*F*/ {Position: math.Point{X: math.Clamp(t.X-a, b.Min.X, b.Max.X-a), Y: b.Max.Y}, RoundedRadius: 0},
				/*G*/ {Position: b.BL(), RoundedRadius: 5},
			}
			// fmt.Printf("D: %+v\n", p)
		}
		c.DrawPolygon(p, o.pen, o.brush)
	}
	o.PaintChildren.Paint(c)
}
