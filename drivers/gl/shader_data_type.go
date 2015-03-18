// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"

	"github.com/go-gl/gl/v3.2-core/gl"
)

type ShaderDataType int

const (
	FLOAT_MAT2x3 ShaderDataType = gl.FLOAT_MAT2x3
	FLOAT_MAT2x4 ShaderDataType = gl.FLOAT_MAT2x4
	FLOAT_MAT2   ShaderDataType = gl.FLOAT_MAT2
	FLOAT_MAT3x2 ShaderDataType = gl.FLOAT_MAT3x2
	FLOAT_MAT3x4 ShaderDataType = gl.FLOAT_MAT3x4
	FLOAT_MAT3   ShaderDataType = gl.FLOAT_MAT3
	FLOAT_MAT4x2 ShaderDataType = gl.FLOAT_MAT4x2
	FLOAT_MAT4x3 ShaderDataType = gl.FLOAT_MAT4x3
	FLOAT_MAT4   ShaderDataType = gl.FLOAT_MAT4
	FLOAT_VEC1   ShaderDataType = gl.FLOAT
	FLOAT_VEC2   ShaderDataType = gl.FLOAT_VEC2
	FLOAT_VEC3   ShaderDataType = gl.FLOAT_VEC3
	FLOAT_VEC4   ShaderDataType = gl.FLOAT_VEC4
	SAMPLER_2D   ShaderDataType = gl.SAMPLER_2D
)

func (s ShaderDataType) String() string {
	switch s {
	case FLOAT_MAT2x3:
		return "mat2x3"
	case FLOAT_MAT2x4:
		return "mat2x4"
	case FLOAT_MAT2:
		return "mat2"
	case FLOAT_MAT3x2:
		return "mat3x2"
	case FLOAT_MAT3x4:
		return "mat3x4"
	case FLOAT_MAT3:
		return "mat3"
	case FLOAT_MAT4x2:
		return "mat4x2"
	case FLOAT_MAT4x3:
		return "mat4x3"
	case FLOAT_MAT4:
		return "mat4"
	case FLOAT_VEC1:
		return "float"
	case FLOAT_VEC2:
		return "vec2"
	case FLOAT_VEC3:
		return "vec3"
	case FLOAT_VEC4:
		return "vec4"
	case SAMPLER_2D:
		return "sampler2D"
	default:
		return "unknown"
	}
}

func (s ShaderDataType) SizeInBytes() int {
	return s.VectorElementCount() * s.VectorElementType().SizeInBytes()
}

func (s ShaderDataType) VectorElementCount() int {
	switch s {
	case FLOAT_MAT2x3:
		return 2 * 3
	case FLOAT_MAT2x4:
		return 2 * 4
	case FLOAT_MAT2:
		return 2 * 2
	case FLOAT_MAT3x2:
		return 3 * 2
	case FLOAT_MAT3x4:
		return 3 * 4
	case FLOAT_MAT3:
		return 3 * 3
	case FLOAT_MAT4x2:
		return 4 * 2
	case FLOAT_MAT4x3:
		return 4 * 3
	case FLOAT_MAT4:
		return 4 * 4
	case FLOAT_VEC1:
		return 1
	case FLOAT_VEC2:
		return 2
	case FLOAT_VEC3:
		return 3
	case FLOAT_VEC4:
		return 4
	case SAMPLER_2D:
		return 1
	default:
		panic(fmt.Errorf("Unknown ShaderDataType 0x%.4x", s))
	}
}

func (s ShaderDataType) VectorElementType() PrimitiveType {
	switch s {
	case FLOAT_MAT2x3:
		return FLOAT
	case FLOAT_MAT2x4:
		return FLOAT
	case FLOAT_MAT2:
		return FLOAT
	case FLOAT_MAT3x2:
		return FLOAT
	case FLOAT_MAT3x4:
		return FLOAT
	case FLOAT_MAT3:
		return FLOAT
	case FLOAT_MAT4x2:
		return FLOAT
	case FLOAT_MAT4x3:
		return FLOAT
	case FLOAT_MAT4:
		return FLOAT
	case FLOAT_VEC1:
		return FLOAT
	case FLOAT_VEC2:
		return FLOAT
	case FLOAT_VEC3:
		return FLOAT
	case FLOAT_VEC4:
		return FLOAT
	case SAMPLER_2D:
		return INT
	default:
		panic(fmt.Errorf("Unknown ShaderDataType 0x%.4x", s))
	}
}
