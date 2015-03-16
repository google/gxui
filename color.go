// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import "github.com/google/gxui/math"

var Transparent = Color{0.0, 0.0, 0.0, 0.0}

var Black = Color{0.0, 0.0, 0.0, 1.0}

var Red10 = Color{0.1, 0.0, 0.0, 1.0}
var Red20 = Color{0.2, 0.0, 0.0, 1.0}
var Red30 = Color{0.3, 0.0, 0.0, 1.0}
var Red40 = Color{0.4, 0.0, 0.0, 1.0}
var Red50 = Color{0.5, 0.0, 0.0, 1.0}
var Red60 = Color{0.6, 0.0, 0.0, 1.0}
var Red70 = Color{0.7, 0.0, 0.0, 1.0}
var Red80 = Color{0.8, 0.0, 0.0, 1.0}
var Red90 = Color{0.9, 0.0, 0.0, 1.0}
var Red = Color{1.0, 0.0, 0.0, 1.0}

var Green10 = Color{0.0, 0.1, 0.0, 1.0}
var Green20 = Color{0.0, 0.2, 0.0, 1.0}
var Green30 = Color{0.0, 0.3, 0.0, 1.0}
var Green40 = Color{0.0, 0.4, 0.0, 1.0}
var Green50 = Color{0.0, 0.5, 0.0, 1.0}
var Green60 = Color{0.0, 0.6, 0.0, 1.0}
var Green70 = Color{0.0, 0.7, 0.0, 1.0}
var Green80 = Color{0.0, 0.8, 0.0, 1.0}
var Green90 = Color{0.0, 0.9, 0.0, 1.0}
var Green = Color{0.0, 1.0, 0.0, 1.0}

var Blue10 = Color{0.0, 0.0, 0.1, 1.0}
var Blue20 = Color{0.0, 0.0, 0.2, 1.0}
var Blue30 = Color{0.0, 0.0, 0.3, 1.0}
var Blue40 = Color{0.0, 0.0, 0.4, 1.0}
var Blue50 = Color{0.0, 0.0, 0.5, 1.0}
var Blue60 = Color{0.0, 0.0, 0.6, 1.0}
var Blue70 = Color{0.0, 0.0, 0.7, 1.0}
var Blue80 = Color{0.0, 0.0, 0.8, 1.0}
var Blue90 = Color{0.0, 0.0, 0.9, 1.0}
var Blue = Color{0.0, 0.0, 1.0, 1.0}

var Gray10 = Color{0.1, 0.1, 0.1, 1.0}
var Gray15 = Color{0.15, 0.15, 0.15, 1.0}
var Gray20 = Color{0.2, 0.2, 0.2, 1.0}
var Gray30 = Color{0.3, 0.3, 0.3, 1.0}
var Gray40 = Color{0.4, 0.4, 0.4, 1.0}
var Gray50 = Color{0.5, 0.5, 0.5, 1.0}
var Gray60 = Color{0.6, 0.6, 0.6, 1.0}
var Gray70 = Color{0.7, 0.7, 0.7, 1.0}
var Gray80 = Color{0.8, 0.8, 0.8, 1.0}
var Gray90 = Color{0.9, 0.9, 0.9, 1.0}
var White = Color{1.0, 1.0, 1.0, 1.0}

var Yellow = Color{1.0, 1.0, 0.0, 1.0}

type Color struct {
	R, G, B, A float32
}

func ColorFromHex(hex uint32) Color {
	return Color{
		A: float32((hex>>24)&0xFF) / 255.0,
		R: float32((hex>>16)&0xFF) / 255.0,
		G: float32((hex>>8)&0xFF) / 255.0,
		B: float32(hex&0xFF) / 255.0,
	}
}

func (c Color) MulRGB(s float32) Color {
	return Color{c.R * s, c.G * s, c.B * s, c.A}
}

func (c Color) Saturate() Color {
	return Color{math.Saturate(c.R), math.Saturate(c.G), math.Saturate(c.B), math.Saturate(c.A)}
}
