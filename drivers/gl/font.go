// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"
	"unicode"

	"github.com/golang/freetype/truetype"
	"github.com/google/gxui"
	"github.com/google/gxui/math"
	fnt "golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type font struct {
	size             int
	scale            fixed.Int26_6
	glyphMaxSizeDips math.Size
	ascentDips       int
	ttf              *truetype.Font
	resolutions      map[resolution]*glyphTable
	glyphAdvanceDips map[rune]int
}

func newFont(data []byte, size int) (*font, error) {
	ttf, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}

	scale := fixed.Int26_6(size << 6)
	bounds := rectangle26_6toRect(ttf.Bounds(scale))
	ascentDips := bounds.Max.Y

	return &font{
		size:             size,
		scale:            scale,
		glyphMaxSizeDips: bounds.Size(),
		ascentDips:       ascentDips,
		ttf:              ttf,
		resolutions:      make(map[resolution]*glyphTable),
		glyphAdvanceDips: make(map[rune]int),
	}, nil
}

func (f *font) advanceDips(r rune) int {
	if g, found := f.glyphAdvanceDips[r]; found {
		return g
	}
	idx := f.ttf.Index(r)
	gb := &truetype.GlyphBuf{}
	err := gb.Load(f.ttf, f.scale, idx, fnt.HintingFull)
	if err != nil {
		panic(err)
	}

	advance := int((gb.AdvanceWidth + 0x3f) >> 6)
	f.glyphAdvanceDips[r] = advance
	return advance
}

func (f *font) glyphTable(resolution resolution) *glyphTable {
	t, found := f.resolutions[resolution]
	if !found {
		opt := truetype.Options{
			Size:              float64(f.size),
			DPI:               float64(resolution.intDipsToPixels(72)),
			Hinting:           fnt.HintingFull,
			GlyphCacheEntries: 1,
			SubPixelsX:        1,
			SubPixelsY:        1,
		}
		t = newGlyphTable(truetype.NewFace(f.ttf, &opt))
		f.resolutions[resolution] = t
	}
	return t
}

func (f *font) align(rect math.Rect, size math.Size, ascent int, h gxui.HorizontalAlignment, v gxui.VerticalAlignment) math.Point {
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

func (f *font) DrawRunes(ctx *context, runes []rune, offsets []math.Point, col gxui.Color, ds *drawState) {
	if len(runes) != len(offsets) {
		panic(fmt.Errorf("There must be the same number of runes to offsets. Got %d runes and %d offsets",
			len(runes), len(offsets)))
	}
	resolution := ctx.resolution
	table := f.glyphTable(resolution)

	for i, r := range runes {
		if unicode.IsSpace(r) {
			continue
		}
		page := table.get(r)
		texture := page.texture()
		entry := page.get(r)
		srcRect := entry.bounds.Offset(entry.offset)
		dstRect := entry.bounds.Offset(resolution.pointDipsToPixels(offsets[i]))
		tc := ctx.getOrCreateTextureContext(texture)
		ctx.blitter.blitGlyph(ctx, tc, col, srcRect, dstRect, ds)
	}
}

func (f *font) Size() int {
	return f.size
}

func (f *font) Measure(fl *gxui.TextBlock) math.Size {
	size := math.Size{W: 0, H: f.glyphMaxSizeDips.H}
	var offset math.Point
	for _, r := range fl.Runes {
		if r == '\n' {
			offset.X = 0
			offset.Y += f.glyphMaxSizeDips.H
			continue
		}
		offset.X += f.advanceDips(r)
		size = size.Max(math.Size{W: offset.X, H: offset.Y + f.glyphMaxSizeDips.H})
	}
	return size
}

func (f *font) Layout(fl *gxui.TextBlock) (offsets []math.Point) {
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
		offset.X += f.advanceDips(r)
		sizeDips = sizeDips.Max(math.Size{W: offset.X, H: offset.Y + f.glyphMaxSizeDips.H})
	}

	origin := f.align(fl.AlignRect, sizeDips, f.ascentDips, fl.H, fl.V)
	for i, p := range offsets {
		offsets[i] = p.Add(origin)
	}
	return offsets
}

func (f *font) LoadGlyphs(first, last rune) {
	if first > last {
		first, last = last, first
	}
	for r := first; r < last; r++ {
		f.advanceDips(r)
	}
}

func (f *font) GlyphMaxSize() math.Size {
	return f.glyphMaxSizeDips
}
