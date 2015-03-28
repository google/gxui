// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"
	"reflect"

	"github.com/go-gl/gl/v2.1/gl"
)

type vertexStream struct {
	refCounted
	name  string
	data  interface{}
	ty    shaderDataType
	count int
}

type vertexStreamContext struct {
	buffer uint32
}

func newVertexStream(name string, ty shaderDataType, data interface{}) *vertexStream {
	dataVal := reflect.ValueOf(data)
	dataLen := dataVal.Len()

	if dataLen%ty.vectorElementCount() != 0 {
		panic(fmt.Errorf("Incorrect multiple of elements. Got: %d, Requires multiple of %d",
			dataLen, ty.vectorElementCount()))
	}
	if !ty.vectorElementType().isArrayOfType(data) {
		panic("Data is not of the specified type")
	}

	vs := &vertexStream{
		name:  name,
		data:  data,
		ty:    ty,
		count: dataLen / ty.vectorElementCount(),
	}
	vs.init()
	globalStats.vertexStreamCount.inc()
	return vs
}

func (s *vertexStream) release() bool {
	if !s.refCounted.release() {
		return false
	}
	globalStats.vertexStreamCount.dec()
	return true
}

func (s *vertexStream) newContext() *vertexStreamContext {
	dataVal := reflect.ValueOf(s.data)
	dataLen := dataVal.Len()
	size := dataLen * s.ty.vectorElementType().sizeInBytes()

	var buffer uint32
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.BufferData(gl.ARRAY_BUFFER, size, gl.Ptr(s.data), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	checkError()

	return &vertexStreamContext{buffer}
}

func (c vertexStreamContext) bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, c.buffer)
}

func (c *vertexStreamContext) destroy() {
	gl.DeleteBuffers(1, &c.buffer)
	c.buffer = 0
}
