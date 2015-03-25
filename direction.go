// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import "fmt"

type Direction int

const (
	TopToBottom Direction = iota
	LeftToRight
	BottomToTop
	RightToLeft
)

func (d Direction) LeftToRight() bool { return d == LeftToRight }
func (d Direction) RightToLeft() bool { return d == RightToLeft }
func (d Direction) TopToBottom() bool { return d == TopToBottom }
func (d Direction) BottomToTop() bool { return d == BottomToTop }

func (d Direction) Flip() Direction {
	switch d {
	case TopToBottom:
		return BottomToTop
	case LeftToRight:
		return RightToLeft
	case BottomToTop:
		return TopToBottom
	case RightToLeft:
		return LeftToRight
	default:
		panic(fmt.Errorf("Unknown direction %d", d))
	}
}

func (d Direction) Orientation() Orientation {
	switch d {
	case TopToBottom, BottomToTop:
		return Vertical
	case LeftToRight, RightToLeft:
		return Horizontal
	default:
		panic(fmt.Errorf("Unknown direction %d", d))
	}
}
