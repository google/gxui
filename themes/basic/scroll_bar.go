// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package basic

import (
	"github.com/google/gxui"
	"github.com/google/gxui/mixins"
)

type ScrollBar struct {
	mixins.ScrollBar
}

func CreateScrollBar(theme *Theme) gxui.ScrollBar {
	s := &ScrollBar{}
	s.ScrollBar.Init(s, theme)
	s.SetBarBrush(theme.ScrollBarBarDefaultStyle.Brush)
	s.SetBarPen(theme.ScrollBarBarDefaultStyle.Pen)
	s.SetRailBrush(theme.ScrollBarRailDefaultStyle.Brush)
	s.SetRailPen(theme.ScrollBarRailDefaultStyle.Pen)
	updateColors := func() {
		switch {
		case s.IsMouseOver():
			s.SetBarBrush(theme.ScrollBarBarOverStyle.Brush)
			s.SetBarPen(theme.ScrollBarBarOverStyle.Pen)
			s.SetRailBrush(theme.ScrollBarRailOverStyle.Brush)
			s.SetRailPen(theme.ScrollBarRailOverStyle.Pen)
		default:
			s.SetBarBrush(theme.ScrollBarBarDefaultStyle.Brush)
			s.SetBarPen(theme.ScrollBarBarDefaultStyle.Pen)
			s.SetRailBrush(theme.ScrollBarRailDefaultStyle.Brush)
			s.SetRailPen(theme.ScrollBarRailDefaultStyle.Pen)
		}
		s.Redraw()
	}
	s.OnMouseEnter(func(gxui.MouseEvent) { updateColors() })
	s.OnMouseExit(func(gxui.MouseEvent) { updateColors() })
	return s
}
