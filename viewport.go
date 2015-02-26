// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"gxui/math"
)

type Viewport interface {
	SizeDips() math.Size
	SizePixels() math.Size
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

	SetContinuousRedraw(bool)
	ContinuousRedraw() bool
	Stats() string
}
