// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"github.com/google/gxui/math"
)

type glyphTable struct {
	index map[rune]int
	pages []*glyphPage
}

func createGlyphTable(resolution resolution, glyphMaxSizeDips math.Size) *glyphTable {
	glyphMaxSizePixels := resolution.sizeDipsToPixels(glyphMaxSizeDips)
	return &glyphTable{
		index: make(map[rune]int),
		pages: []*glyphPage{createGlyphPage(resolution, glyphMaxSizePixels)},
	}
}

func (t *glyphTable) get(r rune, g *glyph) *glyphPage {
	if i, found := t.index[r]; found {
		return t.pages[i]
	}
	if page := t.pages[len(t.pages)-1]; !page.add(r, g) {
		page = createGlyphPage(page.resolution, page.glyphMaxSizePixels)
		page.add(r, g)
		t.pages = append(t.pages, page)
	}
	index := len(t.pages) - 1
	t.index[r] = index
	return t.pages[index]
}
