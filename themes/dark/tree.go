// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dark

import (
	"github.com/google/gxui"
	"github.com/google/gxui/math"
	"github.com/google/gxui/mixins"
)

type Tree struct {
	mixins.Tree
	theme *Theme
}

var expandedPoly = gxui.Polygon{
	gxui.PolygonVertex{Position: math.Point{X: 2, Y: 3}},
	gxui.PolygonVertex{Position: math.Point{X: 8, Y: 3}},
	gxui.PolygonVertex{Position: math.Point{X: 5, Y: 8}},
}

var collapsedPoly = gxui.Polygon{
	gxui.PolygonVertex{Position: math.Point{X: 3, Y: 2}},
	gxui.PolygonVertex{Position: math.Point{X: 8, Y: 5}},
	gxui.PolygonVertex{Position: math.Point{X: 3, Y: 8}},
}

func CreateTree(theme *Theme) gxui.Tree {
	t := &Tree{}
	t.Init(t, theme)
	t.SetPadding(math.Spacing{L: 3, T: 3, R: 3, B: 3})
	t.SetBorderPen(gxui.TransparentPen)
	t.theme = theme

	return t
}

// mixins.Tree overrides
func (t *Tree) Paint(c gxui.Canvas) {
	r := t.Bounds().Size().Rect()

	t.Tree.Paint(c)

	if t.HasFocus() {
		s := t.theme.FocusedStyle
		c.DrawRoundedRect(r, 3, 3, 3, 3, s.Pen, s.Brush)
	}
}

func (t *Tree) PaintMouseOverBackground(c gxui.Canvas, r math.Rect) {
	c.DrawRoundedRect(r, 2.0, 2.0, 2.0, 2.0, gxui.TransparentPen, gxui.CreateBrush(gxui.Gray15))
}

func (t *Tree) CreateExpandButton(theme gxui.Theme, node *mixins.TreeInternalNode) gxui.Button {
	img := theme.CreateImage()
	imgSize := math.Size{W: 10, H: 10}

	btn := theme.CreateButton()
	btn.SetBackgroundBrush(gxui.TransparentBrush)
	btn.SetBorderPen(gxui.CreatePen(1, gxui.Gray30))
	btn.SetMargin(math.Spacing{L: 1, R: 1, T: 1, B: 1})
	btn.OnClick(func(ev gxui.MouseEvent) {
		if ev.Button == gxui.MouseButtonLeft {
			if node.IsExpanded() {
				node.Collapse()
			} else {
				node.Expand()
			}
		}
	})
	btn.AddChild(img)

	updateStyle := func() {
		canvas := theme.Driver().CreateCanvas(imgSize)
		switch {
		case !btn.IsMouseDown(gxui.MouseButtonLeft) && node.IsExpanded():
			canvas.DrawPolygon(expandedPoly, gxui.TransparentPen, gxui.CreateBrush(gxui.Gray70))
		case !btn.IsMouseDown(gxui.MouseButtonLeft) && !node.IsExpanded():
			canvas.DrawPolygon(collapsedPoly, gxui.TransparentPen, gxui.CreateBrush(gxui.Gray70))
		case node.IsExpanded():
			canvas.DrawPolygon(expandedPoly, gxui.TransparentPen, gxui.CreateBrush(gxui.Gray30))
		case !node.IsExpanded():
			canvas.DrawPolygon(collapsedPoly, gxui.TransparentPen, gxui.CreateBrush(gxui.Gray30))
		}
		canvas.Complete()
		img.SetCanvas(canvas)
	}
	btn.OnMouseDown(func(gxui.MouseEvent) { updateStyle() })
	btn.OnMouseUp(func(gxui.MouseEvent) { updateStyle() })
	node.OnExpandedChanged(func(e bool) { updateStyle() })
	updateStyle()
	return btn
}

// mixins.List overrides
func (l *Tree) PaintSelection(c gxui.Canvas, r math.Rect) {
	s := l.theme.HighlightStyle
	c.DrawRoundedRect(r, 2.0, 2.0, 2.0, 2.0, s.Pen, s.Brush)
}
