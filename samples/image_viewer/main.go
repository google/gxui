// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/samples/flags"
)

func appMain(driver gxui.Driver) {
	args := flag.Args()
	if len(args) != 1 {
		fmt.Print("usage: image_viewer image-path\n")
		os.Exit(1)
	}

	file := args[0]
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("Failed to open image '%s': %v\n", file, err)
		os.Exit(1)
	}

	source, _, err := image.Decode(f)
	if err != nil {
		fmt.Printf("Failed to read image '%s': %v\n", file, err)
		os.Exit(1)
	}

	theme := flags.CreateTheme(driver)
	img := theme.CreateImage()

	mx := source.Bounds().Max
	window := theme.CreateWindow(mx.X, mx.Y, "Image viewer")
	window.SetScale(flags.DefaultScaleFactor)
	window.AddChild(img)

	// Copy the image to a RGBA format before handing to a gxui.Texture
	rgba := image.NewRGBA(source.Bounds())
	draw.Draw(rgba, source.Bounds(), source, image.ZP, draw.Src)
	texture := driver.CreateTexture(rgba, 1)
	img.SetTexture(texture)

	window.OnClose(driver.Terminate)
}

func main() {
	flag.Parse()
	gl.StartDriver(appMain)
}
