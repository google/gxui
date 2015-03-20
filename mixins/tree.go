// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/google/gxui"
	"github.com/google/gxui/math"
	"github.com/google/gxui/mixins/parts"
)

type TreeOuter interface {
	ListOuter
	CreateExpandButton(theme gxui.Theme, node *TreeInternalNode) gxui.Button
	PaintUnexpandedSelection(c gxui.Canvas, r math.Rect)
}

type Tree struct {
	List
	parts.Focusable
	outer        TreeOuter
	adapterInner gxui.TreeAdapter
	adapterOuter *TreeToListAdapter
}

func (t *Tree) Init(outer TreeOuter, theme gxui.Theme) {
	t.List.Init(outer, theme)
	t.Focusable.Init(outer)
	t.outer = outer

	// Interface compliance test
	_ = gxui.Tree(t)
}

// gxui.Tree complaince
func (t *Tree) SetAdapter(adapter gxui.TreeAdapter) {
	if t.adapterInner == adapter {
		return
	}
	if adapter != nil {
		t.adapterInner = adapter
		t.adapterOuter = CreateTreeToListAdapter(adapter, t.outer.CreateExpandButton)
		t.List.SetAdapter(t.adapterOuter)
	} else {
		t.adapterOuter = nil
		t.adapterInner = nil
		t.List.SetAdapter(nil)
	}
}

func (t *Tree) Adapter() gxui.TreeAdapter {
	return t.adapterInner
}

func (t *Tree) Show(id gxui.AdapterItemId) {
	t.adapterOuter.ExpandAllParents(id)
	t.List.ScrollTo(id)
}

func (t *Tree) ExpandAll() {
	t.adapterOuter.root.ExpandAll()
	t.DataChanged()
}

func (t *Tree) CollapseAll() {
	t.adapterOuter.root.CollapseAll()
	t.DataChanged()
}

func (t *Tree) CreateExpandButton(theme gxui.Theme, node *TreeInternalNode) gxui.Button {
	btn := theme.CreateButton()
	btn.SetMargin(math.Spacing{L: 2, R: 2, T: 1, B: 1})
	btn.OnClick(func(ev gxui.MouseEvent) {
		if ev.Button == gxui.MouseButtonLeft {
			if node.IsExpanded() {
				node.Collapse()
			} else {
				node.Expand()
			}
		}
	})
	node.OnExpandedChanged(func(e bool) {
		if e {
			btn.SetText("-")
		} else {
			btn.SetText("+")
		}
	})
	if node.IsExpanded() {
		btn.SetText("-")
	} else {
		btn.SetText("+")
	}
	return btn
}

func (t *Tree) PaintUnexpandedSelection(c gxui.Canvas, r math.Rect) {
	c.DrawRoundedRect(r, 2.0, 2.0, 2.0, 2.0, gxui.CreatePen(1, gxui.Gray50), gxui.TransparentBrush)
}

// List override
func (t *Tree) PaintChild(c gxui.Canvas, child gxui.Control, idx int) {
	t.List.PaintChild(c, child, idx)
	if t.selectedId != gxui.InvalidAdapterItemId {
		id := t.adapterOuter.DeepestVisibleAncestor(t.selectedId)
		if id != t.selectedId {
			// The selected item is hidden by an unexpanded node.
			// Highlight the deepest visible node instead.
			if item, found := t.items[id]; found {
				if child == item.Control {
					b := child.Bounds().Expand(child.Margin())
					t.outer.PaintUnexpandedSelection(c, b)
				}
			}
		}
	}
}

// InputEventHandler override
func (t *Tree) KeyPress(ev gxui.KeyboardEvent) (consume bool) {
	id := t.Selected()
	switch ev.Key {
	case gxui.KeyLeft:
		newId := t.adapterOuter.Collapse(id)
		if newId != id {
			t.Select(newId)
			return true
		}
	case gxui.KeyRight:
		if t.adapterOuter.Expand(id) {
			return true
		}
	}
	return t.List.KeyPress(ev)
}
