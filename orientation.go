// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type Orientation int

const (
	Vertical Orientation = iota
	Horizontal
)

func (o Orientation) Horizontal() bool { return o == Horizontal }
func (o Orientation) Vertical() bool   { return o == Vertical }

func (o Orientation) Flip() Orientation {
	if o == Horizontal {
		return Vertical
	} else {
		return Horizontal
	}
}

func (o Orientation) Major(x, y int) int {
	if o == Horizontal {
		return x
	} else {
		return y
	}
}

func (o Orientation) Minor(x, y int) int {
	if o == Horizontal {
		return y
	} else {
		return x
	}
}
