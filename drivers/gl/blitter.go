// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"github.com/google/gxui"
	"github.com/google/gxui/math"

	"github.com/goxjs/gl"
)

const (
	debugWireframePolygons = false

	vsCopySrc = `
  attribute vec2 aPosition;
  varying vec2 vTexcoords;
  uniform mat3 mPos;
  uniform mat3 mUV;
  void main() {
    vec3 pos3 = vec3(aPosition, 1.0);
    gl_Position = vec4(mPos * pos3, 1.0);
    vTexcoords = (mUV * pos3).xy;
  }`

	fsCopySrc = `
  #ifdef GL_ES
    precision mediump float;
  #endif

  uniform sampler2D source;
  varying vec2 vTexcoords;
  void main() {
    gl_FragColor = texture2D(source, vTexcoords);
  }`

	vsColorSrc = `
  attribute vec2 aPosition;
  uniform mat3 mPos;
  void main() {
    vec3 pos3 = vec3(aPosition, 1.0);
    gl_Position = vec4((mPos * pos3).xy, 0.0, 1.0);
  }`

	fsColorSrc = `
  #ifdef GL_ES
    precision mediump float;
  #endif

  uniform vec4 Color;
  void main() {
    gl_FragColor = Color;
    gl_FragColor *= gl_FragColor.a; // PMA
  }`

	vsFontSrc = `
  attribute vec2 aSrc;
  attribute vec2 aDst;
  attribute vec4 aClp;
  attribute vec4 aCol;
  varying vec2 vSrc;
  varying vec4 vCol;
  varying vec2 vClp;
  uniform mat3 mSrc;
  uniform mat3 mDst;
  void main() {
    vec2 vClipMin = (mDst * vec3(aClp.xy, 1.0)).xy;
    vec2 vClipMax = (mDst * vec3(aClp.zw, 1.0)).xy;
    gl_Position = vec4(mDst * vec3(aDst, 1.0), 1.0);
    vSrc = (mSrc * vec3(aSrc, 1.0)).xy;
    vClp = (gl_Position.xy - vClipMin) / (vClipMax - vClipMin);
    vCol = aCol;
  }`

	fsFontSrc = `
  #ifdef GL_ES
    precision mediump float;
  #endif

  uniform sampler2D source;
  varying vec2 vSrc;
  varying vec4 vCol;
  varying vec2 vClp;
  void main() {
    vec2 clipping = step(vec2(0.0, 0.0), vClp) * step(vClp, vec2(1.0, 1.0));
    gl_FragColor  = vCol * texture2D(source, vSrc).aaaa;
    gl_FragColor *= clipping.x * clipping.y;
  }`
)

type glyphBatch struct {
	DstRects  []float32
	SrcRects  []float32
	Colors    []float32
	ClipRects []float32
	Indices   []uint16
	GlyphPage *textureContext
}

type blitter struct {
	stats       *contextStats
	quad        *shape
	copyShader  *shaderProgram
	colorShader *shaderProgram
	fontShader  *shaderProgram
	glyphBatch  glyphBatch
}

func newBlitter(ctx *context, stats *contextStats) *blitter {
	return &blitter{
		stats:       stats,
		quad:        newQuadShape(),
		copyShader:  newShaderProgram(ctx, vsCopySrc, fsCopySrc),
		colorShader: newShaderProgram(ctx, vsColorSrc, fsColorSrc),
		fontShader:  newShaderProgram(ctx, vsFontSrc, fsFontSrc),
	}
}

func (b *blitter) destroy(ctx *context) {
	b.copyShader.destroy(ctx)
	b.colorShader.destroy(ctx)
	b.fontShader.destroy(ctx)
}

func (b *blitter) blit(ctx *context, tc *textureContext, srcRect, dstRect math.Rect, ds *drawState) {
	b.commitGlyphs(ctx)

	dstRect = dstRect.Offset(ds.OriginPixels)
	sw, sh := tc.sizePixels.WH()
	dw, dh := ctx.sizePixels.WH()

	var mUV math.Mat3
	if tc.flipY {
		mUV = math.CreateMat3(
			float32(srcRect.W())/float32(sw), 0, 0,
			0, -float32(srcRect.H())/float32(sh), 0,
			float32(srcRect.Min.X)/float32(sw),
			1.0-float32(srcRect.Min.Y)/float32(sh), 1,
		)
	} else {
		mUV = math.CreateMat3(
			float32(srcRect.W())/float32(sw), 0, 0,
			0, float32(srcRect.H())/float32(sh), 0,
			float32(srcRect.Min.X)/float32(sw),
			float32(srcRect.Min.Y)/float32(sh), 1,
		)
	}
	mPos := math.CreateMat3(
		+2.0*float32(dstRect.W())/float32(dw), 0, 0,
		0, -2.0*float32(dstRect.H())/float32(dh), 0,
		-1.0+2.0*float32(dstRect.Min.X)/float32(dw),
		+1.0-2.0*float32(dstRect.Min.Y)/float32(dh), 1,
	)
	if !tc.pma {
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	}
	b.quad.draw(ctx, b.copyShader, uniformBindings{
		"source": tc,
		"mUV":    mUV,
		"mPos":   mPos,
	})
	if !tc.pma {
		gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	}
	b.stats.drawCallCount++
}

