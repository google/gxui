// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/themes/dark"
)

var data = flag.String("data", "", "path to data")
var file = flag.String("file", "", "path to file")
var width = flag.Int("width", 0, "width of the image")
var height = flag.Int("height", 0, "height of the image")
var imageType = flag.String("type", "rgba", "The type of the image (rgba or depth)")

func appMain(driver gxui.Driver) {
	f, err := os.Open(*file)
	if err != nil {
		panic(err)
	}
	bmp, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	raw := img2rgba(bmp)

	theme := dark.CreateTheme(driver)
	img := theme.CreateImage()

	mx := raw.Bounds().Max

	window := theme.CreateWindow(mx.X, mx.Y, "Image viewer")
	window.AddChild(img)

	texture := driver.CreateTexture(raw, 1)
	img.SetTexture(texture)

	window.OnClose(driver.Terminate)
	gxui.EventLoop(driver)
}

func img2rgba(bmp image.Image) *image.RGBA {
	mx := bmp.Bounds().Max
	raw := image.NewRGBA(bmp.Bounds())

	for y := 0; y < mx.Y; y++ {
		for x := 0; x < mx.X; x++ {
			raw.Set(x, mx.Y-y-1, bmp.At(x, y))
		}
	}

	return raw
}

func main() {
	flag.Parse()
	gl.StartDriver(*data, appMain)
}
