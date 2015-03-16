// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parts

import (
	"github.com/google/gxui/mixins/outer"
)

type VisibleOuter interface {
	outer.Redrawer
	outer.Parenter
}

type Visible struct {
	outer   VisibleOuter
	visible bool
}

func (v *Visible) Init(outer VisibleOuter) {
	v.outer = outer
	v.visible = true
}

func (v *Visible) IsVisible() bool {
	return v.visible
}

func (v *Visible) SetVisible(visible bool) {
	if v.visible != visible {
		v.visible = visible
		if p := v.outer.Parent(); p != nil {
			p.Redraw()
		}
	}
}
