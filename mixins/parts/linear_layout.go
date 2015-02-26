// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parts

import (
	"gxui"
	"gxui/math"
	"gxui/mixins/outer"
)

type LinearLayoutOuter interface {
	gxui.Container
	outer.Bounds
}

type LinearLayout struct {
	outer               LinearLayoutOuter
	orientation         gxui.Orientation
	horizontalAlignment gxui.HorizontalAlignment
	verticalAlignment   gxui.VerticalAlignment
}

func (l *LinearLayout) Init(outer LinearLayoutOuter) {
	l.outer = outer
}

func (l *LinearLayout) LayoutChildren() {
	s := l.outer.Bounds().Size().Contract(l.outer.Padding())
	o := l.outer.Padding().LT()

	d := 0
	children := l.outer.Children()
	for _, c := range children {
		cm := c.Margin()
		cs := c.DesiredSize(math.ZeroSize, s.Contract(cm).Max(math.ZeroSize))
		if l.orientation.Horizontal() {
			var y int
			switch l.verticalAlignment {
			case gxui.AlignTop:
				y = cm.T
			case gxui.AlignMiddle:
				y = (s.H - cs.H) / 2
			case gxui.AlignBottom:
				y = cm.B - cs.H
			}
			d += cm.L
			c.Layout(math.CreateRect(d, y, d+cs.W, y+cs.H).Offset(o))
			d += cs.W
			d += cm.R
			s.W -= cs.W + cm.W()
		} else {
			var x int
			switch l.horizontalAlignment {
			case gxui.AlignLeft:
				x = cm.L
			case gxui.AlignCenter:
				x = (s.W - cs.W) / 2
			case gxui.AlignRight:
				x = cm.R - cs.W
			}
			d += cm.T
			c.Layout(math.CreateRect(x, d, x+cs.W, d+cs.H).Offset(o))
			d += cs.H
			d += cm.B
			s.H -= cs.H + cm.H()
		}
	}
}

func (l *LinearLayout) DesiredSize(min, max math.Size) math.Size {
	bounds := min.Rect()
	children := l.outer.Children()

	offset := math.Point{X: 0, Y: 0}
	for _, c := range children {
		cs := c.DesiredSize(math.ZeroSize, max)
		cm := c.Margin()
		cb := cs.Expand(cm).Rect().Offset(offset)
		if l.orientation.Horizontal() {
			offset.X += cb.W()
		} else {
			offset.Y += cb.H()
		}
		bounds = bounds.Union(cb)
	}

	return bounds.Size().Expand(l.outer.Padding()).Clamp(min, max)
}

func (l *LinearLayout) Orientation() gxui.Orientation {
	return l.orientation
}

func (l *LinearLayout) SetOrientation(o gxui.Orientation) {
	if l.orientation != o {
		l.orientation = o
		l.LayoutChildren()
	}
}

func (l *LinearLayout) HorizontalAlignment() gxui.HorizontalAlignment {
	return l.horizontalAlignment
}

func (l *LinearLayout) SetHorizontalAlignment(alignment gxui.HorizontalAlignment) {
	l.horizontalAlignment = alignment
}

func (l *LinearLayout) VerticalAlignment() gxui.VerticalAlignment {
	return l.verticalAlignment
}

func (l *LinearLayout) SetVerticalAlignment(alignment gxui.VerticalAlignment) {
	l.verticalAlignment = alignment
}
