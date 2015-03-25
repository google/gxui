// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import test "github.com/google/gxui/testing"
import (
	"github.com/google/gxui"
	"github.com/google/gxui/math"
	"testing"
)

type testTreeAdapterNode struct {
	gxui.AdapterBase
	item     gxui.AdapterItem
	children []*testTreeAdapterNode
}

func (n *testTreeAdapterNode) Count() int {
	return len(n.children)
}

func (n *testTreeAdapterNode) NodeAt(index int) gxui.TreeNode {
	return n.children[index]
}

func (n *testTreeAdapterNode) ItemAt(index int) gxui.AdapterItem {
	return n.children[index].item
}

func (n *testTreeAdapterNode) ItemIndex(item gxui.AdapterItem) int {
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

func (n *testTreeAdapterNode) Create(theme gxui.Theme, index int) gxui.Control {
	return nil
}

func (n *testTreeAdapterNode) Size(gxui.Theme) math.Size {
	return math.ZeroSize
}

func createTestTreeNode(item gxui.AdapterItem, subnodes ...*testTreeAdapterNode) *testTreeAdapterNode {
	return &testTreeAdapterNode{item: item, children: subnodes}
}

func TestTINFlatSimple(t *testing.T) {
	N := createTestTreeNode
	root := CreateTreeInternalRoot(N(0, N(10), N(20), N(30)))

	test.AssertEquals(t, 3, root.descendants)
	test.AssertEquals(t, gxui.AdapterItem(10), root.ItemAt(0))
	test.AssertEquals(t, gxui.AdapterItem(20), root.ItemAt(1))
	test.AssertEquals(t, gxui.AdapterItem(30), root.ItemAt(2))
}

func TestTINDeepNoExpand(t *testing.T) {
	N := createTestTreeNode
	root := CreateTreeInternalRoot(N(0,
		N(10,
			N(11), N(12), N(13), N(14)),
		N(20,
			N(21), N(22), N(23), N(24)),
		N(30,
			N(31), N(32), N(33), N(34)),
	))

	test.AssertEquals(t, 3, root.descendants)
	test.AssertEquals(t, gxui.AdapterItem(10), root.ItemAt(0))
	test.AssertEquals(t, gxui.AdapterItem(20), root.ItemAt(1))
	test.AssertEquals(t, gxui.AdapterItem(30), root.ItemAt(2))
}

func TestTINFindByIndex(t *testing.T) {
	N := createTestTreeNode
	root := CreateTreeInternalRoot(N(0,
		/*0*/ N(100,
			/*1*/ N(110),
			/*2*/ N(120,
				/*3*/ N(121),
				/*4*/ N(122),
				/*5*/ N(123)),
			/*6*/ N(130),
			/*7*/ N(140,
				/*8*/ N(141),
				/*9*/ N(141)))))

	root.ExpandAll()
	test.AssertEquals(t, 10, root.descendants)

	n, i, d := root.FindByIndex(0)
	test.AssertEquals(t, root, n)
	test.AssertEquals(t, 0, i)
	test.AssertEquals(t, 0, d)

	n, i, d = root.FindByIndex(4)
	test.AssertEquals(t, root.Child(0).Child(1), n)
	test.AssertEquals(t, 1, i)
	test.AssertEquals(t, 2, d)

	n, i, d = root.FindByIndex(9)
	test.AssertEquals(t, root.Child(0).Child(3), n)
	test.AssertEquals(t, 1, i)
	test.AssertEquals(t, 2, d)
}

func TestTINFindByItem(t *testing.T) {
	N := createTestTreeNode
	root := CreateTreeInternalRoot(N(0,
		/*0*/ N(100,
			/*1*/ N(110),
			/*2*/ N(120,
				/*3*/ N(121),
				/*4*/ N(122),
				/*5*/ N(123)),
			/*6*/ N(130),
			/*7*/ N(140,
				/*8*/ N(141),
				/*9*/ N(141)))))

	root.ExpandAll()
	test.AssertEquals(t, 10, root.descendants)

	n, i, d := root.FindByItem(100)
	test.AssertEquals(t, root, n)
	test.AssertEquals(t, 0, i)
	test.AssertEquals(t, 0, d)

	n, i, d = root.FindByItem(122)
	test.AssertEquals(t, root.Child(0).Child(1), n)
	test.AssertEquals(t, 1, i)
	test.AssertEquals(t, 2, d)

	n, i, d = root.FindByItem(141)
	test.AssertEquals(t, root.Child(0).Child(3), n)
	test.AssertEquals(t, 0, i)
	test.AssertEquals(t, 2, d)
}

func TestTINExpandExpandFirst(t *testing.T) {
	N := createTestTreeNode
	root := CreateTreeInternalRoot(N(0,
		N(10,
			N(11), N(12), N(13), N(14)),
		N(20,
			N(21), N(22), N(23), N(24)),
		N(30,
			N(31), N(32), N(33), N(34)),
	))

	root.Child(0).Expand()
	test.AssertEquals(t, 7, root.descendants)
	test.AssertEquals(t, gxui.AdapterItem(10), root.ItemAt(0))
	test.AssertEquals(t, gxui.AdapterItem(11), root.ItemAt(1))
	test.AssertEquals(t, gxui.AdapterItem(12), root.ItemAt(2))
	test.AssertEquals(t, gxui.AdapterItem(13), root.ItemAt(3))
	test.AssertEquals(t, gxui.AdapterItem(14), root.ItemAt(4))
	test.AssertEquals(t, gxui.AdapterItem(20), root.ItemAt(5))

	root.Child(0).Collapse()
	test.AssertEquals(t, 3, root.descendants)
	test.AssertEquals(t, gxui.AdapterItem(10), root.ItemAt(0))
	test.AssertEquals(t, gxui.AdapterItem(20), root.ItemAt(1))
	test.AssertEquals(t, gxui.AdapterItem(30), root.ItemAt(2))
}

func TestTINExpandCollapseOne(t *testing.T) {
	N := createTestTreeNode
	root := CreateTreeInternalRoot(N(0,
		N(10,
			N(11), N(12), N(13), N(14)),
		N(20,
			N(21), N(22), N(23), N(24)),
		N(30,
			N(31), N(32), N(33), N(34)),
	))

	root.Child(1).Expand()
	test.AssertEquals(t, root.descendants, 7)
	test.AssertEquals(t, gxui.AdapterItem(10), root.ItemAt(0))
	test.AssertEquals(t, gxui.AdapterItem(20), root.ItemAt(1))
	test.AssertEquals(t, gxui.AdapterItem(21), root.ItemAt(2))
	test.AssertEquals(t, gxui.AdapterItem(24), root.ItemAt(5))
	test.AssertEquals(t, gxui.AdapterItem(30), root.ItemAt(6))

	root.Child(1).Collapse()
	test.AssertEquals(t, root.descendants, 3)
	test.AssertEquals(t, gxui.AdapterItem(10), root.ItemAt(0))
	test.AssertEquals(t, gxui.AdapterItem(20), root.ItemAt(1))
	test.AssertEquals(t, gxui.AdapterItem(30), root.ItemAt(2))
}

func TestTINGExpandAll(t *testing.T) {
	N := createTestTreeNode
	root := CreateTreeInternalRoot(N(0,
		/*0*/ N(100,
			/*1*/ N(110),
			/*2*/ N(120,
				/*3*/ N(121),
				/*4*/ N(122),
				/*5*/ N(123)),
			/*6*/ N(130),
			/*7*/ N(140,
				/*8*/ N(141),
				/*9*/ N(141)))))

	root.ExpandAll()
	test.AssertEquals(t, root.descendants, 10)
}

func TestTINItemIndex(t *testing.T) {
	N := createTestTreeNode
	root := CreateTreeInternalRoot(N(0,
		/*0*/ N(100,
			/*1*/ N(110),
			/*2*/ N(120,
				/*3*/ N(121),
				/*4*/ N(122),
				/*5*/ N(123)),
			/*6*/ N(130),
			/*7*/ N(140,
				/*8*/ N(141),
				/*9*/ N(141)))))

	root.ExpandAll()
	test.AssertEquals(t, 0, root.ItemIndex(100))
	test.AssertEquals(t, 2, root.ItemIndex(120))
	test.AssertEquals(t, 7, root.ItemIndex(140))
}
