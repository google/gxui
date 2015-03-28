// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"github.com/google/gxui/math"
	"image"

	"github.com/go-gl/gl/v2.1/gl"
)

type texture struct {
	refCounted
	image        image.Image
	pixelsPerDip float32
	flipY        bool
}

func newTexture(img image.Image, pixelsPerDip float32) *texture {
	t := &texture{
		image:        img,
		pixelsPerDip: pixelsPerDip,
	}
	t.init()
	return t
}

// gxui.Texture compliance
func (t *texture) Image() image.Image {
	return t.image
}

func (t *texture) Size() math.Size {
	return t.SizePixels().ScaleS(1.0 / t.pixelsPerDip)
}

func (t *texture) SizePixels() math.Size {
	s := t.image.Bounds().Size()
	return math.Size{W: s.X, H: s.Y}
}

func (t *texture) FlipY() bool {
	return t.flipY
}

func (t *texture) SetFlipY(flipY bool) {
	t.flipY = flipY
}

func (t *texture) Release() {
	t.release()
}

func (t *texture) newContext() *textureContext {
	var fmt uint32
	var data interface{}
	var pma bool

	switch ty := t.image.(type) {
	case *image.RGBA:
		fmt = gl.RGBA
		data = ty.Pix
		pma = true
	case *image.NRGBA:
		fmt = gl.RGBA
		data = ty.Pix
	case *image.Gray:
		fmt = gl.RED
		data = ty.Pix
	case *image.Alpha:
		fmt = gl.ALPHA
		data = ty.Pix
	default:
		panic("Unsupported image type")
	}

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	w, h := t.SizePixels().WH()
	gl.TexImage2D(gl.TEXTURE_2D, 0, int32(fmt), int32(w), int32(h), 0, fmt, gl.UNSIGNED_BYTE, gl.Ptr(data))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	checkError()

	return &textureContext{
		texture:    texture,
		sizePixels: t.Size(),
		flipY:      t.flipY,
		pma:        pma,
	}
}

type textureContext struct {
	texture    uint32
	sizePixels math.Size
	flipY      bool
	pma        bool
}

func (c *textureContext) destroy() {
	gl.DeleteTextures(1, &c.texture)
	c.texture = 0
}
