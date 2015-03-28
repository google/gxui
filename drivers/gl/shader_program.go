// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
)

type uniformBindings map[string]interface{}

type shaderProgram struct {
	program    uint32
	uniforms   []shaderUniform
	attributes []shaderAttribute
}

func compile(source string, ty uint32) uint32 {
	shader := gl.CreateShader(ty)
	c := gl.Str(source + "\x00")
	gl.ShaderSource(shader, 1, &c, nil)

	var status int32
	gl.CompileShader(shader)
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status != gl.TRUE {
		var l int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &l)
		log := strings.Repeat("\x00", int(l+1))
		gl.GetShaderInfoLog(shader, l, nil, gl.Str(log))
		panic(log)
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

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status != gl.TRUE {
		var l int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &l)
		log := strings.Repeat("\x00", int(l+1))
		gl.GetProgramInfoLog(program, l, nil, gl.Str(log))
		panic(log)
	}

	gl.UseProgram(program)
	checkError()

	var uniformCount int32
	gl.GetProgramiv(program, gl.ACTIVE_UNIFORMS, &uniformCount)
	uniforms := make([]shaderUniform, uniformCount)
	textureUnit := 0
	for i := range uniforms {
		var length, size int32
		var ty uint32
		name := strings.Repeat("\x00", 256)
		cname := gl.Str(name)
		gl.GetActiveUniform(program, uint32(i), int32(len(name)-1), &length, &size, &ty, cname)
		location := gl.GetUniformLocation(program, cname)
		name = name[:strings.IndexRune(name, 0)]
		uniforms[i] = shaderUniform{
			name:        name,
			size:        int(size),
			ty:          shaderDataType(ty),
			location:    location,
			textureUnit: textureUnit,
		}
		if ty == gl.SAMPLER_2D {
			textureUnit++
		}
	}

	var attributeCount int32
	gl.GetProgramiv(program, gl.ACTIVE_ATTRIBUTES, &attributeCount)
	attributes := make([]shaderAttribute, attributeCount)
	for i := range attributes {
		var length, size int32
		var ty uint32
		name := strings.Repeat("\x00", 256)
		cname := gl.Str(name)
		gl.GetActiveAttrib(program, uint32(i), int32(len(name)-1), &length, &size, &ty, cname)
		name = name[:strings.IndexRune(name, 0)]
		location := gl.GetAttribLocation(program, cname)
		attributes[i] = shaderAttribute{
			name:     name,
			size:     int(size),
			ty:       shaderDataType(ty),
			location: uint32(location),
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
	s.program = 0
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
		a.attribPointer(int32(elementCount), uint32(elementTy), false, 0, nil)
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
