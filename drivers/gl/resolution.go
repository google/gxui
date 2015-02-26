// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"gaze/gxui/math"
)

// 16:16 fixed point ratio of DIPs to pixels
type Resolution uint32

func (r Resolution) IntDipsToPixels(s int) int {
	return (s * int(r)) >> 16
}

func (r Resolution) PointDipsToPixels(s math.Point) math.Point {
	return math.Point{
		X: r.IntDipsToPixels(s.X),
		Y: r.IntDipsToPixels(s.Y),
	}
}

func (r Resolution) SizeDipsToPixels(s math.Size) math.Size {
	return math.Size{
		W: r.IntDipsToPixels(s.W),
		H: r.IntDipsToPixels(s.H),
	}
}

func (r Resolution) RectDipsToPixels(s math.Rect) math.Rect {
	return math.Rect{
		Min: r.PointDipsToPixels(s.Min),
		Max: r.PointDipsToPixels(s.Max),
	}
}
