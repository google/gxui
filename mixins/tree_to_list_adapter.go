// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/google/gxui"
	"github.com/google/gxui/math"
)

// Interface used to visualize tree nodes in as a list.
type TreeControlCreator interface {
	// Create returns a Control that contains control (returned by the backing
	// TreeNode) and visualizes the expanded state of node.
	Create(theme gxui.Theme, control gxui.Control, node *TreeToListNode) gxui.Control

	// Size returns the size that each of the controls returned by Create will
	// be displayed at for the given theme.
	// treeControlSize is the size returned the backing TreeNode.
	Size(theme gxui.Theme, treeControlSize math.Size) math.Size
}

// TreeToListAdapter converts a TreeAdapter to a ListAdapter so that the
// tree can be visualized with a List.
type TreeToListAdapter struct {
	TreeToListNode
	gxui.AdapterBase
	adapter gxui.TreeAdapter
	creator TreeControlCreator
}

// CreateTreeToListAdapter wraps the provided TreeAdapter with an adapter
// conforming to the ListAdapter interface.
func CreateTreeToListAdapter(treeAdapter gxui.TreeAdapter, creator TreeControlCreator) *TreeToListAdapter {
	listAdapter := &TreeToListAdapter{}
	listAdapter.depth = -1 // The adapter is not a node, just a container.
	listAdapter.adapter = treeAdapter
	listAdapter.container = treeAdapter
	listAdapter.creator = creator
	treeAdapter.OnDataReplaced(func() {
		listAdapter.reset()
		listAdapter.DataReplaced()
	})
	treeAdapter.OnDataChanged(func() {
		listAdapter.update()
		listAdapter.DataChanged()
	})
	listAdapter.reset()
	return listAdapter
}

func (a *TreeToListAdapter) adjustDescendants(delta int) {
	if delta != 0 {
		a.descendants += delta
		a.DataChanged()
	}
}

// reset clears the current state of the tree.
func (a *TreeToListAdapter) reset() {
	a.descendants = a.adapter.Count()
	a.children = make([]*TreeToListNode, a.descendants)
	for i := range a.children {
		node := a.adapter.NodeAt(i)
		item := node.Item()
		a.children[i] = &TreeToListNode{container: node, item: item, parent: a}
	}
}

// Count returns the total number of expanded nodes in the tree.
func (a *TreeToListAdapter) Count() int {
	return a.descendants
}

// Create returns a Control visualizing the item at the specified index in the
// list of all the expanded nodes treated as as a flattened list.
func (a *TreeToListAdapter) Create(theme gxui.Theme, index int) gxui.Control {
	n := a.NodeAt(index)
	c := n.container.(gxui.TreeNode).Create(theme)
	return a.creator.Create(theme, c, n)
}

// Size returns the size that each of the item's controls will be displayed
// at for the given theme.
func (a *TreeToListAdapter) Size(theme gxui.Theme) math.Size {
	return a.creator.Size(theme, a.adapter.Size(theme))
}

// DeepestNode returns the deepest expanded node to represent item.
// If the item is not found in the adapter, then nil is returned.
func (a *TreeToListAdapter) DeepestNode(item gxui.AdapterItem) *TreeToListNode {
	n := &a.TreeToListNode
	for {
		if i := n.DirectItemIndex(item); i >= 0 {
			n = n.children[i]
		} else {
			return nil
		}
		if item == n.item || !n.IsExpanded() {
			return n
		}
	}
}
