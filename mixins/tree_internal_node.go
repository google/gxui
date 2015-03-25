// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"fmt"
	"github.com/google/gxui"
)

type TreeInternalNode struct {
	item            gxui.AdapterItem    // the item this wraps as AdapterItem
	node            gxui.TreeNode       // the item this wraps as TreeNode
	descendants     int                 // total number of descendants visible
	children        []*TreeInternalNode // if expanded the child nodes, or nil if collapsed
	parent          *TreeInternalNode
	onExpandChanged gxui.Event
}

func (n *TreeInternalNode) findByIndex(absIdx int, d int) (parent *TreeInternalNode, relIndex int, depth int) {
	i := absIdx
	for j, c := range n.children {
		switch {
		case i == 0:
			return n, j, d
		case i <= c.descendants:
			return c.findByIndex(i-1, d+1)
		default:
			i -= c.descendants + 1
		}
	}
	panic(fmt.Errorf("Tree node does not contain index %d", absIdx))
}

func (n *TreeInternalNode) findByItem(item gxui.AdapterItem, baseIdx, depth int) (parent *TreeInternalNode, relIdx, absIdx, d int) {
	relIdx = n.node.ItemIndex(item)
	if relIdx < 0 {
		return nil, -1, -1, -1 // Not found
	}
	if relIdx >= len(n.children) {
		panic(fmt.Errorf(
			"%T.ItemIndex(%v) returned out of bounds index %v. Acceptable range: [%d - %d]",
			n.node, item, relIdx, 0, len(n.children)-1))
	}

	absIdx = baseIdx + 1 // +1 for n
	for i := 0; i < relIdx; i++ {
		absIdx += n.children[i].descendants + 1
	}

	if child := n.children[relIdx]; child.item == item {
		return n, relIdx, absIdx, depth
	} else {
		if child.IsExpanded() {
			return child.findByItem(item, absIdx, depth+1)
		} else {
			return n, relIdx, absIdx, depth
		}
	}
}

func CreateTreeInternalRoot(node gxui.TreeAdapter) *TreeInternalNode {
	root := &TreeInternalNode{node: node}
	root.Expand()
	return root
}

func (n *TreeInternalNode) OnExpandedChanged(f func(bool)) gxui.EventSubscription {
	if n.onExpandChanged == nil {
		n.onExpandChanged = gxui.CreateEvent(func(bool) {})
	}
	return n.onExpandChanged.Listen(f)
}

func (n *TreeInternalNode) Item() gxui.AdapterItem {
	return n.item
}

func (n *TreeInternalNode) IsExpanded() bool {
	return n.children != nil
}

func (n *TreeInternalNode) IsLeaf() bool {
	return n.node == nil || n.node.Count() == 0
}

func (n *TreeInternalNode) Expand() bool {
	if n.IsExpanded() || n.IsLeaf() {
		return false
	}
	n.descendants = n.node.Count()
	if n.descendants < 0 {
		panic(fmt.Errorf("%T.Count() returned a negative value %d", n.node, n.descendants))
	}
	n.children = make([]*TreeInternalNode, n.descendants)
	for i := range n.children {
		item := n.node.ItemAt(i)
		node := n.node.NodeAt(i)
		n.children[i] = &TreeInternalNode{node: node, item: item, parent: n}
	}
	p := n.parent
	for p != nil {
		p.descendants += n.descendants
		p = p.parent
	}
	if n.onExpandChanged != nil {
		n.onExpandChanged.Fire(true)
	}
	return true
}

func (n *TreeInternalNode) Collapse() bool {
	if !n.IsExpanded() || n.IsLeaf() || n.parent == nil {
		return false
	}
	p := n.parent
	for p != nil {
		p.descendants -= n.descendants
		p = p.parent
	}
	n.descendants = 0
	n.children = nil
	if n.onExpandChanged != nil {
		n.onExpandChanged.Fire(false)
	}
	return true
}

func (n *TreeInternalNode) ExpandAll() {
	n.Expand()
	for _, c := range n.children {
		c.ExpandAll()
	}
}

func (n *TreeInternalNode) CollapseAll() {
	if !n.Collapse() { // The root cannot be collapsed
		for _, c := range n.children {
			c.CollapseAll()
		}
	}
}

func (n *TreeInternalNode) Child(i int) (parent *TreeInternalNode) {
	return n.children[i]
}

func (n *TreeInternalNode) FindByIndex(idx int) (parent *TreeInternalNode, childIndex int, depth int) {
	return n.findByIndex(idx, 0)
}

func (n *TreeInternalNode) FindByItem(item gxui.AdapterItem) (parent *TreeInternalNode, childIndex int, depth int) {
	c, relIdx, _, depth := n.findByItem(item, -1, 0)
	return c, relIdx, depth
}

func (n *TreeInternalNode) ItemAt(idx int) gxui.AdapterItem {
	p, i, _ := n.FindByIndex(idx)
	return p.node.ItemAt(i)
}

func (n *TreeInternalNode) ItemIndex(item gxui.AdapterItem) int {
	_, _, absIdx, _ := n.findByItem(item, -1, 0)
	return absIdx
}
