// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"
	"sync"
	"unicode"

	"github.com/go-gl-legacy/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl/platform"
	"github.com/google/gxui/math"
)

const kClearColorR = 0.5
const kClearColorG = 0.5
const kClearColorB = 0.5

type Viewport struct {
	sync.Mutex

	driver                  *Driver
	context                 *Context
	window                  *glfw.Window
	canvas                  *Canvas
	sizeDips, sizePixels    math.Size
	title                   string
	pendingMouseMoveEvent   *gxui.MouseEvent
	pendingMouseScrollEvent *gxui.MouseEvent
	scrollAccumX            float64
	scrollAccumY            float64
	destroyed               bool
	continuousRedraw        bool
	redrawCount             int

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

func CreateViewport(driver *Driver, width, height int, title string) *Viewport {
	v := &Viewport{}

	glfw.DefaultWindowHints()
	glfw.WindowHint(glfw.Samples, 4)
	wnd, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)
	}
	wnd.MakeContextCurrent()
	if gl.Init() != gl.GLenum(0) {
		panic("Failed to initialize gl")
	}

	v.context = CreateContext()

	cursorPoint := func(x, y float64) math.Point {
		// HACK: xpos is off by 1 and ypos is off by 3 on OSX.
		// Compensate until real fix is found.
		x -= 1.0
		y -= 3.0
		return math.Point{X: int(x), Y: int(y)}
	}
	wnd.SetCloseCallback(func(*glfw.Window) {
		v.onClose.Fire()
		v.Close()
	})
	wnd.SetSizeCallback(func(_ *glfw.Window, w, h int) {
		v.Lock()
		v.sizeDips = math.Size{W: w, H: h}
		v.Unlock()
		v.onResize.Fire()
	})
	wnd.SetFramebufferSizeCallback(func(_ *glfw.Window, w, h int) {
		v.Lock()
		v.sizePixels = math.Size{W: w, H: h}
		v.Unlock()
		gl.Viewport(0, 0, w, h)
		gl.ClearColor(kClearColorR, kClearColorG, kClearColorB, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
	})
	wnd.SetCursorPosCallback(func(w *glfw.Window, x, y float64) {
		p := cursorPoint(w.GetCursorPos())
		v.Lock()
		if v.pendingMouseMoveEvent == nil {
			v.pendingMouseMoveEvent = &gxui.MouseEvent{}
			driver.Events() <- func() {
				v.Lock()
				v.onMouseMove.Fire(*v.pendingMouseMoveEvent)
				v.pendingMouseMoveEvent = nil
				v.Unlock()
			}
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
			driver.Events() <- func() {
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
			}
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
	gl.Viewport(0, 0, fw, fh)
	gl.Scissor(0, 0, fw, fh)
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
	v.sizeDips = math.Size{W: width, H: height}
	v.sizePixels = math.Size{W: fw, H: fh}
	v.title = title

	return v
}

// Driver methods
// These methods are all called on the driver thread
func (v *Viewport) render() {
	if v.destroyed {
		panic("Attempting to render a destroyed Viewport")
	}

	v.window.MakeContextCurrent()

	ctx := v.context
	ctx.BeginDraw(v.SizeDips(), v.SizePixels())

	dss := DrawStateStack{DrawState{
		ClipPixels: v.sizePixels.Rect(),
	}}

	v.canvas.draw(ctx, &dss)
	if len(dss) != 1 {
		panic("DrawStateStack count was not 1 after calling Canvas.Draw")
	}

	ctx.Apply(dss.Head())
	ctx.Blitter.Commit(ctx)

	if v.driver.debugEnabled {
		v.drawFrameUpdate(ctx)
	}

	ctx.EndDraw()

	v.window.SwapBuffers()

	if v.continuousRedraw {
		if ctx.Stats().FrameCount%60 == 0 {
			print(ctx.Stats().String() + "\n")
		}
		glfw.PollEvents()
		v.driver.flush()
		v.driver.asyncDriver(v.render)
	}
}

func (v *Viewport) drawFrameUpdate(ctx *Context) {
	dx := (ctx.Stats().FrameCount * 10) & 0xFF
	r := math.CreateRect(dx-5, 0, dx+5, 3)
	ds := &DrawState{}
	ctx.Blitter.BlitRect(ctx, r, gxui.White, ds)
}

// gxui.Viewport compliance
// These methods are all called on the application thread
func (v *Viewport) SetCanvas(canvas gxui.Canvas) {
	if v.destroyed {
		panic("Attempting to set the canvas on a destroyed Viewport")
	}
	v.redrawCount++
	cnt := v.redrawCount
	c := canvas.(*Canvas)
	if c != nil {
		c.AddRef()
	}
	v.driver.asyncDriver(func() {
		if v.redrawCount == cnt {
			if v.canvas != nil {
				v.canvas.Release()
			}
			v.canvas = c
			if v.canvas != nil {
				v.render()
			}
		} else if c != nil {
			c.Release()
		}
	})
}

func (v *Viewport) SizeDips() math.Size {
	v.Lock()
	defer v.Unlock()
	return v.sizeDips
}

func (v *Viewport) SizePixels() math.Size {
	v.Lock()
	defer v.Unlock()
	return v.sizePixels
}

func (v *Viewport) Title() string {
	v.Lock()
	defer v.Unlock()
	return v.title
}

func (v *Viewport) SetTitle(title string) {
	v.Lock()
	v.title = title
	v.Unlock()
	v.driver.asyncDriver(func() {
		v.window.SetTitle(title)
	})
}

func (v *Viewport) Show() {
	v.driver.asyncDriver(func() { v.window.Show() })
}

func (v *Viewport) Hide() {
	v.driver.asyncDriver(func() { v.window.Hide() })
}

func (v *Viewport) Close() {
	v.Destroy()
}

func (v *Viewport) OnResize(f func()) gxui.EventSubscription {
	return v.onResize.Listen(f)
}

func (v *Viewport) OnClose(f func()) gxui.EventSubscription {
	return v.onClose.Listen(f)
}

func (v *Viewport) OnMouseMove(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return v.onMouseMove.Listen(f)
}

func (v *Viewport) OnMouseEnter(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return v.onMouseEnter.Listen(f)
}

func (v *Viewport) OnMouseExit(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return v.onMouseExit.Listen(f)
}

func (v *Viewport) OnMouseDown(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return v.onMouseDown.Listen(f)
}

func (v *Viewport) OnMouseUp(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return v.onMouseUp.Listen(f)
}

func (v *Viewport) OnMouseScroll(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return v.onMouseScroll.Listen(f)
}

func (v *Viewport) OnKeyDown(f func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return v.onKeyDown.Listen(f)
}

func (v *Viewport) OnKeyUp(f func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return v.onKeyUp.Listen(f)
}

func (v *Viewport) OnKeyRepeat(f func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return v.onKeyRepeat.Listen(f)
}

func (v *Viewport) OnKeyStroke(f func(gxui.KeyStrokeEvent)) gxui.EventSubscription {
	return v.onKeyStroke.Listen(f)
}

func (v *Viewport) Destroy() {
	v.driver.asyncDriver(func() {
		if !v.destroyed {
			v.window.MakeContextCurrent()
			v.context.Destroy()
			v.window.Destroy()
			v.onDestroy.Fire()
			v.canvas.Release()
			v.canvas = nil
			v.destroyed = true
		}
	})
}

func (v *Viewport) SetContinuousRedraw(continuousRedraw bool) {
	v.continuousRedraw = continuousRedraw
	if continuousRedraw {
		v.driver.asyncDriver(func() {
			v.render()
		})
	}
}

func (v *Viewport) ContinuousRedraw() bool {
	return v.continuousRedraw
}

func (v *Viewport) Stats() string {
	return fmt.Sprintf("Global: %s\nViewport: %s", globalStats.String(), v.context.LastStats().String())
}
