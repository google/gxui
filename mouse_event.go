// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/google/gxui/math"
)

type MouseEvent struct {
	Button           MouseButton
	State            MouseState
	Point            math.Point // Local to the event receiver
	WindowPoint      math.Point
	Window           Window
	ScrollX, ScrollY int
	Modifier         KeyboardModifier
}
