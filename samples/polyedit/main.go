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
	"github.com/google/gxui/themes/dark"
)

func vertexAt(p gxui.Polygon, at math.Point) int {
	for i, v := range p {
		if v.Position.Sub(at).Len() < 10 {
			return i
		}
	}
	return -1
}

func nearestEdge(p gxui.Polygon, to math.Point) int {
	best := -1
	bestDist := float32(1e10)
	for i := 0; i < len(p); i++ {
		j := (i + 1) % len(p)
		a := p[i].Position.Vec2()
		b := p[j].Position.Vec2()
		ab := b.Sub(a)
		// Find the nearest point on the line ab
		var nearest math.Point
		{
			dir := ab.Normalize()
			plane := dir.Vec3(-dir.Dot(a))
			dist := math.Saturate(plane.Dot(to.Vec3(1)) / ab.Len())
			nearest = a.Add(ab.MulS(dist)).Point()
		}
		dist := nearest.Sub(to).Len()
		if dist < bestDist {
			best = j
			bestDist = dist
		}
	}
	return best
}

func parsePolygon(s string) gxui.Polygon {
	parse := func(t string) bool {
		if len(s) >= len(t) && s[:len(t)] == t {
			s = s[len(t):]
			return true
		}
		return false
	}
	parseInt := func() int {
		i := 0
		for s[0] >= '0' && s[0] <= '9' {
			i *= 10
			i += int(s[0]) - '0'
			s = s[1:]
		}
		return i
	}
	p := gxui.Polygon{}
	parse("gxui.Polygon{")
	for parse("gxui.PolygonVertex") {
		parse("{")
		v := gxui.PolygonVertex{}
		parse("Position:math.Point{")
		parse("X:")
		v.Position.X = parseInt()
		parse(", Y:")
		v.Position.Y = parseInt()
		parse("}")
		parse(", ")
		parse("RoundedRadius:")
		v.RoundedRadius = float32(parseInt())
		parse("}")
		parse(", ")
		p = append(p, v)
	}
	return p
}

func appMain(driver gxui.Driver) {
	theme := dark.CreateTheme(driver)
	image := theme.CreateImage()
	p := gxui.Polygon{
		gxui.PolygonVertex{
			Position:      math.Point{X: 100, Y: 200},
			RoundedRadius: 0,
		},
		gxui.PolygonVertex{
			Position:      math.Point{X: 100, Y: 100},
			RoundedRadius: 0,
		},
		gxui.PolygonVertex{
			Position:      math.Point{X: 200, Y: 100},
			RoundedRadius: 0,
		},
	}
	update := func() {
		image.SetPolygon(p, gxui.CreatePen(3, gxui.White), gxui.CreateBrush(gxui.Gray50))
	}
	dragging := -1
	snap := func(p math.Point) math.Point {
		p.X = ((p.X + 5) / 10) * 10
		p.Y = ((p.Y + 5) / 10) * 10
		return p
	}
	image.OnMouseDown(func(ev gxui.MouseEvent) {
		switch {
		case ev.Modifier.Control() && ev.Button == gxui.MouseButtonLeft:
			i := nearestEdge(p, ev.Point)
			p = append(p, gxui.PolygonVertex{})
			copy(p[i+1:], p[i:])
			p[i] = gxui.PolygonVertex{
				Position:      snap(ev.Point),
				RoundedRadius: 0,
			}
		case ev.Modifier.Control() && ev.Button == gxui.MouseButtonRight:
			i := vertexAt(p, ev.Point)
			if i >= -1 {
				p = append(p[:i], p[i+1:]...)
			}
		}
		dragging = vertexAt(p, ev.Point)
	})
	image.OnMouseUp(func(ev gxui.MouseEvent) {
		dragging = -1
	})
	image.OnMouseMove(func(ev gxui.MouseEvent) {
		if dragging >= 0 {
			p[dragging].Position = snap(ev.Point)
			update()
		}
	})

	update()
	window := theme.CreateWindow(800, 600, "Polygon editor")
	window.SetScale(flags.DefaultScaleFactor)
	window.AddChild(image)
	window.OnClose(driver.Terminate)
	window.OnKeyDown(func(ev gxui.KeyboardEvent) {
		if ev.Modifier.Control() {
			switch ev.Key {
			case gxui.KeyC:
				driver.SetClipboard(fmt.Sprintf("%#v\n", p))
			case gxui.KeyV:
				s, err := driver.GetClipboard()
				if err == nil {
					p = parsePolygon(s)
					update()
				}
			}
		}
	})
}

func main() {
	gl.StartDriver(appMain)
}
