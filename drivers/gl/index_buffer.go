// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"
	"reflect"

	"github.com/go-gl/gl/v2.1/gl"
)

type indexBuffer struct {
	refCounted
	data interface{}
	ty   primitiveType
}

type indexBufferContext struct {
	buffer uint32
	ty     primitiveType
	length int
}

func newIndexBuffer(ty primitiveType, data interface{}) *indexBuffer {
	switch ty {
	case ptUbyte, ptUshort, ptUint:
		if !ty.isArrayOfType(data) {
			panic(fmt.Errorf("Index data is not of type %v", ty))
		}
	default:
		panic(fmt.Errorf("Index type must be either UBYTE, USHORT or UINT. Got: %v", ty))
	}

	ib := &indexBuffer{
		data: data,
		ty:   ty,
	}
	ib.init()
	globalStats.indexBufferCount.inc()
	return ib
}

func (b *indexBuffer) release() bool {
	if !b.refCounted.release() {
		return false
	}
	globalStats.indexBufferCount.dec()
	return true
}

func (b *indexBuffer) newContext() *indexBufferContext {
	dataVal := reflect.ValueOf(b.data)
	length := dataVal.Len()
	size := length * b.ty.sizeInBytes()

	var buffer uint32
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, size, gl.Ptr(b.data), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	checkError()

	return &indexBufferContext{
		buffer: buffer,
		ty:     b.ty,
		length: length,
	}
}

func (c *indexBufferContext) destroy() {
	gl.DeleteBuffers(1, &c.buffer)
	c.buffer = 0
}

func (c *indexBufferContext) render(drawMode drawMode) {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, c.buffer)
	gl.DrawElements(uint32(drawMode), int32(c.length), uint32(c.ty), nil)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	checkError()
}
