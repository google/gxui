// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parts

import (
	"gaze/gxui"
	"gaze/gxui/assert"
	"gaze/gxui/mixins/outer"
	"runtime"
)

const debugVerifyDetachOnGC = false

type DrawPaintOuter interface {
	outer.Attachable
	outer.Bounds
	outer.Painter
	outer.Parenter
}

type DrawPaint struct {
	outer           DrawPaintOuter
	theme           gxui.Theme
	canvas          gxui.Canvas
	dirty           bool
	redrawRequested bool
}

func verifyDetach(o DrawPaintOuter) {
	assert.False(o.Attached(), "%T garbage collected while still attached", o)
}

func (d *DrawPaint) Init(outer DrawPaintOuter, theme gxui.Theme) {
	d.outer = outer
	d.theme = theme
	outer.OnDetach(func() {
		if d.canvas != nil {
			d.canvas.Release()
			d.canvas = nil
		}
	})

	if debugVerifyDetachOnGC {
		runtime.SetFinalizer(d.outer, verifyDetach)
	}
}

func (d *DrawPaint) Redraw() {
	if !d.redrawRequested {
		if p := d.outer.Parent(); p != nil {
			d.redrawRequested = true
			p.Redraw()
		}
	}
}

func (d *DrawPaint) Draw() gxui.Canvas {
	assert.True(d.outer.Attached(), "Attempting to draw a non-attached control")

	s := d.outer.Bounds().Size()
	if s.Area() == 0 {
		return nil // No area to draw in
	}
	if d.canvas == nil || d.canvas.Size() != s || d.redrawRequested {
		if d.canvas != nil {
			d.canvas.Release()
		}
		d.canvas = d.theme.Driver().CreateCanvas(s)
		d.redrawRequested = false
		d.outer.Paint(d.canvas)
		d.canvas.Complete()
	}
	return d.canvas
}
