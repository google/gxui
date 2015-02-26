// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type KeyboardController struct {
	window Window
}

func CreateKeyboardController(w Window) *KeyboardController {
	c := &KeyboardController{
		window: w,
	}
	w.OnKeyDown(c.keyDown)
	w.OnKeyUp(c.keyUp)
	w.OnKeyRepeat(c.keyPress)
	w.OnKeyStroke(c.keyStroke)
	return c
}

func (c *KeyboardController) keyDown(ev KeyboardEvent) {
	f := Control(c.window.Focus())
	for f != nil {
		f.KeyDown(ev)
		f, _ = f.Parent().(Control)
	}
	c.keyPress(ev)
}

func (c *KeyboardController) keyUp(ev KeyboardEvent) {
	f := Control(c.window.Focus())
	for f != nil {
		f.KeyUp(ev)
		f, _ = f.Parent().(Control)
	}
}

func (c *KeyboardController) keyPress(ev KeyboardEvent) {
	f := Control(c.window.Focus())
	for f != nil {
		if f.KeyPress(ev) {
			return
		}
		f, _ = f.Parent().(Control)
	}
	c.window.KeyPress(ev)
}

func (c *KeyboardController) keyStroke(ev KeyStrokeEvent) {
	f := Control(c.window.Focus())
	for f != nil {
		if f.KeyStroke(ev) {
			return
		}
		f, _ = f.Parent().(Control)
	}
	c.window.KeyStroke(ev)
}
