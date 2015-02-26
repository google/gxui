// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

var WhiteBrush = CreateBrush(White)
var TransparentBrush = CreateBrush(Transparent)
var BlackBrush = CreateBrush(Black)
var DefaultBrush = WhiteBrush

type Brush struct {
	Color Color
}

func CreateBrush(color Color) Brush {
	return Brush{color}
}
