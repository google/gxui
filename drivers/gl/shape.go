// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import "github.com/go-gl/gl/v2.1/gl"

type shape struct {
	refCounted
	vb       *vertexBuffer
	ib       *indexBuffer
	drawMode drawMode
}

func newShape(vb *vertexBuffer, ib *indexBuffer, drawMode drawMode) *shape {
	if vb == nil {
		panic("VertexBuffer cannot be nil")
	}

	s := &shape{
		vb:       vb,
		ib:       ib,
		drawMode: drawMode,
	}
	s.init()
	globalStats.shapeCount.inc()
	return s
}

func (s *shape) release() bool {
	if !s.refCounted.release() {
		return false
	}
	if s.vb != nil {
		s.vb.release()
		s.vb = nil
	}
	if s.ib != nil {
		s.ib.release()
		s.ib = nil
	}
	globalStats.shapeCount.dec()
	return true
}

func newQuadShape() *shape {
	pos := newVertexStream("aPosition", stFloatVec2, []float32{
		0.0, 0.0,
		1.0, 0.0,
		0.0, 1.0,
		1.0, 1.0,
	})
	vb := newVertexBuffer(pos)
	ib := newIndexBuffer(ptUint, []uint32{
		0, 1, 2,
		2, 1, 3,
	})
	return newShape(vb, ib, dmTriangles)
}

func (s shape) draw(ctx *context, shader *shaderProgram, ub uniformBindings) {
	s.assertAlive("draw")

	shader.bind(ctx, s.vb, ub)
	if s.ib != nil {
		ctx.getOrCreateIndexBufferContext(s.ib).render(s.drawMode)
	} else {
		gl.DrawArrays(uint32(s.drawMode), 0, int32(s.vb.count))
	}
	shader.unbind(ctx)
	checkError()
}
