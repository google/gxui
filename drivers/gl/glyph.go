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
	w := int(((g.B.XMax-g.B.XMin)*int32(r) + 0x3FFFFF) >> 22)
	h := int(((g.B.YMax-g.B.YMin)*int32(r) + 0x3FFFFF) >> 22)
	return math.Size{W: w, H: h}
}

func (g *glyph) rect(r resolution) math.Rect {
	x := int((g.B.XMin * int32(r)) >> 22)
	y := -int((g.B.YMax * int32(r)) >> 22)
	return g.size(r).Rect().Offset(math.Point{X: x, Y: y})
}

func (g *glyph) advance(r resolution) int {
	return int((g.AdvanceWidth*int32(r) + 0x3FFFFF) >> 22)
}

func (g *glyph) advanceDips() int {
	return int((g.AdvanceWidth + 0x3f) >> 6)
}
