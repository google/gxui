// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import "github.com/goxjs/gl"

type shape struct {
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
	return s
}

func newQuadShape() *shape {
	pos := newVertexStream("aPosition", stFloatVec2, []float32{
		0.0, 0.0,
		1.0, 0.0,
		0.0, 1.0,
		1.0, 1.0,
	})
	vb := newVertexBuffer(pos)
	ib := newIndexBuffer(ptUshort, []uint16{
		0, 1, 2,
		2, 1, 3,
	})
	return newShape(vb, ib, dmTriangles)
}

func (s shape) draw(ctx *context, shader *shaderProgram, ub uniformBindings) {
	shader.bind(ctx, s.vb, ub)
	if s.ib != nil {
		ctx.getOrCreateIndexBufferContext(s.ib).render(s.drawMode)
	} else {
		gl.DrawArrays(gl.Enum(s.drawMode), 0, s.vb.count)
	}
	shader.unbind(ctx)
	checkError()
}
