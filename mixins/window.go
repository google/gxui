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
	outer.IsVisibler
	outer.LayoutChildren
	outer.PaintChilder
	outer.Painter
	outer.Parenter
	outer.Sized
}

type Window struct {
	parts.Attachable
	parts.BackgroundBorderPainter
	parts.Container
	parts.Paddable
	parts.PaintChildren

	driver             gxui.Driver
	outer              WindowOuter
	viewport           gxui.Viewport
	windowedSize       math.Size
	mouseController    *gxui.MouseController
	keyboardController *gxui.KeyboardController
	focusController    *gxui.FocusController
	layoutPending      bool
	drawPending        bool
	updatePending      bool
	onClose            gxui.Event // Raised by viewport
	onResize           gxui.Event // Raised by viewport
	onMouseMove        gxui.Event // Raised by viewport
	onMouseEnter       gxui.Event // Raised by viewport
	onMouseExit        gxui.Event // Raised by viewport
	onMouseDown        gxui.Event // Raised by viewport
	onMouseUp          gxui.Event // Raised by viewport
	onMouseScroll      gxui.Event // Raised by viewport
	onKeyDown          gxui.Event // Raised by viewport
	onKeyUp            gxui.Event // Raised by viewport
	onKeyRepeat        gxui.Event // Raised by viewport
	onKeyStroke        gxui.Event // Raised by viewport

	onClick       gxui.Event // Raised by MouseController
	onDoubleClick gxui.Event // Raised by MouseController

	viewportSubscriptions []gxui.EventSubscription
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
	w.BackgroundBorderPainter.Init(outer)
	w.Container.Init(outer)
	w.Paddable.Init(outer)
	w.PaintChildren.Init(outer)
	w.outer = outer
	w.driver = driver

	w.onClose = gxui.CreateEvent(func() {})
	w.onResize = gxui.CreateEvent(func() {})
	w.onMouseMove = gxui.CreateEvent(func(gxui.MouseEvent) {})
	w.onMouseEnter = gxui.CreateEvent(func(gxui.MouseEvent) {})
	w.onMouseExit = gxui.CreateEvent(func(gxui.MouseEvent) {})
	w.onMouseDown = gxui.CreateEvent(func(gxui.MouseEvent) {})
	w.onMouseUp = gxui.CreateEvent(func(gxui.MouseEvent) {})
	w.onMouseScroll = gxui.CreateEvent(func(gxui.MouseEvent) {})
	w.onKeyDown = gxui.CreateEvent(func(gxui.KeyboardEvent) {})
	w.onKeyUp = gxui.CreateEvent(func(gxui.KeyboardEvent) {})
	w.onKeyRepeat = gxui.CreateEvent(func(gxui.KeyboardEvent) {})
	w.onKeyStroke = gxui.CreateEvent(func(gxui.KeyStrokeEvent) {})

	w.onClick = gxui.CreateEvent(func(gxui.MouseEvent) {})
	w.onDoubleClick = gxui.CreateEvent(func(gxui.MouseEvent) {})

	w.focusController = gxui.CreateFocusController(outer)
	w.mouseController = gxui.CreateMouseController(outer, w.focusController)
	w.keyboardController = gxui.CreateKeyboardController(outer)

	w.onResize.Listen(func() {
		w.outer.LayoutChildren()
		w.Draw()
	})

	w.SetBorderPen(gxui.TransparentPen)

	w.setViewport(driver.CreateWindowedViewport(width, height, title))

	// Window starts shown
	w.Attach()

	// Interface compliance test
	_ = gxui.Window(w)
}

func (w *Window) Draw() gxui.Canvas {
	if s := w.viewport.SizeDips(); s != math.ZeroSize {
		c := w.driver.CreateCanvas(s)
		w.outer.Paint(c)
		c.Complete()
		w.viewport.SetCanvas(c)
		return c
	} else {
		return nil
	}
}

func (w *Window) Paint(c gxui.Canvas) {
	w.PaintBackground(c, c.Size().Rect())
	w.PaintChildren.Paint(c)
	w.PaintBorder(c, c.Size().Rect())
}

func (w *Window) LayoutChildren() {
	s := w.Size().Contract(w.Padding()).Max(math.ZeroSize)
	o := w.Padding().LT()
	for _, c := range w.outer.Children() {
		c.Layout(c.Control.DesiredSize(math.ZeroSize, s).Rect().Offset(o))
	}
}

func (w *Window) Size() math.Size {
	return w.viewport.SizeDips()
}

func (w *Window) SetSize(size math.Size) {
	w.viewport.SetSizeDips(size)
}

