// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/google/gxui/math"
)

type Window interface {
	Container

	// Title returns the title of the window.
	// This is usually the text displayed at the top of the window.
	Title() string

	// SetTitle changes the title of the window.
	SetTitle(string)

	// Scale returns the display scaling for this window.
	// A scale of 1 is unscaled, 2 is twice the regular scaling.
	Scale() float32

	// SetScale alters the display scaling for this window.
	// A scale of 1 is unscaled, 2 is twice the regular scaling.
	SetScale(float32)

	// Position returns position of the window.
	Position() math.Point

	// SetPosition changes position of the window.
	SetPosition(math.Point)

	// Fullscreen returns true if the window is currently full-screen.
	Fullscreen() bool

	// SetFullscreen makes the window either full-screen or windowed.
	SetFullscreen(bool)

	// Show makes the window visible.
	Show()

	// Hide makes the window invisible.
	Hide()

	// Close destroys the window.
	// Once the window is closed, no further calls should be made to it.
	Close()

	// Focus returns the control currently with focus.
	Focus() Focusable

	// SetFocus gives the specified control Focus, returning true on success or
	// false if the control cannot be given focus.
	SetFocus(Control) bool

	// BackgroundBrush returns the brush used to draw the window background.
	BackgroundBrush() Brush

	// SetBackgroundBrush sets the brush used to draw the window background.
	SetBackgroundBrush(Brush)

	// BorderPen returns the pen used to draw the window border.
	BorderPen() Pen

	// SetBorderPen sets the pen used to draw the window border.
	SetBorderPen(Pen)

	Click(MouseEvent)
	DoubleClick(MouseEvent)
	KeyPress(KeyboardEvent)
	KeyStroke(KeyStrokeEvent)

	// Events
	OnClose(func()) EventSubscription
	OnResize(func()) EventSubscription
	OnClick(func(MouseEvent)) EventSubscription
	OnDoubleClick(func(MouseEvent)) EventSubscription
	OnMouseMove(func(MouseEvent)) EventSubscription
	OnMouseEnter(func(MouseEvent)) EventSubscription
	OnMouseExit(func(MouseEvent)) EventSubscription
	OnMouseDown(func(MouseEvent)) EventSubscription
	OnMouseUp(func(MouseEvent)) EventSubscription
	OnMouseScroll(func(MouseEvent)) EventSubscription
	OnKeyDown(func(KeyboardEvent)) EventSubscription
	OnKeyUp(func(KeyboardEvent)) EventSubscription
	OnKeyRepeat(func(KeyboardEvent)) EventSubscription
	OnKeyStroke(func(KeyStrokeEvent)) EventSubscription
}
