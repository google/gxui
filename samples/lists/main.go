// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/math"
	"github.com/google/gxui/samples/flags"
)

// Number picker uses the gxui.DefaultAdapter for driving a list
func numberPicker(theme gxui.Theme, overlay gxui.BubbleOverlay) gxui.Control {
	items := []string{
		"zero", "one", "two", "three", "four", "five",
		"six", "seven", "eight", "nine", "ten",
		"eleven", "twelve", "thirteen", "fourteen", "fifteen",
		"sixteen", "seventeen", "eighteen", "nineteen", "twenty",
	}

	adapter := gxui.CreateDefaultAdapter()
	adapter.SetItems(items)

	layout := theme.CreateLinearLayout()
	layout.SetDirection(gxui.TopToBottom)

	label0 := theme.CreateLabel()
	label0.SetText("Numbers:")
	layout.AddChild(label0)

	dropList := theme.CreateDropDownList()
	dropList.SetAdapter(adapter)
	dropList.SetBubbleOverlay(overlay)
	layout.AddChild(dropList)

	list := theme.CreateList()
	list.SetAdapter(adapter)
	list.SetOrientation(gxui.Vertical)
	layout.AddChild(list)

	label1 := theme.CreateLabel()
	label1.SetMargin(math.Spacing{T: 30})
	label1.SetText("Selected number:")
	layout.AddChild(label1)

	selected := theme.CreateLabel()
	layout.AddChild(selected)

	dropList.OnSelectionChanged(func(item gxui.AdapterItem) {
		if list.Selected() != item {
			list.Select(item)
		}
	})

	list.OnSelectionChanged(func(item gxui.AdapterItem) {
		if dropList.Selected() != item {
			dropList.Select(item)
		}
		selected.SetText(fmt.Sprintf("%s - %d", item, adapter.ItemIndex(item)))
	})

	return layout
}

type customAdapter struct {
	gxui.AdapterBase
}

func (a *customAdapter) Count() int {
	return 1000
}

func (a *customAdapter) ItemAt(index int) gxui.AdapterItem {
	return index // This adapter uses integer indices as AdapterItems
}

func (a *customAdapter) ItemIndex(item gxui.AdapterItem) int {
	return item.(int) // Inverse of ItemAt()
}

func (a *customAdapter) Size(theme gxui.Theme) math.Size {
	return math.Size{W: 100, H: 100}
}

func (a *customAdapter) Create(theme gxui.Theme, index int) gxui.Control {
	phase := float32(index) / 1000
	c := gxui.Color{
		R: 0.5 + 0.5*math.Sinf(math.TwoPi*(phase+0.000)),
		G: 0.5 + 0.5*math.Sinf(math.TwoPi*(phase+0.333)),
		B: 0.5 + 0.5*math.Sinf(math.TwoPi*(phase+0.666)),
		A: 1.0,
	}
	i := theme.CreateImage()
	i.SetBackgroundBrush(gxui.CreateBrush(c))
	i.SetMargin(math.Spacing{L: 3, T: 3, R: 3, B: 3})
	i.OnMouseEnter(func(ev gxui.MouseEvent) {
		i.SetBorderPen(gxui.CreatePen(2, gxui.Gray80))
	})
	i.OnMouseExit(func(ev gxui.MouseEvent) {
		i.SetBorderPen(gxui.TransparentPen)
	})
	i.OnMouseDown(func(ev gxui.MouseEvent) {
		i.SetBackgroundBrush(gxui.CreateBrush(c.MulRGB(0.7)))
	})
	i.OnMouseUp(func(ev gxui.MouseEvent) {
		i.SetBackgroundBrush(gxui.CreateBrush(c))
	})
	return i
}

// Color picker uses the customAdapter for driving a list
func colorPicker(theme gxui.Theme) gxui.Control {
	layout := theme.CreateLinearLayout()
	layout.SetDirection(gxui.TopToBottom)

	label0 := theme.CreateLabel()
	label0.SetText("Color palette:")
	layout.AddChild(label0)

	adapter := &customAdapter{}

	list := theme.CreateList()
	list.SetAdapter(adapter)
	list.SetOrientation(gxui.Horizontal)
	layout.AddChild(list)

	label1 := theme.CreateLabel()
	label1.SetMargin(math.Spacing{T: 30})
	label1.SetText("Selected color:")
	layout.AddChild(label1)

	selected := theme.CreateImage()
	selected.SetExplicitSize(math.Size{W: 32, H: 32})
	layout.AddChild(selected)

	list.OnSelectionChanged(func(item gxui.AdapterItem) {
		if item != nil {
			control := list.ItemControl(item)
			selected.SetBackgroundBrush(control.(gxui.Image).BackgroundBrush())
		}
	})

	return layout
}

func appMain(driver gxui.Driver) {
	theme := flags.CreateTheme(driver)

	overlay := theme.CreateBubbleOverlay()

	holder := theme.CreatePanelHolder()
	holder.AddPanel(numberPicker(theme, overlay), "Default adapter")
	holder.AddPanel(colorPicker(theme), "Custom adapter")

	window := theme.CreateWindow(800, 600, "Lists")
	window.SetScale(flags.DefaultScaleFactor)
	window.AddChild(holder)
	window.AddChild(overlay)
	window.OnClose(driver.Terminate)
	window.SetPadding(math.Spacing{L: 10, T: 10, R: 10, B: 10})
}

func main() {
	gl.StartDriver(appMain)
}
