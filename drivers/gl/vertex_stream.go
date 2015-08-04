// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"
	"math"
	"reflect"

	"github.com/goxjs/gl"
)

type vertexStream struct {
	name  string
	data  []byte
	ty    shaderDataType
	count int
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
	data := float32Bytes(data32...)

	vs := &vertexStream{
		name:  name,
		data:  data,
		ty:    ty,
		count: dataLen / ty.vectorElementCount(),
	}
	return vs
}

func (s *vertexStream) newContext() *vertexStreamContext {
	buffer := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.BufferData(gl.ARRAY_BUFFER, s.data, gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, gl.Buffer{})
	checkError()

	globalStats.vertexStreamContextCount.inc()
	return &vertexStreamContext{buffer: buffer}
}

type vertexStreamContext struct {
	contextResource
	buffer gl.Buffer
}

func (c *vertexStreamContext) bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, c.buffer)
}

func (c *vertexStreamContext) destroy() {
	globalStats.vertexStreamContextCount.dec()
	gl.DeleteBuffer(c.buffer)
	c.buffer = gl.Buffer{}
}

// float32Bytes returns the byte representation of float32 values in little endian byte order.
func float32Bytes(values ...float32) []byte {
	b := make([]byte, 4*len(values))
	for i, v := range values {
		u := math.Float32bits(v)
		b[4*i+0] = byte(u >> 0)
		b[4*i+1] = byte(u >> 8)
		b[4*i+2] = byte(u >> 16)
		b[4*i+3] = byte(u >> 24)
	}
	return b
}
