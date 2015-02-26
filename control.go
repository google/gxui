// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"gaze/gxui/math"
)

type Control interface {
	Layout(math.Rect)
	Draw() Canvas
	Bounds() math.Rect
	Parent() Container
	SetParent(Container)
	Attached() bool
	Attach()
	Detach()
	DesiredSize(min, max math.Size) math.Size
	Margin() math.Spacing
	SetMargin(math.Spacing)
	IsVisible() bool
	SetVisible(bool)
	ContainsPoint(math.Point) bool
	IsMouseOver() bool
	IsMouseDown(button MouseButton) bool

	Click(MouseEvent) (consume bool)
	DoubleClick(MouseEvent) (consume bool)
	KeyPress(KeyboardEvent) (consume bool)
	KeyStroke(KeyStrokeEvent) (consume bool)
	MouseScroll(MouseEvent) (consume bool)

	MouseMove(MouseEvent)
	MouseEnter(MouseEvent)
	MouseExit(MouseEvent)
	MouseDown(MouseEvent)
	MouseUp(MouseEvent)
	KeyDown(KeyboardEvent)
	KeyUp(KeyboardEvent)
	KeyRepeat(KeyboardEvent)

	OnAttach(func()) EventSubscription
	OnDetach(func()) EventSubscription
	OnClick(func(MouseEvent)) EventSubscription
	OnDoubleClick(func(MouseEvent)) EventSubscription
	OnKeyPress(func(KeyboardEvent)) EventSubscription
	OnKeyStroke(func(KeyStrokeEvent)) EventSubscription
	OnMouseMove(func(MouseEvent)) EventSubscription
	OnMouseEnter(func(MouseEvent)) EventSubscription
	OnMouseExit(func(MouseEvent)) EventSubscription
	OnMouseDown(func(MouseEvent)) EventSubscription
	OnMouseUp(func(MouseEvent)) EventSubscription
	OnMouseScroll(func(MouseEvent)) EventSubscription
	OnKeyDown(func(KeyboardEvent)) EventSubscription
	OnKeyUp(func(KeyboardEvent)) EventSubscription
	OnKeyRepeat(func(KeyboardEvent)) EventSubscription
}
