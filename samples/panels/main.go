// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/samples/flags"
)

// Create a PanelHolder with a 3 panels
func panelHolder(name string, theme gxui.Theme) gxui.PanelHolder {
	label := func(text string) gxui.Label {
		label := theme.CreateLabel()
		label.SetText(text)
		return label
	}

	holder := theme.CreatePanelHolder()
	holder.AddPanel(label(name+" 0 content"), name+" 0 panel")
	holder.AddPanel(label(name+" 1 content"), name+" 1 panel")
	holder.AddPanel(label(name+" 2 content"), name+" 2 panel")
	return holder
}

func appMain(driver gxui.Driver) {
	theme := flags.CreateTheme(driver)

	// ┌───────┐║┌───────┐
	// │       │║│       │
	// │   A   │║│   B   │
	// │       │║│       │
	// └───────┘║└───────┘
	// ═══════════════════
	// ┌───────┐║┌───────┐
	// │       │║│       │
	// │   C   │║│   D   │
	// │       │║│       │
	// └───────┘║└───────┘

	splitterAB := theme.CreateSplitterLayout()
	splitterAB.SetOrientation(gxui.Horizontal)
	splitterAB.AddChild(panelHolder("A", theme))
	splitterAB.AddChild(panelHolder("B", theme))

	splitterCD := theme.CreateSplitterLayout()
	splitterCD.SetOrientation(gxui.Horizontal)
	splitterCD.AddChild(panelHolder("C", theme))
	splitterCD.AddChild(panelHolder("D", theme))

	vSplitter := theme.CreateSplitterLayout()
	vSplitter.SetOrientation(gxui.Vertical)
	vSplitter.AddChild(splitterAB)
	vSplitter.AddChild(splitterCD)

	window := theme.CreateWindow(800, 600, "Panels")
	window.SetScale(flags.DefaultScaleFactor)
	window.AddChild(vSplitter)
	window.OnClose(driver.Terminate)
}

func main() {
	gl.StartDriver(appMain)
}
