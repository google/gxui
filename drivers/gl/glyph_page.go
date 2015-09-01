// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/google/gxui/math"
	fnt "golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	dumpGlyphPages     = false
	glyphPageWidth     = 512
	glyphPageHeight    = 512
	glyphSizeAlignment = 8
	glyphPadding       = 1
)

type glyphEntry struct {
	offset math.Point
	bounds math.Rect
}

type glyphPage struct {
	image     *image.Alpha
	size      math.Size // in pixels
	entries   map[rune]glyphEntry
	rowHeight int
	tex       *texture
	nextPoint math.Point
}

func point26_6toPoint(p fixed.Point26_6) math.Point {
	return math.Point{X: int(p.X) >> 6, Y: int(p.Y) >> 6}
}

func rectangle26_6toRect(p fixed.Rectangle26_6) math.Rect {
	return math.Rect{Min: point26_6toPoint(p.Min), Max: point26_6toPoint(p.Max)}
}

func align(v, pot int) int {
	return (v + pot - 1) & ^(pot - 1)
}

func newGlyphPage(face fnt.Face, r rune) *glyphPage {
	// Start the page big enough to hold the initial rune.
	b, _, _ := face.GlyphBounds(r)
	bounds := rectangle26_6toRect(b)
	size := math.Size{W: glyphPageWidth, H: glyphPageHeight}.Max(bounds.Size())
	size.W = align(size.W, glyphSizeAlignment)
	size.H = align(size.H, glyphSizeAlignment)

	page := &glyphPage{
		image:     image.NewAlpha(image.Rect(0, 0, size.W, size.H)),
		size:      size,
		entries:   make(map[rune]glyphEntry),
		rowHeight: 0,
	}
	page.add(face, r)
	return page
}

func (p *glyphPage) commit() {
	if p.tex != nil {
		return
	}
	p.tex = newTexture(p.image, 1.0)
	if dumpGlyphPages {
		f, _ := os.Create("glyph-page.png")
		defer f.Close()
		png.Encode(f, p.image)
	}
}

func (p *glyphPage) add(face fnt.Face, r rune) bool {
	if _, found := p.entries[r]; found {
		panic("Glyph already added to glyph page")
	}

	b, mask, maskp, _, _ := face.Glyph(fixed.Point26_6{}, r)
	bounds := math.CreateRect(b.Min.X, b.Min.Y, b.Max.X, b.Max.Y)

	w, h := bounds.Size().WH()
	x, y := p.nextPoint.X, p.nextPoint.Y

	if x+w > p.size.W {
		// Row full, start new line
		x = 0
		y += p.rowHeight + glyphPadding
		p.rowHeight = 0
	}

	if y+h > p.size.H {
		return false // Page full
	}

	draw.Draw(p.image, image.Rect(x, y, x+w, y+h), mask, maskp, draw.Src)

	p.entries[r] = glyphEntry{
		offset: math.Point{X: x, Y: y}.Sub(bounds.Min),
		bounds: bounds,
	}
	p.nextPoint = math.Point{X: x + w + glyphPadding, Y: y}
	if h > p.rowHeight {
		p.rowHeight = h
	}
	p.tex = nil

	return true
}

func (p *glyphPage) texture() *texture {
	if p.tex == nil {
		p.commit()
	}
	return p.tex
}

func (p *glyphPage) get(rune rune) glyphEntry {
	return p.entries[rune]
}
