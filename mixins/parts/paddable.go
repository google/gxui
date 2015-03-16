// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parts

import (
	"github.com/google/gxui/math"
	"github.com/google/gxui/mixins/outer"
)

type PaddableOuter interface {
	outer.LayoutChildren
	outer.Redrawer
}

type Paddable struct {
	outer   PaddableOuter
	padding math.Spacing
}

func (p *Paddable) Init(outer PaddableOuter) {
	p.outer = outer
}

func (p *Paddable) SetPadding(m math.Spacing) {
	p.padding = m
	p.outer.LayoutChildren()
	p.outer.Redraw()
}

func (p *Paddable) Padding() math.Spacing {
	return p.padding
}
