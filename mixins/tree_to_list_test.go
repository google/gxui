// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"testing"

	"github.com/google/gxui"
	"github.com/google/gxui/math"
)

type testTreeNode struct {
	item     gxui.AdapterItem
	children []*testTreeNode
}

func (n *testTreeNode) Count() int                           { return len(n.children) }
func (n *testTreeNode) NodeAt(index int) gxui.TreeNode       { return n.children[index] }
func (n *testTreeNode) Item() gxui.AdapterItem               { return n.item }
func (n *testTreeNode) Create(theme gxui.Theme) gxui.Control { return nil }

func (n *testTreeNode) ItemIndex(item gxui.AdapterItem) int {
	for i, c := range n.children {
		if item == c.item {
			return i
		}
		if idx := c.ItemIndex(item); idx >= 0 {
			return i
		}
	}
	return -1
}

type testTreeAdapter struct {
	gxui.AdapterBase
	testTreeNode
}

func (n *testTreeNode) Size(theme gxui.Theme) math.Size { return math.ZeroSize }

// n creates and returns a testTreeNode with the item i and children c.
func n(i gxui.AdapterItem, c ...*testTreeNode) *testTreeNode {
	return &testTreeNode{item: i, children: c}
}

// a creates and returns a list and tree adapters for the children c.
func a(c ...*testTreeNode) (list_adapter *TreeToListAdapter, tree_adapter *testTreeAdapter) {
	adapter := &testTreeAdapter{}
	adapter.children = c
	return CreateTreeToListAdapter(adapter, nil), adapter
}

func test(t *testing.T, name string, adapter *TreeToListAdapter, expected ...gxui.AdapterItem) {
	if len(expected) != adapter.Count() {
		t.Errorf("%s: Count was not as expected.\nExpected: %v\nGot:      %v",
			name, len(expected), adapter.Count())
	}
	for expected_index, expected_item := range expected {
		got_item := adapter.ItemAt(expected_index)
		got_index := adapter.ItemIndex(expected_item)
		if expected_item != got_item {
			t.Errorf("%s: Item at index %v was not as expected.\nExpected: %v\nGot:      %v",
				name, expected_index, expected_item, got_item)
		}
		if expected_index != got_index {
			t.Errorf("%s: Index of item %v was not as expected.\nExpected: %v\nGot:      %v",
				name, expected_item, expected_item, got_item)
		}
	}
}

func TestTreeToListNodeFlat(t *testing.T) {
	list_adapter, _ := a(n(10), n(20), n(30))
	test(t, "flat", list_adapter,
		gxui.AdapterItem(10),
		gxui.AdapterItem(20),
		gxui.AdapterItem(30),
	)
}

func TestTreeToListNodeDeep(t *testing.T) {

	list_adapter, tree_adapter := a(
		n(100,
			n(110),
			n(120,
				n(121),
				n(122),
				n(123)),
			n(130),
			n(140,
				n(141),
				n(142))))

	test(t, "unexpanded", list_adapter,
		gxui.AdapterItem(100),
	)

	list_adapter.node.children[0].Expand()
	test(t, "single expanded", list_adapter,
		gxui.AdapterItem(100), // (0) 100
		gxui.AdapterItem(110), // (1)  ╠══ 110
		gxui.AdapterItem(120), // (2)  ╠══ 120
		gxui.AdapterItem(130), // (3)  ╠══ 130
		gxui.AdapterItem(140), // (4)  ╚══ 140
	)

	list_adapter.ExpandAll()
	test(t, "fully expanded", list_adapter,
		gxui.AdapterItem(100), // (0) 100
		gxui.AdapterItem(110), // (1)  ╠══ 110
		gxui.AdapterItem(120), // (2)  ╠══ 120
		gxui.AdapterItem(121), // (3)  ║    ╠══ 121
		gxui.AdapterItem(122), // (4)  ║    ╠══ 122
		gxui.AdapterItem(123), // (5)  ║    ╚══ 123
		gxui.AdapterItem(130), // (6)  ╠══ 130
		gxui.AdapterItem(140), // (7)  ╚══ 140
		gxui.AdapterItem(141), // (8)       ╠══ 141
		gxui.AdapterItem(142), // (9)       ╚══ 142
	)

	list_adapter.node.NodeAt(2).Collapse()
	test(t, "one collapsed", list_adapter,
		gxui.AdapterItem(100), // (0) 100
		gxui.AdapterItem(110), // (1)  ╠══ 110
		gxui.AdapterItem(120), // (2)  ╠══ 120
		gxui.AdapterItem(130), // (3)  ╠══ 130
		gxui.AdapterItem(140), // (4)  ╚══ 140
		gxui.AdapterItem(141), // (5)       ╠══ 141
		gxui.AdapterItem(142), // (6)       ╚══ 142
	)

	tree_adapter.children[0].children = append(tree_adapter.children[0].children, n(150))
	test(t, "mutate, no data-changed", list_adapter,
		gxui.AdapterItem(100), // (0) 100
		gxui.AdapterItem(110), // (1)  ╠══ 110
		gxui.AdapterItem(120), // (2)  ╠══ 120
		gxui.AdapterItem(130), // (3)  ╠══ 130
		gxui.AdapterItem(140), // (4)  ╚══ 140
		gxui.AdapterItem(141), // (5)       ╠══ 141
		gxui.AdapterItem(142), // (6)       ╚══ 142
	)

	tree_adapter.DataChanged(false)
	test(t, "data-changed", list_adapter,
		gxui.AdapterItem(100), // (0) 100
		gxui.AdapterItem(110), // (1)  ╠══ 110
		gxui.AdapterItem(120), // (2)  ╠══ 120
		gxui.AdapterItem(130), // (3)  ╠══ 130
		gxui.AdapterItem(140), // (4)  ╠══ 140
		gxui.AdapterItem(141), // (5)  ║    ╠══ 141
		gxui.AdapterItem(142), // (6)  ║    ╚══ 142
		gxui.AdapterItem(150), // (7)  ╚══ 150
	)
}
