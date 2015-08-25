// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import "github.com/google/gxui/math"

// Tree is the interface of all controls that visualize a hierarchical tree
// structure of items.
type Tree interface {
	Focusable

	// SetAdapter binds the specified TreeAdapter to this Tree control, replacing
	// any previously bound adapter.
	SetAdapter(TreeAdapter)

	// TreeAdapter returns the currently bound adapter.
	Adapter() TreeAdapter

	// Show makes the specified item visible, expanding the tree if necessary.
	Show(AdapterItem)

	// ExpandAll expands all tree nodes.
	ExpandAll()

	// CollapseAll collapses all tree nodes.
	CollapseAll()

	// Selected returns the currently selected item.
	Selected() AdapterItem

	// Select makes the specified item selected. The tree will not automatically
	// expand to the newly selected item. If the Tree does not contain the
	// specified item, then Select returns false and the previous selection
	// remains unaltered.
	Select(AdapterItem) bool

	// OnSelectionChanged registers the function f to be called when the selection
	// changes.
	OnSelectionChanged(f func(AdapterItem)) EventSubscription
}

// TreeNodeContainer is the interface used by nodes that can hold sub-nodes in the tree.
type TreeNodeContainer interface {
	// Count returns the number of immediate child nodes.
	Count() int

	// Node returns the i'th child TreeNode.
	NodeAt(i int) TreeNode

	// ItemIndex returns the index of the child equal to item, or the index of the
	// child that indirectly contains item, or if the item is not found under this
	// node, -1.
	ItemIndex(item AdapterItem) int
}

// TreeNode is the interface used by nodes in the tree.
type TreeNode interface {
	TreeNodeContainer

	// Item returns the AdapterItem this node.
	// It is important for the TreeNode to return consistent AdapterItems for
	// the same data, so that selections can be persisted, or re-ordering
	// animations can be played when the dataset changes.
	// The AdapterItem returned must be equality-unique across the entire Adapter.
	Item() AdapterItem

	// Create returns a Control visualizing this node.
	Create(theme Theme) Control
}

// TreeAdapter is an interface used to visualize a set of hierarchical items.
// Users of the TreeAdapter should presume the data is unchanged until the
// OnDataChanged or OnDataReplaced events are fired.
type TreeAdapter interface {
	TreeNodeContainer

	// Size returns the size that each of the item's controls will be displayed
	// at for the given theme.
	Size(Theme) math.Size

	// OnDataChanged registers f to be called when there is a partial change in
	// the items of the adapter. Scroll positions and selections should be
	// preserved if possible.
	// If recreateControls is true then each of the visible controls should be
	// recreated by re-calling Create().
	OnDataChanged(f func(recreateControls bool)) EventSubscription

	// OnDataReplaced registers f to be called when there is a complete
	// replacement of items in the adapter.
	OnDataReplaced(f func()) EventSubscription
}
