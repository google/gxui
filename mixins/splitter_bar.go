// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/google/gxui"
	"github.com/google/gxui/math"
	"github.com/google/gxui/mixins/base"
)

type SplitterBarOuter interface {
	base.ControlOuter
}

type SplitterBar struct {
	base.Control

	onDrag          func(wndPnt math.Point)
	outer           SplitterBarOuter
	theme           gxui.Theme
	onDragStart     gxui.Event
	onDragEnd       gxui.Event
	backgroundColor gxui.Color
	foregroundColor gxui.Color
	isDragging      bool
}

func (b *SplitterBar) Init(outer SplitterBarOuter, theme gxui.Theme) {
	b.Control.Init(outer, theme)

	b.outer = outer
	b.theme = theme
	b.onDragStart = gxui.CreateEvent(func(gxui.MouseEvent) {})
	b.onDragEnd = gxui.CreateEvent(func(gxui.MouseEvent) {})
	b.backgroundColor = gxui.Red
	b.foregroundColor = gxui.Green
}

func (b *SplitterBar) SetBackgroundColor(c gxui.Color) {
	b.backgroundColor = c
}

func (b *SplitterBar) SetForegroundColor(c gxui.Color) {
	b.foregroundColor = c
}

func (b *SplitterBar) OnSplitterDragged(f func(wndPnt math.Point)) {
	b.onDrag = f
}

func (b *SplitterBar) IsDragging() bool {
	return b.isDragging
}

func (b *SplitterBar) OnDragStart(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return b.onDragStart.Listen(f)
}

func (b *SplitterBar) OnDragEnd(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return b.onDragEnd.Listen(f)
}

// parts.DrawPaint overrides
func (b *SplitterBar) Paint(c gxui.Canvas) {
	r := b.outer.Size().Rect()
	c.DrawRect(r, gxui.CreateBrush(b.backgroundColor))
	if b.foregroundColor != b.backgroundColor {
		c.DrawRect(r.ContractI(1), gxui.CreateBrush(b.foregroundColor))
	}
}

// InputEventHandler overrides
func (b *SplitterBar) MouseDown(e gxui.MouseEvent) {
	b.isDragging = true
	b.onDragStart.Fire(e)
	var mms, mus gxui.EventSubscription
	mms = e.Window.OnMouseMove(func(we gxui.MouseEvent) {
		if b.onDrag != nil {
			b.onDrag(we.WindowPoint)
		}
	})
	mus = e.Window.OnMouseUp(func(we gxui.MouseEvent) {
		mms.Unlisten()
		mus.Unlisten()
		b.isDragging = false
		b.onDragEnd.Fire(we)
	})

	b.InputEventHandler.MouseDown(e)
}
