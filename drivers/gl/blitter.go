// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"github.com/google/gxui"
	"github.com/google/gxui/math"

	"github.com/go-gl/gl/v3.2-core/gl"
)

const debugWireframePolygons = false

var COPY_VS_SOURCE = `
  attribute vec2 aPosition;
  varying vec2 vTexcoords;
  uniform mat3 mPos;
  uniform mat3 mUV;
  void main() {
  	vec3 pos3 = vec3(aPosition, 1.0);
    gl_Position = vec4(mPos * pos3, 1.0);
    vTexcoords = (mUV * pos3).xy;
  }
`
var COPY_FS_SOURCE = `
  uniform sampler2D source;
  varying vec2 vTexcoords;
  void main() {
    gl_FragColor = texture2D(source, vTexcoords);
  }
`

var COLOR_VS_SOURCE = `
  attribute vec2 aPosition;
  uniform mat3 mPos;
  void main() {
  	vec3 pos3 = vec3(aPosition, 1.0);
    gl_Position = vec4((mPos * pos3).xy, 0.0, 1.0);
  }
`

var COLOR_FS_SOURCE = `
  uniform vec4 Color;
  void main() {
    gl_FragColor = Color;
    gl_FragColor *= gl_FragColor.a; // PMA
  }
`

var FONT_VS_SOURCE = `
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
  }
`

var FONT_FS_SOURCE = `
  uniform sampler2D source;
  varying vec2 vSrc;
  varying vec4 vCol;
  varying vec2 vClp;
  void main() {
    vec2 clipping = step(vec2(0.0, 0.0), vClp) * step(vClp, vec2(1.0, 1.0));
    gl_FragColor  = vCol * texture2D(source, vSrc).aaaa;
    gl_FragColor *= clipping.x * clipping.y;
  }
`

type GlyphBatch struct {
	DstRects  []float32
	SrcRects  []float32
	Colors    []float32
	ClipRects []float32
	Indices   []uint32
	GlyphPage SamplerSource
}

type Blitter struct {
	stats       *Stats
	quad        *Shape
	copyShader  *ShaderProgram
	colorShader *ShaderProgram
	fontShader  *ShaderProgram
	glyphBatch  GlyphBatch
}

func CreateBlitter(ctx *Context, stats *Stats) *Blitter {
	return &Blitter{
		stats:       stats,
		quad:        CreateQuadShape(),
		copyShader:  CreateShaderProgram(ctx, COPY_VS_SOURCE, COPY_FS_SOURCE),
		colorShader: CreateShaderProgram(ctx, COLOR_VS_SOURCE, COLOR_FS_SOURCE),
		fontShader:  CreateShaderProgram(ctx, FONT_VS_SOURCE, FONT_FS_SOURCE),
	}
}

func (b *Blitter) Destroy(ctx *Context) {
	b.quad.Release()
	b.copyShader.Destroy(ctx)
	b.colorShader.Destroy(ctx)
	b.fontShader.Destroy(ctx)
}

func (b *Blitter) Blit(ctx *Context, ss SamplerSource, srcRect, dstRect math.Rect, ds *DrawState) {
	b.CommitGlyphs(ctx)

	dstRect = dstRect.Offset(ds.OriginPixels)
	sw, sh := ss.SizePixels().WH()
	dw, dh := ctx.RenderTargetSizePixels().WH()

	var mUV math.Mat3
	if ss.FlipY() {
		mUV = math.CreateMat3(
			float32(srcRect.W())/float32(sw), 0, 0,
			0, float32(srcRect.H())/float32(sh), 0,
			float32(srcRect.Min.X)/float32(sw),
			float32(srcRect.Min.Y)/float32(sh), 1,
		)
	} else {
		mUV = math.CreateMat3(
			float32(srcRect.W())/float32(sw), 0, 0,
			0, -float32(srcRect.H())/float32(sh), 0,
			float32(srcRect.Min.X)/float32(sw),
			1.0-float32(srcRect.Min.Y)/float32(sh), 1,
		)
	}
	mPos := math.CreateMat3(
		+2.0*float32(dstRect.W())/float32(dw), 0, 0,
		0, -2.0*float32(dstRect.H())/float32(dh), 0,
		-1.0+2.0*float32(dstRect.Min.X)/float32(dw),
		+1.0-2.0*float32(dstRect.Min.Y)/float32(dh), 1,
	)
	if !ss.PremultipliedAlpha() {
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	}
	b.quad.Draw(ctx, b.copyShader, UniformBindings{
		"source": ss,
		"mUV":    mUV,
		"mPos":   mPos,
	})
	if !ss.PremultipliedAlpha() {
		gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	}
	b.stats.DrawCallCount++
}

