// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/themes/dark"
)

var data = flag.String("data", "", "path to data")
var file = flag.String("file", "", "path to file")

func appMain(driver gxui.Driver) {
	f, err := os.Open(*file)
	if err != nil {
		panic(err)
	}
	source, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	theme := dark.CreateTheme(driver)
	img := theme.CreateImage()

	mx := source.Bounds().Max
	window := theme.CreateWindow(mx.X, mx.Y, "Image viewer")
	window.AddChild(img)

	// Copy the image to a RGBA format before handing to a gxui.Texture
	rgba := image.NewRGBA(source.Bounds())
	draw.Draw(rgba, source.Bounds(), source, image.ZP, draw.Src)
	texture := driver.CreateTexture(rgba, 1)
	texture.SetFlipY(true)
	img.SetTexture(texture)

	window.OnClose(driver.Terminate)
	gxui.EventLoop(driver)
}

func main() {
	flag.Parse()
	gl.StartDriver(appMain)
}
