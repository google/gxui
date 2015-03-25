// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"container/list"
	"image"
	"runtime"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/google/gxui"
	"github.com/google/gxui/math"
)

func init() {
	runtime.LockOSThread()
}

type driver struct {
	pendingDriver chan func()
	pendingApp    chan func()
	viewports     *list.List
	debugEnabled  bool
}

func StartDriver(appRoutine func(driver gxui.Driver)) {
	if runtime.GOMAXPROCS(-1) < 2 {
		runtime.GOMAXPROCS(2)
	}

	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	driver := &driver{
		pendingDriver: make(chan func(), 256),
		pendingApp:    make(chan func(), 256),
		viewports:     list.New(),
	}

	go appRoutine(driver)
	driver.run()
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

func (d *driver) flush() {
	for {
		select {
		case ev := <-d.pendingDriver:
			ev()
		default:
			return
		}
	}
}

func (d *driver) run() {
	for {
		select {
		case ev, open := <-d.pendingDriver:
			if open {
				ev()
			} else {
				return // closed channel represents driver shutdown
			}
		default:
			glfw.WaitEvents()
		}
	}
}

func (d *driver) wake() {
	glfw.PostEmptyEvent()
}

func (d *driver) EnableDebug(enabled bool) {
	d.debugEnabled = enabled
}

// gxui.Driver compliance
func (d *driver) Events() chan func() {
	return d.pendingApp
}

func (d *driver) Terminate() {
	d.asyncDriver(func() {
		for v := d.viewports.Front(); v != nil; v = v.Next() {
			v.Value.(*viewport).Destroy()
		}
		d.flush()
		close(d.pendingDriver)
		close(d.pendingApp)
		d.viewports.Init()
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

func (d *driver) CreateViewport(width, height int, name string) gxui.Viewport {
	var v *viewport
	d.syncDriver(func() {
		v = newViewport(d, width, height, name)
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
