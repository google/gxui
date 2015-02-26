// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dark

import (
	"gaze/gxui"
	"gaze/gxui/mixins"
)

type Window struct {
	mixins.Window
	theme *Theme
}

func CreateWindow(theme *Theme, width, height int, title string) gxui.Window {
	w := &Window{}
	w.Window.Init(w, theme.Driver(), width, height, title)
	w.theme = theme
	return w
}

func (w *Window) Paint(c gxui.Canvas) {
	c.Clear(w.theme.WindowBackground)
	w.PaintChildren.Paint(c)
}
