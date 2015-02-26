// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import (
	"math"
)

type Size struct {
	W, H int
}

func (s Size) Point() Point {
	return Point{s.W, s.H}
}

func (s Size) Vec2() Vec2 {
	return Vec2{float32(s.W), float32(s.H)}
}

func (s Size) Rect() Rect {
	return CreateRect(0, 0, s.W, s.H)
}

func (s Size) CenteredRect() Rect {
	return CreateRect(-s.W/2, -s.H/2, s.W/2, s.H/2)
}

func (s Size) Scale(v Vec2) Size {
	return Size{
		int(math.Ceil(float64(s.W) * float64(v.X))),
		int(math.Ceil(float64(s.H) * float64(v.Y))),
	}
}
func (s Size) ScaleS(v float32) Size {
	return Size{
		int(math.Ceil(float64(s.W) * float64(v))),
		int(math.Ceil(float64(s.H) * float64(v))),
	}
}

func (s Size) Expand(sp Spacing) Size {
	return Size{s.W + sp.W(), s.H + sp.H()}
}

func (s Size) Contract(sp Spacing) Size {
	return Size{s.W - sp.W(), s.H - sp.H()}
}

func (s Size) Add(o Size) Size {
	return Size{s.W + o.W, s.H + o.H}
}

func (s Size) Sub(o Size) Size {
	return Size{s.W - o.W, s.H - o.H}
}

func (s Size) Min(o Size) Size {
	return Size{Min(s.W, o.W), Min(s.H, o.H)}
}

func (s Size) Max(o Size) Size {
	return Size{Max(s.W, o.W), Max(s.H, o.H)}
}

func (s Size) Clamp(min, max Size) Size {
	return Size{Clamp(s.W, min.W, max.W), Clamp(s.H, min.H, max.H)}
}

func (s Size) WH() (w, h int) {
	return s.W, s.H
}

func (s Size) Area() int {
	return s.W * s.H
}

func (s Size) EdgeAlignedFit(outer Rect, edgePoint Point) Rect {
	r := s.CenteredRect().Offset(edgePoint).Constrain(outer)
	if topFits := edgePoint.Y+s.H < outer.Max.Y; topFits {
		return r.OffsetY(edgePoint.Y - r.Min.Y)
	}
	if bottomFits := edgePoint.Y-s.H >= outer.Min.Y; bottomFits {
		return r.OffsetY(edgePoint.Y - r.Max.Y)
	}
	if leftFits := edgePoint.X+s.W < outer.Max.X; leftFits {
		return r.OffsetX(edgePoint.X - r.Min.X)
	}
	if rightFits := edgePoint.X-s.W >= outer.Min.X; rightFits {
		return r.OffsetX(edgePoint.X - r.Max.X)
	}
	return r
}
