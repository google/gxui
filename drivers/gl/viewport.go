// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"sync"
	"sync/atomic"
	"unicode"

	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl/platform"
	"github.com/google/gxui/math"
	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"
)

const viewportDebugEnabled = false

const clearColorR = 0.5
const clearColorG = 0.5
const clearColorB = 0.5

type viewport struct {
	sync.Mutex

	driver                  *driver
	context                 *context
	window                  *glfw.Window
	canvas                  *canvas
	fullscreen              bool
	scaling                 float32
	sizeDipsUnscaled        math.Size
	sizeDips                math.Size
	sizePixels              math.Size
	position                math.Point
	title                   string
	pendingMouseMoveEvent   *gxui.MouseEvent
	pendingMouseScrollEvent *gxui.MouseEvent
	scrollAccumX            float64
	scrollAccumY            float64
	destroyed               bool
	redrawCount             uint32

	// Broadcasts to application thread
	onClose       gxui.Event // ()
	onResize      gxui.Event // ()
	onMouseMove   gxui.Event // (gxui.MouseEvent)
	onMouseEnter  gxui.Event // (gxui.MouseEvent)
	onMouseExit   gxui.Event // (gxui.MouseEvent)
	onMouseDown   gxui.Event // (gxui.MouseEvent)
	onMouseUp     gxui.Event // (gxui.MouseEvent)
	onMouseScroll gxui.Event // (gxui.MouseEvent)
	onKeyDown     gxui.Event // (gxui.KeyboardEvent)
	onKeyUp       gxui.Event // (gxui.KeyboardEvent)
	onKeyRepeat   gxui.Event // (gxui.KeyboardEvent)
	onKeyStroke   gxui.Event // (gxui.KeyStrokeEvent)
	// Broadcasts to driver thread
	onDestroy gxui.Event
}

