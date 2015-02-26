// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

type Point struct {
	X, Y int
}

func NewPoint(X, Y int) Point {
	return Point{X: X, Y: Y}
}

func (p Point) Add(o Point) Point {
	return Point{p.X + o.X, p.Y + o.Y}
}

func (p Point) AddX(o int) Point {
	return Point{p.X + o, p.Y}
}

func (p Point) AddY(o int) Point {
	return Point{p.X, p.Y + o}
}

func (p Point) Sub(o Point) Point {
	return Point{p.X - o.X, p.Y - o.Y}
}

func (p Point) Neg() Point {
	return Point{-p.X, -p.Y}
}

func (p Point) SqrLen() int {
	return p.Dot(p)
}

func (p Point) Len() float32 {
	return Sqrtf(float32(p.SqrLen()))
}

func (p Point) Dot(o Point) int {
	return p.X*o.X + p.Y*o.Y
}

func (p Point) XY() (x, y int) {
	return p.X, p.Y
}

func (p Point) Vec2() Vec2 {
	return Vec2{float32(p.X), float32(p.Y)}
}

func (p Point) Vec3(z float32) Vec3 {
	return Vec3{float32(p.X), float32(p.Y), z}
}

func (p Point) Scale(s Vec2) Point {
	return Point{int(float32(p.X) * s.X), int(float32(p.Y) * s.Y)}
}

func (p Point) ScaleS(s float32) Point {
	return Point{int(float32(p.X) * s), int(float32(p.Y) * s)}
}

func (p Point) ScaleX(s float32) Point {
	return Point{int(float32(p.X) * s), p.Y}
}

func (p Point) ScaleY(s float32) Point {
	return Point{p.X, int(float32(p.Y) * s)}
}

func (p Point) Size() Size {
	return Size{p.X, p.Y}
}

func (p Point) Min(o Point) Point {
	return Point{Min(p.X, o.X), Min(p.Y, o.Y)}
}

func (p Point) Max(o Point) Point {
	return Point{Max(p.X, o.X), Max(p.Y, o.Y)}
}

func (p Point) Clamp(min, max Point) Point {
	return p.Min(max).Max(min)
}

func (p Point) Remap(from, to Rect) Point {
	return p.Sub(from.Min).
		ScaleX(float32(to.W()) / float32(from.W())).
		ScaleY(float32(to.H()) / float32(from.H())).
		Add(to.Min)
}
