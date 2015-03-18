// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"
	"reflect"

	"github.com/go-gl/gl/v3.2-core/gl"
)

type IndexBuffer struct {
	refCounted
	data interface{}
	ty   PrimitiveType
}

type IndexBufferContext struct {
	buffer uint32
	ty     PrimitiveType
	length int
}

func CreateIndexBuffer(ty PrimitiveType, data interface{}) *IndexBuffer {
	switch ty {
	case UBYTE, USHORT, UINT:
		if !ty.IsArrayOfType(data) {
			panic(fmt.Errorf("Index data is not of type %v", ty))
		}
	default:
		panic(fmt.Errorf("Index type must be either UBYTE, USHORT or UINT. Got: %v", ty))
	}

	ib := &IndexBuffer{
		data: data,
		ty:   ty,
	}
	ib.init()
	globalStats.IndexBufferCount++
	return ib
}

func (b *IndexBuffer) Release() {
	if b.release() {
		globalStats.IndexBufferCount--
	}
}

func (b IndexBuffer) CreateContext() IndexBufferContext {
	dataVal := reflect.ValueOf(b.data)
	length := dataVal.Len()
	size := length * b.ty.SizeInBytes()

	var buffer uint32
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, size, gl.Ptr(b.data), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	CheckError()

	return IndexBufferContext{
		buffer: buffer,
		ty:     b.ty,
		length: length,
	}
}

func (c *IndexBufferContext) Destroy() {
	gl.DeleteBuffers(1, &c.buffer)
	c.buffer = 0
}

func (c IndexBufferContext) Render(drawMode DrawMode) {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, c.buffer)
	gl.DrawElements(uint32(drawMode), int32(c.length), uint32(c.ty), nil)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	CheckError()
}
