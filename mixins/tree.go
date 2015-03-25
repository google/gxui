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
	outer       TreeOuter
	treeAdapter gxui.TreeAdapter
	listAdapter *TreeToListAdapter
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
	if t.treeAdapter == adapter {
		return
	}
	if adapter != nil {
		t.treeAdapter = adapter
		t.listAdapter = CreateTreeToListAdapter(adapter, t.outer.CreateExpandButton)
		t.List.SetAdapter(t.listAdapter)
	} else {
		t.listAdapter = nil
		t.treeAdapter = nil
		t.List.SetAdapter(nil)
	}
}

func (t *Tree) Adapter() gxui.TreeAdapter {
	return t.treeAdapter
}

func (t *Tree) Show(item gxui.AdapterItem) {
	t.listAdapter.ExpandAllParents(item)
	t.List.ScrollTo(item)
}

func (t *Tree) ExpandAll() {
	t.listAdapter.root.ExpandAll()
	t.DataChanged()
}

func (t *Tree) CollapseAll() {
	t.listAdapter.root.CollapseAll()
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
	if t.selectedItem != nil {
		item := t.listAdapter.DeepestVisibleAncestor(t.selectedItem)
		if item != t.selectedItem {
			// The selected item is hidden by an unexpanded node.
			// Highlight the deepest visible node instead.
			if details, found := t.details[item]; found {
				if child == details.control {
					b := child.Bounds().Expand(child.Margin())
					t.outer.PaintUnexpandedSelection(c, b)
				}
			}
		}
	}
}

// InputEventHandler override
func (t *Tree) KeyPress(ev gxui.KeyboardEvent) (consume bool) {
	item := t.Selected()
	switch ev.Key {
	case gxui.KeyLeft:
		newItem := t.listAdapter.Collapse(item)
		if newItem != item {
			t.Select(newItem)
			return true
		}
	case gxui.KeyRight:
		if t.listAdapter.Expand(item) {
			return true
		}
	}
	return t.List.KeyPress(ev)
}
