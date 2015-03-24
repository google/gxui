// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/google/gxui/math"
)

// A Font represents a TrueType font loaded by the GXUI driver.
type Font interface {
	LoadGlyphs(first, last rune)
	Size() int
	GlyphMaxSize() math.Size
	Measure(*TextBlock) math.Size
	Layout(*TextBlock) (offsets []math.Point)
}

// TextBlock is a sequence of runes to be laid out.
type TextBlock struct {
	Runes     []rune
	AlignRect math.Rect
	H         HorizontalAlignment
	V         VerticalAlignment
}
