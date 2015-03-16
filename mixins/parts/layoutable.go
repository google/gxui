// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parts

import (
	"github.com/google/gxui/assert"
	"github.com/google/gxui/math"
	"github.com/google/gxui/mixins/outer"
)

type LayoutableOuter interface {
	outer.Parenter
	outer.Redrawer
}

type Layoutable struct {
	outer             LayoutableOuter
	margin            math.Spacing
	rect              math.Rect
	relayoutRequested bool
}

func (l *Layoutable) Init(outer LayoutableOuter) {
	l.outer = outer
}

func (l *Layoutable) SetMargin(m math.Spacing) {
	l.margin = m
	if p := l.outer.Parent(); p != nil {
		p.Relayout()
	}
}

func (l *Layoutable) Margin() math.Spacing {
	return l.margin
}

func (l *Layoutable) Bounds() math.Rect {
	return l.rect
}

func (l *Layoutable) Layout(rect math.Rect) {
	assert.False(rect.W() < 0, "Layout() called with a negative width. Rect: %v", rect)
	assert.False(rect.H() < 0, "Layout() called with a negative height. Rect: %v", rect)

	boundsChanged := l.rect != rect
	l.rect = rect
	if l.relayoutRequested || boundsChanged {
		callLayoutChildrenIfSupported(l.outer)
		l.outer.Redraw()
		l.relayoutRequested = false
	}
}

func (l *Layoutable) Relayout() {
	if !l.relayoutRequested {
		if p := l.outer.Parent(); p != nil {
			l.relayoutRequested = true
			p.Relayout()
		}
	}
}
