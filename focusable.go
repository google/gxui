// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

// Focusable is the optional interface implmented by controls that have the
// ability to acquire focus. A control with focus will receive keyboard input
// first.
type Focusable interface {
	Control

	// IsFocusable returns true if the control is currently in a state where it
	// can acquire focus.
	IsFocusable() bool

	// HasFocus returns true when the control has focus.
	HasFocus() bool

	// GainedFocus is called when the Focusable gains focus.
	// This method is called by the FocusManager should not be called by the user.
	GainedFocus()

	// LostFocus is called when the Focusable loses focus.
	// This method is called by the FocusManager should not be called by the user.
	LostFocus()

	// OnGainedFocus subscribes f to be called whenever the control gains focus.
	OnGainedFocus(f func()) EventSubscription

	// OnLostFocus subscribes f to be called whenever the control loses focus.
	OnLostFocus(f func()) EventSubscription
}
