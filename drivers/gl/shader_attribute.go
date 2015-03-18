// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import "github.com/go-gl/gl/v3.2-core/gl"

type shaderAttribute struct {
	name     string
	size     int
	ty       ShaderDataType
	location uint32
}

func (a shaderAttribute) EnableArray() {
	gl.EnableVertexAttribArray(a.location)
}

func (a shaderAttribute) DisableArray() {
	gl.DisableVertexAttribArray(a.location)
}

func (a shaderAttribute) AttribPointer(size int32, ty uint32, normalized bool, stride int32, data interface{}) {
	gl.VertexAttribPointer(a.location, size, ty, normalized, stride, gl.Ptr(data))
}
