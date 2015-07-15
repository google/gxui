// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package basic

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
	t.SetControlCreator(treeControlCreator{})

	return t
}

// mixins.Tree overrides
func (t *Tree) Paint(c gxui.Canvas) {
	r := t.Size().Rect()

	t.Tree.Paint(c)

	if t.HasFocus() {
		s := t.theme.FocusedStyle
		c.DrawRoundedRect(r, 3, 3, 3, 3, s.Pen, s.Brush)
	}
}

func (t *Tree) PaintMouseOverBackground(c gxui.Canvas, r math.Rect) {
	c.DrawRoundedRect(r, 2.0, 2.0, 2.0, 2.0, gxui.TransparentPen, gxui.CreateBrush(gxui.Gray15))
}

// mixins.List overrides
func (l *Tree) PaintSelection(c gxui.Canvas, r math.Rect) {
	s := l.theme.HighlightStyle
	c.DrawRoundedRect(r, 2.0, 2.0, 2.0, 2.0, s.Pen, s.Brush)
}

type treeControlCreator struct{}

func (treeControlCreator) Create(theme gxui.Theme, control gxui.Control, node *mixins.TreeToListNode) gxui.Control {
	img := theme.CreateImage()
	imgSize := math.Size{W: 10, H: 10}

	ll := theme.CreateLinearLayout()
	ll.SetDirection(gxui.LeftToRight)

	btn := theme.CreateButton()
	btn.SetBackgroundBrush(gxui.TransparentBrush)
	btn.SetBorderPen(gxui.CreatePen(1, gxui.Gray30))
	btn.SetMargin(math.Spacing{L: 1, R: 1, T: 1, B: 1})
	btn.OnClick(func(ev gxui.MouseEvent) {
		if ev.Button == gxui.MouseButtonLeft {
			node.ToggleExpanded()
		}
	})
	btn.AddChild(img)

	update := func() {
		expanded := node.IsExpanded()
		canvas := theme.Driver().CreateCanvas(imgSize)
		btn.SetVisible(!node.IsLeaf())
		switch {
		case !btn.IsMouseDown(gxui.MouseButtonLeft) && expanded:
			canvas.DrawPolygon(expandedPoly, gxui.TransparentPen, gxui.CreateBrush(gxui.Gray70))
		case !btn.IsMouseDown(gxui.MouseButtonLeft) && !expanded:
			canvas.DrawPolygon(collapsedPoly, gxui.TransparentPen, gxui.CreateBrush(gxui.Gray70))
		case expanded:
			canvas.DrawPolygon(expandedPoly, gxui.TransparentPen, gxui.CreateBrush(gxui.Gray30))
		case !expanded:
			canvas.DrawPolygon(collapsedPoly, gxui.TransparentPen, gxui.CreateBrush(gxui.Gray30))
		}
		canvas.Complete()
		img.SetCanvas(canvas)
	}
	btn.OnMouseDown(func(gxui.MouseEvent) { update() })
	btn.OnMouseUp(func(gxui.MouseEvent) { update() })
	update()

	gxui.WhileAttached(btn, node.OnChange, update)

	ll.AddChild(btn)
	ll.AddChild(control)
	ll.SetPadding(math.Spacing{L: 16 * node.Depth()})
	return ll
}

func (treeControlCreator) Size(theme gxui.Theme, treeControlSize math.Size) math.Size {
	return treeControlSize
}
