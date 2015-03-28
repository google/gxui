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

	Terminate()
	SetClipboard(str string)
	GetClipboard() (string, error)

	// CreateFont loads a font from the provided TrueType bytes.
	CreateFont(data []byte, size int) (Font, error)

	CreateViewport(width, height int, name string) Viewport
	CreateCanvas(math.Size) Canvas
	CreateTexture(img image.Image, pixelsPerDip float32) Texture
}
