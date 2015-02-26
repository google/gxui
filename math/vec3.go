// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import (
	"fmt"
)

type Vec3 struct {
	X, Y, Z float32
}

func (v Vec3) String() string {
	return fmt.Sprintf("(%.5f, %.5f, %.5f)", v.X, v.Y, v.Z)
}

func (v Vec3) SqrLen() float32 {
	return v.Dot(v)
}

func (v Vec3) Len() float32 {
	return Sqrtf(v.SqrLen())
}

func (v Vec3) Normalize() Vec3 {
	l := v.Len()
	if l == 0 {
		return Vec3{0, 0, 0}
	} else {
		return v.MulS(1.0 / v.Len())
	}
}

func (v Vec3) Neg() Vec3 {
	return Vec3{-v.X, -v.Y, -v.Z}
}

func (v Vec3) XY() Vec2 {
	return Vec2{v.X, v.Y}
}

func (v Vec3) Add(o Vec3) Vec3 {
	return Vec3{v.X + o.X, v.Y + o.Y, v.Z + o.Z}
}

func (v Vec3) Sub(o Vec3) Vec3 {
	return Vec3{v.X - o.X, v.Y - o.Y, v.Z - o.Z}
}

func (v Vec3) Mul(o Vec3) Vec3 {
	return Vec3{v.X * o.X, v.Y * o.Y, v.Z * o.Z}
}

func (v Vec3) Div(o Vec3) Vec3 {
	return Vec3{v.X / o.X, v.Y / o.Y, v.Z / o.Z}
}

func (v Vec3) Dot(o Vec3) float32 {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}

func (v Vec3) Cross(o Vec3) Vec3 {
	return Vec3{
		v.Y*o.Z - v.Z*o.Y,
		v.Z*o.X - v.X*o.Z,
		v.X*o.Y - v.Y*o.X,
	}
}

//                ╭          ╮
//                │ M₀ M₁ M₂ │
// [V₀, V₁, V₂] ⨯ │ M₃ M₄ M₅ │ = [R₀, R₁, R₂]
//                │ M₆ M₇ M₈ │
//                ╰          ╯
// R₀ = V₀ • M₀ + V₁ • M₃ + V₂ • M₆
// R₁ = V₀ • M₁ + V₁ • M₄ + V₂ • M₇
// R₂ = V₀ • M₂ + V₁ • M₅ + V₂ • M₈
func (v Vec3) MulM(m Mat3) Vec3 {
	a := m.Row(0).MulS(v.X)
	b := m.Row(1).MulS(v.Y)
	c := m.Row(2).MulS(v.Z)
	return a.Add(b).Add(c)
}

func (v Vec3) MulS(s float32) Vec3 {
	return Vec3{v.X * s, v.Y * s, v.Z * s}
}

func (v Vec3) DivS(s float32) Vec3 {
	return Vec3{v.X / s, v.Y / s, v.Z / s}
}
