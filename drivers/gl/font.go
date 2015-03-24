// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"

	"github.com/google/gxui"
	"github.com/google/gxui/math"

	"code.google.com/p/freetype-go/freetype/truetype"
)

type Quad struct {
	Texture          *Texture
	SrcRect, DstRect math.Rect
}

type Font struct {
	size             int
	scale            int32
	glyphMaxSizeDips math.Size
	ascentDips       int
	ttf              *truetype.Font
	resolutions      map[Resolution]*glyphTable
	glyphs           map[rune]*glyph
	quads            []Quad // Reused each call to Draw()
}

func createFont(data []byte, size int) (*Font, error) {
	ttf, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}

	scale := int32(size << 6)
	bounds := ttf.Bounds(scale)
	glyphMaxSizeDips := math.Size{
		W: int(bounds.XMax-bounds.XMin) >> 6,
		H: int(bounds.YMax-bounds.YMin) >> 6,
	}
	ascentDips := int(bounds.YMax >> 6)

	return &Font{
		size:             size,
		scale:            scale,
		glyphMaxSizeDips: glyphMaxSizeDips,
		ascentDips:       ascentDips,
		ttf:              ttf,
		resolutions:      make(map[Resolution]*glyphTable),
		glyphs:           make(map[rune]*glyph),
		quads:            []Quad{},
	}, nil
}

func (f *Font) glyph(r rune) *glyph {
	if g, found := f.glyphs[r]; found {
		return g
	}
	idx := f.ttf.Index(r)
	gb := truetype.NewGlyphBuf()
	err := gb.Load(f.ttf, f.scale, idx, truetype.Hinting(truetype.FullHinting))
	if err != nil {
		panic(err)
	}

	g := glyph(*gb)
	f.glyphs[r] = &g
	return &g
}

func (f *Font) glyphTable(resolution Resolution) *glyphTable {
	t, found := f.resolutions[resolution]
	if !found {
		t = createGlyphTable(resolution, f.glyphMaxSizeDips)
		f.resolutions[resolution] = t
	}
	return t
}

func (f *Font) align(rect math.Rect, size math.Size, ascent int, h gxui.HorizontalAlignment, v gxui.VerticalAlignment) math.Point {
	var origin math.Point
	switch h {
	case gxui.AlignLeft:
		origin.X = rect.Min.X
	case gxui.AlignCenter:
		origin.X = rect.Mid().X - (size.W / 2)
	case gxui.AlignRight:
		origin.X = rect.Max.X - size.W
	}
	switch v {
	case gxui.AlignTop:
		origin.Y = rect.Min.Y + ascent
	case gxui.AlignMiddle:
		origin.Y = rect.Mid().Y - (size.H / 2) + ascent
	case gxui.AlignBottom:
		origin.Y = rect.Max.Y - size.H + ascent
	}
	return origin
}

func (f *Font) DrawRunes(ctx *Context, runes []rune, col gxui.Color, offsets []math.Point, origin math.Point, ds *DrawState) {
	if len(runes) != len(offsets) {
		panic(fmt.Errorf("There must be the same number of runes to offsets. Got %d runes and %d offsets",
			len(runes), len(offsets)))
	}
	resolution := ctx.Resolution()
	table := f.glyphTable(resolution)

	for i, r := range runes {
		if r == '\t' {
			continue
		}
		glyph := f.glyph(r)
		page := table.get(r, glyph)
		texture := page.texture()
		srcRect := glyph.size(resolution).Rect().Offset(page.offset(r))
		dstRect := glyph.rect(resolution).
			Offset(resolution.PointDipsToPixels(offsets[i])).
			Offset(origin)
		tc := ctx.GetOrCreateTextureContext(texture)
		ctx.Blitter.BlitGlyph(ctx, tc, col, srcRect, dstRect, ds)
	}
}

func (f *Font) Size() int {
	return f.size
}

func (f *Font) Measure(fl *gxui.TextBlock) math.Size {
	size := math.Size{W: 0, H: f.glyphMaxSizeDips.H}
	var offset math.Point
	for _, r := range fl.Runes {
		if r == '\n' {
			offset.X = 0
			offset.Y += f.glyphMaxSizeDips.H
			continue
		}
		offset.X += f.glyph(r).advanceDips()
		size = size.Max(math.Size{W: offset.X, H: offset.Y + f.glyphMaxSizeDips.H})
	}
	return size
}

func (f *Font) Layout(fl *gxui.TextBlock) (offsets []math.Point) {
	sizeDips := math.Size{}
	offsets = make([]math.Point, len(fl.Runes))
	var offset math.Point
	for i, r := range fl.Runes {
		if r == '\n' {
			offset.X = 0
			offset.Y += f.glyphMaxSizeDips.H
			continue
		}

		offsets[i] = offset
		offset.X += f.glyph(r).advanceDips()
		sizeDips = sizeDips.Max(math.Size{W: offset.X, H: offset.Y + f.glyphMaxSizeDips.H})
	}

	origin := f.align(fl.AlignRect, sizeDips, f.ascentDips, fl.H, fl.V)
	for i, p := range offsets {
		offsets[i] = p.Add(origin)
	}
	return offsets
}

func (f *Font) LoadGlyphs(first, last rune) {
	if first > last {
		first, last = last, first
	}
	for r := first; r < last; r++ {
		f.glyph(r)
	}
}

func (f *Font) GlyphMaxSize() math.Size {
	return f.glyphMaxSizeDips
}
