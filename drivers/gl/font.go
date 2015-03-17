// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"
	"unicode/utf8"

	"github.com/google/gxui"
	"github.com/google/gxui/math"

	"code.google.com/p/freetype-go/freetype/truetype"
)

type Quad struct {
	Texture          *Texture
	SrcRect, DstRect math.Rect
}

type Font struct {
	name             string
	size             int
	scale            int32
	glyphMaxSizeDips math.Size
	ascentDips       int
	ttf              *truetype.Font
	resolutions      map[Resolution]*glyphTable
	glyphs           map[rune]*glyph
	quads            []Quad // Reused each call to Draw()
}

func CreateFont(name string, data []byte, size int) (*Font, error) {
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
		name:             name,
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

func (f *Font) Draw(ctx *Context, str string, col gxui.Color, alignRectDips math.Rect, h gxui.HorizontalAlignment, v gxui.VerticalAlignment, ds *DrawState) {
	resolution := ctx.Resolution()
	glyphMaxSizePixels := resolution.SizeDipsToPixels(f.glyphMaxSizeDips)
	ascentPixels := resolution.IntDipsToPixels(f.ascentDips)
	alignRectPixels := resolution.RectDipsToPixels(alignRectDips)
	table := f.glyphTable(resolution)
	quads := f.quads[:0]
	sizePixels := math.Size{}
	offset := math.Point{}
	for i := 0; i < len(str); {
		r, l := utf8.DecodeRuneInString(str[i:])
		i += l

		if r == '\n' {
			offset.X = 0
			offset.Y += glyphMaxSizePixels.H
			continue
		}

		glyph := f.glyph(r)
		page := table.get(r, glyph)
		srcRect := glyph.size(resolution).Rect().Offset(page.offset(r))
		dstRect := glyph.rect(resolution).Offset(offset)
		offset.X += glyph.advance(resolution)
		quads = append(quads, Quad{page.texture(), srcRect, dstRect})
		sizePixels = sizePixels.Max(math.Size{
			W: offset.X,
			H: offset.Y + glyphMaxSizePixels.H,
		})
	}

	origin := f.align(alignRectPixels, sizePixels, ascentPixels, h, v)
	for _, q := range quads {
		tc := ctx.GetOrCreateTextureContext(q.Texture)
		ctx.Blitter.BlitGlyph(ctx, tc, col, q.SrcRect, q.DstRect.Offset(origin), ds)
	}
}

func (f *Font) DrawRunes(ctx *Context, runes []rune, col gxui.Color, points []math.Point, origin math.Point, ds *DrawState) {
	if len(runes) != len(points) {
		panic(fmt.Errorf("There must be the same number of runes to points. Got %d runes and %d points",
			len(runes), len(points)))
	}
	resolution := ctx.Resolution()
	table := f.glyphTable(resolution)

	for i, r := range runes {
		glyph := f.glyph(r)
		page := table.get(r, glyph)
		texture := page.texture()
		srcRect := glyph.size(resolution).Rect().Offset(page.offset(r))
		dstRect := glyph.rect(resolution).
			Offset(resolution.PointDipsToPixels(points[i])).
			Offset(origin)
		tc := ctx.GetOrCreateTextureContext(texture)
		ctx.Blitter.BlitGlyph(ctx, tc, col, srcRect, dstRect, ds)
	}
}

func (f *Font) Name() string {
	return f.name
}

func (f *Font) Size() int {
	return f.size
}

func (f *Font) Measure(s string) math.Size {
	size := math.Size{W: 0, H: f.glyphMaxSizeDips.H}
	var offset math.Point
	for i := 0; i < len(s); {
		r, l := utf8.DecodeRuneInString(s[i:])
		i += l

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

func (f *Font) MeasureRunes(runes []rune) math.Size {
	size := math.Size{W: 0, H: f.glyphMaxSizeDips.H}
	var offset math.Point
	for _, r := range runes {
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

func (f *Font) LayoutRunes(points []math.Point, runes []rune, alignRectDips math.Rect, h gxui.HorizontalAlignment, v gxui.VerticalAlignment) {
	sizeDips := math.Size{}
	var offset math.Point
	for i, r := range runes {
		if r == '\n' {
			offset.X = 0
			offset.Y += f.glyphMaxSizeDips.H
			continue
		}

		points[i] = offset
		offset.X += f.glyph(r).advanceDips()
		sizeDips = sizeDips.Max(math.Size{W: offset.X, H: offset.Y + f.glyphMaxSizeDips.H})
	}

	origin := f.align(alignRectDips, sizeDips, f.ascentDips, h, v)
	for i, p := range points {
		points[i] = p.Add(origin)
	}
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
