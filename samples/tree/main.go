// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/math"
	"github.com/google/gxui/themes/dark"
)

type treeAdapterNode struct {
	children []treeAdapterNode
	item     string
}

func (n treeAdapterNode) Count() int {
	return len(n.children)
}

func (n treeAdapterNode) ItemAt(index int) gxui.AdapterItem {
	return n.children[index].item
}

func (n treeAdapterNode) ItemIndex(item gxui.AdapterItem) int {
	for i, c := range n.children {
		if c.item == item {
			return i
		}
		if c.ItemIndex(item) >= 0 {
			return i
		}
	}
	return -1
}

func (n treeAdapterNode) Create(theme gxui.Theme, index int) gxui.Control {
	l := theme.CreateLabel()
	l.SetText(n.children[index].item)
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

func (a treeAdapter) Size(theme gxui.Theme) math.Size {
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

	node := func(item string, children ...treeAdapterNode) treeAdapterNode {
		return treeAdapterNode{
			children: children,
			item:     item,
		}
	}

	layout := theme.CreateLinearLayout()
	layout.SetDirection(gxui.TopToBottom)

	adapter := treeAdapter{}
	adapter.children = []treeAdapterNode{
		node("Animals",
			node("Mammals",
				node("Cats"),
				node("Dogs"),
				node("Horses"),
				node("Duck-billed platypuses"),
			),
			node("Birds",
				node("Peacocks"),
				node("Doves"),
			),
			node("Reptiles",
				node("Lizards"),
				node("Turtles"),
				node("Crocodiles"),
				node("Snakes"),
			),
			node("Amphibians",
				node("Frogs"),
				node("Toads"),
			),
			node("Arthropods",
				node("Crustaceans",
					node("Crabs"),
					node("Lobsters"),
				),
				node("Insects",
					node("Ants"),
					node("Bees"),
				),
				node("Arachnids",
					node("Spiders"),
					node("Scorpions"),
				),
			),
		),
	}

	tree := theme.CreateTree()
	tree.SetAdapter(adapter)
	tree.Select("Duck-billed platypuses")
	tree.Show(tree.Selected())

	layout.AddChild(tree)

	row := theme.CreateLinearLayout()
	row.SetDirection(gxui.LeftToRight)
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
	gl.StartDriver(appMain)
}