func newViewport(driver *driver, width, height int, title string, fullscreen bool) *viewport {
	v := &viewport{
		fullscreen: fullscreen,
		scaling:    1,
		title:      title,
	}

	glfw.DefaultWindowHints()
	glfw.WindowHint(glfw.Samples, 4)
	var monitor *glfw.Monitor
	if fullscreen {
		monitor = glfw.GetPrimaryMonitor()
		if width == 0 || height == 0 {
			vm := monitor.GetVideoMode()
			width, height = vm.Width, vm.Height
		}
	}
	wnd, err := glfw.CreateWindow(width, height, v.title, monitor, nil)
	if err != nil {
		panic(err)
	}
	width, height = wnd.GetSize() // At this time, width and height serve as a "hint" for glfw.CreateWindow, so get actual values from window.

	wnd.MakeContextCurrent()

	v.context = newContext()

	cursorPoint := func(x, y float64) math.Point {
		// HACK: xpos is off by 1 and ypos is off by 3 on OSX.
		// Compensate until real fix is found.
		x -= 1.0
		y -= 3.0
		return math.Point{X: int(x), Y: int(y)}.ScaleS(1 / v.scaling)
	}
	wnd.SetCloseCallback(func(*glfw.Window) {
		v.Close()
	})
	wnd.SetPosCallback(func(w *glfw.Window, x, y int) {
		v.Lock()
		v.position = math.NewPoint(x, y)
		v.Unlock()
	})
	wnd.SetSizeCallback(func(_ *glfw.Window, w, h int) {
		v.Lock()
		v.sizeDipsUnscaled = math.Size{W: w, H: h}
		v.sizeDips = v.sizeDipsUnscaled.ScaleS(1 / v.scaling)
		v.Unlock()
		v.onResize.Fire()
	})
	wnd.SetFramebufferSizeCallback(func(_ *glfw.Window, w, h int) {
		v.Lock()
		v.sizePixels = math.Size{W: w, H: h}
		v.Unlock()
		gl.Viewport(0, 0, w, h)
		gl.ClearColor(clearColorR, clearColorG, clearColorB, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
	})
	wnd.SetCursorPosCallback(func(w *glfw.Window, x, y float64) {
		p := cursorPoint(w.GetCursorPos())
		v.Lock()
		if v.pendingMouseMoveEvent == nil {
			v.pendingMouseMoveEvent = &gxui.MouseEvent{}
			driver.Call(func() {
				v.Lock()
				ev := *v.pendingMouseMoveEvent
				v.pendingMouseMoveEvent = nil
				v.Unlock()
				v.onMouseMove.Fire(ev)
			})
		}
		v.pendingMouseMoveEvent.Point = p
		v.pendingMouseMoveEvent.State = getMouseState(w)
		v.Unlock()
	})
	wnd.SetCursorEnterCallback(func(w *glfw.Window, entered bool) {
		p := cursorPoint(w.GetCursorPos())
		ev := gxui.MouseEvent{
			Point: p,
		}
		ev.State = getMouseState(w)
		if entered {
			v.onMouseEnter.Fire(ev)
		} else {
			v.onMouseExit.Fire(ev)
		}
	})
	wnd.SetScrollCallback(func(w *glfw.Window, xoff, yoff float64) {
		p := cursorPoint(w.GetCursorPos())
		v.Lock()
		if v.pendingMouseScrollEvent == nil {
			v.pendingMouseScrollEvent = &gxui.MouseEvent{}
			driver.Call(func() {
				v.Lock()
				ev := *v.pendingMouseScrollEvent
				v.pendingMouseScrollEvent = nil
				ev.ScrollX, ev.ScrollY = int(v.scrollAccumX), int(v.scrollAccumY)
				if ev.ScrollX != 0 || ev.ScrollY != 0 {
					v.scrollAccumX -= float64(ev.ScrollX)
					v.scrollAccumY -= float64(ev.ScrollY)
					v.Unlock()
					v.onMouseScroll.Fire(ev)
				} else {
					v.Unlock()
				}
			})
		}
		v.pendingMouseScrollEvent.Point = p
		v.scrollAccumX += xoff * platform.ScrollSpeed
		v.scrollAccumY += yoff * platform.ScrollSpeed
		v.pendingMouseScrollEvent.State = getMouseState(w)
		v.Unlock()
	})
	wnd.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		p := cursorPoint(w.GetCursorPos())
		ev := gxui.MouseEvent{
			Point:    p,
			Modifier: translateKeyboardModifier(mod),
		}
		ev.Button = translateMouseButton(button)
		ev.State = getMouseState(w)
		if action == glfw.Press {
			v.onMouseDown.Fire(ev)
		} else {
			v.onMouseUp.Fire(ev)
		}
	})
	wnd.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		ev := gxui.KeyboardEvent{
			Key:      translateKeyboardKey(key),
			Modifier: translateKeyboardModifier(mods),
		}
		switch action {
		case glfw.Press:
			v.onKeyDown.Fire(ev)
		case glfw.Release:
			v.onKeyUp.Fire(ev)
		case glfw.Repeat:
			v.onKeyRepeat.Fire(ev)
		}
	})
	wnd.SetCharModsCallback(func(w *glfw.Window, char rune, mods glfw.ModifierKey) {
		if !unicode.IsControl(char) &&
			!unicode.IsGraphic(char) &&
			!unicode.IsLetter(char) &&
			!unicode.IsMark(char) &&
			!unicode.IsNumber(char) &&
			!unicode.IsPunct(char) &&
			!unicode.IsSpace(char) &&
			!unicode.IsSymbol(char) {
			return // Weird unicode character. Ignore
		}

		ev := gxui.KeyStrokeEvent{
			Character: char,
			Modifier:  translateKeyboardModifier(mods),
		}
		v.onKeyStroke.Fire(ev)
	})
	wnd.SetRefreshCallback(func(w *glfw.Window) {
		if v.canvas != nil {
			v.render()
		}
	})

	fw, fh := wnd.GetFramebufferSize()
	posX, posY := wnd.GetPos()

	// Pre-multiplied alpha blending
	gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.SCISSOR_TEST)
	gl.Viewport(0, 0, fw, fh)
	gl.Scissor(0, 0, int32(fw), int32(fh))
	gl.ClearColor(clearColorR, clearColorG, clearColorB, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	wnd.SwapBuffers()

	v.window = wnd
	v.driver = driver
	v.onClose = driver.createAppEvent(func() {})
	v.onResize = driver.createAppEvent(func() {})
	v.onMouseMove = gxui.CreateEvent(func(gxui.MouseEvent) {})
	v.onMouseEnter = driver.createAppEvent(func(gxui.MouseEvent) {})
	v.onMouseExit = driver.createAppEvent(func(gxui.MouseEvent) {})
	v.onMouseDown = driver.createAppEvent(func(gxui.MouseEvent) {})
	v.onMouseUp = driver.createAppEvent(func(gxui.MouseEvent) {})
	v.onMouseScroll = gxui.CreateEvent(func(gxui.MouseEvent) {})
	v.onKeyDown = driver.createAppEvent(func(gxui.KeyboardEvent) {})
	v.onKeyUp = driver.createAppEvent(func(gxui.KeyboardEvent) {})
	v.onKeyRepeat = driver.createAppEvent(func(gxui.KeyboardEvent) {})
	v.onKeyStroke = driver.createAppEvent(func(gxui.KeyStrokeEvent) {})
	v.onDestroy = driver.createDriverEvent(func() {})
	v.sizeDipsUnscaled = math.Size{W: width, H: height}
	v.sizeDips = v.sizeDipsUnscaled.ScaleS(1 / v.scaling)
	v.sizePixels = math.Size{W: fw, H: fh}
	v.position = math.Point{X: posX, Y: posY}
	return v
}

