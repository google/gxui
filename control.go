// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import "github.com/google/gxui/math"

// Control is the interface exposed by all UI control elements.
type Control interface {
	// Size returns the size of the control. If the control is not attached, then
	// the returned size is undefined.
	Size() math.Size

	// SetSize sets the size of the control to the specified value.
	// SetSize should only be called by the parent of the control during layout.
	SetSize(math.Size)

	// Draw draws the control's visual apperance into the returned, new canvas.
	// Draw is typically called by the parent of the control - calling Draw will
	// not issue a re-draw of an attached control.
	Draw() Canvas

	// Parent returns the parent of the control.
	Parent() Parent

	// SetParent sets the parent of the control.
	// SetParent should only be called by the new parent of the control.
	SetParent(Parent)

	// Attached returns true if the control is directly or indirectly attached
	// to a window.
	Attached() bool

	// Attach is called when the control is directly or indirectly attached to a
	// window.
	// Attach should only be called by the parent of the control.
	Attach()

	// Detach is called when the control is directly or indirectly detached from a
	// window.
	// Detach should only be called by the parent of the control.
	Detach()

	// DesiredSize returns the desired size of the control based on the min and
	// max size limits. The parent control may ignore the desired size.
	DesiredSize(min, max math.Size) math.Size

	// Margin returns the desired spacing between sibling controls.
	Margin() math.Spacing

	// SetMargin set the desired spacing between sibling controls, issuing a
	// relayout if the margin has changed.
	SetMargin(math.Spacing)

	// IsVisible returns true if the control is visible.
	IsVisible() bool

	// SetVisible sets the visibility of the control.
	SetVisible(bool)

	// ContainsPoint returns true if the specified local-space point is considered
	// within the control.
	ContainsPoint(math.Point) bool

	// IsMouseOver returns true if the mouse cursor was last reported within the
	// control.
	IsMouseOver() bool

	// IsMouseDown returns true if button was last reported pressed on the
	// control.
	IsMouseDown(button MouseButton) bool

	// Click is called when the mouse is pressed and released on the control.
	// If Click returns true, then the click event is consumed by the control,
	// otherwise the next control below the should be considered for the click
	// event.
	Click(MouseEvent) (consume bool)

	// DoubleClick is called when the mouse is double-clicked on the control.
	// If DoubleClick returns true, then the double-click event is consumed by the
	// control, otherwise the next control below the should be considered for the
	// double-click event.
	DoubleClick(MouseEvent) (consume bool)

	// KeyPress is called when a keyboard key is pressed while the control (or
	// non-consuming child) has focus. If KeyPress returns true, then the
	// key-press event is consumed by the control, otherwise the parent control
	// should be considered for the key-press event.
	KeyPress(KeyboardEvent) (consume bool)

	// KeyStroke is called when a key-storke is made while the control (or
	// non-consuming child) has focus. If KeyStroke returns true, then the
	// key-stroke event is consumed by the control, otherwise the parent control
	// should be considered for the key-stroke event.
	KeyStroke(KeyStrokeEvent) (consume bool)

	// MouseScroll is called when a mouse scroll is made while the control (or
	// non-consuming child) has focus. If MouseScroll returns true, then the
	// mouse-scroll event is consumed by the control, otherwise the parent control
	// should be considered for the key-stroke event.
	MouseScroll(MouseEvent) (consume bool)

	// MouseMove is called when the mouse cursor moves over the control.
	MouseMove(MouseEvent)

	// MouseEnter is called when the mouse cursor transitions from outside to
	// inside the bounds of the control.
	MouseEnter(MouseEvent)

	// MouseExit is called when the mouse cursor transitions from inside to
	// outside the bounds of the control.
	MouseExit(MouseEvent)

	// MouseDown is called when a mouse button is pressed while the mouse cursor
	// is over the control.
	MouseDown(MouseEvent)

	// MouseUp is called when a mouse button is released while the mouse cursor
	// is over the control.
	MouseUp(MouseEvent)

	// KeyDown is called when a keyboard button is pressed while the control (or
	// child control) has focus.
	KeyDown(KeyboardEvent)

	// KeyUp is called when a keyboard button is released while the control (or
	// child control) has focus.
	KeyUp(KeyboardEvent)

	// KeyRepeat is called when a keyboard button held long enough for a
	// repeat-key event while the control (or child control) has focus.
	KeyRepeat(KeyboardEvent)

	// OnAttach subscribes f to be called whenever the control is attached.
	OnAttach(f func()) EventSubscription

	// OnDetach subscribes f to be called whenever the control is detached.
	OnDetach(f func()) EventSubscription

	// OnKeyPress subscribes f to be called whenever the control receives a
	// key-press event.
	OnKeyPress(f func(KeyboardEvent)) EventSubscription

	// OnKeyStroke subscribes f to be called whenever the control receives a
	// key-stroke event.
	OnKeyStroke(f func(KeyStrokeEvent)) EventSubscription

	// OnClick subscribes f to be called whenever the control receives a click
	// event.
	OnClick(f func(MouseEvent)) EventSubscription

	// OnDoubleClick subscribes f to be called whenever the control receives a
	// double-click event.
	OnDoubleClick(f func(MouseEvent)) EventSubscription

	// OnMouseMove subscribes f to be called whenever the control receives a
	// mouse-move event.
	OnMouseMove(f func(MouseEvent)) EventSubscription

	// OnMouseEnter subscribes f to be called whenever the control receives a
	// mouse-enter event.
	OnMouseEnter(f func(MouseEvent)) EventSubscription

	// OnMouseExit subscribes f to be called whenever the control receives a
	// mouse-exit event.
	OnMouseExit(f func(MouseEvent)) EventSubscription

	// OnMouseDown subscribes f to be called whenever the control receives a
	// mouse-down event.
	OnMouseDown(f func(MouseEvent)) EventSubscription

	// OnMouseUp subscribes f to be called whenever the control receives a
	// mouse-up event.
	OnMouseUp(f func(MouseEvent)) EventSubscription

	// OnMouseScroll subscribes f to be called whenever the control receives a
	// mouse-scroll event.
	OnMouseScroll(f func(MouseEvent)) EventSubscription

	// OnKeyDown subscribes f to be called whenever the control receives a
	// key-down event.
	OnKeyDown(f func(KeyboardEvent)) EventSubscription

	// OnKeyUp subscribes f to be called whenever the control receives a
	// key-up event.
	OnKeyUp(f func(KeyboardEvent)) EventSubscription

	// OnKeyRepeat subscribes f to be called whenever the control receives a
	// key-repeat event.
	OnKeyRepeat(f func(KeyboardEvent)) EventSubscription
}
