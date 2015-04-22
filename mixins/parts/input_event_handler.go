// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parts

import (
	"github.com/google/gxui"
)

type InputEventHandlerOuter interface{}

type InputEventHandler struct {
	outer         InputEventHandlerOuter
	isMouseOver   bool
	isMouseDown   map[gxui.MouseButton]bool
	onClick       gxui.Event
	onDoubleClick gxui.Event
	onKeyPress    gxui.Event
	onKeyStroke   gxui.Event
	onMouseMove   gxui.Event
	onMouseEnter  gxui.Event
	onMouseExit   gxui.Event
	onMouseDown   gxui.Event
	onMouseUp     gxui.Event
	onMouseScroll gxui.Event
	onKeyDown     gxui.Event
	onKeyUp       gxui.Event
	onKeyRepeat   gxui.Event
}

func (m *InputEventHandler) getOnClick() gxui.Event {
	if m.onClick == nil {
		m.onClick = gxui.CreateEvent(m.Click)
	}
	return m.onClick
}

func (m *InputEventHandler) getOnDoubleClick() gxui.Event {
	if m.onDoubleClick == nil {
		m.onDoubleClick = gxui.CreateEvent(m.DoubleClick)
	}
	return m.onDoubleClick
}

func (m *InputEventHandler) getOnKeyPress() gxui.Event {
	if m.onKeyPress == nil {
		m.onKeyPress = gxui.CreateEvent(m.KeyPress)
	}
	return m.onKeyPress
}

func (m *InputEventHandler) getOnKeyStroke() gxui.Event {
	if m.onKeyStroke == nil {
		m.onKeyStroke = gxui.CreateEvent(m.KeyStroke)
	}
	return m.onKeyStroke
}

func (m *InputEventHandler) getOnMouseMove() gxui.Event {
	if m.onMouseMove == nil {
		m.onMouseMove = gxui.CreateEvent(m.MouseMove)
	}
	return m.onMouseMove
}

func (m *InputEventHandler) getOnMouseEnter() gxui.Event {
	if m.onMouseEnter == nil {
		m.onMouseEnter = gxui.CreateEvent(m.MouseEnter)
	}
	return m.onMouseEnter
}

func (m *InputEventHandler) getOnMouseExit() gxui.Event {
	if m.onMouseExit == nil {
		m.onMouseExit = gxui.CreateEvent(m.MouseExit)
	}
	return m.onMouseExit
}

func (m *InputEventHandler) getOnMouseDown() gxui.Event {
	if m.onMouseDown == nil {
		m.onMouseDown = gxui.CreateEvent(m.MouseDown)
	}
	return m.onMouseDown
}

func (m *InputEventHandler) getOnMouseUp() gxui.Event {
	if m.onMouseUp == nil {
		m.onMouseUp = gxui.CreateEvent(m.MouseUp)
	}
	return m.onMouseUp
}

func (m *InputEventHandler) getOnMouseScroll() gxui.Event {
	if m.onMouseScroll == nil {
		m.onMouseScroll = gxui.CreateEvent(m.MouseScroll)
	}
	return m.onMouseScroll
}

func (m *InputEventHandler) getOnKeyDown() gxui.Event {
	if m.onKeyDown == nil {
		m.onKeyDown = gxui.CreateEvent(m.KeyDown)
	}
	return m.onKeyDown
}

func (m *InputEventHandler) getOnKeyUp() gxui.Event {
	if m.onKeyUp == nil {
		m.onKeyUp = gxui.CreateEvent(m.KeyUp)
	}
	return m.onKeyUp
}

func (m *InputEventHandler) getOnKeyRepeat() gxui.Event {
	if m.onKeyRepeat == nil {
		m.onKeyRepeat = gxui.CreateEvent(m.KeyRepeat)
	}
	return m.onKeyRepeat
}

func (m *InputEventHandler) Init(outer InputEventHandlerOuter) {
	m.outer = outer
	m.isMouseDown = make(map[gxui.MouseButton]bool)
}

