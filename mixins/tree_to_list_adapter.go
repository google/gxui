// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"gxui"
	"gxui/math"
)

type CreateExpandButton func(theme gxui.Theme, node *TreeInternalNode) gxui.Button

type TreeToListAdapter struct {
	gxui.AdapterBase
	createExpandButton CreateExpandButton
	adapter            gxui.TreeAdapter
	root               *TreeInternalNode
}

func CreateTreeToListAdapter(inner gxui.TreeAdapter, ceb CreateExpandButton) *TreeToListAdapter {
	outer := &TreeToListAdapter{
		createExpandButton: ceb,
		adapter:            inner,
		root:               CreateTreeInternalRoot(inner),
	}
	inner.OnDataReplaced(func() {
		outer.root = CreateTreeInternalRoot(inner)
		outer.DataReplaced()
	})
	inner.OnDataChanged(outer.DataChanged)

	return outer
}

func (a TreeToListAdapter) Collapse(id gxui.AdapterItemId) gxui.AdapterItemId {
	n, i, _ := a.root.FindById(id)
	if n.Child(i).Collapse() {
		a.DataChanged()
		return n.Child(i).Id()
	}
	if n != a.root && n.Collapse() {
		a.DataChanged()
		return n.Id()
	}
	return id
}

func (a TreeToListAdapter) Expand(id gxui.AdapterItemId) bool {
	n, i, _ := a.root.FindById(id)
	if n.Child(i).Expand() {
		a.DataChanged()
		return true
	}
	return false
}

func (a TreeToListAdapter) ExpandAllParents(id gxui.AdapterItemId) bool {
	for a.Expand(a.DeepestVisibleAncestor(id)) {
	}
	return false
}

func (a TreeToListAdapter) DeepestVisibleAncestor(id gxui.AdapterItemId) gxui.AdapterItemId {
	n, i, _ := a.root.FindById(id)
	child := n.Child(i)
	return child.Id()
}

// Adapter compliance
func (a TreeToListAdapter) ItemSize(theme gxui.Theme) math.Size {
	return a.adapter.ItemSize(theme)
}

func (a TreeToListAdapter) Count() int {
	return a.root.childCount
}

func (a TreeToListAdapter) ItemId(index int) gxui.AdapterItemId {
	return a.root.ItemId(index)
}

func (a TreeToListAdapter) ItemIndex(id gxui.AdapterItemId) int {
	return a.root.ItemIndex(id)
}

func (a TreeToListAdapter) Create(theme gxui.Theme, index int) gxui.Control {
	n, i, d := a.root.FindByIndex(index)
	child := n.Child(i)
	toggle := a.createExpandButton(theme, child)
	toggle.SetVisible(!child.IsLeaf())

	child.OnExpandedChanged(func(e bool) {
		a.DataChanged()
	})

	control := n.adapterNode.Create(theme, i)

	layout := theme.CreateLinearLayout()
	layout.SetPadding(math.Spacing{L: d * 16, T: 0, R: 0, B: 0})
	layout.SetOrientation(gxui.Horizontal)
	layout.AddChild(toggle)
	layout.AddChild(control)
	return layout
}
