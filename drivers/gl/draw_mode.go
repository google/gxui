// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"

	"github.com/go-gl-legacy/gl"
)

type DrawMode int

const (
	POINTS         DrawMode = gl.POINTS
	LINE_STRIP     DrawMode = gl.LINE_STRIP
	LINE_LOOP      DrawMode = gl.LINE_LOOP
	LINES          DrawMode = gl.LINES
	TRIANGLE_STRIP DrawMode = gl.TRIANGLE_STRIP
	TRIANGLE_FAN   DrawMode = gl.TRIANGLE_FAN
	TRIANGLES      DrawMode = gl.TRIANGLES
)

func (d DrawMode) PrimativeCount(vertexCount int) int {
	switch d {
	case POINTS:
		return vertexCount
	case LINE_STRIP:
		return vertexCount - 1
	case LINE_LOOP:
		return vertexCount
	case LINES:
		return vertexCount / 2
	case TRIANGLE_STRIP:
		return vertexCount - 2
	case TRIANGLE_FAN:
		return vertexCount - 2
	case TRIANGLES:
		return vertexCount / 3
	default:
		panic(fmt.Errorf("Unknown DrawMode 0x%.4x", d))
	}
}
