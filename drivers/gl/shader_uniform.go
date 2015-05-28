// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"

	"github.com/google/gxui"
	"github.com/google/gxui/math"
	"github.com/goxjs/gl"
)

type shaderUniform struct {
	name        string
	size        int
	ty          shaderDataType
	location    gl.Uniform
	textureUnit int
}

func (u *shaderUniform) bind(context *context, v interface{}) {
	switch u.ty {
	case stFloatMat2:
		gl.UniformMatrix2fv(u.location, v.([]float32))
	case stFloatMat3:
		switch m := v.(type) {
		case math.Mat3:
			gl.UniformMatrix3fv(u.location, m[:])
		case []float32:
			gl.UniformMatrix3fv(u.location, m)
		}
	case stFloatMat4:
		gl.UniformMatrix4fv(u.location, v.([]float32))
	case stFloatVec1:
		switch v := v.(type) {
		case float32:
			gl.Uniform1f(u.location, v)
		case []float32:
			gl.Uniform1fv(u.location, v)
		}
	case stFloatVec2:
		switch v := v.(type) {
		case math.Vec2:
			gl.Uniform2fv(u.location, []float32{v.X, v.Y})
		case []float32:
			if len(v)%2 != 0 {
				panic(fmt.Errorf("Uniform '%s' of type vec2 should be an float32 array with a multiple of two length", u.name))
			}
			gl.Uniform2fv(u.location, v)
		}
	case stFloatVec3:
		switch v := v.(type) {
		case math.Vec3:
			gl.Uniform3fv(u.location, []float32{v.X, v.Y, v.Z})
		case []float32:
			if len(v)%3 != 0 {
				panic(fmt.Errorf("Uniform '%s' of type vec3 should be an float32 array with a multiple of three length", u.name))
			}
			gl.Uniform3fv(u.location, v)
		}
	case stFloatVec4:
		switch v := v.(type) {
		case math.Vec4:
			gl.Uniform4fv(u.location, []float32{v.X, v.Y, v.Z, v.W})
		case gxui.Color:
			gl.Uniform4fv(u.location, []float32{v.R, v.G, v.B, v.A})
		case []float32:
			if len(v)%4 != 0 {
				panic(fmt.Errorf("Uniform '%s' of type vec4 should be an float32 array with a multiple of four length", u.name))
			}
			gl.Uniform4fv(u.location, v)
		}
	case stSampler2d:
		tc := v.(*textureContext)
		gl.ActiveTexture(gl.Enum(gl.TEXTURE0 + u.textureUnit))
		gl.BindTexture(gl.TEXTURE_2D, tc.texture)
		gl.Uniform1i(u.location, u.textureUnit)
	default:
		panic(fmt.Errorf("Uniform of unsupported type %s", u.ty))
	}
}