func (b *blitter) blitGlyph(ctx *context, tc *textureContext, c gxui.Color, srcRect, dstRect math.Rect, ds *drawState) {
	dstRect = dstRect.Offset(ds.OriginPixels)

	if b.glyphBatch.GlyphPage != tc {
		b.commitGlyphs(ctx)
		b.glyphBatch.GlyphPage = tc
	}
	i := uint16(len(b.glyphBatch.DstRects)) / 2
	clip := []float32{
		float32(ds.ClipPixels.Min.X),
		float32(ds.ClipPixels.Min.Y),
		float32(ds.ClipPixels.Max.X),
		float32(ds.ClipPixels.Max.Y),
	}
	b.glyphBatch.DstRects = append(b.glyphBatch.DstRects,
		float32(dstRect.Min.X), float32(dstRect.Min.Y),
		float32(dstRect.Max.X), float32(dstRect.Min.Y),
		float32(dstRect.Min.X), float32(dstRect.Max.Y),
		float32(dstRect.Max.X), float32(dstRect.Max.Y),
	)
	b.glyphBatch.SrcRects = append(b.glyphBatch.SrcRects,
		float32(srcRect.Min.X), float32(srcRect.Min.Y),
		float32(srcRect.Max.X), float32(srcRect.Min.Y),
		float32(srcRect.Min.X), float32(srcRect.Max.Y),
		float32(srcRect.Max.X), float32(srcRect.Max.Y),
	)
	b.glyphBatch.ClipRects = append(b.glyphBatch.ClipRects,
		clip[0], clip[1], clip[2], clip[3],
		clip[0], clip[1], clip[2], clip[3],
		clip[0], clip[1], clip[2], clip[3],
		clip[0], clip[1], clip[2], clip[3],
	)

	c = c.MulRGB(c.A) // PMA
	b.glyphBatch.Colors = append(b.glyphBatch.Colors,
		c.R, c.G, c.B, c.A,
		c.R, c.G, c.B, c.A,
		c.R, c.G, c.B, c.A,
		c.R, c.G, c.B, c.A,
	)
	b.glyphBatch.Indices = append(b.glyphBatch.Indices,
		i, i+1, i+2,
		i+2, i+1, i+3,
	)
}

func (b *blitter) blitShape(ctx *context, shape shape, color gxui.Color, ds *drawState) {
	b.commitGlyphs(ctx)
	dipsToPixels := ctx.resolution.dipsToPixels()
	dw, dh := ctx.sizePixels.WH()
	mPos := math.CreateMat3(
		+2.0*dipsToPixels/float32(dw), 0, 0,
		0, -2.0*dipsToPixels/float32(dh), 0,
		-1.0+2.0*float32(ds.OriginPixels.X)/float32(dw),
		+1.0-2.0*float32(ds.OriginPixels.Y)/float32(dh), 1,
	)

	shape.draw(ctx, b.colorShader, uniformBindings{
		"mPos":  mPos,
		"Color": color,
	})

	if debugWireframePolygons {
		// glPolygonMode is not available in OpenGL ES/WebGL (since its implementation is very inefficient; a shame because it's useful for debugging).
		//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		shape.draw(ctx, b.colorShader, uniformBindings{
			"mPos":  mPos,
			"Color": gxui.Blue,
		})
		//gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}
	b.stats.drawCallCount++
}

func (b *blitter) blitRect(ctx *context, dstRect math.Rect, color gxui.Color, ds *drawState) {
	b.commitGlyphs(ctx)
	dstRect = dstRect.Offset(ds.OriginPixels)
	dw, dh := ctx.sizePixels.WH()
	mPos := math.CreateMat3(
		+2.0*float32(dstRect.W())/float32(dw), 0, 0,
		0, -2.0*float32(dstRect.H())/float32(dh), 0,
		-1.0+2.0*float32(dstRect.Min.X)/float32(dw),
		+1.0-2.0*float32(dstRect.Min.Y)/float32(dh), 1,
	)

	b.quad.draw(ctx, b.colorShader, uniformBindings{
		"mPos":  mPos,
		"Color": color,
	})
	b.stats.drawCallCount++
}

func (b *blitter) commit(ctx *context) {
	b.commitGlyphs(ctx)
}

func (b *blitter) commitGlyphs(ctx *context) {
	tc := b.glyphBatch.GlyphPage
	if tc == nil {
		return
	}
	sw, sh := tc.sizePixels.WH()
	dw, dh := ctx.sizePixels.WH()

	mSrc := math.CreateMat3(
		1.0/float32(sw), 0, 0,
		0, 1.0/float32(sh), 0,
		0.0, 0.0, 1,
	)
	mDst := math.CreateMat3(
		+2.0/float32(dw), 0, 0,
		0, -2.0/float32(dh), 0,
		-1.0, +1.0, 1,
	)
	vb := newVertexBuffer(
		newVertexStream("aDst", stFloatVec2, b.glyphBatch.DstRects),
		newVertexStream("aSrc", stFloatVec2, b.glyphBatch.SrcRects),
		newVertexStream("aClp", stFloatVec4, b.glyphBatch.ClipRects),
		newVertexStream("aCol", stFloatVec4, b.glyphBatch.Colors),
	)
	ib := newIndexBuffer(ptUshort, b.glyphBatch.Indices)
	s := newShape(vb, ib, dmTriangles)
	gl.Disable(gl.SCISSOR_TEST)
	s.draw(ctx, b.fontShader, uniformBindings{
		"source": tc,
		"mDst":   mDst,
		"mSrc":   mSrc,
	})
	gl.Enable(gl.SCISSOR_TEST)
	b.glyphBatch.GlyphPage = nil
	b.glyphBatch.DstRects = b.glyphBatch.DstRects[:0]
	b.glyphBatch.SrcRects = b.glyphBatch.SrcRects[:0]
	b.glyphBatch.ClipRects = b.glyphBatch.ClipRects[:0]
	b.glyphBatch.Colors = b.glyphBatch.Colors[:0]
	b.glyphBatch.Indices = b.glyphBatch.Indices[:0]
	b.stats.drawCallCount++
}
