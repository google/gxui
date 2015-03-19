// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/google/gxui/math"
	"image"
)

type Driver interface {
	// Events returns the event queue for the UI. The application should pull each
	// func from the returned chan and execute it on the main UI go-routine. The
	// application is free to write additional funcs to this chan in order for
	// them to be executed on the main UI go-routine.
	Events() chan func()
	Terminate()
	SetClipboard(str string)
	GetClipboard() (string, error)
	LoadFont(name string, size int) (Font, error)
	CreateViewport(width, height int, name string) Viewport
	CreateCanvas(math.Size) Canvas
	CreateTexture(img image.Image, pixelsPerDip float32) Texture
}
