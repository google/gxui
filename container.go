// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import "gaze/gxui/math"

type Parent interface {
	Children() []Control
}

type Container interface {
	Parent
	Relayout()
	Redraw()
	ChildCount() int
	ChildIndex(child Control) int
	ChildAt(int) Control
	AddChild(child Control)
	AddChildAt(index int, child Control)
	RemoveChild(child Control)
	RemoveChildAt(index int)
	RemoveAll()
	Padding() math.Spacing
	SetPadding(math.Spacing)
}
