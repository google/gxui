// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"github.com/google/gxui/math"
	"image"

	"github.com/go-gl-legacy/gl"
)

type Texture struct {
	refCounted
	image        image.Image
	pixelsPerDip float32
	flipY        bool
}

func CreateTexture(img image.Image, pixelsPerDip float32) *Texture {
	t := &Texture{
		image:        img,
		pixelsPerDip: pixelsPerDip,
	}
	t.init()
	return t
}

// gxui.Texture compliance
func (t *Texture) Image() image.Image {
	return t.image
}

func (t *Texture) Size() math.Size {
	return t.SizePixels().ScaleS(1.0 / t.pixelsPerDip)
}

func (t *Texture) SizePixels() math.Size {
	s := t.image.Bounds().Size()
	return math.Size{W: s.X, H: s.Y}
}

func (t *Texture) FlipY() bool {
	return t.flipY
}

func (t *Texture) SetFlipY(flipY bool) {
	t.flipY = flipY
}

func (t *Texture) CreateContext() TextureContext {
	var fmt gl.GLenum
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

	glTex := gl.GenTexture()
	glTex.Bind(gl.TEXTURE_2D)
	w, h := t.SizePixels().WH()
	gl.TexImage2D(gl.TEXTURE_2D, 0, int(fmt), w, h, 0, fmt, gl.UNSIGNED_BYTE, data)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	glTex.Unbind(gl.TEXTURE_2D)
	CheckError()

	return TextureContext{
		texture:    glTex,
		sizePixels: t.Size(),
		flipY:      t.flipY,
		pma:        pma,
	}
}

type TextureContext struct {
	texture    gl.Texture
	sizePixels math.Size
	flipY      bool
	pma        bool
}

func (c TextureContext) Texture() gl.Texture {
	return c.texture
}

func (c TextureContext) SizePixels() math.Size {
	return c.sizePixels
}

func (c TextureContext) FlipY() bool {
	return c.flipY
}

func (c TextureContext) PremultipliedAlpha() bool {
	return c.pma
}

func (c *TextureContext) Destroy() {
	c.texture.Delete()
}
