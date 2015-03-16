// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/google/gxui"
	"github.com/google/gxui/math"
	"github.com/google/gxui/mixins/parts"
)

type ButtonOuter interface {
	LinearLayoutOuter
	IsChecked() bool
	SetChecked(bool)
}

type Button struct {
	LinearLayout
	parts.Focusable

	outer      ButtonOuter
	theme      gxui.Theme
	label      gxui.Label
	buttonType gxui.ButtonType
	checked    bool
}

func (b *Button) Init(outer ButtonOuter, theme gxui.Theme) {
	b.LinearLayout.Init(outer, theme)
	b.Focusable.Init(outer)

	b.buttonType = gxui.PushButton
	b.theme = theme
	b.outer = outer

	// Interface compliance test
	_ = gxui.Button(b)
}

func (b *Button) Label() gxui.Label {
	return b.label
}

func (b *Button) Text() string {
	if b.label != nil {
		return b.label.Text()
	} else {
		return ""
	}
}

func (b *Button) SetText(text string) {
	if b.Text() == text {
		return
	}
	if text == "" {
		if b.label != nil {
			b.RemoveChild(b.label)
			b.label = nil
		}
	} else {
		if b.label == nil {
			b.label = b.theme.CreateLabel()
			b.label.SetMargin(math.ZeroSpacing)
			b.AddChild(b.label)
		}
		b.label.SetText(text)
	}
}

func (b *Button) Type() gxui.ButtonType {
	return b.buttonType
}

func (b *Button) SetType(buttonType gxui.ButtonType) {
	if buttonType != b.buttonType {
		b.buttonType = buttonType
		b.outer.Redraw()
	}
}

func (b *Button) IsChecked() bool {
	return b.checked
}

func (b *Button) SetChecked(checked bool) {
	if checked != b.checked {
		b.checked = checked
		b.outer.Redraw()
	}
}

// InputEventHandler override
func (b *Button) Click(ev gxui.MouseEvent) (consume bool) {
	if ev.Button == gxui.MouseButtonLeft {
		if b.buttonType == gxui.ToggleButton {
			b.outer.SetChecked(!b.outer.IsChecked())
		}
		b.LinearLayout.Click(ev)
		return true
	}
	return b.LinearLayout.Click(ev)
}

func (b *Button) KeyPress(ev gxui.KeyboardEvent) (consume bool) {
	consume = b.LinearLayout.KeyPress(ev)
	if ev.Key == gxui.KeySpace || ev.Key == gxui.KeyEnter {
		me := gxui.MouseEvent{
			Button: gxui.MouseButtonLeft,
		}
		return b.Click(me)
	}
	return
}
