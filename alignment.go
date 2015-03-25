// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type HorizontalAlignment int

const (
	AlignLeft HorizontalAlignment = iota
	AlignCenter
	AlignRight
)

func (a HorizontalAlignment) AlignLeft() bool   { return a == AlignLeft }
func (a HorizontalAlignment) AlignCenter() bool { return a == AlignCenter }
func (a HorizontalAlignment) AlignRight() bool  { return a == AlignRight }

type VerticalAlignment int

const (
	AlignTop VerticalAlignment = iota
	AlignMiddle
	AlignBottom
)

func (a VerticalAlignment) AlignTop() bool    { return a == AlignTop }
func (a VerticalAlignment) AlignMiddle() bool { return a == AlignMiddle }
func (a VerticalAlignment) AlignBottom() bool { return a == AlignBottom }
