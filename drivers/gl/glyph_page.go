// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"github.com/google/gxui/math"
	"image"
	"image/png"
	"os"

	"code.google.com/p/freetype-go/freetype/raster"
	"code.google.com/p/freetype-go/freetype/truetype"
)

const (
	dumpGlyphPages  = false
	glyphPageWidth  = 512
	glyphPageHeight = 512
	glyphPadding    = 1
)

type glyphPage struct {
	resolution         resolution
	glyphMaxSizePixels math.Size
	image              *image.Alpha
	offsets            map[rune]math.Point
	rowHeight          int
	rast               *raster.Rasterizer
	tex                *texture
	nextPoint          math.Point
}

func createGlyphPage(resolution resolution, glyphMaxSizePixels math.Size) *glyphPage {
	return &glyphPage{
		resolution:         resolution,
		glyphMaxSizePixels: glyphMaxSizePixels,
		image:              image.NewAlpha(image.Rect(0, 0, glyphPageWidth, glyphPageHeight)),
		offsets:            make(map[rune]math.Point),
		rowHeight:          0,
		rast:               raster.NewRasterizer(glyphMaxSizePixels.W, glyphMaxSizePixels.H),
	}
}

// drawContour draws the given closed contour with the given offset.
func (p *glyphPage) drawContour(ps []truetype.Point, dx, dy raster.Fix32) {
	if len(ps) == 0 {
		return
	}
	rast := p.rast
	resolution := p.resolution
	// ps[0] is a truetype.Point measured in FUnits and positive Y going upwards.
	// start is the same thing measured in fixed point units and positive Y
	// going downwards, and offset by (dx, dy)
	start := raster.Point{
		X: dx + raster.Fix32(ps[0].X*int32(resolution)>>14),
		Y: dy - raster.Fix32(ps[0].Y*int32(resolution)>>14),
	}
	rast.Start(start)
	q0, on0 := start, true
	for _, p := range ps[1:] {
		q := raster.Point{
			X: dx + raster.Fix32(p.X*int32(resolution)>>14),
			Y: dy - raster.Fix32(p.Y*int32(resolution)>>14),
		}
		on := p.Flags&0x01 != 0
		if on {
			if on0 {
				rast.Add1(q)
			} else {
				rast.Add2(q0, q)
			}
		} else {
			if on0 {
				// No-op.
			} else {
				mid := raster.Point{
					X: (q0.X + q.X) / 2,
					Y: (q0.Y + q.Y) / 2,
				}
				rast.Add2(q0, mid)
			}
		}
		q0, on0 = q, on
	}
	// Close the curve.
	if on0 {
		rast.Add1(start)
	} else {
		rast.Add2(q0, start)
	}
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

func (p *glyphPage) add(rune rune, g *glyph) bool {
	if _, found := p.offsets[rune]; found {
		panic("Glyph already added to glyph page")
	}

	w, h := g.size(p.resolution).WH()
	x, y := p.nextPoint.X, p.nextPoint.Y

	if x+w > glyphPageWidth {
		// Row full, start new line
		x = 0
		y += p.rowHeight + glyphPadding
		p.rowHeight = 0
	}

	if y+h > glyphPageHeight {
		return false // Page full
	}

	// Build the raster contours
	p.rast.Clear()
	fx := -raster.Fix32(g.B.XMin * int32(p.resolution) >> 14)
	fy := +raster.Fix32(g.B.YMax * int32(p.resolution) >> 14)
	e0 := 0
	for _, e1 := range g.End {
		p.drawContour(g.Point[e0:e1], fx, fy)
		e0 = e1
	}

	// Perform the rasterization
	a := &image.Alpha{
		Pix:    p.image.Pix[x+y*p.image.Stride:],
		Stride: p.image.Stride,
		Rect:   image.Rect(0, 0, w, h),
	}
	p.rast.Rasterize(raster.NewAlphaSrcPainter(a))

	p.offsets[rune] = math.Point{X: x, Y: y}
	p.nextPoint = math.Point{X: x + w + glyphPadding, Y: y}
	if h > p.rowHeight {
		p.rowHeight = h
	}

	if p.tex != nil {
		p.tex.Release()
		p.tex = nil
	}

	return true
}

func (p *glyphPage) texture() *texture {
	if p.tex == nil {
		p.commit()
	}
	return p.tex
}

func (p *glyphPage) offset(rune rune) math.Point {
	return p.offsets[rune]
}
