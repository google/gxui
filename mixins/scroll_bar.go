// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/google/gxui"
	"github.com/google/gxui/math"
	"github.com/google/gxui/mixins/base"
)

type ScrollBarOuter interface {
	base.ControlOuter
}

type ScrollBar struct {
	base.Control
	outer ScrollBarOuter

	orientation         gxui.Orientation
	thickness           int
	minBarLength        int
	scrollPositionFrom  int
	scrollPositionTo    int
	scrollLimit         int
	railBrush, barBrush gxui.Brush
	railPen, barPen     gxui.Pen
	barRect             math.Rect
	onScroll            gxui.Event
	autoHide            bool
}

func (s *ScrollBar) positionAt(p math.Point) int {
	o := s.orientation
	frac := float32(o.Major(p.XY())) / float32(o.Major(s.Size().WH()))
	max := s.ScrollLimit()
	return int(float32(max) * frac)
}

func (s *ScrollBar) rangeAt(p math.Point) (from, to int) {
	width := s.scrollPositionTo - s.scrollPositionFrom
	from = math.Clamp(s.positionAt(p), 0, s.scrollLimit-width)
	to = from + width
	return
}

func (s *ScrollBar) updateBarRect() {
	sf, st := s.ScrollFraction()
	size := s.Size()
	b := size.Rect()
	halfMinLen := s.minBarLength / 2
	if s.orientation.Horizontal() {
		b.Min.X = math.Lerp(0, size.W, sf)
		b.Max.X = math.Lerp(0, size.W, st)
		if b.W() < s.minBarLength {
			c := (b.Min.X + b.Max.X) / 2
			c = math.Clamp(c, b.Min.X+halfMinLen, b.Max.X-halfMinLen)
			b.Min.X, b.Max.X = c-halfMinLen, c+halfMinLen
		}
	} else {
		b.Min.Y = math.Lerp(0, size.H, sf)
		b.Max.Y = math.Lerp(0, size.H, st)
		if b.H() < s.minBarLength {
			c := (b.Min.Y + b.Max.Y) / 2
			c = math.Clamp(c, b.Min.Y+halfMinLen, b.Max.Y-halfMinLen)
			b.Min.Y, b.Max.Y = c-halfMinLen, c+halfMinLen
		}
	}
	s.barRect = b
}

func (s *ScrollBar) Init(outer ScrollBarOuter, theme gxui.Theme) {
	s.Control.Init(outer, theme)

	s.outer = outer
	s.thickness = 10
	s.minBarLength = 10
	s.scrollPositionFrom = 0
	s.scrollPositionTo = 100
	s.scrollLimit = 100
	s.onScroll = gxui.CreateEvent(s.SetScrollPosition)

	// Interface compliance test
	_ = gxui.ScrollBar(s)
}

func (s *ScrollBar) OnScroll(f func(from, to int)) gxui.EventSubscription {
	return s.onScroll.Listen(f)
}

func (s *ScrollBar) ScrollFraction() (from, to float32) {
	from = float32(s.scrollPositionFrom) / float32(s.scrollLimit)
	to = float32(s.scrollPositionTo) / float32(s.scrollLimit)
	return
}

func (s *ScrollBar) DesiredSize(min, max math.Size) math.Size {
	if s.orientation.Horizontal() {
		return math.Size{W: max.W, H: s.thickness}.Clamp(min, max)
	} else {
		return math.Size{W: s.thickness, H: max.H}.Clamp(min, max)
	}
}

func (s *ScrollBar) Paint(c gxui.Canvas) {
	c.DrawRoundedRect(s.outer.Size().Rect(), 3, 3, 3, 3, s.railPen, s.railBrush)
	c.DrawRoundedRect(s.barRect, 3, 3, 3, 3, s.barPen, s.barBrush)
}

func (s *ScrollBar) RailBrush() gxui.Brush {
	return s.railBrush
}

func (s *ScrollBar) SetRailBrush(b gxui.Brush) {
	if s.railBrush != b {
		s.railBrush = b
		s.Redraw()
	}
}

func (s *ScrollBar) BarBrush() gxui.Brush {
	return s.barBrush
}

func (s *ScrollBar) SetBarBrush(b gxui.Brush) {
	if s.barBrush != b {
		s.barBrush = b
		s.Redraw()
	}
}

func (s *ScrollBar) RailPen() gxui.Pen {
	return s.railPen
}

func (s *ScrollBar) SetRailPen(b gxui.Pen) {
	if s.railPen != b {
		s.railPen = b
		s.Redraw()
	}
}

func (s *ScrollBar) BarPen() gxui.Pen {
	return s.barPen
}

func (s *ScrollBar) SetBarPen(b gxui.Pen) {
	if s.barPen != b {
		s.barPen = b
		s.Redraw()
	}
}

func (s *ScrollBar) ScrollPosition() (from, to int) {
	return s.scrollPositionFrom, s.scrollPositionTo
}

func (s *ScrollBar) SetScrollPosition(from, to int) {
	if s.scrollPositionFrom != from || s.scrollPositionTo != to {
		s.scrollPositionFrom, s.scrollPositionTo = from, to
		s.updateBarRect()
		s.Redraw()
		s.onScroll.Fire(from, to)
	}
}

func (s *ScrollBar) ScrollLimit() int {
	return s.scrollLimit
}

func (s *ScrollBar) SetScrollLimit(l int) {
	if s.scrollLimit != l {
		s.scrollLimit = l
		s.updateBarRect()
		s.Redraw()
	}
}

func (s *ScrollBar) AutoHide() bool {
	return s.autoHide
}

func (s *ScrollBar) SetAutoHide(autoHide bool) {
	if s.autoHide != autoHide {
		s.autoHide = autoHide
		s.Redraw()
	}
}

func (s *ScrollBar) IsVisible() bool {
	if s.autoHide && s.scrollPositionFrom == 0 && s.scrollPositionTo == s.scrollLimit {
		return false
	}
	return s.Control.IsVisible()
}

func (s *ScrollBar) Orientation() gxui.Orientation {
	return s.orientation
}

func (s *ScrollBar) SetOrientation(o gxui.Orientation) {
	if s.orientation != o {
		s.orientation = o
		s.Redraw()
	}
}

// InputEventHandler overrides
func (s *ScrollBar) Click(ev gxui.MouseEvent) (consume bool) {
	if !s.barRect.Contains(ev.Point) {
		p := s.positionAt(ev.Point)
		from, to := s.scrollPositionFrom, s.scrollPositionTo
		switch {
		case p < from:
			width := to - from
			from = math.Max(from-width, 0)
			s.SetScrollPosition(from, from+width)
		case p > to:
			width := to - from
			to = math.Min(to+width, s.scrollLimit)
			s.SetScrollPosition(to-width, to)
		}
	}
	return true
}

func (s *ScrollBar) MouseDown(ev gxui.MouseEvent) {
	if s.barRect.Contains(ev.Point) {
		initialOffset := ev.Point.Sub(s.barRect.Min)
		var mms, mus gxui.EventSubscription
		mms = ev.Window.OnMouseMove(func(we gxui.MouseEvent) {
			p := gxui.WindowToChild(we.WindowPoint, s.outer)
			s.SetScrollPosition(s.rangeAt(p.Sub(initialOffset)))
		})
		mus = ev.Window.OnMouseUp(func(we gxui.MouseEvent) {
			mms.Unlisten()
			mus.Unlisten()
		})
	}
	s.InputEventHandler.MouseDown(ev)
}
