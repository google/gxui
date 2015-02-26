// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

type Spacing struct {
	L, T, R, B int
}

func CreateSpacing(s int) Spacing {
	return Spacing{s, s, s, s}
}

func (s Spacing) LT() Point {
	return Point{s.L, s.T}
}

func (s Spacing) W() int {
	return s.L + s.R
}

func (s Spacing) H() int {
	return s.T + s.B
}

func (s Spacing) Size() Size {
	return Size{s.W(), s.H()}
}

func (s Spacing) Add(o Spacing) Spacing {
	return Spacing{s.L + o.L, s.T + o.T, s.R + o.R, s.B + o.B}
}

func (s Spacing) Sub(o Spacing) Spacing {
	return Spacing{s.L - o.L, s.T - o.T, s.R - o.R, s.B - o.B}
}

func (s Spacing) Min(o Spacing) Spacing {
	return Spacing{Min(s.L, o.L), Min(s.T, o.T), Min(s.R, o.R), Min(s.B, o.B)}
}

func (s Spacing) Max(o Spacing) Spacing {
	return Spacing{Max(s.L, o.L), Max(s.T, o.T), Max(s.R, o.R), Max(s.B, o.B)}
}
