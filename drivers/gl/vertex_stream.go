// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"golang.org/x/mobile/f32"
	"golang.org/x/mobile/gl"
)

type vertexStream struct {
	refCounted
	name  string
	data  []byte
	ty    shaderDataType
	count int
}

type vertexStreamContext struct {
	buffer gl.Buffer
}

func newVertexStream(name string, ty shaderDataType, data32 []float32) *vertexStream {
	dataVal := reflect.ValueOf(data32)
	dataLen := dataVal.Len()

	if dataLen%ty.vectorElementCount() != 0 {
		panic(fmt.Errorf("Incorrect multiple of elements. Got: %d, Requires multiple of %d",
			dataLen, ty.vectorElementCount()))
	}
	if !ty.vectorElementType().isArrayOfType(data32) {
		panic("Data is not of the specified type")
	}

	// HACK.
	data := f32.Bytes(binary.LittleEndian, data32...)

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
	buffer := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.BufferData(gl.ARRAY_BUFFER, s.data, gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, gl.Buffer{})
	checkError()

	return &vertexStreamContext{buffer}
}

func (c vertexStreamContext) bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, c.buffer)
}

func (c *vertexStreamContext) destroy() {
	gl.DeleteBuffer(c.buffer)
	c.buffer = gl.Buffer{}
}
