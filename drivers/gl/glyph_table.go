// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import fnt "golang.org/x/image/font"

type glyphTable struct {
	face  fnt.Face
	index map[rune]int
	pages []*glyphPage
}

func newGlyphTable(face fnt.Face) *glyphTable {
	return &glyphTable{face: face, index: make(map[rune]int)}
}

func (t *glyphTable) get(r rune) *glyphPage {
	if i, found := t.index[r]; found {
		return t.pages[i]
	}
	if len(t.pages) == 0 {
		t.pages = append(t.pages, newGlyphPage(t.face, r))
	} else {
		page := t.pages[len(t.pages)-1]
		if !page.add(t.face, r) {
			page = newGlyphPage(t.face, r)
			t.pages = append(t.pages, page)
		}
	}
	index := len(t.pages) - 1
	t.index[r] = index
	return t.pages[index]
}
