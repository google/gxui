// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

type Rect struct {
	Min, Max Point
}

func CreateRect(minX, minY, maxX, maxY int) Rect {
	return Rect{Point{minX, minY}, Point{maxX, maxY}}
}

func (r Rect) Mid() Point {
	return Point{
		(r.Min.X + r.Max.X) / 2,
		(r.Min.Y + r.Max.Y) / 2,
	}
}

func (r Rect) W() int {
	return r.Max.X - r.Min.X
}

func (r Rect) H() int {
	return r.Max.Y - r.Min.Y
}

func (r Rect) TL() Point {
	return r.Min
}

func (r Rect) TC() Point {
	return Point{(r.Min.X + r.Max.X) / 2, r.Min.Y}
}

func (r Rect) TR() Point {
	return Point{r.Max.X, r.Min.Y}
}

func (r Rect) BL() Point {
	return Point{r.Min.X, r.Max.Y}
}

func (r Rect) BC() Point {
	return Point{(r.Min.X + r.Max.X) / 2, r.Max.Y}
}

func (r Rect) BR() Point {
	return r.Max
}

func (r Rect) ML() Point {
	return Point{r.Min.X, (r.Min.Y + r.Max.Y) / 2}
}

func (r Rect) MR() Point {
	return Point{r.Max.X, (r.Min.Y + r.Max.Y) / 2}
}

func (r Rect) Size() Size {
	return Size{r.Max.X - r.Min.X, r.Max.Y - r.Min.Y}
}

func (r Rect) ScaleAt(p Point, s Vec2) Rect {
	return Rect{
		p.Add(r.Min.Sub(p).Scale(s)),
		p.Add(r.Max.Sub(p).Scale(s)),
	}
}

func (r Rect) ScaleS(s float32) Rect {
	return Rect{r.Min.ScaleS(s), r.Max.ScaleS(s)}
}

func (r Rect) Offset(p Point) Rect {
	return Rect{r.Min.Add(p), r.Max.Add(p)}
}

func (r Rect) OffsetX(i int) Rect {
	return r.Offset(Point{i, 0})
}

func (r Rect) OffsetY(i int) Rect {
	return r.Offset(Point{0, i})
}

func (r Rect) ClampXY(x, y int) (int, int) {
	return Clamp(x, r.Min.X, r.Max.X), Clamp(y, r.Min.Y, r.Max.Y)
}

func (r Rect) Lerp(v Vec2) Point {
	return r.Min.Add(r.Size().Scale(v).Point())
}

func (r Rect) Frac(v Point) Vec2 {
	return v.Sub(r.Min).Vec2().Div(r.Size().Vec2())
}

func (r Rect) Remap(from, to Rect) Rect {
	return Rect{r.Min.Remap(from, to), r.Max.Remap(from, to)}
}

func (r Rect) Expand(s Spacing) Rect {
	return Rect{
		Point{r.Min.X - s.L, r.Min.Y - s.T},
		Point{r.Max.X + s.R, r.Max.Y + s.B},
	}.Canon()
}

func (r Rect) ExpandI(s int) Rect {
	return Rect{
		Point{r.Min.X - s, r.Min.Y - s},
		Point{r.Max.X + s, r.Max.Y + s},
	}.Canon()
}

func (r Rect) Contract(s Spacing) Rect {
	return Rect{
		Point{r.Min.X + s.L, r.Min.Y + s.T},
		Point{r.Max.X - s.R, r.Max.Y - s.B},
	}.Canon()
}

func (r Rect) ContractI(s int) Rect {
	return Rect{
		Point{r.Min.X + s, r.Min.Y + s},
		Point{r.Max.X - s, r.Max.Y - s},
	}.Canon()
}

func (r Rect) Union(o Rect) Rect {
	return Rect{r.Min.Min(o.Min), r.Max.Max(o.Max)}
}

func (r Rect) Intersect(o Rect) Rect {
	return Rect{
		r.Min.Max(o.Min),
		r.Max.Min(o.Max),
	}.Canon()
}

func (r Rect) Constrain(o Rect) Rect {
	overflowMin := o.Min.Sub(r.Min).Max(ZeroPoint)
	overflowMax := o.Max.Sub(r.Max).Min(ZeroPoint)
	return Rect{
		r.Min.Add(overflowMax).Max(o.Min),
		r.Max.Add(overflowMin).Min(o.Max),
	}
}

func (r Rect) Canon() Rect {
	return Rect{
		r.Min.Min(r.Max),
		r.Min.Max(r.Max),
	}
}

func (r Rect) Contains(p Point) bool {
	return r.Min.X <= p.X && r.Min.Y <= p.Y &&
		r.Max.X > p.X && r.Max.Y > p.Y
}
