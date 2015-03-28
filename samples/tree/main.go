// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/math"
	"github.com/google/gxui/samples/flags"
	"github.com/google/gxui/themes/dark"
)

type node struct {
	name     string
	children []node
}

func (n node) Count() int {
	return len(n.children)
}

func (n node) NodeAt(index int) gxui.TreeNode {
	return n.children[index]
}

func (n node) ItemAt(index int) gxui.AdapterItem {
	return n.children[index].name
}

func (n node) ItemIndex(item gxui.AdapterItem) int {
	name := item.(string)
	for i, c := range n.children {
		if c.name == name || c.ItemIndex(item) >= 0 {
			return i
		}
	}
	return -1
}

func (n node) Create(theme gxui.Theme, index int) gxui.Control {
	label := theme.CreateLabel()
	label.SetText(n.children[index].name)
	return label
}

type adapter struct {
	gxui.AdapterBase
	node
}

func (a *adapter) Size(t gxui.Theme) math.Size {
	return math.Size{W: math.MaxSize.W, H: 18}
}

func appMain(driver gxui.Driver) {
	theme := dark.CreateTheme(driver)

	layout := theme.CreateLinearLayout()
	layout.SetDirection(gxui.TopToBottom)

	animals := &adapter{
		node: node{
			name: "Animals",
			children: []node{
				node{
					name: "Mammals",
					children: []node{
						node{name: "Cats"},
						node{name: "Dogs"},
						node{name: "Horses"},
						node{name: "Duck-billed platypuses"},
					},
				},
				node{
					name: "Birds",
					children: []node{
						node{name: "Peacocks"},
						node{name: "Doves"},
					},
				},
				node{
					name: "Reptiles",
					children: []node{
						node{name: "Lizards"},
						node{name: "Turtles"},
						node{name: "Crocodiles"},
						node{name: "Snakes"},
					},
				},
				node{
					name: "Amphibians",
					children: []node{
						node{name: "Frogs"},
						node{name: "Toads"},
					},
				},
				node{
					name: "Arthropods",
					children: []node{
						node{
							name: "Crustaceans",
							children: []node{
								node{name: "Crabs"},
								node{name: "Lobsters"},
							},
						},
						node{
							name: "Insects",
							children: []node{
								node{name: "Ants"},
								node{name: "Bees"},
							},
						},
						node{
							name: "Arachnids",
							children: []node{
								node{name: "Spiders"},
								node{name: "Scorpions"},
							},
						},
					},
				},
			},
		},
	}

	tree := theme.CreateTree()
	tree.SetAdapter(animals)
	tree.Select("Doves")
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
