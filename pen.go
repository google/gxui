// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

var DefaultPen Pen = CreatePen(1.0, Black)
var TransparentPen Pen = CreatePen(0.0, Transparent)
var WhitePen Pen = CreatePen(1.0, White)

type Pen struct {
	Width float32
	Color Color
}

func CreatePen(width float32, color Color) Pen {
	return Pen{width, color}
}
