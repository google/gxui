// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package basic

import (
	"github.com/google/gxui"
	"github.com/google/gxui/math"
	"github.com/google/gxui/mixins"
)

type BubbleOverlay struct {
	mixins.BubbleOverlay
	theme *Theme
}

func CreateBubbleOverlay(theme *Theme) gxui.BubbleOverlay {
	b := &BubbleOverlay{}
	b.Init(b, theme)
	b.SetMargin(math.Spacing{L: 3, T: 3, R: 3, B: 3})
	b.SetPadding(math.Spacing{L: 5, T: 5, R: 5, B: 5})
	b.SetPen(theme.BubbleOverlayStyle.Pen)
	b.SetBrush(theme.BubbleOverlayStyle.Brush)
	b.theme = theme
	return b
}
