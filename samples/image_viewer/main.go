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
	"image"
	"image/color"
	"io/ioutil"
	gomath "math"
)

var data = flag.String("data", "data", "path to data")
var file = flag.String("file", "file", "path to file")
var width = flag.Int("width", 0, "width of the image")
var height = flag.Int("height", 0, "height of the image")
var imageType = flag.String("type", "rgba", "The type of the image (rgba or depth)")

func appMain(driver gxui.Driver) {
	theme := dark.CreateTheme(driver)
	img := theme.CreateImage()
	window := theme.CreateWindow(800, 600, "Image viewer")
	window.AddChild(img)

	raw, _ := ioutil.ReadFile(*file)
	bmp := image.NewRGBA(image.Rect(0, 0, *width, *height))
	if *imageType == "rgba" {
		bmp.Pix = raw
	} else if *imageType == "depth" {
		depthToImage(bmp, *width, *height, raw)
	}

	texture := driver.CreateTexture(bmp, 1)
	img.SetTexture(texture)

	window.OnClose(driver.Terminate)
	gxui.EventLoop(driver)
}

func depthToImage(img *image.RGBA, w int, h int, buffer []byte) {
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			bits := (uint32(buffer[3]) << 24) | (uint32(buffer[2]) << 16) | (uint32(buffer[1]) << 8) | (uint32(buffer[0]) << 0)
			depth := gomath.Float32frombits(bits)
			buffer = buffer[4:]

			d := 0.01 / (1.0 - depth)
			c := color.RGBA{
				R: byte(math.Cosf(d+math.TwoPi*0.000)*127.0 + 128.0),
				G: byte(math.Cosf(d+math.TwoPi*0.333)*127.0 + 128.0),
				B: byte(math.Cosf(d+math.TwoPi*0.666)*127.0 + 128.0),
				A: byte(0xFF),
			}
			img.Set(x, y, c)
		}
	}
}

func main() {
	flag.Parse()
	gl.StartDriver(*data, appMain)
}
