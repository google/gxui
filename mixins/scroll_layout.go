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

type ScrollLayoutOuter interface {
	base.ContainerOuter
}

type ScrollLayout struct {
	base.Container
	parts.BackgroundBorderPainter

	outer                  ScrollLayoutOuter
	theme                  gxui.Theme
	scrollOffset           math.Point
	canScrollX, canScrollY bool
	scrollBarX, scrollBarY *gxui.Child
	child                  *gxui.Child
	innerSize              math.Size
}

func (l *ScrollLayout) Init(outer ScrollLayoutOuter, theme gxui.Theme) {
	l.Container.Init(outer, theme)
	l.BackgroundBorderPainter.Init(outer)

	l.outer = outer
	l.theme = theme
	l.canScrollX = true
	l.canScrollY = true
	scrollBarX := theme.CreateScrollBar()
	scrollBarX.SetOrientation(gxui.Horizontal)
	scrollBarX.OnScroll(func(from, to int) { l.SetScrollOffset(math.Point{X: from, Y: l.scrollOffset.Y}) })
	scrollBarY := theme.CreateScrollBar()
	scrollBarY.SetOrientation(gxui.Vertical)
	scrollBarY.OnScroll(func(from, to int) { l.SetScrollOffset(math.Point{X: l.scrollOffset.X, Y: from}) })
	l.scrollBarX = l.AddChild(scrollBarX)
	l.scrollBarY = l.AddChild(scrollBarY)
	l.SetMouseEventTarget(true)

	// Interface compliance test
	_ = gxui.ScrollLayout(l)
}

func (l *ScrollLayout) LayoutChildren() {
	s := l.outer.Size().Contract(l.Padding())
	o := l.Padding().LT()

	var sxs, sys math.Size
	if l.canScrollX {
		sxs = l.scrollBarX.Control.DesiredSize(math.ZeroSize, s)
	}
	if l.canScrollY {
		sys = l.scrollBarY.Control.DesiredSize(math.ZeroSize, s)
	}

	l.scrollBarX.Layout(math.CreateRect(0, s.H-sxs.H, s.W-sys.W, s.H).Canon().Offset(o))
	l.scrollBarY.Layout(math.CreateRect(s.W-sys.W, 0, s.W, s.H-sxs.H).Canon().Offset(o))

	l.innerSize = s.Contract(math.Spacing{R: sys.W, B: sxs.H})

	if l.child != nil {
		max := l.innerSize
		if l.canScrollX {
			max.W = math.MaxSize.W
		}
		if l.canScrollY {
			max.H = math.MaxSize.H
		}
		cs := l.child.Control.DesiredSize(math.ZeroSize, max)
		l.child.Layout(cs.Rect().Offset(l.scrollOffset.Neg()).Offset(o))
		l.scrollBarX.Control.(gxui.ScrollBar).SetScrollLimit(cs.W)
		l.scrollBarY.Control.(gxui.ScrollBar).SetScrollLimit(cs.H)
	}

	l.SetScrollOffset(l.scrollOffset)
}

func (l *ScrollLayout) DesiredSize(min, max math.Size) math.Size {
	return max
}

func (l *ScrollLayout) SetScrollOffset(scrollOffset math.Point) bool {
	var cs math.Size
	if l.child != nil {
		cs = l.child.Control.Size()
	}

	s := l.innerSize
	scrollOffset = scrollOffset.Min(cs.Sub(s).Point()).Max(math.Point{})

	l.scrollBarX.Control.SetVisible(l.canScrollX && cs.W > s.W)
	l.scrollBarY.Control.SetVisible(l.canScrollY && cs.H > s.H)
	l.scrollBarX.Control.(gxui.ScrollBar).SetScrollPosition(l.scrollOffset.X, l.scrollOffset.X+s.W)
	l.scrollBarY.Control.(gxui.ScrollBar).SetScrollPosition(l.scrollOffset.Y, l.scrollOffset.Y+s.H)

	if l.scrollOffset != scrollOffset {
		l.scrollOffset = scrollOffset
		l.Relayout()
		return true
	}

	return false
}

// InputEventHandler override
func (l *ScrollLayout) MouseScroll(ev gxui.MouseEvent) (consume bool) {
	if ev.ScrollY == 0 {
		return l.InputEventHandler.MouseScroll(ev)
	}
	switch {
	case l.canScrollY:
		return l.SetScrollOffset(l.scrollOffset.AddY(-ev.ScrollY))
	case l.canScrollX:
		return l.SetScrollOffset(l.scrollOffset.AddX(-ev.ScrollY))
	default:
		return false
	}
}

// gxui.ScrollLayout complaince
func (l *ScrollLayout) SetChild(control gxui.Control) {
	if l.child != nil {
		l.RemoveChild(l.child.Control)
	}
	if control != nil {
		l.child = l.AddChildAt(0, control)
	}
}

func (l *ScrollLayout) Child() gxui.Control {
	return l.child.Control
}

func (l *ScrollLayout) SetScrollAxis(horizontal, vertical bool) {
	if l.canScrollX != horizontal || l.canScrollY != vertical {
		l.canScrollX, l.canScrollY = horizontal, vertical
		l.scrollBarX.Control.SetVisible(horizontal)
		l.scrollBarY.Control.SetVisible(vertical)
		l.Relayout()
	}
}

func (l *ScrollLayout) ScrollAxis() (horizontal, vertical bool) {
	return l.canScrollX, l.canScrollY
}
