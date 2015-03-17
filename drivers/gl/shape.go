// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import "github.com/go-gl-legacy/gl"

type Shape struct {
	refCounted
	vb       *VertexBuffer
	ib       *IndexBuffer
	drawMode DrawMode
}

func CreateShape(vb *VertexBuffer, ib *IndexBuffer, drawMode DrawMode) *Shape {
	if vb == nil {
		panic("VertexBuffer cannot be nil")
	}

	s := &Shape{
		vb:       vb,
		ib:       ib,
		drawMode: drawMode,
	}
	s.init()
	globalStats.ShapeCount++
	return s
}

func (s *Shape) Release() {
	if s.release() {
		if s.vb != nil {
			s.vb.Release()
			s.vb = nil
		}
		if s.ib != nil {
			s.ib.Release()
			s.ib = nil
		}
		globalStats.ShapeCount--
	}
}

func CreateQuadShape() *Shape {
	pos := CreateVertexStream("aPosition", FLOAT_VEC2, []float32{
		0.0, 0.0,
		1.0, 0.0,
		0.0, 1.0,
		1.0, 1.0,
	})
	vb := CreateVertexBuffer(pos)
	ib := CreateIndexBuffer(UINT, []uint32{
		0, 1, 2,
		2, 1, 3,
	})
	return CreateShape(vb, ib, TRIANGLES)
}

func (s Shape) Draw(ctx *Context, shader *ShaderProgram, ub UniformBindings) {
	s.AssertAlive("Draw")

	shader.Bind(ctx, s.vb, ub)
	if s.ib != nil {
		ctx.GetOrCreateIndexBufferContext(s.ib).Render(s.drawMode)
	} else {
		gl.DrawArrays(gl.GLenum(s.drawMode), 0, s.vb.VertexCount)
	}
	CheckError()
	shader.Unbind(ctx)
}
