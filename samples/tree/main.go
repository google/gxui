// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"

	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/math"
	"github.com/google/gxui/themes/dark"
)

var data = flag.String("data", "", "path to data")

type treeAdapterNode struct {
	children []treeAdapterNode
	data     string
	id       gxui.AdapterItemId
}

func (n treeAdapterNode) Count() int {
	return len(n.children)
}

func (n treeAdapterNode) ItemId(index int) gxui.AdapterItemId {
	return n.children[index].id
}

func (n treeAdapterNode) ItemIndex(id gxui.AdapterItemId) int {
	for i, c := range n.children {
		if c.id == id {
			return i
		}
		if c.ItemIndex(id) >= 0 {
			return i
		}
	}
	return -1
}

func (n treeAdapterNode) Create(theme gxui.Theme, index int) gxui.Control {
	l := theme.CreateLabel()
	l.SetText(n.children[index].data)
	return l
}

func (n treeAdapterNode) CreateNode(index int) gxui.TreeAdapterNode {
	if len(n.children[index].children) > 0 {
		return n.children[index]
	} else {
		return nil // This child is a leaf.
	}
}

type treeAdapter struct {
	treeAdapterNode
	onDataChanged  gxui.Event
	onDataReplaced gxui.Event
}

func (a treeAdapter) ItemSize(theme gxui.Theme) math.Size {
	return math.Size{W: math.MaxSize.W, H: 20}
}

func (a treeAdapter) OnDataChanged(f func()) gxui.EventSubscription {
	if a.onDataChanged == nil {
		a.onDataChanged = gxui.CreateEvent(f)
	}
	return a.onDataChanged.Listen(f)
}

func (a treeAdapter) OnDataReplaced(f func()) gxui.EventSubscription {
	if a.onDataReplaced == nil {
		a.onDataReplaced = gxui.CreateEvent(f)
	}
	return a.onDataReplaced.Listen(f)
}

func appMain(driver gxui.Driver) {
	theme := dark.CreateTheme(driver)

	node := func(id gxui.AdapterItemId, data string, children ...treeAdapterNode) treeAdapterNode {
		return treeAdapterNode{
			children: children,
			data:     data,
			id:       id,
		}
	}

	layout := theme.CreateLinearLayout()
	layout.SetOrientation(gxui.Vertical)

	adapter := treeAdapter{}
	adapter.children = []treeAdapterNode{
		node(0x000, "Animals",
			node(0x100, "Mammals",
				node(0x110, "Cats"),
				node(0x120, "Dogs"),
				node(0x130, "Horses"),
				node(0x140, "Duck-billed platypuses"),
			),
			node(0x200, "Birds",
				node(0x210, "Peacocks"),
				node(0x220, "Doves"),
			),
			node(0x300, "Reptiles",
				node(0x310, "Lizards"),
				node(0x320, "Turtles"),
				node(0x330, "Crocodiles"),
				node(0x340, "Snakes"),
			),
			node(0x400, "Amphibians",
				node(0x410, "Frogs"),
				node(0x420, "Toads"),
			),
			node(0x500, "Arthropods",
				node(0x510, "Crustaceans",
					node(0x511, "Crabs"),
					node(0x512, "Lobsters"),
				),
				node(0x520, "Insects",
					node(0x521, "Ants"),
					node(0x522, "Bees"),
				),
				node(0x530, "Arachnids",
					node(0x531, "Spiders"),
					node(0x532, "Scorpions"),
				),
			),
		),
	}

	tree := theme.CreateTree()
	tree.SetAdapter(adapter)
	tree.Select(0x140) // Duck-billed platypuses
	tree.Show(tree.Selected())

	layout.AddChild(tree)

	row := theme.CreateLinearLayout()
	row.SetOrientation(gxui.Horizontal)
	layout.AddChild(row)

	expandAll := theme.CreateButton()
	expandAll.SetText("Expand All")
	expandAll.OnClick(func(gxui.MouseEvent) { tree.ExpandAll() })
	row.AddChild(expandAll)

	collapseAll := theme.CreateButton()
	collapseAll.SetText("Collapse All")
	collapseAll.OnClick(func(gxui.MouseEvent) { tree.CollapseAll() })
	row.AddChild(collapseAll)

	window := theme.CreateWindow(800, 600, "Tree view")
	window.AddChild(layout)
	window.OnClose(driver.Terminate)
	window.SetPadding(math.Spacing{L: 10, T: 10, R: 10, B: 10})
	gxui.EventLoop(driver)
}

func main() {
	flag.Parse()
	gl.StartDriver(*data, appMain)
}
