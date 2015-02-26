// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"time"
)

var doubleClickTime = time.Millisecond * 300

type MouseController struct {
	window          Window
	focusController *FocusController
	lastOver        ControlPointList
	lastDown        map[MouseButton]ControlPointList
	lastUpTime      map[MouseButton]time.Time
}

func CreateMouseController(w Window, focusController *FocusController) *MouseController {
	c := &MouseController{
		window:          w,
		focusController: focusController,
		lastDown:        make(map[MouseButton]ControlPointList),
		lastUpTime:      make(map[MouseButton]time.Time),
	}
	w.OnMouseMove(c.mouseMove)
	w.OnMouseEnter(c.mouseMove)
	w.OnMouseExit(c.mouseMove)
	w.OnMouseDown(c.mouseDown)
	w.OnMouseUp(c.mouseUp)
	w.OnMouseScroll(c.mouseScroll)
	return c
}

func (m *MouseController) updatePosition(ev MouseEvent) {
	ValidateHierarchy(m.window)

	nowOver := TopControlsUnder(ev.Point, m.window)

	for _, cp := range m.lastOver {
		if !nowOver.Contains(cp.C) {
			e := ev
			e.Point = cp.P
			cp.C.MouseExit(e)
		}
	}

	for _, cp := range nowOver {
		if !m.lastOver.Contains(cp.C) {
			e := ev
			e.Point = cp.P
			cp.C.MouseEnter(e)
		}
	}

	m.lastOver = nowOver
}

func (m *MouseController) mouseMove(ev MouseEvent) {
	m.updatePosition(ev)
	for _, cp := range m.lastOver {
		e := ev
		e.Point = cp.P
		cp.C.MouseMove(e)
	}
}

func (m *MouseController) mouseDown(ev MouseEvent) {
	m.updatePosition(ev)

	for _, cp := range m.lastOver {
		e := ev
		e.Point = cp.P
		cp.C.MouseDown(e)
	}

	m.lastDown[ev.Button] = m.lastOver
}

func (m *MouseController) mouseUp(ev MouseEvent) {
	m.updatePosition(ev)

	for _, cp := range m.lastDown[ev.Button] {
		e := ev
		e.Point = cp.P
		cp.C.MouseUp(e)
	}

	setFocusCount := m.focusController.SetFocusCount()

	dblClick := time.Since(m.lastUpTime[ev.Button]) < doubleClickTime
	clickConsumed := false
	for i := len(m.lastDown[ev.Button]) - 1; i >= 0; i-- {
		cp := m.lastDown[ev.Button][i]
		if p, found := m.lastOver.Find(cp.C); found {
			ev.Point = p
			if (dblClick && cp.C.DoubleClick(ev)) || (!dblClick && cp.C.Click(ev)) {
				clickConsumed = true
				break
			}
		}
	}

	if !clickConsumed {
		ev.Point = ev.WindowPoint
		if dblClick {
			m.window.DoubleClick(ev)
		} else {
			m.window.Click(ev)
		}
	}

	focusSet := setFocusCount != m.focusController.SetFocusCount()
	if !focusSet {
		for i := len(m.lastDown[ev.Button]) - 1; i >= 0; i-- {
			cp := m.lastDown[ev.Button][i]
			if m.lastOver.Contains(cp.C) && m.window.SetFocus(cp.C) {
				focusSet = true
				break
			}
		}

		if !focusSet {
			m.window.SetFocus(nil)
		}
	}

	delete(m.lastDown, ev.Button)
	m.lastUpTime[ev.Button] = time.Now()
}

func (m *MouseController) mouseScroll(ev MouseEvent) {
	m.updatePosition(ev)

	for i := len(m.lastOver) - 1; i >= 0; i-- {
		cp := m.lastOver[i]
		e := ev
		e.Point = cp.P
		cp.C.MouseScroll(e)
	}
}
