// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"gaze/gxui/math"
	"image"
)

type Driver interface {
	Events() chan func()
	Call(func())
	Terminate()
	SetClipboard(str string)
	GetClipboard() (string, error)
	CreateViewport(width, height int, name string) Viewport
	CreateCanvas(math.Size) Canvas
	CreateFont(name string, size int) Font
	CreateTexture(img image.Image, pixelsPerDip float32) Texture
}
