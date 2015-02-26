// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"gaze/gxui"
	"gaze/gxui/drivers/gl"
	"gaze/gxui/math"
	"gaze/gxui/themes/dark"
	"time"
)

var data = flag.String("data", "data", "path to data")

func appMain(driver gxui.Driver) {
	theme := dark.CreateTheme(driver)

	label := theme.CreateLabel()
	label.SetText("This is a progress bar:")

	progressBar := theme.CreateProgressBar()
	progressBar.SetDesiredSize(math.Size{W: 400, H: 20})
	progressBar.SetTarget(100)

	layout := theme.CreateLinearLayout()
	layout.AddChild(label)
	layout.AddChild(progressBar)
	layout.SetHorizontalAlignment(gxui.AlignCenter)

	window := theme.CreateWindow(800, 600, "Progress bar")
	window.AddChild(layout)
	window.OnClose(driver.Terminate)

	progress := 0
	var timer *time.Timer
	timer = time.AfterFunc(time.Millisecond*500, func() {
		driver.Events() <- func() {
			progress = (progress + 3) % progressBar.Target()
			progressBar.SetProgress(progress)
			timer.Reset(time.Millisecond * 500)
		}
	})

	gxui.EventLoop(driver)
}

func main() {
	flag.Parse()
	gl.StartDriver(*data, appMain)
}
