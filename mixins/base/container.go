// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"github.com/google/gxui"
	"github.com/google/gxui/mixins/outer"
	"github.com/google/gxui/mixins/parts"
)

type ContainerNoControlOuter interface {
	gxui.Container
	outer.PaintChilder
	outer.Painter
	outer.LayoutChildren
}

type ContainerOuter interface {
	ContainerNoControlOuter
	gxui.Control
}

type Container struct {
	parts.Attachable
	parts.Container
	parts.DrawPaint
	parts.InputEventHandler
	parts.Layoutable
	parts.Paddable
	parts.PaintChildren
	parts.Parentable
	parts.Visible
}

func (c *Container) Init(outer ContainerOuter, theme gxui.Theme) {
	c.Attachable.Init(outer)
	c.Container.Init(outer)
	c.DrawPaint.Init(outer, theme)
	c.InputEventHandler.Init(outer)
	c.Layoutable.Init(outer, theme)
	c.Paddable.Init(outer)
	c.PaintChildren.Init(outer)
	c.Parentable.Init(outer)
	c.Visible.Init(outer)

	// Interface compliance test
	_ = gxui.Container(c)
}
