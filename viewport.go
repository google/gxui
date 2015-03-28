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

	// SizePixels returns the size of the viewport in pixels.
	SizePixels() math.Size

	// Scale returns the display scaling for this window.
	// A scale of 1 is unscaled, 2 is twice the regular scaling.
	Scale() float32

	// SetScale alters the display scaling for this window.
	// A scale of 1 is unscaled, 2 is twice the regular scaling.
	SetScale(float32)

	Title() string
	SetTitle(string)
	Show()
	Hide()
	Close()
	SetCanvas(Canvas)
	OnClose(func()) EventSubscription
	OnResize(func()) EventSubscription
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
