// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dark

import (
	"gaze/gxui"
	"gaze/gxui/mixins"
)

func CreateScrollLayout(theme *Theme) gxui.ScrollLayout {
	l := &mixins.ScrollLayout{}
	l.Init(l, theme)
	return l
}
