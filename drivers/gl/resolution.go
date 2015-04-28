// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"

	"github.com/google/gxui/math"
)

// 16:16 fixed point ratio of DIPs to pixels
type resolution uint32

func (r resolution) String() string {
	return fmt.Sprintf("%f", r.dipsToPixels())
}

func (r resolution) dipsToPixels() float32 {
	return float32(r) / 65536.0
}

func (r resolution) intDipsToPixels(s int) int {
	return (s * int(r)) >> 16
}

func (r resolution) pointDipsToPixels(s math.Point) math.Point {
	return math.Point{
		X: r.intDipsToPixels(s.X),
		Y: r.intDipsToPixels(s.Y),
	}
}

func (r resolution) sizeDipsToPixels(s math.Size) math.Size {
	return math.Size{
		W: r.intDipsToPixels(s.W),
		H: r.intDipsToPixels(s.H),
	}
}

func (r resolution) rectDipsToPixels(s math.Rect) math.Rect {
	return math.Rect{
		Min: r.pointDipsToPixels(s.Min),
		Max: r.pointDipsToPixels(s.Max),
	}
}
