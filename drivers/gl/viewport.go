// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unicode"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl/platform"
	"github.com/google/gxui/math"
)

const viewportDebugEnabled = false

const kClearColorR = 0.5
const kClearColorG = 0.5
const kClearColorB = 0.5

type viewport struct {
	sync.Mutex

	driver                  *driver
	context                 *context
	window                  *glfw.Window
	canvas                  *canvas
	scaling                 float32
	sizeDipsUnscaled        math.Size
	sizeDips                math.Size
	sizePixels              math.Size
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
	onMouseScroll gxui.Event // (dx, dy int, p math.Point)
	onKeyDown     gxui.Event // (gxui.KeyboardEvent)
	onKeyUp       gxui.Event // (gxui.KeyboardEvent)
	onKeyRepeat   gxui.Event // (gxui.KeyboardEvent)
	onKeyStroke   gxui.Event // (gxui.KeyStrokeEvent)
	// Broadcasts to driver thread
	onDestroy gxui.Event
}

func newViewport(driver *driver, width, height int, title string) *viewport {
	v := &viewport{}

	glfw.DefaultWindowHints()
	glfw.WindowHint(glfw.Samples, 4)
	wnd, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)
	}
	wnd.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		panic(fmt.Errorf("Failed to initialize gl: %v", err))
	}

	v.context = newContext()
	v.scaling = 1

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
		gl.Viewport(0, 0, int32(w), int32(h))
		gl.ClearColor(kClearColorR, kClearColorG, kClearColorB, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
	})
	wnd.SetCursorPosCallback(func(w *glfw.Window, x, y float64) {
		p := cursorPoint(w.GetCursorPos())
		v.Lock()
		if v.pendingMouseMoveEvent == nil {
			v.pendingMouseMoveEvent = &gxui.MouseEvent{}
			driver.Call(func() {
				v.Lock()
				v.onMouseMove.Fire(*v.pendingMouseMoveEvent)
				v.pendingMouseMoveEvent = nil
				v.Unlock()
			})
		}
		v.pendingMouseMoveEvent.Point = p
		v.Unlock()
	})
	wnd.SetCursorEnterCallback(func(w *glfw.Window, entered bool) {
		p := cursorPoint(w.GetCursorPos())
		ev := gxui.MouseEvent{
			Point: p,
		}
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
				dx, dy := int(v.scrollAccumX), int(v.scrollAccumY)
				if dx != 0 || dy != 0 {
					v.pendingMouseScrollEvent.ScrollX = dx
					v.pendingMouseScrollEvent.ScrollY = dy
					v.onMouseScroll.Fire(*v.pendingMouseScrollEvent)
					v.scrollAccumX -= float64(dx)
					v.scrollAccumY -= float64(dy)
				}
				v.pendingMouseScrollEvent = nil
				v.Unlock()
			})
		}
		v.pendingMouseScrollEvent.Point = p
		v.scrollAccumX += xoff * platform.ScrollSpeed
		v.scrollAccumY += yoff * platform.ScrollSpeed
		v.Unlock()
	})
	wnd.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		p := cursorPoint(w.GetCursorPos())
		ev := gxui.MouseEvent{
			Point:    p,
			Modifier: translateKeyboardModifier(mod),
		}
		switch button {
		case glfw.MouseButtonLeft:
			ev.Button = gxui.MouseButtonLeft
		case glfw.MouseButtonMiddle:
			ev.Button = gxui.MouseButtonMiddle
		case glfw.MouseButtonRight:
			ev.Button = gxui.MouseButtonRight
		}
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

	// Pre-multiplied alpha blending
	gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.SCISSOR_TEST)
	gl.Viewport(0, 0, int32(fw), int32(fh))
	gl.Scissor(0, 0, int32(fw), int32(fh))
	gl.ClearColor(kClearColorR, kClearColorG, kClearColorB, 1.0)
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
	v.title = title

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
	if c != nil {
		c.addRef()
	}
	v.driver.asyncDriver(func() {
		// Only use the canvas of the most recent SetCanvas call.
		if atomic.LoadUint32(&v.redrawCount) == cnt {
			if v.canvas != nil {
				v.canvas.release()
			}
			v.canvas = c
			if v.canvas != nil {
				v.render()
			}
		} else if c != nil {
			c.release()
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
			v.context.destroy()
			v.window.Destroy()
			v.onDestroy.Fire()
			v.canvas.Release()
			v.canvas = nil
			v.destroyed = true
		}
	})
}
