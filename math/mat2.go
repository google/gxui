// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import (
	"fmt"
)

// ╭       ╮
// │ M₀ M₁ │
// │ M₂ M₃ │
// ╰       ╯
type Mat2 [4]float32

func CreateMat2(r0c0, r0c1, r1c0, r1c1 float32) Mat2 {
	return Mat2{
		r0c0, r0c1,
		r1c0, r1c1,
	}
}

// ╭    ╮
// │ R₀ │
// │ R₁ │
// ╰    ╯
func CreateMat2FromRows(r0, r1 Vec2) Mat2 {
	return Mat2{
		r0.X, r0.Y,
		r1.X, r1.Y,
	}
}

func (m Mat2) String() string {
	s := make([]string, 4)
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
		"\n╭ %s %s ╮"+
			"\n│ %s %s │"+
			"\n│ %s %s │"+
			"\n╰ %s %s ╯",
		p, p,
		s[0], s[1],
		s[2], s[3],
		p, p,
	)
}

func (m Mat2) Rows() (r0, r1 Vec2) {
	return Vec2{m[0], m[1]}, Vec2{m[2], m[3]}
}

func (m Mat2) Row(i int) Vec2 {
	i *= 2
	return Vec2{m[i+0], m[i+1]}
}

func (m Mat2) Invert() Mat2 {

	// ╭       ╮⁻¹        ╭         ╮
	// │ M₀ M₁ │    =  1  │  M₃ -M₁ │
	// │ M₂ M₃ │      ─── │ -M₂  M₀ │
	// ╰       ╯      det ╰         ╯
	//
	// Where: det = M₀ • M₃ - M₁ • M₂
	//
	det := m[0]*m[3] - m[1]*m[2]
	return DivM2S(CreateMat2(m[3], -m[1], -m[2], m[0]), det)
}

func (m Mat2) Transpose() Mat2 {
	return CreateMat2(
		m[0], m[2],
		m[1], m[3],
	)
}

//            ╭       ╮
//            │ M₀ M₁ │
// [V₀, V₁] ⨯ │ M₂ M₃ │ = [V₀ • M₀ + V₁ • M₂, V₀ • M₁ + V₁ • M₃]
//            ╰       ╯
func MulVM2(v Vec2, m Mat2) Vec2 {
	a := m.Row(0).MulS(v.X)
	b := m.Row(1).MulS(v.Y)
	return a.Add(b)
}

func DivM2S(m Mat2, s float32) Mat2 {
	return CreateMat2FromRows(
		m.Row(0).DivS(s),
		m.Row(1).DivS(s),
	)
}
