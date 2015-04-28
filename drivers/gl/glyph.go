// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"github.com/google/gxui/math"

	"code.google.com/p/freetype-go/freetype/truetype"
)

//            ╾──────w──────╼
//         min╔═════════════╗
//           ╿║    ▒███▒    ║ ╿
//           │║  █▒     ▒█  ║ │
//           │║ █░       ▒█ ║ │ ascent
//           │║ █░       ░█ ║ │
//           │║  █▒      ▒█ ║ │
//       │   h║   ░█████▒▒█ ║ ╽
//     ──┼───┼╫───────────█─╫───── +x
// origin│   │║          ▒█ ║
//       │   │║  ██▒   ░██  ║
//       │   ╽║   ░█████░   ║
//       │    ╚═════════════╝max
//       ╾────────────────────╼
//       │     advance
//       │
//       │
//       +y
//
// y-axis is flipped from freetype's.
// See: http://www.freetype.org/freetype2/docs/glyphs/glyphs-3.html#section-4
type glyph truetype.GlyphBuf

func (g *glyph) size(r resolution) math.Size {
	w := int((int64(g.B.XMax-g.B.XMin)*int64(r) + 0x3FFFFF) >> 22)
	h := int((int64(g.B.YMax-g.B.YMin)*int64(r) + 0x3FFFFF) >> 22)
	return math.Size{W: w, H: h}
}

func (g *glyph) sizeDips() math.Size {
	w := int(((g.B.XMax - g.B.XMin) + 0x1F) >> 6)
	h := int(((g.B.YMax - g.B.YMin) + 0x1F) >> 6)
	return math.Size{W: w, H: h}
}

func (g *glyph) rect(r resolution) math.Rect {
	x := int((int64(g.B.XMin) * int64(r)) >> 22)
	y := -int((int64(g.B.YMax) * int64(r)) >> 22)
	return g.size(r).Rect().Offset(math.Point{X: x, Y: y})
}

func (g *glyph) rectDips() math.Rect {
	x := int(g.B.XMin >> 6)
	y := -int(g.B.YMax >> 6)
	return g.sizeDips().Rect().Offset(math.Point{X: x, Y: y})
}

func (g *glyph) advance(r resolution) int {
	return int((int64(g.AdvanceWidth)*int64(r) + 0x3FFFFF) >> 22)
}

func (g *glyph) advanceDips() int {
	return int((g.AdvanceWidth + 0x3f) >> 6)
}
