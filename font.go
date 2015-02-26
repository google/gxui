// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"gxui/math"
)

type Font interface {
	LoadGlyphs(first, last rune)
	Name() string
	Size() int
	GlyphMaxSize() math.Size
	Measure(string) math.Size
	MeasureRunes([]rune) math.Size
	LayoutRunes(offsets []math.Point, runes []rune, alignRect math.Rect, h HorizontalAlignment, v VerticalAlignment)
}
