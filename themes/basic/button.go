// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package basic

import (
	"github.com/google/gxui"
	"github.com/google/gxui/math"
	"github.com/google/gxui/mixins"
)

type Button struct {
	mixins.Button
	theme *Theme
}

func CreateButton(theme *Theme) gxui.Button {
	b := &Button{}
	b.Init(b, theme)
	b.theme = theme
	b.SetPadding(math.Spacing{L: 3, T: 3, R: 3, B: 3})
	b.SetMargin(math.Spacing{L: 3, T: 3, R: 3, B: 3})
	b.SetBackgroundBrush(b.theme.ButtonDefaultStyle.Brush)
	b.SetBorderPen(b.theme.ButtonDefaultStyle.Pen)
	b.OnMouseEnter(func(gxui.MouseEvent) { b.Redraw() })
	b.OnMouseExit(func(gxui.MouseEvent) { b.Redraw() })
	b.OnMouseDown(func(gxui.MouseEvent) { b.Redraw() })
	b.OnMouseUp(func(gxui.MouseEvent) { b.Redraw() })
	b.OnGainedFocus(b.Redraw)
	b.OnLostFocus(b.Redraw)
	return b
}

// Button internal overrides
func (b *Button) Paint(c gxui.Canvas) {
	pen := b.Button.BorderPen()
	brush := b.Button.BackgroundBrush()
	fontColor := b.theme.ButtonDefaultStyle.FontColor

	switch {
	case b.IsMouseDown(gxui.MouseButtonLeft) && b.IsMouseOver():
		pen = b.theme.ButtonPressedStyle.Pen
		brush = b.theme.ButtonPressedStyle.Brush
		fontColor = b.theme.ButtonPressedStyle.FontColor
	case b.IsMouseOver():
		pen = b.theme.ButtonOverStyle.Pen
		brush = b.theme.ButtonOverStyle.Brush
		fontColor = b.theme.ButtonOverStyle.FontColor
	}

	if l := b.Label(); l != nil {
		l.SetColor(fontColor)
	}

	r := b.Size().Rect()

	c.DrawRoundedRect(r, 2, 2, 2, 2, gxui.TransparentPen, brush)

	b.PaintChildren.Paint(c)

	c.DrawRoundedRect(r, 2, 2, 2, 2, pen, gxui.TransparentBrush)

	if b.IsChecked() {
		pen = b.theme.HighlightStyle.Pen
		brush = b.theme.HighlightStyle.Brush
		c.DrawRoundedRect(r, 2.0, 2.0, 2.0, 2.0, pen, brush)
	}

	if b.HasFocus() {
		pen = b.theme.FocusedStyle.Pen
		brush = b.theme.FocusedStyle.Brush
		c.DrawRoundedRect(r.ContractI(int(pen.Width)), 3.0, 3.0, 3.0, 3.0, pen, brush)
	}
}
