// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"
	"reflect"

	"github.com/go-gl/gl/v3.2-core/gl"
)

type VertexStream struct {
	refCounted
	name  string
	data  interface{}
	ty    ShaderDataType
	count int
}

type VertexStreamContext struct {
	buffer uint32
}

func CreateVertexStream(name string, ty ShaderDataType, data interface{}) *VertexStream {
	dataVal := reflect.ValueOf(data)
	dataLen := dataVal.Len()

	if dataLen%ty.VectorElementCount() != 0 {
		panic(fmt.Errorf("Incorrect multiple of elements. Got: %d, Requires multiple of %d",
			dataLen, ty.VectorElementCount()))
	}
	if !ty.VectorElementType().IsArrayOfType(data) {
		panic("Data is not of the specified type")
	}

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

	var buffer uint32
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.BufferData(gl.ARRAY_BUFFER, size, gl.Ptr(s.data), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	CheckError()

	return VertexStreamContext{buffer}
}

func (c VertexStreamContext) Bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, c.buffer)
}

func (c *VertexStreamContext) Destroy() {
	gl.DeleteBuffers(1, &c.buffer)
	c.buffer = 0
}
