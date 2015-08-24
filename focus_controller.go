// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type FocusController struct {
	window             Window
	focus              Focusable
	setFocusCount      int
	detachSubscription EventSubscription
}

func CreateFocusController(window Window) *FocusController {
	return &FocusController{
		window: window,
	}
}

func (c *FocusController) SetFocus(f Focusable) {
	c.setFocusCount++
	if c.focus == f {
		return
	}
	if c.focus != nil {
		o := c.focus
		c.focus = nil
		c.detachSubscription.Unlisten()
		o.LostFocus()
		if c.focus != nil {
			return // Something in LostFocus() called SetFocus(). Respect their call.
		}
	}
	c.focus = f
	if c.focus != nil {
		c.detachSubscription = c.focus.OnDetach(func() { c.SetFocus(nil) })
		c.focus.GainedFocus()
	}
}

func (c *FocusController) SetFocusCount() int {
	return c.setFocusCount
}

func (c *FocusController) Focus() Focusable {
	return c.focus
}

func (c *FocusController) FocusNext() {
	c.SetFocus(c.NextFocusable(c.focus, true))
}

func (c *FocusController) FocusPrev() {
	c.SetFocus(c.NextFocusable(c.focus, false))
}

func (c *FocusController) NextFocusable(after Control, forwards bool) Focusable {
	container, _ := after.(Container)
	if container != nil {
		f := c.NextChildFocusable(container, nil, forwards)
		if f != nil {
			return f
		}
	}

	for after != nil {
		parent := after.Parent()
		if parent != nil {
			f := c.NextChildFocusable(parent, after, forwards)
			if f != nil {
				return f
			}
		}
		after, _ = parent.(Control)
	}

	return c.NextChildFocusable(c.window, nil, forwards)
}

func (c *FocusController) NextChildFocusable(p Parent, after Control, forwards bool) Focusable {
	examineNext := after == nil
	children := p.Children()

	i := 0
	e := len(children)
	if !forwards {
		i = len(children) - 1
		e = -1
	}

	for i != e {
		f := children[i]
		if forwards {
			i++
		} else {
			i--
		}

		if !examineNext {
			if f.Control == after {
				examineNext = true
			}
			continue
		}

		if focusable := c.Focusable(f.Control); focusable != nil {
			return focusable
		}

		if container, ok := f.Control.(Container); ok {
			focusable := c.NextChildFocusable(container, nil, forwards)
			if focusable != nil {
				return focusable
			}
		}
	}
	return nil
}

func (c *FocusController) Focusable(ctrl Control) Focusable {
	focusable, _ := ctrl.(Focusable)
	if focusable != nil && focusable.IsFocusable() {
		return focusable
	}
	return nil
}
