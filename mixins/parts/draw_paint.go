// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parts

import (
	"fmt"
	"runtime"

	"github.com/google/gxui"
	"github.com/google/gxui/mixins/outer"
)

const debugVerifyDetachOnGC = false

type DrawPaintOuter interface {
	outer.Attachable
	outer.Painter
	outer.Parenter
	outer.Sized
}

type DrawPaint struct {
	outer           DrawPaintOuter
	driver          gxui.Driver
	canvas          gxui.Canvas
	dirty           bool
	redrawRequested bool
}

func verifyDetach(o DrawPaintOuter) {
	if o.Attached() {
		panic(fmt.Errorf("%T garbage collected while still attached", o))
	}
}

func (d *DrawPaint) Init(outer DrawPaintOuter, theme gxui.Theme) {
	d.outer = outer
	d.driver = theme.Driver()

	if debugVerifyDetachOnGC {
		runtime.SetFinalizer(d.outer, verifyDetach)
	}
}

func (d *DrawPaint) Redraw() {
	d.driver.AssertUIGoroutine()
	if !d.redrawRequested {
		if p := d.outer.Parent(); p != nil {
			d.redrawRequested = true
			p.Redraw()
		}
	}
}

func (d *DrawPaint) Draw() gxui.Canvas {
	if !d.outer.Attached() {
		panic(fmt.Errorf("Attempting to draw a non-attached control %T", d.outer))
	}

	s := d.outer.Size()
	if s.Area() == 0 {
		return nil // No area to draw in
	}
	if d.canvas == nil || d.canvas.Size() != s || d.redrawRequested {
		d.canvas = d.driver.CreateCanvas(s)
		d.redrawRequested = false
		d.outer.Paint(d.canvas)
		d.canvas.Complete()
	}
	return d.canvas
}
