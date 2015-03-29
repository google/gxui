// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/google/gxui"
	"github.com/google/gxui/math"
	"github.com/google/gxui/mixins/outer"
	"github.com/google/gxui/mixins/parts"
)

type WindowOuter interface {
	gxui.Window
	outer.Attachable
	outer.Bounds
	outer.IsVisibler
	outer.LayoutChildren
	outer.PaintChilder
	outer.Painter
	outer.Parenter
}

type Window struct {
	parts.Attachable
	parts.Container
	parts.Paddable
	parts.PaintChildren

	driver             gxui.Driver
	outer              WindowOuter
	viewport           gxui.Viewport
	mouseController    *gxui.MouseController
	keyboardController *gxui.KeyboardController
	focusController    *gxui.FocusController
	layoutPending      bool
	drawPending        bool
	updatePending      bool
	onAttach           gxui.Event
	onDetach           gxui.Event
}

func (w *Window) requestUpdate() {
	if !w.updatePending {
		w.updatePending = true
		w.driver.Call(w.update)
	}
}

func (w *Window) update() {
	if !w.Attached() {
		// Window was detached between requestUpdate() and update()
		w.updatePending = false
		w.layoutPending = false
		w.drawPending = false
		return
	}
	w.updatePending = false
	if w.layoutPending {
		w.layoutPending = false
		w.drawPending = true
		w.outer.LayoutChildren()
	}
	if w.drawPending {
		w.drawPending = false
		w.Draw()
	}
}

func (w *Window) Init(outer WindowOuter, driver gxui.Driver, width, height int, title string) {
	w.Attachable.Init(outer)
	w.Container.Init(outer)
	w.Paddable.Init(outer)
	w.PaintChildren.Init(outer)
	w.outer = outer
	w.driver = driver
	w.viewport = driver.CreateViewport(width, height, title)
	w.focusController = gxui.CreateFocusController(outer)
	w.mouseController = gxui.CreateMouseController(outer, w.focusController)
	w.keyboardController = gxui.CreateKeyboardController(outer)
	w.viewport.OnResize(func() {
		w.outer.LayoutChildren()
		w.Draw()
	})

	// Window starts shown
	w.Attach()

	// Interface compliance test
	_ = gxui.Window(w)
}

func (w *Window) Draw() gxui.Canvas {
	if s := w.viewport.SizeDips(); s != math.ZeroSize {
		c := w.driver.CreateCanvas(w.viewport.SizeDips())
		w.outer.Paint(c)
		c.Complete()
		w.viewport.SetCanvas(c)
		c.Release()
		return c
	} else {
		return nil
	}
}

func (w *Window) LayoutChildren() {
	s := w.Bounds().Size().Contract(w.Padding())
	o := w.Padding().LT()
	for _, c := range w.outer.Children() {
		c.Layout(c.DesiredSize(math.ZeroSize, s).Rect().Offset(o))
	}
}

func (w *Window) Bounds() math.Rect {
	s := w.viewport.SizeDips()
	return math.CreateRect(0, 0, s.W, s.H)
}

func (w *Window) Parent() gxui.Container {
	return nil
}

func (w *Window) Viewport() gxui.Viewport {
	return w.viewport
}

func (w *Window) Title() string {
	return w.viewport.Title()
}

func (w *Window) SetTitle(t string) {
	w.viewport.SetTitle(t)
}

func (w *Window) Scale() float32 {
	return w.viewport.Scale()
}

func (w *Window) SetScale(scale float32) {
	w.viewport.SetScale(scale)
}

func (w *Window) Show() {
	w.Attach()
	w.viewport.Show()
}

func (w *Window) Hide() {
	w.Detach()
	w.viewport.Hide()
}

func (w *Window) Close() {
	w.Detach()
	w.viewport.Close()
}

func (w *Window) Focus() gxui.Focusable {
	return w.focusController.Focus()
}

func (w *Window) SetFocus(c gxui.Control) bool {
	fc := w.focusController
	if c == nil {
		fc.SetFocus(nil)
		return true
	}
	if f := fc.Focusable(c); f != nil {
		fc.SetFocus(f)
		return true
	}
	return false
}

func (w *Window) IsVisible() bool {
	return true
}

func (w *Window) OnClose(f func()) gxui.EventSubscription {
	return w.viewport.OnClose(f)
}

func (w *Window) OnResize(f func()) gxui.EventSubscription {
	return w.viewport.OnResize(f)
}

func (w *Window) OnMouseMove(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return w.viewport.OnMouseMove(func(ev gxui.MouseEvent) {
		ev.Window = w
		ev.WindowPoint = ev.Point
		f(ev)
	})
}

func (w *Window) OnMouseEnter(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return w.viewport.OnMouseEnter(func(ev gxui.MouseEvent) {
		ev.Window = w
		ev.WindowPoint = ev.Point
		f(ev)
	})
}

func (w *Window) OnMouseExit(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return w.viewport.OnMouseExit(func(ev gxui.MouseEvent) {
		ev.Window = w
		ev.WindowPoint = ev.Point
		f(ev)
	})
}

func (w *Window) OnMouseDown(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return w.viewport.OnMouseDown(func(ev gxui.MouseEvent) {
		ev.Window = w
		ev.WindowPoint = ev.Point
		f(ev)
	})
}

func (w *Window) OnMouseUp(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return w.viewport.OnMouseUp(func(ev gxui.MouseEvent) {
		ev.Window = w
		ev.WindowPoint = ev.Point
		f(ev)
	})
}

func (w *Window) OnMouseScroll(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return w.viewport.OnMouseScroll(func(ev gxui.MouseEvent) {
		ev.Window = w
		ev.WindowPoint = ev.Point
		f(ev)
	})
}

func (w *Window) OnKeyDown(f func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return w.viewport.OnKeyDown(f)
}

func (w *Window) OnKeyUp(f func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return w.viewport.OnKeyUp(f)
}

func (w *Window) OnKeyRepeat(f func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return w.viewport.OnKeyRepeat(f)
}

func (w *Window) OnKeyStroke(f func(gxui.KeyStrokeEvent)) gxui.EventSubscription {
	return w.viewport.OnKeyStroke(f)
}

func (w *Window) Relayout() {
	w.layoutPending = true
	w.requestUpdate()
}

func (w *Window) Redraw() {
	w.drawPending = true
	w.requestUpdate()
}

func (w *Window) Click(gxui.MouseEvent)       {}
func (w *Window) DoubleClick(gxui.MouseEvent) {}

func (w *Window) KeyPress(ev gxui.KeyboardEvent) {
	if ev.Key == gxui.KeyTab {
		if ev.Modifier&gxui.ModShift != 0 {
			w.focusController.FocusPrev()
		} else {
			w.focusController.FocusNext()
		}
	}
}
func (w *Window) KeyStroke(gxui.KeyStrokeEvent) {}
