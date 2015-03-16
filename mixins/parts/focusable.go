// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parts

import (
	"github.com/google/gxui"
)

type FocusableOuter interface{}

type Focusable struct {
	outer         FocusableOuter
	focusable     bool
	hasFocus      bool
	onGainedFocus gxui.Event
	onLostFocus   gxui.Event
}

func (f *Focusable) Init(outer FocusableOuter) {
	f.outer = outer
	f.focusable = true
}

// gxui.Control compliance
func (f *Focusable) IsFocusable() bool {
	return f.focusable
}

func (f *Focusable) HasFocus() bool {
	return f.hasFocus
}

func (f *Focusable) SetFocusable(bool) {
	f.focusable = true
}

func (f *Focusable) OnGainedFocus(l func()) gxui.EventSubscription {
	if f.onGainedFocus == nil {
		f.onGainedFocus = gxui.CreateEvent(f.GainedFocus)
	}
	return f.onGainedFocus.Listen(l)
}

func (f *Focusable) OnLostFocus(l func()) gxui.EventSubscription {
	if f.onLostFocus == nil {
		f.onLostFocus = gxui.CreateEvent(f.LostFocus)
	}
	return f.onLostFocus.Listen(l)
}

func (f *Focusable) GainedFocus() {
	f.hasFocus = true
	if f.onGainedFocus != nil {
		f.onGainedFocus.Fire()
	}
}

func (f *Focusable) LostFocus() {
	f.hasFocus = false
	if f.onLostFocus != nil {
		f.onLostFocus.Fire()
	}
}
