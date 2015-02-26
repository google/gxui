// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"
	"gxui/assert"
	"reflect"

	"github.com/go-gl/gl"
)

type IndexBuffer struct {
	refCounted
	data interface{}
	ty   PrimitiveType
}

type IndexBufferContext struct {
	buffer gl.Buffer
	ty     PrimitiveType
	length int
}

func CreateIndexBuffer(ty PrimitiveType, data interface{}) *IndexBuffer {
	switch ty {
	case UBYTE, USHORT, UINT:
		assert.True(ty.IsArrayOfType(data), "Index data is not of the specified type")
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

	buffer := gl.GenBuffer()
	buffer.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, size, b.data, gl.STATIC_DRAW)
	CheckError()
	buffer.Unbind(gl.ELEMENT_ARRAY_BUFFER)

	return IndexBufferContext{
		buffer: buffer,
		ty:     b.ty,
		length: length,
	}
}

func (c *IndexBufferContext) Destroy() {
	c.buffer.Delete()
	c.buffer = gl.Buffer(0)
}

func (c IndexBufferContext) Render(drawMode DrawMode) {
	c.buffer.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.DrawElements(gl.GLenum(drawMode), c.length, gl.GLenum(c.ty), nil)
	c.buffer.Unbind(gl.ELEMENT_ARRAY_BUFFER)
	CheckError()
}
