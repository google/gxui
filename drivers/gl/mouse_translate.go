// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"github.com/google/gxui"

	"github.com/go-gl/glfw/v3.1/glfw"
)

func getMouseState(w *glfw.Window) gxui.MouseState {
	var s gxui.MouseState
	for _, button := range []glfw.MouseButton{glfw.MouseButtonLeft, glfw.MouseButtonMiddle, glfw.MouseButtonRight} {
		if w.GetMouseButton(button) == glfw.Press {
			switch button {
			case glfw.MouseButtonLeft:
				s |= 1 << uint(gxui.MouseButtonLeft)
			case glfw.MouseButtonMiddle:
				s |= 1 << uint(gxui.MouseButtonMiddle)
			case glfw.MouseButtonRight:
				s |= 1 << uint(gxui.MouseButtonRight)
			}
		}
	}
	return s
}
