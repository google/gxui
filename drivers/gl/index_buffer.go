// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"
	"reflect"

	"github.com/goxjs/gl"
)

type indexBuffer struct {
	data []byte
	ty   primitiveType
}

func newIndexBuffer(ty primitiveType, data16 []uint16) *indexBuffer {
	switch ty {
	case ptUbyte, ptUshort, ptUint:
		if !ty.isArrayOfType(data16) {
			panic(fmt.Errorf("Index data is not of type %v", ty))
		}
	default:
		panic(fmt.Errorf("Index type must be either UBYTE, USHORT or UINT. Got: %v", ty))
	}

	// HACK: Hardcode support for only ptUshort.
	data := make([]byte, len(data16)*2)
	for i, v := range data16 {
		data[2*i+0] = byte(v >> 0)
		data[2*i+1] = byte(v >> 8)
	}

	ib := &indexBuffer{
		data: data,
		ty:   ty,
	}
	return ib
}

func (b *indexBuffer) newContext() *indexBufferContext {
	dataVal := reflect.ValueOf(b.data)
	length := dataVal.Len() / 2 // HACK: Hardcode support for only ptUshort.

	buffer := gl.CreateBuffer()
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, b.data, gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, gl.Buffer{})
	checkError()

	globalStats.indexBufferContextCount.inc()
	return &indexBufferContext{
		buffer: buffer,
		ty:     b.ty,
		length: length,
	}
}

type indexBufferContext struct {
	contextResource
	buffer gl.Buffer
	ty     primitiveType
	length int
}

func (c *indexBufferContext) destroy() {
	globalStats.indexBufferContextCount.dec()
	gl.DeleteBuffer(c.buffer)
	c.buffer = gl.Buffer{}
}

func (c *indexBufferContext) render(drawMode drawMode) {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, c.buffer)
	gl.DrawElements(gl.Enum(drawMode), c.length, gl.Enum(c.ty), 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, gl.Buffer{})
	checkError()
}
