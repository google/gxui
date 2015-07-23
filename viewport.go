// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/google/gxui/math"
)

type Viewport interface {
	// SizeDips returns the size of the viewport in device-independent pixels.
	// The ratio of pixels to DIPs is based on the screen density and scale
	// adjustments made with the SetScale method.
	SizeDips() math.Size

	// SetSizeDips sets the size of the viewport in device-independent pixels.
	// The ratio of pixels to DIPs is based on the screen density and scale
	// adjustments made with the SetScale method.
	SetSizeDips(math.Size)

	// SizePixels returns the size of the viewport in pixels.
	SizePixels() math.Size

	// Scale returns the display scaling for this viewport.
	// A scale of 1 is unscaled, 2 is twice the regular scaling.
	Scale() float32

	// SetScale alters the display scaling for this viewport.
	// A scale of 1 is unscaled, 2 is twice the regular scaling.
	SetScale(float32)

	// Fullscreen returns true if the viewport was created full-screen.
	Fullscreen() bool

	// Title returns the title of the window.
	// This is usually the text displayed at the top of the window.
	Title() string

	// SetTitle changes the title of the window.
	SetTitle(string)

	// Position returns position of the window.
	Position() math.Point

	// SetPosition changes position of the window.
	SetPosition(math.Point)

	// Show makes the window visible.
	Show()

	// Hide makes the window invisible.
	Hide()

	// Close destroys the window.
	// Once the window is closed, no further calls should be made to it.
	Close()

	// SetCanvas changes the displayed content of the viewport to the specified
	// Canvas. As canvases are immutable once completed, every visual update of a
	// viewport will require a call to SetCanvas.
	SetCanvas(Canvas)

	// OnClose subscribes f to be called when the viewport closes.
	OnClose(f func()) EventSubscription

	// OnResize subscribes f to be called whenever the viewport changes size.
	OnResize(f func()) EventSubscription

	// OnMouseMove subscribes f to be called whenever the mouse cursor moves over
	// the viewport.
	OnMouseMove(f func(MouseEvent)) EventSubscription

	// OnMouseEnter subscribes f to be called whenever the mouse cursor enters the
	// viewport.
	OnMouseEnter(f func(MouseEvent)) EventSubscription

	// OnMouseEnter subscribes f to be called whenever the mouse cursor leaves the
	// viewport.
	OnMouseExit(f func(MouseEvent)) EventSubscription

	// OnMouseDown subscribes f to be called whenever a mouse button is pressed
	// while the cursor is inside the viewport.
	OnMouseDown(f func(MouseEvent)) EventSubscription

	// OnMouseUp subscribes f to be called whenever a mouse button is released
	// while the cursor is inside the viewport.
	OnMouseUp(f func(MouseEvent)) EventSubscription

	// OnMouseScroll subscribes f to be called whenever the mouse scroll wheel
	// turns while the cursor is inside the viewport.
	OnMouseScroll(f func(MouseEvent)) EventSubscription

	// OnKeyDown subscribes f to be called whenever a keyboard key is pressed
	// while the viewport has focus.
	OnKeyDown(f func(KeyboardEvent)) EventSubscription

	// OnKeyUp subscribes f to be called whenever a keyboard key is released
	// while the viewport has focus.
	OnKeyUp(f func(KeyboardEvent)) EventSubscription

	// OnKeyRepeat subscribes f to be called whenever a keyboard key-repeat event
	// is raised while the viewport has focus.
	OnKeyRepeat(f func(KeyboardEvent)) EventSubscription

	// OnKeyStroke subscribes f to be called whenever a keyboard key-stroke event
	// is raised while the viewport has focus.
	OnKeyStroke(f func(KeyStrokeEvent)) EventSubscription
}