func (b *Blitter) BlitGlyph(ctx *Context, ss SamplerSource, c gxui.Color, srcRect, dstRect math.Rect, ds *DrawState) {
	dstRect = dstRect.Offset(ds.OriginPixels)

	if b.glyphBatch.GlyphPage != ss {
		b.CommitGlyphs(ctx)
		b.glyphBatch.GlyphPage = ss
	}
	i := uint32(len(b.glyphBatch.DstRects)) / 2
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

func (b *Blitter) BlitShape(ctx *Context, shape Shape, color gxui.Color, ds *DrawState) {
	b.CommitGlyphs(ctx)
	dipsToPixels := ctx.DipsToPixels()
	dw, dh := ctx.RenderTargetSizePixels().WH()
	mPos := math.CreateMat3(
		+2.0*dipsToPixels/float32(dw), 0, 0,
		0, -2.0*dipsToPixels/float32(dh), 0,
		-1.0+2.0*float32(ds.OriginPixels.X)/float32(dw),
		+1.0-2.0*float32(ds.OriginPixels.Y)/float32(dh), 1,
	)

	shape.Draw(ctx, b.colorShader, UniformBindings{
		"mPos":  mPos,
		"Color": color,
	})

	if debugWireframePolygons {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		shape.Draw(ctx, b.colorShader, UniformBindings{
			"mPos":  mPos,
			"Color": gxui.Blue,
		})
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}
	b.stats.DrawCallCount++
}

func (b *Blitter) BlitRect(ctx *Context, dstRect math.Rect, color gxui.Color, ds *DrawState) {
	b.CommitGlyphs(ctx)
	dstRect = dstRect.Offset(ds.OriginPixels)
	dw, dh := ctx.RenderTargetSizePixels().WH()
	mPos := math.CreateMat3(
		+2.0*float32(dstRect.W())/float32(dw), 0, 0,
		0, -2.0*float32(dstRect.H())/float32(dh), 0,
		-1.0+2.0*float32(dstRect.Min.X)/float32(dw),
		+1.0-2.0*float32(dstRect.Min.Y)/float32(dh), 1,
	)

	b.quad.Draw(ctx, b.colorShader, UniformBindings{
		"mPos":  mPos,
		"Color": color,
	})
	b.stats.DrawCallCount++
}

func (b *Blitter) Commit(ctx *Context) {
	b.CommitGlyphs(ctx)
}

func (b *Blitter) CommitGlyphs(ctx *Context) {
	ss := b.glyphBatch.GlyphPage
	if ss == nil {
		return
	}
	sw, sh := ss.SizePixels().WH()
	dw, dh := ctx.RenderTargetSizePixels().WH()

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
	vb := CreateVertexBuffer(
		CreateVertexStream("aDst", FLOAT_VEC2, b.glyphBatch.DstRects),
		CreateVertexStream("aSrc", FLOAT_VEC2, b.glyphBatch.SrcRects),
		CreateVertexStream("aClp", FLOAT_VEC4, b.glyphBatch.ClipRects),
		CreateVertexStream("aCol", FLOAT_VEC4, b.glyphBatch.Colors),
	)
	ib := CreateIndexBuffer(UINT, b.glyphBatch.Indices)
	s := CreateShape(vb, ib, TRIANGLES)
	gl.Disable(gl.SCISSOR_TEST)
	s.Draw(ctx, b.fontShader, UniformBindings{
		"source": ss,
		"mDst":   mDst,
		"mSrc":   mSrc,
	})
	gl.Enable(gl.SCISSOR_TEST)
	s.Release()
	b.glyphBatch.GlyphPage = nil
	b.glyphBatch.DstRects = b.glyphBatch.DstRects[:0]
	b.glyphBatch.SrcRects = b.glyphBatch.SrcRects[:0]
	b.glyphBatch.ClipRects = b.glyphBatch.ClipRects[:0]
	b.glyphBatch.Colors = b.glyphBatch.Colors[:0]
	b.glyphBatch.Indices = b.glyphBatch.Indices[:0]
	b.stats.DrawCallCount++
}