// Driver methods
// These methods are all called on the driver routine
func (v *viewport) render() {
	if v.destroyed {
		return
	}

	v.window.MakeContextCurrent()

	ctx := v.context
	ctx.beginDraw(v.SizeDips(), v.SizePixels())

	dss := drawStateStack{drawState{
		ClipPixels: v.sizePixels.Rect(),
	}}

	v.canvas.draw(ctx, &dss)
	if len(dss) != 1 {
		panic("DrawStateStack count was not 1 after calling Canvas.Draw")
	}

	ctx.apply(dss.head())
	ctx.blitter.commit(ctx)

	if viewportDebugEnabled {
		v.drawFrameUpdate(ctx)
	}

	ctx.endDraw()

	v.window.SwapBuffers()
}

func (v *viewport) drawFrameUpdate(ctx *context) {
	dx := (ctx.stats.frameCount * 10) & 0xFF
	r := math.CreateRect(dx-5, 0, dx+5, 3)
	ds := &drawState{}
	ctx.blitter.blitRect(ctx, r, gxui.White, ds)
}

// gxui.viewport compliance
// These methods are all called on the application routine
func (v *viewport) SetCanvas(cc gxui.Canvas) {
	cnt := atomic.AddUint32(&v.redrawCount, 1)
	c := cc.(*canvas)
	v.driver.asyncDriver(func() {
		// Only use the canvas of the most recent SetCanvas call.
		v.window.MakeContextCurrent()
		if atomic.LoadUint32(&v.redrawCount) == cnt {
			v.canvas = c
			if v.canvas != nil {
				v.render()
			}
		}
	})
}

func (v *viewport) Scale() float32 {
	v.Lock()
	defer v.Unlock()
	return v.scaling
}

func (v *viewport) SetScale(s float32) {
	v.Lock()
	defer v.Unlock()
	if s != v.scaling {
		v.scaling = s
		v.sizeDips = v.sizeDipsUnscaled.ScaleS(1 / s)
		v.onResize.Fire()
	}
}

func (v *viewport) SizeDips() math.Size {
	v.Lock()
	defer v.Unlock()
	return v.sizeDips
}

func (v *viewport) SetSizeDips(size math.Size) {
	v.driver.syncDriver(func() {
		v.sizeDips = size
		v.sizeDipsUnscaled = size.ScaleS(v.scaling)
		v.window.SetSize(v.sizeDipsUnscaled.W, v.sizeDipsUnscaled.H)
	})
}

func (v *viewport) SizePixels() math.Size {
	v.Lock()
	defer v.Unlock()
	return v.sizePixels
}

func (v *viewport) Title() string {
	v.Lock()
	defer v.Unlock()
	return v.title
}

func (v *viewport) SetTitle(title string) {
	v.Lock()
	v.title = title
	v.Unlock()
	v.driver.asyncDriver(func() {
		v.window.SetTitle(title)
	})
}

func (v *viewport) Position() math.Point {
	v.Lock()
	defer v.Unlock()
	return v.position
}

func (v *viewport) SetPosition(pos math.Point) {
	v.Lock()
	v.position = pos
	v.Unlock()
	v.driver.asyncDriver(func() {
		v.window.SetPos(pos.X, pos.Y)
	})
}

func (v *viewport) Fullscreen() bool {
	return v.fullscreen
}

func (v *viewport) Show() {
	v.driver.asyncDriver(func() { v.window.Show() })
}

func (v *viewport) Hide() {
	v.driver.asyncDriver(func() { v.window.Hide() })
}

func (v *viewport) Close() {
	v.onClose.Fire()
	v.Destroy()
}

func (v *viewport) OnResize(f func()) gxui.EventSubscription {
	return v.onResize.Listen(f)
}

func (v *viewport) OnClose(f func()) gxui.EventSubscription {
	return v.onClose.Listen(f)
}

func (v *viewport) OnMouseMove(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return v.onMouseMove.Listen(f)
}

func (v *viewport) OnMouseEnter(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return v.onMouseEnter.Listen(f)
}

func (v *viewport) OnMouseExit(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return v.onMouseExit.Listen(f)
}

func (v *viewport) OnMouseDown(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return v.onMouseDown.Listen(f)
}

func (v *viewport) OnMouseUp(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return v.onMouseUp.Listen(f)
}

func (v *viewport) OnMouseScroll(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return v.onMouseScroll.Listen(f)
}

func (v *viewport) OnKeyDown(f func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return v.onKeyDown.Listen(f)
}

func (v *viewport) OnKeyUp(f func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return v.onKeyUp.Listen(f)
}

func (v *viewport) OnKeyRepeat(f func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return v.onKeyRepeat.Listen(f)
}

func (v *viewport) OnKeyStroke(f func(gxui.KeyStrokeEvent)) gxui.EventSubscription {
	return v.onKeyStroke.Listen(f)
}

func (v *viewport) Destroy() {
	v.driver.asyncDriver(func() {
		if !v.destroyed {
			v.window.MakeContextCurrent()
			v.canvas = nil
			v.context.destroy()
			v.window.Destroy()
			v.onDestroy.Fire()
			v.destroyed = true
		}
	})
}
