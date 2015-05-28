// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"

	"github.com/goxjs/gl"
)

type drawMode int

const (
	dmPoints        drawMode = gl.POINTS
	dmLineStrip     drawMode = gl.LINE_STRIP
	dmLineLoop      drawMode = gl.LINE_LOOP
	dmLines         drawMode = gl.LINES
	dmTriangleStrip drawMode = gl.TRIANGLE_STRIP
	dmTriangleFan   drawMode = gl.TRIANGLE_FAN
	dmTriangles     drawMode = gl.TRIANGLES
)

func (d drawMode) primitiveCount(vertexCount int) int {
	switch d {
	case dmPoints:
		return vertexCount
	case dmLineStrip:
		return vertexCount - 1
	case dmLineLoop:
		return vertexCount
	case dmLines:
		return vertexCount / 2
	case dmTriangleStrip:
		return vertexCount - 2
	case dmTriangleFan:
		return vertexCount - 2
	case dmTriangles:
		return vertexCount / 3
	default:
		panic(fmt.Errorf("Unknown drawMode 0x%.4x", d))
	}
}
