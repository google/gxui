// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parts

import (
	"github.com/google/gxui"
	"github.com/google/gxui/math"
	"github.com/google/gxui/mixins/outer"
)

type LinearLayoutOuter interface {
	gxui.Container
	outer.Bounds
}

type LinearLayout struct {
	outer               LinearLayoutOuter
	direction           gxui.Direction
	sizeMode            gxui.SizeMode
	horizontalAlignment gxui.HorizontalAlignment
	verticalAlignment   gxui.VerticalAlignment
}

func (l *LinearLayout) Init(outer LinearLayoutOuter) {
	l.outer = outer
}

func (l *LinearLayout) LayoutChildren() {
	s := l.outer.Bounds().Size().Contract(l.outer.Padding())
	o := l.outer.Padding().LT()
	children := l.outer.Children()
	major := 0
	if l.direction.RightToLeft() || l.direction.BottomToTop() {
		if l.direction.RightToLeft() {
			major = s.W
		} else {
			major = s.H
		}
	}
	for _, c := range children {
		cm := c.Margin()
		cs := c.DesiredSize(math.ZeroSize, s.Contract(cm).Max(math.ZeroSize))

		// Calculate minor-axis alignment
		var minor int
		switch l.direction.Orientation() {
		case gxui.Horizontal:
			switch l.verticalAlignment {
			case gxui.AlignTop:
				minor = cm.T
			case gxui.AlignMiddle:
				minor = (s.H - cs.H) / 2
			case gxui.AlignBottom:
				minor = s.H - cs.H
			}
		case gxui.Vertical:
			switch l.horizontalAlignment {
			case gxui.AlignLeft:
				minor = cm.L
			case gxui.AlignCenter:
				minor = (s.W - cs.W) / 2
			case gxui.AlignRight:
				minor = s.W - cs.W
			}
		}

		// Peform layout
		switch l.direction {
		case gxui.LeftToRight:
			major += cm.L
			c.Layout(math.CreateRect(major, minor, major+cs.W, minor+cs.H).Offset(o))
			major += cs.W
			major += cm.R
			s.W -= cs.W + cm.W()
		case gxui.RightToLeft:
			major -= cm.R
			c.Layout(math.CreateRect(major-cs.W, minor, major, minor+cs.H).Offset(o))
			major -= cs.W
			major -= cm.L
			s.W -= cs.W + cm.W()
		case gxui.TopToBottom:
			major += cm.T
			c.Layout(math.CreateRect(minor, major, minor+cs.W, major+cs.H).Offset(o))
			major += cs.H
			major += cm.B
			s.H -= cs.H + cm.H()
		case gxui.BottomToTop:
			major -= cm.B
			c.Layout(math.CreateRect(minor, major-cs.H, minor+cs.W, major).Offset(o))
			major -= cs.H
			major -= cm.T
			s.H -= cs.H + cm.H()
		}
	}
}

func (l *LinearLayout) DesiredSize(min, max math.Size) math.Size {
	if l.sizeMode.Fill() {
		return max
	}

	bounds := min.Rect()
	children := l.outer.Children()

	horizontal := l.direction.Orientation().Horizontal()
	offset := math.Point{X: 0, Y: 0}
	for _, c := range children {
		cs := c.DesiredSize(math.ZeroSize, max)
		cm := c.Margin()
		cb := cs.Expand(cm).Rect().Offset(offset)
		if horizontal {
			offset.X += cb.W()
		} else {
			offset.Y += cb.H()
		}
		bounds = bounds.Union(cb)
	}

	return bounds.Size().Expand(l.outer.Padding()).Clamp(min, max)
}

func (l *LinearLayout) Direction() gxui.Direction {
	return l.direction
}

func (l *LinearLayout) SetDirection(d gxui.Direction) {
	if l.direction != d {
		l.direction = d
		l.outer.Relayout()
	}
}

func (l *LinearLayout) SizeMode() gxui.SizeMode {
	return l.sizeMode
}

func (l *LinearLayout) SetSizeMode(mode gxui.SizeMode) {
	if l.sizeMode != mode {
		l.sizeMode = mode
		l.outer.Relayout()
	}
}

func (l *LinearLayout) HorizontalAlignment() gxui.HorizontalAlignment {
	return l.horizontalAlignment
}

func (l *LinearLayout) SetHorizontalAlignment(alignment gxui.HorizontalAlignment) {
	if l.horizontalAlignment != alignment {
		l.horizontalAlignment = alignment
		l.outer.Relayout()
	}
}

func (l *LinearLayout) VerticalAlignment() gxui.VerticalAlignment {
	return l.verticalAlignment
}

func (l *LinearLayout) SetVerticalAlignment(alignment gxui.VerticalAlignment) {
	if l.verticalAlignment != alignment {
		l.verticalAlignment = alignment
		l.outer.Relayout()
	}
}
