// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v3.2-core/gl"
)

type UniformBindings map[string]interface{}

type ShaderProgram struct {
	Program    uint32
	Uniforms   []shaderUniform
	Attributes []shaderAttribute
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
	CheckError()

	return shader
}

func CreateShaderProgram(ctx *Context, vsSource, fsSource string) *ShaderProgram {
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
	CheckError()

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
			ty:          ShaderDataType(ty),
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
			ty:       ShaderDataType(ty),
			location: uint32(location),
		}
	}

	ctx.Stats().ShaderProgramCount++

	return &ShaderProgram{
		Program:    program,
		Uniforms:   uniforms,
		Attributes: attributes,
	}
}

func (s *ShaderProgram) Destroy(ctx *Context) {
	gl.DeleteProgram(s.Program)
	s.Program = 0
	// TODO: Delete shaders.
	ctx.Stats().ShaderProgramCount--
}

func (s *ShaderProgram) Bind(ctx *Context, vb *VertexBuffer, uniforms UniformBindings) {
	gl.UseProgram(s.Program)
	for _, a := range s.Attributes {
		vs, found := vb.Streams[a.name]
		if !found {
			panic(fmt.Errorf("VertexBuffer missing required stream '%s'", a.name))
		}
		if a.ty != vs.Type() {
			panic(fmt.Errorf("Attribute '%s' type '%s' does not match stream type '%s'",
				a.name, a.ty, vs.Type()))
		}
		elementCount := a.ty.VectorElementCount()
		elementTy := a.ty.VectorElementType()
		ctx.GetOrCreateVertexStreamContext(vs).Bind()
		a.EnableArray()
		a.AttribPointer(int32(elementCount), uint32(elementTy), false, 0, nil)
	}
	for _, u := range s.Uniforms {
		v, found := uniforms[u.name]
		if !found {
			panic(fmt.Errorf("Uniforms missing '%s'", u.name))
		}
		u.Bind(ctx, v)
	}
	CheckError()
}

func (s *ShaderProgram) Unbind(ctx *Context) {
	for _, a := range s.Attributes {
		a.DisableArray()
	}
	CheckError()
}
