// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"github.com/google/gxui/assert"

	"github.com/go-gl-legacy/gl"
)

type UniformBindings map[string]interface{}

type ShaderProgram struct {
	Program    gl.Program
	Uniforms   []ShaderUniform
	Attributes []ShaderAttribute
}

func compile(source string, ty gl.GLenum) gl.Shader {
	s := gl.CreateShader(ty)
	s.Source(source)
	s.Compile()
	err := s.GetInfoLog()
	if err != "" {
		panic(err)
	}
	return s
}

func CreateShaderProgram(ctx *Context, vsSource, fsSource string) *ShaderProgram {
	vs := compile(vsSource, gl.VERTEX_SHADER)
	fs := compile(fsSource, gl.FRAGMENT_SHADER)

	program := gl.CreateProgram()
	program.AttachShader(vs)
	program.AttachShader(fs)
	program.Link()
	program.Validate()

	err := program.GetInfoLog()
	if err != "" {
		panic(err)
	}

	program.Use()
	CheckError()

	uniformCount := program.Get(gl.ACTIVE_UNIFORMS)
	uniforms := make([]ShaderUniform, uniformCount)
	textureUnit := 0
	for i := 0; i < uniformCount; i++ {
		size, ty, name := program.GetActiveUniform(i)
		location := program.GetUniformLocation(name)
		uniforms[i] = CreateShaderUniform(name, size, ShaderDataType(ty), location, textureUnit)
		if ty == gl.SAMPLER_2D {
			textureUnit++
		}
	}

	attributeCount := program.Get(gl.ACTIVE_ATTRIBUTES)
	attributes := make([]ShaderAttribute, attributeCount)
	for i := 0; i < attributeCount; i++ {
		size, ty, name := program.GetActiveAttrib(i)
		location := program.GetAttribLocation(name)
		attributes[i] = CreateShaderAttribute(name, size, ShaderDataType(ty), location)
	}

	ctx.Stats().ShaderProgramCount++

	return &ShaderProgram{
		Program:    program,
		Uniforms:   uniforms,
		Attributes: attributes,
	}
}

func (s *ShaderProgram) Destroy(ctx *Context) {
	s.Program.Delete()
	s.Program = gl.Program(0)
	// TODO: Delete shaders.
	ctx.Stats().ShaderProgramCount--
}

func (s *ShaderProgram) Bind(ctx *Context, vb *VertexBuffer, uniforms UniformBindings) {
	s.Program.Use()
	for _, a := range s.Attributes {
		vs, found := vb.Streams[a.Name]
		assert.True(found, "VertexBuffer missing required stream '%s'", a.Name)
		assert.Equals(a.Type, vs.Type(), "Attribute %s type", a.Name)
		elementCount := a.Type.VectorElementCount()
		elementTy := a.Type.VectorElementType()
		ctx.GetOrCreateVertexStreamContext(vs).Bind()
		a.Location.EnableArray()
		a.Location.AttribPointer(uint(elementCount), gl.GLenum(elementTy), false, 0, nil)
		CheckError()
	}
	for _, u := range s.Uniforms {
		v, found := uniforms[u.Name]
		assert.True(found, "Uniforms missing '%s'", u.Name)
		u.Bind(ctx, v)
	}
}

func (s *ShaderProgram) Unbind(ctx *Context) {
	for _, a := range s.Attributes {
		a.Location.DisableArray()
	}
	CheckError()
}
