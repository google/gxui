// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/google/gxui/math"
)

type MouseEvent struct {
	Button           MouseButton
	Point            math.Point // Local to the event receiver
	WindowPoint      math.Point
	Window           Window
	ScrollX, ScrollY int
	Modifier         KeyboardModifier
}

func (ev MouseEvent) IsLeftDown() bool {
	return ev.Button&MouseButtonLeft != 0
}

func (ev MouseEvent) IsMiddleDown() bool {
	return ev.Button&MouseButtonMiddle != 0
}

func (ev MouseEvent) IsRightDown() bool {
	return ev.Button&MouseButtonRight != 0
}
