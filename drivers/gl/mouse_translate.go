// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"

	"github.com/google/gxui"

	"github.com/goxjs/glfw"
)

func translateMouseButton(button glfw.MouseButton) gxui.MouseButton {
	switch button {
	case glfw.MouseButtonLeft:
		return gxui.MouseButtonLeft
	case glfw.MouseButtonMiddle:
		return gxui.MouseButtonMiddle
	case glfw.MouseButtonRight:
		return gxui.MouseButtonRight
	default:
		panic(fmt.Errorf("Unknown mouse button %v", button))
	}
}

func getMouseState(w *glfw.Window) gxui.MouseState {
	var s gxui.MouseState
	for _, button := range []glfw.MouseButton{glfw.MouseButtonLeft, glfw.MouseButtonMiddle, glfw.MouseButtonRight} {
		if w.GetMouseButton(button) == glfw.Press {
			s |= 1 << uint(translateMouseButton(button))
		}
	}
	return s
}
