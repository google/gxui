// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package basic

import (
	"github.com/google/gxui"
	"github.com/google/gxui/mixins"
)

type Window struct {
	mixins.Window
}

func CreateWindow(theme *Theme, width, height int, title string) gxui.Window {
	w := &Window{}
	w.Window.Init(w, theme.Driver(), width, height, title)
	w.SetBackgroundBrush(gxui.CreateBrush(theme.WindowBackground))
	return w
}
