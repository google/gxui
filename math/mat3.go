// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import (
	"fmt"
)

// A 3x3 matrix:
//   ╭          ╮
//   │ M₀ M₁ M₂ │
//   │ M₃ M₄ M₅ │
//   │ M₆ M₇ M₈ │
//   ╰          ╯
type Mat3 [9]float32

func CreateMat3(r0c0, r0c1, r0c2, r1c0, r1c1, r1c2, r2c0, r2c1, r2c2 float32) Mat3 {
	return Mat3{
		r0c0, r0c1, r0c2,
		r1c0, r1c1, r1c2,
		r2c0, r2c1, r2c2,
	}
}

// Build a 3x3 matrix from 3 row vectors
//  ╭    ╮
//  │ R₀ │
//  │ R₁ │
//  │ R₂ │
//  ╰    ╯
func CreateMat3FromRows(r0, r1, r2 Vec3) Mat3 {
	return Mat3{
		r0.X, r0.Y, r0.Z,
		r1.X, r1.Y, r1.Z,
		r2.X, r2.Y, r2.Z,
	}
}

//      A
//     ╱ ╲
//    ╱___╲
//   B     C
//
// [V₀, V₁, V₂] * M = [λ₀, λ₁, 1]
// λ₂ = 1 - (λ₀ + λ₁)
//
// [V₀, V₁, V₂] = A • λ₀ + B • λ₁ + C • λ₂
//
// A * M = (1, 0, 1)
// B * M = (0, 1, 1)
// C * M = (0, 0, 1)
func CreateMat3PositionToBarycentric(a, b, c Vec2) Mat3 {
	m := CreateMat2FromRows(a.Sub(c), b.Sub(c)).Invert()
	o := MulVM2(c, m)
	return Mat3{
		m[0], m[1], 0,
		m[2], m[3], 0,
		-o.X, -o.Y, 1,
	}
}

func (m Mat3) String() string {
	s := make([]string, 9)
	l := 0
	for i, v := range m {
		s[i] = fmt.Sprintf("%.5f", v)
		l = Max(l, len(s[i]))
	}
	for i, _ := range m {
		for len(s[i]) < l {
			s[i] = " " + s[i]
		}
	}
	p := ""
	for i := 0; i < l; i++ {
		p += " "
	}
	return fmt.Sprintf(
		"\n╭ %s %s %s ╮"+
			"\n│ %s %s %s │"+
			"\n│ %s %s %s │"+
			"\n│ %s %s %s │"+
			"\n╰ %s %s %s ╯",
		p, p, p,
		s[0], s[1], s[2],
		s[3], s[4], s[5],
		s[6], s[7], s[8],
		p, p, p,
	)
}

func (m Mat3) Rows() (r0, r1, r2 Vec3) {
	return Vec3{m[0], m[1], m[2]}, Vec3{m[3], m[4], m[5]}, Vec3{m[6], m[7], m[8]}
}

func (m Mat3) Row(i int) Vec3 {
	i *= 3
	return Vec3{m[i+0], m[i+1], m[i+2]}
}

func (m Mat3) Invert() Mat3 {
	//           ╭         ╮
	//        1  │ C₁ ⨯ C₂ │
	// M⁻¹ = ─── │ C₂ ⨯ C₀ │
	//       det │ C₀ ⨯ C₁ │
	//           ╰         ╯
	//
	// Where: det = C₀ • (C₁ ⨯ C₂)
	//
	C0, C1, C2 := m.Transpose().Rows()
	C0C1, C1C2, C2C0 := C0.Cross(C1), C1.Cross(C2), C2.Cross(C0)
	det := C0.Dot(C1C2)
	inv := CreateMat3FromRows(C1C2, C2C0, C0C1).DivS(det)
	return inv
}

func (m Mat3) Transpose() Mat3 {
	return CreateMat3(
		m[0], m[3], m[6],
		m[1], m[4], m[7],
		m[2], m[5], m[8],
	)
}

func (m Mat3) DivS(s float32) Mat3 {
	return CreateMat3FromRows(
		m.Row(0).DivS(s),
		m.Row(1).DivS(s),
		m.Row(2).DivS(s),
	)
}
