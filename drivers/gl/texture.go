// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"image"

	"github.com/google/gxui/math"
	"github.com/goxjs/gl"
)

type texture struct {
	image        image.Image
	pixelsPerDip float32
	flipY        bool
}

func newTexture(img image.Image, pixelsPerDip float32) *texture {
	t := &texture{
		image:        img,
		pixelsPerDip: pixelsPerDip,
	}
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

func (t *texture) newContext() *textureContext {
	var fmt gl.Enum
	var data []byte
	var pma bool

	switch ty := t.image.(type) {
	case *image.RGBA:
		fmt = gl.RGBA
		data = ty.Pix
		pma = true
	case *image.NRGBA:
		fmt = gl.RGBA
		data = ty.Pix
	case *image.Alpha:
		fmt = gl.ALPHA
		data = ty.Pix
	default:
		panic("Unsupported image type")
	}

	texture := gl.CreateTexture()
	gl.BindTexture(gl.TEXTURE_2D, texture)
	w, h := t.SizePixels().WH()
	gl.TexImage2D(gl.TEXTURE_2D, 0, w, h, fmt, gl.UNSIGNED_BYTE, data)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.BindTexture(gl.TEXTURE_2D, gl.Texture{})
	checkError()

	globalStats.textureContextCount.inc()
	return &textureContext{
		texture:    texture,
		sizePixels: t.Size(),
		flipY:      t.flipY,
		pma:        pma,
	}
}

type textureContext struct {
	contextResource
	texture    gl.Texture
	sizePixels math.Size
	flipY      bool
	pma        bool
}

func (c *textureContext) destroy() {
	globalStats.textureContextCount.dec()
	gl.DeleteTexture(c.texture)
	c.texture = gl.Texture{}
}
