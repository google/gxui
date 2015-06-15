// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"

	"github.com/goxjs/gl"
)

type uniformBindings map[string]interface{}

type shaderProgram struct {
	program    gl.Program
	uniforms   []shaderUniform
	attributes []shaderAttribute
}

func compile(source string, ty int) gl.Shader {
	shader := gl.CreateShader(gl.Enum(ty))
	gl.ShaderSource(shader, source)

	gl.CompileShader(shader)
	if gl.GetShaderi(shader, gl.COMPILE_STATUS) != gl.TRUE {
		panic(gl.GetShaderInfoLog(shader))
	}
	checkError()

	return shader
}

func newShaderProgram(ctx *context, vsSource, fsSource string) *shaderProgram {
	vs := compile(vsSource, gl.VERTEX_SHADER)
	fs := compile(fsSource, gl.FRAGMENT_SHADER)

	program := gl.CreateProgram()
	gl.AttachShader(program, vs)
	gl.AttachShader(program, fs)
	gl.LinkProgram(program)

	if gl.GetProgrami(program, gl.LINK_STATUS) != gl.TRUE {
		panic(gl.GetProgramInfoLog(program))
	}

	gl.UseProgram(program)
	checkError()

	uniformCount := gl.GetProgrami(program, gl.ACTIVE_UNIFORMS)
	uniforms := make([]shaderUniform, uniformCount)
	textureUnit := 0
	for i := range uniforms {
		name, size, ty := gl.GetActiveUniform(program, uint32(i))
		location := gl.GetUniformLocation(program, name)
		uniforms[i] = shaderUniform{
			name:        name,
			size:        size,
			ty:          shaderDataType(ty),
			location:    location,
			textureUnit: textureUnit,
		}
		if ty == gl.SAMPLER_2D {
			textureUnit++
		}
	}

	attributeCount := gl.GetProgrami(program, gl.ACTIVE_ATTRIBUTES)
	attributes := make([]shaderAttribute, attributeCount)
	for i := range attributes {
		name, size, ty := gl.GetActiveAttrib(program, uint32(i))
		location := gl.GetAttribLocation(program, name)
		attributes[i] = shaderAttribute{
			name:     name,
			size:     size,
			ty:       shaderDataType(ty),
			location: location,
		}
	}

	ctx.stats.shaderProgramCount++

	return &shaderProgram{
		program:    program,
		uniforms:   uniforms,
		attributes: attributes,
	}
}

func (s *shaderProgram) destroy(ctx *context) {
	gl.DeleteProgram(s.program)
	s.program = gl.Program{}
	// TODO: Delete shaders.
	ctx.stats.shaderProgramCount--
}

func (s *shaderProgram) bind(ctx *context, vb *vertexBuffer, uniforms uniformBindings) {
	gl.UseProgram(s.program)
	for _, a := range s.attributes {
		vs, found := vb.streams[a.name]
		if !found {
			panic(fmt.Errorf("VertexBuffer missing required stream '%s'", a.name))
		}
		if a.ty != vs.ty {
			panic(fmt.Errorf("Attribute '%s' type '%s' does not match stream type '%s'",
				a.name, a.ty, vs.ty))
		}
		elementCount := a.ty.vectorElementCount()
		elementTy := a.ty.vectorElementType()
		ctx.getOrCreateVertexStreamContext(vs).bind()
		a.enableArray()
		a.attribPointer(int32(elementCount), uint32(elementTy), false, 0, 0)
	}
	for _, u := range s.uniforms {
		v, found := uniforms[u.name]
		if !found {
			panic(fmt.Errorf("Uniforms missing '%s'", u.name))
		}
		u.bind(ctx, v)
	}
	checkError()
}

func (s *shaderProgram) unbind(ctx *context) {
	for _, a := range s.attributes {
		a.disableArray()
	}
	checkError()
}
