// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"github.com/go-gl-legacy/gl"
)

type ShaderAttribute struct {
	Name     string
	Size     int
	Type     ShaderDataType
	Location gl.AttribLocation
}

func CreateShaderAttribute(name string, size int, ty ShaderDataType, location gl.AttribLocation) ShaderAttribute {
	return ShaderAttribute{
		Name:     name,
		Size:     size,
		Type:     ty,
		Location: location,
	}
}
