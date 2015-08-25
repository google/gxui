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
	gxui.AdapterBase
	node    TreeToListNode
	adapter gxui.TreeAdapter
	creator TreeControlCreator
}

// CreateTreeToListAdapter wraps the provided TreeAdapter with an adapter
// conforming to the ListAdapter interface.
func CreateTreeToListAdapter(treeAdapter gxui.TreeAdapter, creator TreeControlCreator) *TreeToListAdapter {
	listAdapter := &TreeToListAdapter{}
	listAdapter.node.depth = -1 // The node is just a container.
	listAdapter.node.container = treeAdapter
	listAdapter.adapter = treeAdapter
	listAdapter.creator = creator
	treeAdapter.OnDataReplaced(func() {
		listAdapter.reset()
		listAdapter.DataReplaced()
	})
	treeAdapter.OnDataChanged(func(recreateControls bool) {
		listAdapter.node.update(listAdapter)
		listAdapter.DataChanged(recreateControls)
	})
	listAdapter.reset()
	return listAdapter
}

func (a *TreeToListAdapter) adjustDescendants(delta int) {
	if delta != 0 {
		a.node.descendants += delta
		a.DataChanged(false)
	}
}

// reset clears the current state of the tree.
func (a *TreeToListAdapter) reset() {
	count := a.adapter.Count()
	a.node.descendants = count
	a.node.children = make([]*TreeToListNode, count)
	for i := range a.node.children {
		node := a.adapter.NodeAt(i)
		item := node.Item()
		a.node.children[i] = &TreeToListNode{container: node, item: item, parent: a}
	}
}

// Count returns the total number of expanded nodes in the tree.
func (a *TreeToListAdapter) Count() int {
	return a.node.descendants
}

// Create returns a Control visualizing the item at the specified index in the
// list of all the expanded nodes treated as as a flattened list.
func (a *TreeToListAdapter) Create(theme gxui.Theme, index int) gxui.Control {
	n := a.node.NodeAt(index)
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
	n := &a.node
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

// ItemAt returns the idx'th item of all the expanded nodes treated as a
// flattened list.
// Index 0 represents the first root node, index 1 may represent the the second
// root node or the first child of the first root node, and so on.
func (a *TreeToListAdapter) ItemAt(idx int) gxui.AdapterItem {
	return a.node.ItemAt(idx)
}

// ItemIndex returns the index of item in the list of all the expanded nodes
// treated as a flattened list.
// Index 0 represents the first root node, index 1 may represent the the second
// root node or the first child of the first root node, and so on.
func (a *TreeToListAdapter) ItemIndex(item gxui.AdapterItem) int {
	return a.node.ItemIndex(item)
}

// ExpandItem expands the tree to show item.
func (a *TreeToListAdapter) ExpandItem(item gxui.AdapterItem) {
	node := &a.node
	for {
		idx := node.DirectItemIndex(item)
		if idx < 0 {
			break
		}
		node = node.children[idx]
		node.Expand()
	}
}

// ExpandAll expands this node and all child nodes.
func (a *TreeToListAdapter) ExpandAll() {
	for _, n := range a.node.children {
		n.ExpandAll()
	}
}

// CollapseAll collapses this node and all child nodes.
func (a *TreeToListAdapter) CollapseAll() {
	for _, n := range a.node.children {
		n.CollapseAll()
	}
}

// Contains returns true if item is part of the tree (regardless of whether it
// is part of the expanded tree or not).
func (a *TreeToListAdapter) Contains(item gxui.AdapterItem) bool {
	return a.node.DirectItemIndex(item) >= 0
}
