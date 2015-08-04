// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/google/gxui/math"
)

type Canvas interface {
	Size() math.Size
	IsComplete() bool
	Complete()
	Push()
	Pop()
	AddClip(math.Rect)
	Clear(Color)
	DrawCanvas(c Canvas, position math.Point)
	DrawTexture(t Texture, bounds math.Rect)
	DrawRunes(font Font, runes []rune, points []math.Point, color Color)
	DrawLines(Polygon, Pen)
	DrawPolygon(Polygon, Pen, Brush)
	DrawRect(math.Rect, Brush)
	DrawRoundedRect(rect math.Rect, tl, tr, bl, br float32, p Pen, b Brush)
}
