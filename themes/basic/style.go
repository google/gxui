// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package basic

import (
	"github.com/google/gxui"
)

type Style struct {
	FontColor gxui.Color
	Brush     gxui.Brush
	Pen       gxui.Pen
}

func CreateStyle(fontColor, brushColor, penColor gxui.Color, penWidth float32) Style {
	return Style{
		FontColor: fontColor,
		Pen:       gxui.CreatePen(penWidth, penColor),
		Brush:     gxui.CreateBrush(brushColor),
	}
}
