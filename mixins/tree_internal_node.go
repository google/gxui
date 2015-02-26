// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"fmt"
	"gaze/gxui"
	"gaze/gxui/assert"
)

type TreeInternalNode struct {
	adapterNode     gxui.TreeAdapterNode
	childCount      int
	children        []*TreeInternalNode
	isExpanded      bool
	id              gxui.AdapterItemId
	parent          *TreeInternalNode
	onExpandChanged gxui.Event
}

func (n *TreeInternalNode) findByIndex(absIdx int, d int) (parent *TreeInternalNode, relIndex int, depth int) {
	i := absIdx
	for j, c := range n.children {
		switch {
		case i == 0:
			return n, j, d
		case i <= c.childCount:
			return c.findByIndex(i-1, d+1)
		default:
			i -= c.childCount + 1
		}
	}
	panic(fmt.Errorf("Node does not contain index %d", absIdx))
	return nil, 0, 0
}

func (n *TreeInternalNode) findById(id gxui.AdapterItemId, baseIdx, d int) (parent *TreeInternalNode, relIdx, absIdx, depth int) {
	relIdx = n.adapterNode.ItemIndex(id)
	if relIdx < 0 {
		panic(fmt.Errorf("Node does not contain id %d", id))
	}
	assert.True(relIdx < len(n.children),
		"TreeAdapterNode.ItemIndex(%v) returned out of bounds index %v. Count = %v",
		id, relIdx, len(n.children))

	absIdx = baseIdx + 1 // +1 for n

	for i := 0; i < relIdx; i++ {
		absIdx += n.children[i].childCount + 1
	}

	if n.adapterNode.ItemId(relIdx) == id {
		return n, relIdx, absIdx, d
	}

	if n.children[relIdx].IsExpanded() {
		return n.children[relIdx].findById(id, absIdx, d+1)
	} else {
		return n, relIdx, absIdx, d
	}
}

func CreateTreeInternalRoot(adapterNode gxui.TreeAdapterNode) *TreeInternalNode {
	root := &TreeInternalNode{
		adapterNode: adapterNode,
		id:          gxui.InvalidAdapterItemId,
	}
	root.Expand()
	return root
}

func (n *TreeInternalNode) OnExpandedChanged(f func(bool)) gxui.EventSubscription {
	if n.onExpandChanged == nil {
		n.onExpandChanged = gxui.CreateEvent(func(bool) {})
	}
	return n.onExpandChanged.Listen(f)
}

func (n *TreeInternalNode) Id() gxui.AdapterItemId {
	return n.id
}

func (n *TreeInternalNode) IsExpanded() bool {
	return n.isExpanded
}

func (n *TreeInternalNode) IsLeaf() bool {
	return n.adapterNode == nil
}

func (n *TreeInternalNode) Expand() bool {
	if n.isExpanded || n.IsLeaf() {
		return false
	}
	n.isExpanded = true
	n.childCount = n.adapterNode.Count()
	assert.True(n.childCount >= 0, "%T.Count() returned a negative value %d", n.adapterNode, n.childCount)
	n.children = make([]*TreeInternalNode, n.childCount)
	for i := range n.children {
		n.children[i] = &TreeInternalNode{
			adapterNode: n.adapterNode.CreateNode(i),
			id:          n.adapterNode.ItemId(i),
			parent:      n,
		}
	}
	p := n.parent
	for p != nil {
		p.childCount += n.childCount
		p = p.parent
	}
	if n.onExpandChanged != nil {
		n.onExpandChanged.Fire(true)
	}
	return true
}

func (n *TreeInternalNode) Collapse() bool {
	if !n.isExpanded || n.IsLeaf() {
		return false
	}
	n.isExpanded = false
	p := n.parent
	for p != nil {
		p.childCount -= n.childCount
		p = p.parent
	}
	n.childCount = 0
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
	n.Collapse()
	for _, c := range n.children {
		c.CollapseAll()
	}
}

func (n *TreeInternalNode) Child(i int) (parent *TreeInternalNode) {
	return n.children[i]
}

func (n *TreeInternalNode) FindByIndex(idx int) (parent *TreeInternalNode, childIndex int, depth int) {
	return n.findByIndex(idx, 0)
}

func (n *TreeInternalNode) FindById(id gxui.AdapterItemId) (parent *TreeInternalNode, childIndex int, depth int) {
	c, relIdx, _, depth := n.findById(id, -1, 0)
	return c, relIdx, depth
}

func (n *TreeInternalNode) ItemId(idx int) gxui.AdapterItemId {
	p, i, _ := n.FindByIndex(idx)
	return p.adapterNode.ItemId(i)
}

func (n *TreeInternalNode) ItemIndex(id gxui.AdapterItemId) int {
	_, _, absIdx, _ := n.findById(id, -1, 0)
	return absIdx
}
