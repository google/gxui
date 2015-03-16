// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dark

import (
	"github.com/google/gxui"
	"github.com/google/gxui/math"
	"github.com/google/gxui/mixins"
	"time"
)

type ProgressBar struct {
	mixins.ProgressBar
	theme        *Theme
	timer        *time.Timer
	chevrons     gxui.Canvas
	chevronWidth int
	scroll       int
}

func CreateProgressBar(theme *Theme) gxui.ProgressBar {
	b := &ProgressBar{}
	b.Init(b, theme)
	b.theme = theme
	b.OnAttach(func() {
		events := theme.Driver().Events()
		b.timer = time.AfterFunc(time.Millisecond*50, func() {
			b.scroll = (b.scroll + 1) % (b.chevronWidth * 2)
			events <- b.Redraw
			b.timer.Reset(time.Millisecond * 50)
		})
	})
	b.OnDetach(func() {
		if b.chevrons != nil {
			b.chevrons.Release()
			b.chevrons = nil
			b.timer.Stop()
			b.timer = nil
		}
	})
	b.SetBackgroundBrush(gxui.CreateBrush(gxui.Gray10))
	b.SetBorderPen(gxui.CreatePen(1, gxui.Gray40))
	return b
}

func (b *ProgressBar) Layout(r math.Rect) {
	b.ProgressBar.Layout(r)

	if b.chevrons != nil {
		b.chevrons.Release()
		b.chevrons = nil
	}
	if r.Size().Area() > 0 {
		b.chevrons = b.theme.Driver().CreateCanvas(r.Size())
		b.chevronWidth = r.H() / 2
		cw := b.chevronWidth
		for x := -cw * 2; x < r.W(); x += cw * 2 {
			// x0    x2
			// |  x1 |  x3
			//    |     |
			// A-----B    - y0
			//  \     \
			//   \     \
			//    F     C - y1
			//   /     /
			//  /     /
			// E-----D    - y2
			y0, y1, y2 := 0, r.H()/2, r.H()
			x0, x1 := x, x+cw/2
			x2, x3 := x0+cw, x1+cw
			var chevron = gxui.Polygon{
				/* A */ gxui.PolygonVertex{Position: math.Point{X: x0, Y: y0}},
				/* B */ gxui.PolygonVertex{Position: math.Point{X: x2, Y: y0}},
				/* C */ gxui.PolygonVertex{Position: math.Point{X: x3, Y: y1}},
				/* D */ gxui.PolygonVertex{Position: math.Point{X: x2, Y: y2}},
				/* E */ gxui.PolygonVertex{Position: math.Point{X: x0, Y: y2}},
				/* F */ gxui.PolygonVertex{Position: math.Point{X: x1, Y: y1}},
			}
			b.chevrons.DrawPolygon(chevron, gxui.TransparentPen, gxui.CreateBrush(gxui.Gray30))
		}
		b.chevrons.Complete()
	}
}

func (b *ProgressBar) PaintProgress(c gxui.Canvas, r math.Rect, frac float32) {
	r.Max.X = math.Lerp(r.Min.X, r.Max.X, frac)
	c.DrawRect(r, gxui.CreateBrush(gxui.Gray50))
	c.Push()
	c.AddClip(r)
	c.DrawCanvas(b.chevrons, math.Point{X: b.scroll})
	c.Pop()
}
