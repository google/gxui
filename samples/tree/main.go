// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/math"
	"github.com/google/gxui/samples/flags"
)

// item is used as an gxui.AdapterItem to identifiy each of the nodes.
// Each node's item must be equality-unique for the entire tree.
type item int

var nextUniqueItem item // the next item to used by node.add

// node is an implementation of gxui.TreeNode.
type node struct {
	name     string  // The name and item for this node.
	item     item    // The unique item for this node.
	changed  func()  // Called when a new item is added to this node.
	children []*node // The list of all child nodes.
}

// add appends a new child node to n with the specified name.
func (n *node) add(name string) *node {
	child := &node{
		name:    name,
		item:    nextUniqueItem,
		changed: n.changed,
	}
	nextUniqueItem++
	n.children = append(n.children, child)
	n.changed()
	return child
}

// Count implements gxui.TreeNodeContainer.
func (n *node) Count() int {
	return len(n.children)
}

// NodeAt implements gxui.TreeNodeContainer.
func (n *node) NodeAt(index int) gxui.TreeNode {
	return n.children[index]
}

// ItemIndex implements gxui.TreeNodeContainer.
func (n *node) ItemIndex(item gxui.AdapterItem) int {
	for i, c := range n.children {
		if c.item == item || c.ItemIndex(item) >= 0 {
			return i
		}
	}
	return -1
}

// Item implements gxui.TreeNode.
func (n *node) Item() gxui.AdapterItem {
	return n.item
}

// Create implements gxui.TreeNode.
func (n *node) Create(theme gxui.Theme) gxui.Control {
	layout := theme.CreateLinearLayout()
	layout.SetDirection(gxui.LeftToRight)

	label := theme.CreateLabel()
	label.SetText(n.name)

	textbox := theme.CreateTextBox()
	textbox.SetText(n.name)
	textbox.SetPadding(math.ZeroSpacing)
	textbox.SetMargin(math.ZeroSpacing)

	addButton := theme.CreateButton()
	addButton.SetText("+")
	addButton.OnClick(func(gxui.MouseEvent) { n.add("<new>") })

	edit := func() {
		layout.RemoveAll()
		layout.AddChild(textbox)
		layout.AddChild(addButton)
		gxui.SetFocus(textbox)
	}

	commit := func() {
		n.name = textbox.Text()
		label.SetText(n.name)
		layout.RemoveAll()
		layout.AddChild(label)
		layout.AddChild(addButton)
	}

	// When the user clicks the label, replace it with an editable text-box
	label.OnClick(func(gxui.MouseEvent) { edit() })

	// When the text-box loses focus, replace it with a label again.
	textbox.OnLostFocus(commit)

	layout.AddChild(label)
	layout.AddChild(addButton)
	return layout
}

// adapter is an implementation of gxui.TreeAdapter.
type adapter struct {
	gxui.AdapterBase
	node
}

// Size implements gxui.TreeAdapter.
func (a *adapter) Size(t gxui.Theme) math.Size {
	return math.Size{W: math.MaxSize.W, H: 18}
}

// addSpecies adds the list of species to animals.
// A map of name to item is returned.
func addSpecies(animals *node) map[string]item {
	items := make(map[string]item)

	add := func(to *node, name string) *node {
		n := to.add(name)
		items[name] = n.item
		return n
	}

	mammals := add(animals, "Mammals")
	add(mammals, "Cats")
	add(mammals, "Dogs")
	add(mammals, "Horses")
	add(mammals, "Duck-billed platypuses")

	birds := add(animals, "Birds")
	add(birds, "Peacocks")
	add(birds, "Doves")

	reptiles := add(animals, "Reptiles")
	add(reptiles, "Lizards")
	add(reptiles, "Turtles")
	add(reptiles, "Crocodiles")
	add(reptiles, "Snakes")

	arthropods := add(animals, "Arthropods")

	crustaceans := add(arthropods, "Crustaceans")
	add(crustaceans, "Crabs")
	add(crustaceans, "Lobsters")

	insects := add(arthropods, "Insects")
	add(insects, "Ants")
	add(insects, "Bees")

	arachnids := add(arthropods, "Arachnids")
	add(arachnids, "Spiders")
	add(arachnids, "Scorpions")

	return items
}

func appMain(driver gxui.Driver) {
	theme := flags.CreateTheme(driver)

	layout := theme.CreateLinearLayout()
	layout.SetDirection(gxui.TopToBottom)

	adapter := &adapter{}

	// hook up node changed function to the adapter OnDataChanged event.
	adapter.changed = func() { adapter.DataChanged(false) }

	// add all the species to the 'Animals' root node.
	items := addSpecies(adapter.add("Animals"))

	tree := theme.CreateTree()
	tree.SetAdapter(adapter)
	tree.Select(items["Doves"])
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
	window.SetScale(flags.DefaultScaleFactor)
	window.AddChild(layout)
	window.OnClose(driver.Terminate)
	window.SetPadding(math.Spacing{L: 10, T: 10, R: 10, B: 10})
}

func main() {
	gl.StartDriver(appMain)
}
