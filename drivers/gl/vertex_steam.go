// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"github.com/google/gxui/assert"
	"reflect"

	"github.com/go-gl-legacy/gl"
)

type VertexStream struct {
	refCounted
	name  string
	data  interface{}
	ty    ShaderDataType
	count int
}

type VertexStreamContext struct {
	buffer gl.Buffer
}

func CreateVertexStream(name string, ty ShaderDataType, data interface{}) *VertexStream {
	dataVal := reflect.ValueOf(data)
	dataLen := dataVal.Len()

	assert.True(dataLen%ty.VectorElementCount() == 0,
		"Incorrect multiple of elements. Got: %d, Requires multiple of %d",
		dataLen, ty.VectorElementCount())
	assert.True(ty.VectorElementType().IsArrayOfType(data), "Data is not of the specified type")

	vs := &VertexStream{
		name:  name,
		data:  data,
		ty:    ty,
		count: dataLen / ty.VectorElementCount(),
	}
	vs.init()
	globalStats.VertexStreamCount++
	return vs
}

func (s *VertexStream) Release() {
	if s.release() {
		globalStats.VertexStreamCount--
	}
}

func (s VertexStream) Name() string {
	return s.name
}

func (s VertexStream) Type() ShaderDataType {
	return s.ty
}

func (s VertexStream) VertexCount() int {
	return s.count
}

func (s VertexStream) CreateContext() VertexStreamContext {
	dataVal := reflect.ValueOf(s.data)
	dataLen := dataVal.Len()
	size := dataLen * s.ty.VectorElementType().SizeInBytes()

	buffer := gl.GenBuffer()
	buffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, size, s.data, gl.STATIC_DRAW)
	CheckError()
	buffer.Unbind(gl.ARRAY_BUFFER)

	return VertexStreamContext{buffer}
}

func (c VertexStreamContext) Bind() {
	c.buffer.Bind(gl.ARRAY_BUFFER)
}

func (c *VertexStreamContext) Destroy() {
	c.buffer.Delete()
	c.buffer = gl.Buffer(0)
}