func (m *InputEventHandler) Click(ev gxui.MouseEvent) (consume bool) {
	m.getOnClick().Fire(ev)
	return false
}

func (m *InputEventHandler) DoubleClick(ev gxui.MouseEvent) (consume bool) {
	m.getOnDoubleClick().Fire(ev)
	return false
}

func (m *InputEventHandler) KeyPress(ev gxui.KeyboardEvent) (consume bool) {
	m.getOnKeyPress().Fire(ev)
	return false
}

func (m *InputEventHandler) KeyStroke(ev gxui.KeyStrokeEvent) (consume bool) {
	m.getOnKeyStroke().Fire(ev)
	return false
}

func (m *InputEventHandler) MouseScroll(ev gxui.MouseEvent) (consume bool) {
	m.getOnMouseScroll().Fire(ev)
	return false
}

func (m *InputEventHandler) MouseMove(ev gxui.MouseEvent) {
	m.getOnMouseMove().Fire(ev)
}

func (m *InputEventHandler) MouseEnter(ev gxui.MouseEvent) {
	m.isMouseOver = true
	m.getOnMouseEnter().Fire(ev)
}

func (m *InputEventHandler) MouseExit(ev gxui.MouseEvent) {
	m.isMouseOver = false
	m.getOnMouseExit().Fire(ev)
}

func (m *InputEventHandler) MouseDown(ev gxui.MouseEvent) {
	m.isMouseDown[ev.Button] = true
	m.getOnMouseDown().Fire(ev)
}

func (m *InputEventHandler) MouseUp(ev gxui.MouseEvent) {
	m.isMouseDown[ev.Button] = false
	m.getOnMouseUp().Fire(ev)
}

func (m *InputEventHandler) KeyDown(ev gxui.KeyboardEvent) {
	m.getOnKeyDown().Fire(ev)
}

func (m *InputEventHandler) KeyUp(ev gxui.KeyboardEvent) {
	m.getOnKeyUp().Fire(ev)
}

func (m *InputEventHandler) KeyRepeat(ev gxui.KeyboardEvent) {
	m.getOnKeyRepeat().Fire(ev)
}

func (m *InputEventHandler) OnClick(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return m.getOnClick().Listen(f)
}

func (m *InputEventHandler) OnDoubleClick(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return m.getOnDoubleClick().Listen(f)
}

func (m *InputEventHandler) OnKeyPress(f func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return m.getOnKeyPress().Listen(f)
}

func (m *InputEventHandler) OnKeyStroke(f func(gxui.KeyStrokeEvent)) gxui.EventSubscription {
	return m.getOnKeyStroke().Listen(f)
}

func (m *InputEventHandler) OnMouseMove(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return m.getOnMouseMove().Listen(f)
}

func (m *InputEventHandler) OnMouseEnter(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return m.getOnMouseEnter().Listen(f)
}

func (m *InputEventHandler) OnMouseExit(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return m.getOnMouseExit().Listen(f)
}

func (m *InputEventHandler) OnMouseDown(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return m.getOnMouseDown().Listen(f)
}

func (m *InputEventHandler) OnMouseUp(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return m.getOnMouseUp().Listen(f)
}

func (m *InputEventHandler) OnMouseScroll(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return m.getOnMouseScroll().Listen(f)
}

func (m *InputEventHandler) OnKeyDown(f func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return m.getOnKeyDown().Listen(f)
}

func (m *InputEventHandler) OnKeyUp(f func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return m.getOnKeyUp().Listen(f)
}

func (m *InputEventHandler) OnKeyRepeat(f func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return m.getOnKeyRepeat().Listen(f)
}

func (m *InputEventHandler) IsMouseOver() bool {
	return m.isMouseOver
}

func (m *InputEventHandler) IsMouseDown(button gxui.MouseButton) bool {
	return m.isMouseDown[button]
}
