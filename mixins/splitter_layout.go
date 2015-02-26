// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"gaze/gxui"
	"gaze/gxui/math"
	"gaze/gxui/mixins/base"
)

type SplitterLayoutOuter interface {
	base.ContainerOuter
	CreateSplitterBar() gxui.Control
}

type SplitterLayout struct {
	base.Container

	outer         SplitterLayoutOuter
	theme         gxui.Theme
	orientation   gxui.Orientation
	splitterWidth int
	weights       map[gxui.Control]float32
}

func (l *SplitterLayout) Init(outer SplitterLayoutOuter, theme gxui.Theme) {
	l.Container.Init(outer, theme)
	l.outer = outer
	l.theme = theme
	l.weights = make(map[gxui.Control]float32)
	l.splitterWidth = 4
	l.SetMouseEventTarget(true)

	// Interface compliance test
	_ = gxui.SplitterLayout(l)
}

func (l *SplitterLayout) LayoutChildren() {
	s := l.outer.Bounds().Size().Contract(l.Padding())
	o := l.Padding().LT()

	children := l.outer.Children()

	splitterCount := len(children) / 2

	splitterWidth := l.splitterWidth
	if l.orientation.Horizontal() {
		s.W -= splitterWidth * splitterCount
	} else {
		s.H -= splitterWidth * splitterCount
	}

	netWeight := float32(0.0)
	for i, c := range children {
		if isSplitter := (i & 1) == 1; !isSplitter {
			netWeight += l.weights[c]
		}
	}

	d := 0
	for i, c := range children {
		var cr math.Rect
		if isSplitter := (i & 1) == 1; !isSplitter {
			cm := c.Margin()
			frac := l.weights[c] / netWeight
			if l.orientation.Horizontal() {
				cw := int(float32(s.W) * frac)
				cr = math.CreateRect(d+cm.L, cm.T, d+cw-cm.R, s.H-cm.B)
				d += cw
			} else {
				ch := int(float32(s.H) * frac)
				cr = math.CreateRect(cm.L, d+cm.T, s.W-cm.R, d+ch-cm.B)
				d += ch
			}
		} else {
			if l.orientation.Horizontal() {
				cr = math.CreateRect(d, 0, d+splitterWidth, s.H)
			} else {
				cr = math.CreateRect(0, d, s.W, d+splitterWidth)
			}
			d += splitterWidth
		}
		c.Layout(cr.Offset(o).Canon())
	}
}

func (l *SplitterLayout) ChildWeight(child gxui.Control) float32 {
	return l.weights[child]
}

func (l *SplitterLayout) SetChildWeight(child gxui.Control, weight float32) {
	if l.weights[child] != weight {
		l.weights[child] = weight
		l.LayoutChildren()
	}
}

func (l *SplitterLayout) DesiredSize(min, max math.Size) math.Size {
	return max
}

func (l *SplitterLayout) Orientation() gxui.Orientation {
	return l.orientation
}

func (l *SplitterLayout) SetOrientation(o gxui.Orientation) {
	if l.orientation != o {
		l.orientation = o
		l.LayoutChildren()
	}
}

func (l *SplitterLayout) CreateSplitterBar() gxui.Control {
	b := &SplitterBar{}
	b.Init(b, l.theme)
	b.OnSplitterDragged(func(wndPnt math.Point) { l.SplitterDragged(b, wndPnt) })
	return b
}

func (l *SplitterLayout) SplitterDragged(splitter gxui.Control, wndPnt math.Point) {
	o := l.orientation
	p := gxui.WindowToChild(wndPnt, l)
	splitterIndex := l.ChildIndex(splitter)
	childA, childB := l.ChildAt(splitterIndex-1), l.ChildAt(splitterIndex+1)
	boundsA, boundsB := childA.Bounds(), childB.Bounds()

	min, max := o.Major(boundsA.Min.XY()), o.Major(boundsB.Max.XY())
	frac := math.RampSat(float32(o.Major(p.XY())), float32(min), float32(max))

	netWeight := l.weights[childA] + l.weights[childB]
	l.weights[childA] = netWeight * frac
	l.weights[childB] = netWeight * (1.0 - frac)
	l.LayoutChildren()
}

// parts.Container overrides
func (l *SplitterLayout) AddChildAt(index int, child gxui.Control) {
	l.weights[child] = 1.0
	if l.Container.ChildCount() > 0 {
		l.Container.AddChildAt(index, l.outer.CreateSplitterBar())
		index++
	}
	l.Container.AddChildAt(index, child)
}

func (l *SplitterLayout) RemoveChildAt(index int) {
	if l.Container.ChildCount() > 1 {
		l.Container.RemoveChildAt(index + 1)
	}
	delete(l.weights, l.ChildAt(index))
	l.Container.RemoveChildAt(index)
}
