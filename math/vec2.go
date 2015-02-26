// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

type Vec2 struct {
	X, Y float32
}

func (v Vec2) SqrLen() float32 {
	return v.Dot(v)
}

func (v Vec2) Len() float32 {
	return Sqrtf(v.SqrLen())
}

func (v Vec2) ZeroLength() bool {
	return v.X == 0 && v.Y == 0
}

func (v Vec2) Normalize() Vec2 {
	l := v.Len()
	if l == 0 {
		return Vec2{0, 0}
	} else {
		return v.MulS(1.0 / v.Len())
	}
}

func (v Vec2) Neg() Vec2 {
	return Vec2{-v.X, -v.Y}
}

func (v Vec2) Tangent() Vec2 {
	return Vec2{-v.Y, v.X}
}

func (v Vec2) Point() Point {
	return Point{Round(v.X), Round(v.Y)}
}

func (v Vec2) Vec3(z float32) Vec3 {
	return Vec3{v.X, v.Y, z}
}

func (v Vec2) Vec4(z, w float32) Vec4 {
	return Vec4{v.X, v.Y, z, w}
}

func (v Vec2) XY() (x, y float32) {
	return v.X, v.Y
}

func (v Vec2) Add(o Vec2) Vec2 {
	return Vec2{v.X + o.X, v.Y + o.Y}
}

func (v Vec2) Sub(o Vec2) Vec2 {
	return Vec2{v.X - o.X, v.Y - o.Y}
}

func (v Vec2) Mul(o Vec2) Vec2 {
	return Vec2{v.X * o.X, v.Y * o.Y}
}

func (v Vec2) Div(o Vec2) Vec2 {
	return Vec2{v.X / o.X, v.Y / o.Y}
}

func (v Vec2) Dot(o Vec2) float32 {
	return v.X*o.X + v.Y*o.Y
}

func (v Vec2) Cross(o Vec2) float32 {
	return v.X*o.Y - v.Y*o.X
}

func (v Vec2) MulS(s float32) Vec2 {
	return Vec2{v.X * s, v.Y * s}
}

func (v Vec2) DivS(s float32) Vec2 {
	return Vec2{v.X / s, v.Y / s}
}
