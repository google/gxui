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
	id       gxui.AdapterItemId
	children []*testTreeAdapterNode
}

func (n *testTreeAdapterNode) Count() int {
	return len(n.children)
}

func (n *testTreeAdapterNode) ItemId(index int) gxui.AdapterItemId {
	return n.children[index].id
}

func (n *testTreeAdapterNode) ItemIndex(id gxui.AdapterItemId) int {
	for i, c := range n.children {
		if id == c.id {
			return i
		}
		if idx := c.ItemIndex(id); idx >= 0 {
			return i
		}
	}
	return -1
}

func (n *testTreeAdapterNode) Create(theme gxui.Theme, index int) gxui.Control {
	return nil
}

func (n *testTreeAdapterNode) CreateNode(index int) gxui.TreeAdapterNode {
	return n.children[index]
}

func (n *testTreeAdapterNode) ItemSize(gxui.Theme) math.Size {
	return math.ZeroSize
}

func createTestTreeNode(id gxui.AdapterItemId, subnodes ...*testTreeAdapterNode) *testTreeAdapterNode {
	return &testTreeAdapterNode{id: id, children: subnodes}
}

func TestTINFlatSimple(t *testing.T) {
	N := createTestTreeNode
	root := CreateTreeInternalRoot(N(0, N(10), N(20), N(30)))

	test.AssertEquals(t, 3, root.childCount)
	test.AssertEquals(t, gxui.AdapterItemId(10), root.ItemId(0))
	test.AssertEquals(t, gxui.AdapterItemId(20), root.ItemId(1))
	test.AssertEquals(t, gxui.AdapterItemId(30), root.ItemId(2))
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

	test.AssertEquals(t, 3, root.childCount)
	test.AssertEquals(t, gxui.AdapterItemId(10), root.ItemId(0))
	test.AssertEquals(t, gxui.AdapterItemId(20), root.ItemId(1))
	test.AssertEquals(t, gxui.AdapterItemId(30), root.ItemId(2))
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
	test.AssertEquals(t, 10, root.childCount)

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

func TestTINFindById(t *testing.T) {
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
	test.AssertEquals(t, 10, root.childCount)

	n, i, d := root.FindById(100)
	test.AssertEquals(t, root, n)
	test.AssertEquals(t, 0, i)
	test.AssertEquals(t, 0, d)

	n, i, d = root.FindById(122)
	test.AssertEquals(t, root.Child(0).Child(1), n)
	test.AssertEquals(t, 1, i)
	test.AssertEquals(t, 2, d)

	n, i, d = root.FindById(141)
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
	test.AssertEquals(t, 7, root.childCount)
	test.AssertEquals(t, gxui.AdapterItemId(10), root.ItemId(0))
	test.AssertEquals(t, gxui.AdapterItemId(11), root.ItemId(1))
	test.AssertEquals(t, gxui.AdapterItemId(12), root.ItemId(2))
	test.AssertEquals(t, gxui.AdapterItemId(13), root.ItemId(3))
	test.AssertEquals(t, gxui.AdapterItemId(14), root.ItemId(4))
	test.AssertEquals(t, gxui.AdapterItemId(20), root.ItemId(5))

	root.Child(0).Collapse()
	test.AssertEquals(t, 3, root.childCount)
	test.AssertEquals(t, gxui.AdapterItemId(10), root.ItemId(0))
	test.AssertEquals(t, gxui.AdapterItemId(20), root.ItemId(1))
	test.AssertEquals(t, gxui.AdapterItemId(30), root.ItemId(2))
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
	test.AssertEquals(t, root.childCount, 7)
	test.AssertEquals(t, gxui.AdapterItemId(10), root.ItemId(0))
	test.AssertEquals(t, gxui.AdapterItemId(20), root.ItemId(1))
	test.AssertEquals(t, gxui.AdapterItemId(21), root.ItemId(2))
	test.AssertEquals(t, gxui.AdapterItemId(24), root.ItemId(5))
	test.AssertEquals(t, gxui.AdapterItemId(30), root.ItemId(6))

	root.Child(1).Collapse()
	test.AssertEquals(t, root.childCount, 3)
	test.AssertEquals(t, gxui.AdapterItemId(10), root.ItemId(0))
	test.AssertEquals(t, gxui.AdapterItemId(20), root.ItemId(1))
	test.AssertEquals(t, gxui.AdapterItemId(30), root.ItemId(2))
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
	test.AssertEquals(t, root.childCount, 10)
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
