// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"container/list"
	"fmt"
	"github.com/google/gxui/math"

	"github.com/go-gl/gl/v3.2-core/gl"
)

type Framebuffer struct {
	dimensions  math.Size
	framebuffer uint32
	texture     uint32
}

func CreateFramebuffer(dimensions math.Size) *Framebuffer {
	var framebuffer uint32
	gl.GenFramebuffers(1, &framebuffer)

	var texture uint32
	gl.GenTextures(1, &texture)

	f := &Framebuffer{
		dimensions:  dimensions,
		framebuffer: framebuffer,
		texture:     texture,
	}

	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(dimensions.W), int32(dimensions.H), 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.BindFramebuffer(gl.FRAMEBUFFER, framebuffer)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, f.texture, 0)
	status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)
	if status != gl.FRAMEBUFFER_COMPLETE {
		panic(fmt.Errorf("CheckFramebufferStatus returned 0x%.4x for size %+v", status, dimensions))
	}
	gl.BindTexture(gl.TEXTURE_2D, 0)

	return f
}

func (f Framebuffer) SizeBytes() int {
	return f.dimensions.W * f.dimensions.H * 4
}

func (f Framebuffer) Delete() {
	gl.DeleteFramebuffers(1, &f.framebuffer)
	gl.DeleteTextures(1, &f.texture)
	f.framebuffer = 0
	f.texture = 0
}

type FramebufferPool struct {
	stats                *Stats
	used                 *list.List
	free                 *list.List
	bytesAllocated       int
	targetBytesAllocated int
}

func CreateFramebufferPool(targetBytesAllocated int, stats *Stats) *FramebufferPool {
	return &FramebufferPool{
		stats:                stats,
		used:                 list.New(),
		free:                 list.New(),
		targetBytesAllocated: targetBytesAllocated,
	}
}

func (p *FramebufferPool) Acquire(dimensions math.Size) *Framebuffer {
	for e := p.free.Front(); e != nil; e = e.Next() {
		f := e.Value.(*Framebuffer)
		if f.dimensions == dimensions {
			p.free.Remove(e)
			p.used.PushFront(f)
			p.updateStats()
			return f
		}
	}

	f := CreateFramebuffer(dimensions)
	p.used.PushFront(f)
	p.bytesAllocated += f.SizeBytes()
	p.reap()
	p.updateStats()
	return f
}

func (p *FramebufferPool) Release(framebuffer *Framebuffer) {
	for e := p.used.Front(); e != nil; e = e.Next() {
		f := e.Value.(*Framebuffer)
		if f == framebuffer {
			p.used.Remove(e)
			p.free.PushFront(f)
			p.reap()
			p.updateStats()
			return
		}
	}
	panic("Framebuffer is not part of the pool")
}

func (p *FramebufferPool) reap() {
	for p.bytesAllocated > p.targetBytesAllocated {
		if p.free.Len() == 0 {
			return
		}
		e := p.free.Back()
		f := e.Value.(*Framebuffer)
		p.bytesAllocated -= f.SizeBytes()
		f.Delete()
		p.free.Remove(e)
	}
}

func (p *FramebufferPool) updateStats() {
	p.stats.FramebufferUsedCount = p.used.Len()
	p.stats.FramebufferFreeCount = p.free.Len()
	p.stats.FramebufferBytesAllocated = p.bytesAllocated
}
