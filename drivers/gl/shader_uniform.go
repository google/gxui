// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"
	"github.com/google/gxui"
	"github.com/google/gxui/assert"
	"github.com/google/gxui/math"

	"github.com/go-gl-legacy/gl"
)

type ShaderUniform struct {
	Name        string
	Size        int
	Type        ShaderDataType
	Location    gl.UniformLocation
	TextureUnit int
}

func CreateShaderUniform(name string, size int, ty ShaderDataType, location gl.UniformLocation, textureUnit int) ShaderUniform {
	return ShaderUniform{
		Name:        name,
		Size:        size,
		Type:        ty,
		Location:    location,
		TextureUnit: textureUnit,
	}
}

func (u *ShaderUniform) Bind(context *Context, v interface{}) {
	transpose := true // UniformMatrix expects column-major, gxui is row-major
	switch u.Type {
	case FLOAT_MAT2x3:
		u.Location.UniformMatrix2x3fv(transpose, v.([6]float32))
	case FLOAT_MAT2x4:
		u.Location.UniformMatrix2x4fv(transpose, v.([8]float32))
	case FLOAT_MAT2:
		u.Location.UniformMatrix2fv(transpose, v.([4]float32))
	case FLOAT_MAT3x2:
		u.Location.UniformMatrix3x2fv(transpose, v.([6]float32))
	case FLOAT_MAT3x4:
		u.Location.UniformMatrix3x4fv(transpose, v.([12]float32))
	case FLOAT_MAT3:
		switch m := v.(type) {
		case math.Mat3:
			u.Location.UniformMatrix3fv(transpose, [9]float32(m))
		case [9]float32:
			u.Location.UniformMatrix3fv(transpose, m)
		}
	case FLOAT_MAT4x2:
		u.Location.UniformMatrix4x2fv(transpose, v.([8]float32))
	case FLOAT_MAT4x3:
		u.Location.UniformMatrix4x3fv(transpose, v.([12]float32))
	case FLOAT_MAT4:
		u.Location.UniformMatrix4fv(transpose, v.([16]float32))
	case FLOAT_VEC1:
		switch v := v.(type) {
		case float32:
			u.Location.Uniform1f(v)
		case []float32:
			u.Location.Uniform1fv(len(v), v)
		}
	case FLOAT_VEC2:
		switch v := v.(type) {
		case math.Vec2:
			u.Location.Uniform2fv(1, []float32{v.X, v.Y})
		case []float32:
			assert.True(len(v)%2 == 0, "Uniform '%s' of type vec2 should be an float32 array with a multiple of two length", u.Name)
			u.Location.Uniform2fv(len(v)/2, v)
		}
	case FLOAT_VEC3:
		switch v := v.(type) {
		case math.Vec3:
			u.Location.Uniform3fv(1, []float32{v.X, v.Y, v.Z})
		case []float32:
			assert.True(len(v)%3 == 0, "Uniform '%s' of type vec3 should be an float32 array with a multiple of three length", u.Name)
			u.Location.Uniform3fv(len(v)/3, v)
		}
	case FLOAT_VEC4:
		switch v := v.(type) {
		case math.Vec4:
			u.Location.Uniform4fv(1, []float32{v.X, v.Y, v.Z, v.W})
		case gxui.Color:
			u.Location.Uniform4fv(1, []float32{v.R, v.G, v.B, v.A})
		case []float32:
			assert.True(len(v)%4 == 0, "Uniform '%s' of type vec4 should be an float32 array with a multiple of four length", u.Name)
			u.Location.Uniform4fv(len(v)/4, v)
		}
	case SAMPLER_2D:
		ss := v.(SamplerSource)
		gl.ActiveTexture(gl.GLenum(gl.TEXTURE0 + u.TextureUnit))
		ss.Texture().Bind(gl.TEXTURE_2D)
		u.Location.Uniform1i(u.TextureUnit)
	default:
		panic(fmt.Errorf("Uniform of unsupported type %s", u.Type))
	}
}
