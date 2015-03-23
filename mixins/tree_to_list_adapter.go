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

func CreateTreeToListAdapter(adapter gxui.TreeAdapter, ceb CreateExpandButton) *TreeToListAdapter {
	outer := &TreeToListAdapter{
		createExpandButton: ceb,
		adapter:            adapter,
		root:               CreateTreeInternalRoot(adapter),
	}
	adapter.OnDataReplaced(func() {
		outer.root = CreateTreeInternalRoot(adapter)
		outer.DataReplaced()
	})
	adapter.OnDataChanged(outer.DataChanged)

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
	if n, i, _ := a.root.FindByItem(item); n != nil {
		return n.Child(i).Expand()
	} else {
		return false
	}
}

func (a TreeToListAdapter) ExpandAllParents(item gxui.AdapterItem) {
	changed := false
	n := a.DeepestVisibleAncestor(item)
	for n != nil && n != item {
		if a.Expand(n) {
			changed = true
			n = a.DeepestVisibleAncestor(item)
		} else {
			break // Already expanded, nothing more to do.
		}
	}
	if changed {
		a.DataChanged()
	}
}

func (a TreeToListAdapter) DeepestVisibleAncestor(item gxui.AdapterItem) gxui.AdapterItem {
	if n, i, _ := a.root.FindByItem(item); n != nil {
		child := n.Child(i)
		return child.Item()
	} else {
		return nil
	}
}

// Adapter compliance
func (a TreeToListAdapter) Count() int {
	return a.root.descendants
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

	control := n.node.Create(theme, i)

	layout := theme.CreateLinearLayout()
	layout.SetPadding(math.Spacing{L: d * 16})
	layout.SetDirection(gxui.LeftToRight)
	layout.AddChild(toggle)
	layout.AddChild(control)
	return layout
}
