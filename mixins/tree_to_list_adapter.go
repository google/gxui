// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/google/gxui"
	"github.com/google/gxui/math"
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

func (a TreeToListAdapter) Collapse(item gxui.AdapterItem) gxui.AdapterItem {
	n, i, _ := a.root.FindByItem(item)
	if n.Child(i).Collapse() {
		return n.Child(i).Item()
	}
	if n != a.root && n.Collapse() {
		return n.Item()
	}
	return item
}

func (a TreeToListAdapter) Expand(item gxui.AdapterItem) bool {
	n, i, _ := a.root.FindByItem(item)
	return n.Child(i).Expand()
}

func (a TreeToListAdapter) ExpandAllParents(item gxui.AdapterItem) bool {
	for a.Expand(a.DeepestVisibleAncestor(item)) {
	}
	a.DataChanged()
	return false
}

func (a TreeToListAdapter) DeepestVisibleAncestor(item gxui.AdapterItem) gxui.AdapterItem {
	n, i, _ := a.root.FindByItem(item)
	child := n.Child(i)
	return child.Item()
}

// Adapter compliance
func (a TreeToListAdapter) Count() int {
	return a.root.childCount
}

func (a TreeToListAdapter) ItemAt(index int) gxui.AdapterItem {
	return a.root.ItemAt(index)
}

func (a TreeToListAdapter) ItemIndex(item gxui.AdapterItem) int {
	return a.root.ItemIndex(item)
}

func (a TreeToListAdapter) Size(theme gxui.Theme) math.Size {
	return a.adapter.Size(theme)
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
	layout.SetDirection(gxui.LeftToRight)
	layout.AddChild(toggle)
	layout.AddChild(control)
	return layout
}
