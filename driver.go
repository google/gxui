// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"image"

	"github.com/google/gxui/math"
)

type Driver interface {
	// Call queues f to be run on the UI go-routine, returning before f may have
	// been called. Call returns false if the driver has been terminated, in which
	// case f may not be called.
	Call(f func()) bool

	// CallSync queues and then blocks for f to be run on the UI go-routine.
	// Call returns false if the driver has been terminated, in which case f may
	// not be called.
	CallSync(f func()) bool

	Terminate()
	SetClipboard(str string)
	GetClipboard() (string, error)

	// CreateFont loads a font from the provided TrueType bytes.
	CreateFont(data []byte, size int) (Font, error)

	// CreateWindowedViewport creates a new windowed Viewport with the specified
	// width and height in device independent pixels.
	CreateWindowedViewport(width, height int, name string) Viewport

	// CreateFullscreenViewport creates a new fullscreen Viewport with the
	// specified width and height in device independent pixels. If width or
	// height is 0, then the viewport adopts the current screen resolution.
	CreateFullscreenViewport(width, height int, name string) Viewport

	CreateCanvas(math.Size) Canvas
	CreateTexture(img image.Image, pixelsPerDip float32) Texture

	// Debug function used to verify that the caller is executing on the UI
	// go-routine. If the caller is not on the UI go-routine then the function
	// panics.
	AssertUIGoroutine()
}
