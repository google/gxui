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

func buildStar(theme gxui.Theme, center math.Point, radius, rotation float32, points int) gxui.Image {
	p := make(gxui.Polygon, points*2)
	for i := 0; i < points*2; i++ {
		frac := float32(i) / float32(points*2)
		α := frac*math.TwoPi + rotation
		r := []float32{radius, radius / 2}[i&1]
		p[i] = gxui.PolygonVertex{
			Position: math.Point{
				X: center.X + int(r*math.Cosf(α)),
				Y: center.Y + int(r*math.Sinf(α)),
			},
			RoundedRadius: []float32{0, 50}[i&1],
		}
	}
	image := theme.CreateImage()
	image.SetPolygon(p, gxui.CreatePen(3, gxui.Red), gxui.CreateBrush(gxui.Yellow))
	return image
}

func buildMoon(theme gxui.Theme, center math.Point, radius float32) gxui.Image {
	c := 40
	p := make(gxui.Polygon, c*2)
	for i := 0; i < c; i++ {
		frac := float32(i) / float32(c)
		α := math.Lerpf(math.Pi*1.2, math.Pi*-0.2, frac)
		p[i] = gxui.PolygonVertex{
			Position: math.Point{
				X: center.X + int(radius*math.Sinf(α)),
				Y: center.Y + int(radius*math.Cosf(α)),
			},
			RoundedRadius: 0,
		}
	}
	for i := 0; i < c; i++ {
		frac := float32(i) / float32(c)
		α := math.Lerpf(math.Pi*-0.2, math.Pi*1.2, frac)
		r := math.Lerpf(radius, radius*0.5, math.Sinf(frac*math.Pi))
		p[i+c] = gxui.PolygonVertex{
			Position: math.Point{
				X: center.X + int(r*math.Sinf(α)),
				Y: center.Y + int(r*math.Cosf(α)),
			},
			RoundedRadius: 0,
		}
	}
	image := theme.CreateImage()
	image.SetPolygon(p, gxui.CreatePen(3, gxui.Gray80), gxui.CreateBrush(gxui.Gray40))
	return image
}

func appMain(driver gxui.Driver) {
	theme := dark.CreateTheme(driver)
	image := theme.CreateImage()
	window := theme.CreateWindow(800, 600, "Polygon")
	window.SetScale(flags.DefaultScaleFactor)
	window.AddChild(image)

	window.AddChild(buildStar(theme, math.Point{X: 100, Y: 100}, 50, 0.2, 6))
	window.AddChild(buildStar(theme, math.Point{X: 650, Y: 170}, 70, 0.5, 7))
	window.AddChild(buildStar(theme, math.Point{X: 40, Y: 300}, 20, 0, 5))
	window.AddChild(buildStar(theme, math.Point{X: 410, Y: 320}, 25, 0.9, 5))
	window.AddChild(buildStar(theme, math.Point{X: 220, Y: 520}, 45, 0, 6))

	window.AddChild(buildMoon(theme, math.Point{X: 400, Y: 300}, 200))

	window.OnClose(driver.Terminate)
}

func main() {
	gl.StartDriver(appMain)
}
