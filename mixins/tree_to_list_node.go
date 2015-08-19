// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"fmt"

	"github.com/google/gxui"
)

type treeToListNodeParent interface {
	adjustDescendants(delta int)
}

type TreeToListNode struct {
	item        gxui.AdapterItem       // The wrapped AdapterItem.
	container   gxui.TreeNodeContainer // The wrapped TreeNode.
	descendants int                    // Total number of descendants.
	children    []*TreeToListNode      // The child nodes if expanded, or nil if collapsed.
	parent      treeToListNodeParent   // The parent of this node.
	depth       int                    // The depth of this node.
	onChange    gxui.Event             // event()
}

func (n *TreeToListNode) adjustDescendants(delta int) {
	n.descendants += delta
	n.parent.adjustDescendants(delta)
}

func (n *TreeToListNode) update(nAsParent treeToListNodeParent) {
	if n.IsExpanded() {
		// Build a map of item -> child for the current state.
		m := make(map[gxui.AdapterItem]*TreeToListNode, len(n.children))
		for _, c := range n.children {
			m[c.item] = c
		}

		// Re-create children, reusing previous nodes if found.
		depth := n.depth + 1
		n.descendants = 0
		n.children = make([]*TreeToListNode, n.container.Count())
		for i := range n.children {
			node := n.container.NodeAt(i)
			item := node.Item()
			if p, ok := m[item]; ok {
				p.container = node
				p.update(p)
				n.children[i] = p
				n.descendants += p.descendants + 1
			} else {
				n.children[i] = &TreeToListNode{container: node, item: item, parent: nAsParent, depth: depth}
				n.descendants++
			}
		}
	}
	if n.onChange != nil {
		n.onChange.Fire()
	}
}

// Item returns the AdapterItem this node represents.
func (n *TreeToListNode) Item() gxui.AdapterItem {
	return n.item
}

// Depth returns the depth of this node.
func (n *TreeToListNode) Depth() int {
	return n.depth
}

// IsExpanded returns true if the node is currently expanded.
func (n *TreeToListNode) IsExpanded() bool {
	return n.children != nil
}

// IsLeaf returns true if the node is a leaf in the tree.
func (n *TreeToListNode) IsLeaf() bool {
	return n.container == nil || n.container.Count() == 0
}

// Expand attempts to expand the node, returning true if the node expands.
// If the node is already expanded or is a leaf then Expand returns false.
func (n *TreeToListNode) Expand() bool {
	if n.parent == nil {
		panic("Expand cannot be called for root nodes")
	}
	if n.IsExpanded() || n.IsLeaf() {
		return false
	}
	depth := n.depth + 1
	n.descendants = n.container.Count()
	n.children = make([]*TreeToListNode, n.descendants)
	for i := range n.children {
		node := n.container.NodeAt(i)
		item := node.Item()
		n.children[i] = &TreeToListNode{container: node, item: item, parent: n, depth: depth}
	}
	n.parent.adjustDescendants(n.descendants)
	if n.onChange != nil {
		n.onChange.Fire()
	}
	return true
}

// Collapse attempts to collapse the node, returning true if the node collapses.
// If the node is already collapsed then Collapse returns false.
func (n *TreeToListNode) Collapse() bool {
	if n.parent == nil {
		panic("Collapse cannot be called for root nodes")
	}
	if !n.IsExpanded() || n.IsLeaf() {
		return false
	}
	n.parent.adjustDescendants(-n.descendants)
	n.descendants = 0
	n.children = nil
	if n.onChange != nil {
		n.onChange.Fire()
	}
	return true
}

// ToggleExpanded attempts to toggles the expanded state of the node, returning
// true if the state changed.
func (n *TreeToListNode) ToggleExpanded() bool {
	if n.IsExpanded() {
		return n.Collapse()
	} else {
		return n.Expand()
	}
}

// ExpandAll expands this node and all child nodes.
func (n *TreeToListNode) ExpandAll() {
	n.Expand()
	for _, c := range n.children {
		c.ExpandAll()
	}
}

// CollapseAll collapses this node and all child nodes.
func (n *TreeToListNode) CollapseAll() {
	n.Collapse()
	for _, c := range n.children {
		c.CollapseAll()
	}
}

// OnChange registers f to be called when the node is expanded, collapsed or has
// a change in the number of children.
func (n *TreeToListNode) OnChange(f func()) gxui.EventSubscription {
	if n.onChange == nil {
		n.onChange = gxui.CreateEvent(f)
	}
	return n.onChange.Listen(f)
}

// Descendants returns the total number of descendants of this node.
func (n *TreeToListNode) Descendants() int {
	return n.descendants
}

// Children returns all the immediate child nodes.
func (n *TreeToListNode) Children() []*TreeToListNode {
	return n.children
}

// Parent returns the parent of this node, or nil if this node is a root.
func (n *TreeToListNode) Parent() *TreeToListNode {
	p, _ := n.parent.(*TreeToListNode)
	return p
}

// NodeAt returns the idx'th TreeToListNode of all the expanded nodes under
// this node treated as a flattened list.
// Index 0 represents the first child of n, index 1 may represent the the second
// child of n or the first grandchild of n, and so on.
func (n *TreeToListNode) NodeAt(idx int) *TreeToListNode {
	for _, c := range n.children {
		switch {
		case idx == 0:
			return c
		case idx <= c.descendants:
			return c.NodeAt(idx - 1)
		default:
			idx -= c.descendants + 1
		}
	}
	panic("Index out of bounds")
}

// ItemAt returns the idx'th item of all the expanded nodes treated as a
// flattened list.
// Index 0 represents the first child of n, index 1 may represent the the second
// child of n or the first grandchild of n, and so on.
func (n *TreeToListNode) ItemAt(idx int) gxui.AdapterItem {
	return n.NodeAt(idx).item
}

// ItemIndex returns the index of item in the list of all the expanded nodes
// treated as a flattened list.
// Index 0 represents the first child of n, index 1 may represent the the second
// child of n or the first grandchild of n, and so on.
func (n *TreeToListNode) ItemIndex(item gxui.AdapterItem) int {
	c := n.DirectItemIndex(item)
	if c < 0 {
		return c
	}
	base := 0
	for i := 0; i < c; i++ {
		base += n.children[i].descendants + 1
	}
	if n.children[c].item == item {
		return base
	} else {
		return base + n.children[c].ItemIndex(item) + 1
	}
}

// DirectItemIndex returns the immediate child index that wraps or indirectly
// contains item. If no children contain item, the function returns -1.
func (n *TreeToListNode) DirectItemIndex(item gxui.AdapterItem) int {
	if !n.IsExpanded() {
		return -1
	}
	childIdx := n.container.ItemIndex(item)
	if childIdx < 0 {
		return -1 // Not found
	}
	if childIdx >= len(n.children) {
		panic(fmt.Errorf(
			"%T.ItemIndex(%v) returned out of bounds index %v. Acceptable range: [%d - %d]",
			n.container, item, childIdx, 0, len(n.children)-1))
	}
	return childIdx
}
