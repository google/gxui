// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package basic

import (
	"github.com/google/gxui"
	"github.com/google/gxui/math"
	"github.com/google/gxui/mixins"
)

type PanelTab struct {
	mixins.Button
	theme  *Theme
	active bool
}

func CreatePanelTab(theme *Theme) mixins.PanelTab {
	t := &PanelTab{}
	t.Button.Init(t, theme)
	t.theme = theme
	t.SetPadding(math.Spacing{L: 5, T: 3, R: 5, B: 3})
	t.OnMouseEnter(func(gxui.MouseEvent) { t.Redraw() })
	t.OnMouseExit(func(gxui.MouseEvent) { t.Redraw() })
	t.OnMouseDown(func(gxui.MouseEvent) { t.Redraw() })
	t.OnMouseUp(func(gxui.MouseEvent) { t.Redraw() })
	t.OnGainedFocus(t.Redraw)
	t.OnLostFocus(t.Redraw)
	return t
}

func (t *PanelTab) SetActive(active bool) {
	t.active = active
	t.Redraw()
}

func (t *PanelTab) Paint(c gxui.Canvas) {
	s := t.Size()
	var style Style
	switch {
	case t.IsMouseDown(gxui.MouseButtonLeft) && t.IsMouseOver():
		style = t.theme.TabPressedStyle
	case t.IsMouseOver():
		style = t.theme.TabOverStyle
	default:
		style = t.theme.TabDefaultStyle
	}
	if l := t.Label(); l != nil {
		l.SetColor(style.FontColor)
	}

	c.DrawRoundedRect(s.Rect(), 5.0, 5.0, 0.0, 0.0, style.Pen, style.Brush)

	if t.HasFocus() {
		style = t.theme.FocusedStyle
		r := math.CreateRect(1, 1, s.W-1, s.H-1)
		c.DrawRoundedRect(r, 4.0, 4.0, 0.0, 0.0, style.Pen, style.Brush)
	}

	if t.active {
		style = t.theme.TabActiveHighlightStyle
		r := math.CreateRect(1, s.H-1, s.W-1, s.H)
		c.DrawRect(r, style.Brush)
	}

	t.Button.Paint(c)
}
