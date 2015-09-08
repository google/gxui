// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gl contains an OpenGL implementation of the gxui.Driver interface.
package gl

import (
	"container/list"
	"image"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/google/gxui"
	"github.com/google/gxui/math"
	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"
)

// Maximum time allowed for application to process events on termination.
const maxFlushTime = time.Second * 3

func init() {
	runtime.LockOSThread()
}

type driver struct {
	pendingDriver chan func()
	pendingApp    chan func()
	terminated    int32 // non-zero represents driver terminations
	viewports     *list.List

	pcs  []uintptr // reusable scratch-buffer for use by runtime.Callers.
	uiPC uintptr   // the program-counter of the applicationLoop function.
}

// StartDriver starts the gl driver with the given appRoutine.
func StartDriver(appRoutine func(driver gxui.Driver)) {
	if runtime.GOMAXPROCS(-1) < 2 {
		runtime.GOMAXPROCS(2)
	}

	if err := glfw.Init(gl.ContextWatcher); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	driver := &driver{
		pendingDriver: make(chan func(), 256),
		pendingApp:    make(chan func(), 256),
		viewports:     list.New(),
		pcs:           make([]uintptr, 256),
	}

	driver.pendingApp <- driver.discoverUIGoRoutine
	driver.pendingApp <- func() { appRoutine(driver) }
	go driver.applicationLoop()
	driver.driverLoop()
}

func (d *driver) asyncDriver(f func()) {
	d.pendingDriver <- f
	d.wake()
}

func (d *driver) syncDriver(f func()) {
	c := make(chan bool, 1)
	d.asyncDriver(func() { f(); c <- true })
	<-c
}

func (d *driver) createDriverEvent(signature interface{}) gxui.Event {
	return gxui.CreateChanneledEvent(signature, d.pendingDriver)
}

func (d *driver) createAppEvent(signature interface{}) gxui.Event {
	return gxui.CreateChanneledEvent(signature, d.pendingApp)
}

// driverLoop pulls and executes funcs from the pendingDriver chan until chan
// close. If there are no funcs enqueued, the driver routine calls and blocks on
// glfw.WaitEvents. All sends on the pendingDriver chan should be paired with a
// call to wake() so that glfw.WaitEvents can return.
func (d *driver) driverLoop() {
	for {
		select {
		case ev, open := <-d.pendingDriver:
			if open {
				ev()
			} else {
				return // termintated
			}
		default:
			glfw.WaitEvents()
		}
	}
}

func (d *driver) wake() {
	glfw.PostEmptyEvent()
}

// applicationLoop pulls and executes funcs from the pendingApp chan until
// the chan is closed.
func (d *driver) applicationLoop() {
	for ev := range d.pendingApp {
		ev()
	}
}

// gxui.Driver compliance
func (d *driver) Call(f func()) bool {
	if f == nil {
		panic("Function must not be nil")
	}
	if atomic.LoadInt32(&d.terminated) != 0 {
		return false // Driver.Terminate has been called
	}
	d.pendingApp <- f
	return true
}

func (d *driver) CallSync(f func()) bool {
	c := make(chan struct{})
	if d.Call(func() { f(); close(c) }) {
		<-c
		return true
	}
	return false
}

func (d *driver) Terminate() {
	d.asyncDriver(func() {
		// Close all viewports. This will notify the application.
		for v := d.viewports.Front(); v != nil; v = v.Next() {
			v.Value.(*viewport).Destroy()
		}

		// Flush all remaining events from the application and driver.
		// This gives the application an opportunity to handle shutdown.
		flushStart := time.Now()
		for time.Since(flushStart) < maxFlushTime {
			done := true

			// Process any application events
			sync := make(chan struct{})
			d.Call(func() {
				select {
				case ev := <-d.pendingApp:
					ev()
					done = false
				default:
				}
				close(sync)
			})
			<-sync

			// Process any driver events
			select {
			case ev := <-d.pendingDriver:
				ev()
				done = false
			default:
			}

			if done {
				break
			}
		}

		// All done.
		atomic.StoreInt32(&d.terminated, 1)
		close(d.pendingApp)
		close(d.pendingDriver)

		d.viewports = nil
	})
}

func (d *driver) SetClipboard(str string) {
	d.asyncDriver(func() {
		v := d.viewports.Front().Value.(*viewport)
		v.window.SetClipboardString(str)
	})
}

func (d *driver) GetClipboard() (str string, err error) {
	d.syncDriver(func() {
		c := d.viewports.Front().Value.(*viewport)
		str, err = c.window.GetClipboardString()
	})
	return
}

func (d *driver) CreateFont(data []byte, size int) (gxui.Font, error) {
	return newFont(data, size)
}

func (d *driver) CreateWindowedViewport(width, height int, name string) gxui.Viewport {
	var v *viewport
	d.syncDriver(func() {
		v = newViewport(d, width, height, name, false)
		e := d.viewports.PushBack(v)
		v.onDestroy.Listen(func() {
			d.viewports.Remove(e)
		})
	})
	return v
}

func (d *driver) CreateFullscreenViewport(width, height int, name string) gxui.Viewport {
	var v *viewport
	d.syncDriver(func() {
		v = newViewport(d, width, height, name, true)
		e := d.viewports.PushBack(v)
		v.onDestroy.Listen(func() {
			d.viewports.Remove(e)
		})
	})
	return v
}

func (d *driver) CreateCanvas(s math.Size) gxui.Canvas {
	return newCanvas(s)
}

func (d *driver) CreateTexture(img image.Image, pixelsPerDip float32) gxui.Texture {
	return newTexture(img, pixelsPerDip)
}
