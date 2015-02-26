// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import (
	"fmt"
)

type Vec4 struct {
	X, Y, Z, W float32
}

func (v Vec4) String() string {
	return fmt.Sprintf("(%.5f, %.5f, %.5f, %.5f)", v.X, v.Y, v.Z, v.W)
}

func (v Vec4) SqrLen() float32 {
	return v.Dot(v)
}

func (v Vec4) Len() float32 {
	return Sqrtf(v.SqrLen())
}

func (v Vec4) Normalize() Vec4 {
	l := v.Len()
	if l == 0 {
		return Vec4{0, 0, 0, 0}
	} else {
		return v.MulS(1.0 / v.Len())
	}
}

func (v Vec4) Neg() Vec4 {
	return Vec4{-v.X, -v.Y, -v.Z, -v.W}
}

func (v Vec4) XY() Vec2 {
	return Vec2{v.X, v.Y}
}

func (v Vec4) Add(o Vec4) Vec4 {
	return Vec4{v.X + o.X, v.Y + o.Y, v.Z + o.Z, v.W + o.W}
}

func (v Vec4) Sub(o Vec4) Vec4 {
	return Vec4{v.X - o.X, v.Y - o.Y, v.Z - o.Z, v.W - o.W}
}

func (v Vec4) Mul(o Vec4) Vec4 {
	return Vec4{v.X * o.X, v.Y * o.Y, v.Z * o.Z, v.W * o.W}
}

func (v Vec4) Div(o Vec4) Vec4 {
	return Vec4{v.X / o.X, v.Y / o.Y, v.Z / o.Z, v.W / o.W}
}

func (v Vec4) Dot(o Vec4) float32 {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z + v.W*o.W
}

func (v Vec4) MulS(s float32) Vec4 {
	return Vec4{v.X * s, v.Y * s, v.Z * s, v.W * s}
}

func (v Vec4) DivS(s float32) Vec4 {
	return Vec4{v.X / s, v.Y / s, v.Z / s, v.W / s}
}
