// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import "fmt"

type HorizontalAlignment int

const (
	AlignLeft HorizontalAlignment = iota
	AlignCenter
	AlignRight
)

func (a HorizontalAlignment) AlignLeft() bool   { return a == AlignLeft }
func (a HorizontalAlignment) AlignCenter() bool { return a == AlignCenter }
func (a HorizontalAlignment) AlignRight() bool  { return a == AlignRight }

func (a HorizontalAlignment) String() string {
	switch a {
	case AlignLeft:
		return "Align Left"
	case AlignCenter:
		return "Align Center"
	case AlignRight:
		return "Align Right"
	default:
		return fmt.Sprintf("HorizontalAlignment<%d>", a)
	}
}

type VerticalAlignment int

const (
	AlignTop VerticalAlignment = iota
	AlignMiddle
	AlignBottom
)

func (a VerticalAlignment) AlignTop() bool    { return a == AlignTop }
func (a VerticalAlignment) AlignMiddle() bool { return a == AlignMiddle }
func (a VerticalAlignment) AlignBottom() bool { return a == AlignBottom }

func (a VerticalAlignment) String() string {
	switch a {
	case AlignTop:
		return "Align Top"
	case AlignMiddle:
		return "Align Middle"
	case AlignBottom:
		return "Align Bottom"
	default:
		return fmt.Sprintf("VerticalAlignment<%d>", a)
	}
}
