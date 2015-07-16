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

func appMain(driver gxui.Driver) {
	theme := flags.CreateTheme(driver)
	layout := theme.CreateLinearLayout()
	layout.SetSizeMode(gxui.Fill)

	buttonState := map[gxui.Button]func() bool{}
	update := func() {
		for button, f := range buttonState {
			button.SetChecked(f())
		}
	}

	button := func(name string, action func(), isSelected func() bool) gxui.Button {
		b := theme.CreateButton()
		b.SetText(name)
		b.OnClick(func(gxui.MouseEvent) { action(); update() })
		layout.AddChild(b)
		buttonState[b] = isSelected
		return b
	}

	button("TopToBottom",
		func() { layout.SetDirection(gxui.TopToBottom) },
		func() bool { return layout.Direction().TopToBottom() },
	)
	button("LeftToRight",
		func() { layout.SetDirection(gxui.LeftToRight) },
		func() bool { return layout.Direction().LeftToRight() },
	)
	button("BottomToTop",
		func() { layout.SetDirection(gxui.BottomToTop) },
		func() bool { return layout.Direction().BottomToTop() },
	)
	button("RightToLeft",
		func() { layout.SetDirection(gxui.RightToLeft) },
		func() bool { return layout.Direction().RightToLeft() },
	)

	button("AlignLeft",
		func() { layout.SetHorizontalAlignment(gxui.AlignLeft) },
		func() bool { return layout.HorizontalAlignment().AlignLeft() },
	)
	button("AlignCenter",
		func() { layout.SetHorizontalAlignment(gxui.AlignCenter) },
		func() bool { return layout.HorizontalAlignment().AlignCenter() },
	)
	button("AlignRight",
		func() { layout.SetHorizontalAlignment(gxui.AlignRight) },
		func() bool { return layout.HorizontalAlignment().AlignRight() },
	)

	button("AlignTop",
		func() { layout.SetVerticalAlignment(gxui.AlignTop) },
		func() bool { return layout.VerticalAlignment().AlignTop() },
	)
	button("AlignMiddle",
		func() { layout.SetVerticalAlignment(gxui.AlignMiddle) },
		func() bool { return layout.VerticalAlignment().AlignMiddle() },
	)
	button("AlignBottom",
		func() { layout.SetVerticalAlignment(gxui.AlignBottom) },
		func() bool { return layout.VerticalAlignment().AlignBottom() },
	)

	update()

	window := theme.CreateWindow(800, 600, "Linear layout")
	window.SetScale(flags.DefaultScaleFactor)
	window.AddChild(layout)
	window.OnClose(driver.Terminate)
	window.SetPadding(math.Spacing{L: 10, T: 10, R: 10, B: 10})
}

func main() {
	gl.StartDriver(appMain)
}