func (w *Window) Parent() gxui.Parent {
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

func (w *Window) Position() math.Point {
	return w.viewport.Position()
}

func (w *Window) SetPosition(pos math.Point) {
	w.viewport.SetPosition(pos)
}

func (w *Window) Fullscreen() bool {
	return w.viewport.Fullscreen()
}

func (w *Window) SetFullscreen(fullscreen bool) {
	title := w.viewport.Title()
	if fullscreen != w.Fullscreen() {
		old := w.viewport
		if fullscreen {
			w.windowedSize = old.SizeDips()
			w.setViewport(w.driver.CreateFullscreenViewport(0, 0, title))
		} else {
			width, height := w.windowedSize.WH()
			w.setViewport(w.driver.CreateWindowedViewport(width, height, title))
		}
		old.Close()
	}
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
	return w.onClose.Listen(f)
}

func (w *Window) OnResize(f func()) gxui.EventSubscription {
	return w.onResize.Listen(f)
}

func (w *Window) OnClick(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return w.onClick.Listen(f)
}

func (w *Window) OnDoubleClick(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return w.onDoubleClick.Listen(f)
}

func (w *Window) OnMouseMove(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return w.onMouseMove.Listen(func(ev gxui.MouseEvent) {
		ev.Window = w
		ev.WindowPoint = ev.Point
		f(ev)
	})
}

func (w *Window) OnMouseEnter(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return w.onMouseEnter.Listen(func(ev gxui.MouseEvent) {
		ev.Window = w
		ev.WindowPoint = ev.Point
		f(ev)
	})
}

func (w *Window) OnMouseExit(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return w.onMouseExit.Listen(func(ev gxui.MouseEvent) {
		ev.Window = w
		ev.WindowPoint = ev.Point
		f(ev)
	})
}

func (w *Window) OnMouseDown(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return w.onMouseDown.Listen(func(ev gxui.MouseEvent) {
		ev.Window = w
		ev.WindowPoint = ev.Point
		f(ev)
	})
}

func (w *Window) OnMouseUp(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return w.onMouseUp.Listen(func(ev gxui.MouseEvent) {
		ev.Window = w
		ev.WindowPoint = ev.Point
		f(ev)
	})
}

func (w *Window) OnMouseScroll(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return w.onMouseScroll.Listen(func(ev gxui.MouseEvent) {
		ev.Window = w
		ev.WindowPoint = ev.Point
		f(ev)
	})
}

func (w *Window) OnKeyDown(f func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return w.onKeyDown.Listen(f)
}

func (w *Window) OnKeyUp(f func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return w.onKeyUp.Listen(f)
}

func (w *Window) OnKeyRepeat(f func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return w.onKeyRepeat.Listen(f)
}

func (w *Window) OnKeyStroke(f func(gxui.KeyStrokeEvent)) gxui.EventSubscription {
	return w.onKeyStroke.Listen(f)
}

func (w *Window) Relayout() {
	w.layoutPending = true
	w.requestUpdate()
}

func (w *Window) Redraw() {
	w.drawPending = true
	w.requestUpdate()
}

func (w *Window) Click(ev gxui.MouseEvent) {
	w.onClick.Fire(ev)
}

func (w *Window) DoubleClick(ev gxui.MouseEvent) {
	w.onDoubleClick.Fire(ev)
}

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

func (w *Window) setViewport(v gxui.Viewport) {
	for _, s := range w.viewportSubscriptions {
		s.Unlisten()
	}
	w.viewport = v
	w.viewportSubscriptions = []gxui.EventSubscription{
		v.OnClose(func() { w.onClose.Fire() }),
		v.OnResize(func() { w.onResize.Fire() }),
		v.OnMouseMove(func(ev gxui.MouseEvent) { w.onMouseMove.Fire(ev) }),
		v.OnMouseEnter(func(ev gxui.MouseEvent) { w.onMouseEnter.Fire(ev) }),
		v.OnMouseExit(func(ev gxui.MouseEvent) { w.onMouseExit.Fire(ev) }),
		v.OnMouseDown(func(ev gxui.MouseEvent) { w.onMouseDown.Fire(ev) }),
		v.OnMouseUp(func(ev gxui.MouseEvent) { w.onMouseUp.Fire(ev) }),
		v.OnMouseScroll(func(ev gxui.MouseEvent) { w.onMouseScroll.Fire(ev) }),
		v.OnKeyDown(func(ev gxui.KeyboardEvent) { w.onKeyDown.Fire(ev) }),
		v.OnKeyUp(func(ev gxui.KeyboardEvent) { w.onKeyUp.Fire(ev) }),
		v.OnKeyRepeat(func(ev gxui.KeyboardEvent) { w.onKeyRepeat.Fire(ev) }),
		v.OnKeyStroke(func(ev gxui.KeyStrokeEvent) { w.onKeyStroke.Fire(ev) }),
	}
	w.Relayout()
}
