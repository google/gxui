// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"gaze/gxui/math"
	"github.com/go-gl/gl"
	"sync"
)

type SamplerSource interface {
	Texture() gl.Texture
	SizePixels() math.Size
	FlipY() bool
	PremultipliedAlpha() bool
}

type Context struct {
	lastStatsMutex       sync.Mutex
	Blitter              *Blitter
	currStats, lastStats Stats
	textureContexts      map[*Texture]TextureContext
	vertexStreamContexts map[*VertexStream]VertexStreamContext
	indexBufferContexts  map[*IndexBuffer]IndexBufferContext
	renderTarget         *Canvas
	framebufferPool      *FramebufferPool
	sizeDips, sizePixels math.Size
	clip                 math.Rect
}

func CreateContext() *Context {
	ctx := &Context{}
	ctx.Blitter = CreateBlitter(ctx, &ctx.currStats)
	ctx.textureContexts = make(map[*Texture]TextureContext)
	ctx.vertexStreamContexts = make(map[*VertexStream]VertexStreamContext)
	ctx.indexBufferContexts = make(map[*IndexBuffer]IndexBufferContext)
	ctx.framebufferPool = CreateFramebufferPool(256*1024*1024, &ctx.currStats)

	return ctx
}

func (c *Context) Destroy() {
	for texture, tc := range c.textureContexts {
		delete(c.textureContexts, texture)
		tc.Destroy()
		c.currStats.TextureCount--
	}
	for stream, sc := range c.vertexStreamContexts {
		delete(c.vertexStreamContexts, stream)
		sc.Destroy()
		c.currStats.VertexStreamCount--
	}
	for buffer, ic := range c.indexBufferContexts {
		delete(c.indexBufferContexts, buffer)
		ic.Destroy()
		c.currStats.IndexBufferCount--
	}
	c.Blitter.Destroy(c)
	c.Blitter = nil
}

func (c *Context) Stats() *Stats {
	return &c.currStats
}

func (c *Context) LastStats() Stats {
	c.lastStatsMutex.Lock()
	defer c.lastStatsMutex.Unlock()
	return c.lastStats
}

func (c *Context) BeginDraw(sizeDips, sizePixels math.Size) {
	// Reap any dead textures
	for texture, tc := range c.textureContexts {
		if !texture.Alive() {
			delete(c.textureContexts, texture)
			tc.Destroy()
			c.currStats.TextureCount--
		}
	}
	for stream, sc := range c.vertexStreamContexts {
		if !stream.Alive() {
			delete(c.vertexStreamContexts, stream)
			sc.Destroy()
			c.currStats.VertexStreamCount--
		}
	}
	for buffer, ic := range c.indexBufferContexts {
		if !buffer.Alive() {
			delete(c.indexBufferContexts, buffer)
			ic.Destroy()
			c.currStats.IndexBufferCount--
		}
	}

	c.sizeDips = sizeDips
	c.sizePixels = sizePixels
	c.currStats.DrawCallCount = 0
	c.currStats.Timer("Frame").Start()
}

func (c *Context) EndDraw() {
	c.currStats.Timer("Frame").Stop()
	c.currStats.FrameCount++

	c.lastStatsMutex.Lock()
	c.lastStats = c.currStats
	c.lastStatsMutex.Unlock()
}

func (c *Context) GetOrCreateTextureContext(t *Texture) TextureContext {
	tc, found := c.textureContexts[t]
	if !found {
		tc = t.CreateContext()
		c.textureContexts[t] = tc
		c.currStats.TextureCount++
	}
	return tc
}

func (c *Context) GetOrCreateVertexStreamContext(vs *VertexStream) VertexStreamContext {
	vc, found := c.vertexStreamContexts[vs]
	if !found {
		vc = vs.CreateContext()
		c.vertexStreamContexts[vs] = vc
		c.currStats.VertexStreamCount++
	}
	return vc
}

func (c *Context) GetOrCreateIndexBufferContext(ib *IndexBuffer) IndexBufferContext {
	ic, found := c.indexBufferContexts[ib]
	if !found {
		ic = ib.CreateContext()
		c.indexBufferContexts[ib] = ic
		c.currStats.IndexBufferCount++
	}
	return ic
}

func (c *Context) Apply(ds *DrawState) {
	r := ds.ClipPixels
	o := c.clip
	if o != r {
		c.clip = r
		vs := c.sizePixels
		rs := r.Size()
		gl.Scissor(r.Min.X, vs.H-r.Max.Y, rs.W, rs.H)
	}
}

func (c *Context) RenderTargetSizePixels() math.Size {
	return c.sizePixels
}

func (c *Context) DipsToPixels() float32 {
	wd := c.sizeDips.W
	wp := c.sizePixels.W
	return float32(wp) / float32(wd)
}

func (c *Context) PointDipsToPixels(dips math.Point) math.Point {
	sd := c.sizeDips
	sp := c.sizePixels
	return math.Point{
		X: (dips.X * sp.W) / sd.W,
		Y: (dips.Y * sp.H) / sd.H,
	}
}

func (c *Context) RectDipsToPixels(dips math.Rect) math.Rect {
	return math.Rect{
		Min: c.PointDipsToPixels(dips.Min),
		Max: c.PointDipsToPixels(dips.Max),
	}
}

func (c *Context) Resolution() Resolution {
	return Resolution(c.DipsToPixels()*0xffff + 0.5)
}

func (c *Context) SizeDipsToPixels(dips math.Size) math.Size {
	sd := c.sizeDips
	sp := c.sizePixels
	return math.Size{
		W: (dips.W * sp.W) / sd.W,
		H: (dips.H * sp.H) / sd.H,
	}
}

func (c *Context) SizePixelsToDips(dips math.Size) math.Size {
	sd := c.sizeDips
	sp := c.sizePixels
	return math.Size{
		W: (dips.W * sd.W) / sp.W,
		H: (dips.H * sd.H) / sp.H,
	}
}
