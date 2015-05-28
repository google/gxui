// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import "github.com/goxjs/gl"

type shaderAttribute struct {
	name     string
	size     int
	ty       shaderDataType
	location gl.Attrib
}

func (a shaderAttribute) enableArray() {
	gl.EnableVertexAttribArray(a.location)
}

func (a shaderAttribute) disableArray() {
	gl.DisableVertexAttribArray(a.location)
}

func (a shaderAttribute) attribPointer(size int32, ty uint32, normalized bool, stride int32, offset int) {
	gl.VertexAttribPointer(a.location, int(size), gl.Enum(ty), normalized, int(stride), offset)
}
