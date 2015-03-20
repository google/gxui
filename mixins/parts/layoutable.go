// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parts

import (
	"fmt"

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
	inLayoutChildren  bool // True when calling LayoutChildren
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
	if rect.W() < 0 {
		panic(fmt.Errorf("Layout() called with a negative width. Rect: %v", rect))
	}
	if rect.H() < 0 {
		panic(fmt.Errorf("Layout() called with a negative height. Rect: %v", rect))
	}

	boundsChanged := l.rect != rect
	l.rect = rect
	if l.relayoutRequested || boundsChanged {
		l.relayoutRequested = false
		l.inLayoutChildren = true
		callLayoutChildrenIfSupported(l.outer)
		l.inLayoutChildren = false
		l.outer.Redraw()
	}
}

func (l *Layoutable) Relayout() {
	if l.inLayoutChildren {
		panic("Cannot call Relayout() while in LayoutChildren")
	}
	if !l.relayoutRequested {
		if p := l.outer.Parent(); p != nil {
			l.relayoutRequested = true
			p.Relayout()
		}
	}
}
