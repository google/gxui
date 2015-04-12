// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parts

import (
	"fmt"

	"github.com/google/gxui"
	"github.com/google/gxui/math"
	"github.com/google/gxui/mixins/outer"
)

type ContainerOuter interface {
	gxui.Container
	outer.Attachable
	outer.IsVisibler
	outer.LayoutChildren
	outer.Parenter
	outer.Sized
}

type Container struct {
	outer              ContainerOuter
	children           gxui.Children
	isMouseEventTarget bool
	relayoutSuspended  bool
}

func (c *Container) Init(outer ContainerOuter) {
	c.outer = outer
	c.children = gxui.Children{}
	outer.OnAttach(func() {
		for _, v := range c.children {
			v.Control.Attach()
		}
	})
	outer.OnDetach(func() {
		for _, v := range c.children {
			v.Control.Detach()
		}
	})
}

func (c *Container) SetMouseEventTarget(mouseEventTarget bool) {
	c.isMouseEventTarget = mouseEventTarget
}

func (c *Container) IsMouseEventTarget() bool {
	return c.isMouseEventTarget
}

// RelayoutSuspended returns true if adding or removing a child Control to this
// Container will not trigger a relayout of this Container. The default is false
// where any mutation will trigger a relayout.
func (c *Container) RelayoutSuspended() bool {
	return c.relayoutSuspended
}

// SetRelayoutSuspended enables or disables relayout of the Container on
// adding or removing a child Control to this Container.
func (c *Container) SetRelayoutSuspended(enable bool) {
	c.relayoutSuspended = true
}

// gxui.Parent compliance
func (c *Container) Children() gxui.Children {
	return c.children
}

// gxui.Container compliance
func (c *Container) AddChild(control gxui.Control) *gxui.Child {
	return c.outer.AddChildAt(len(c.children), control)
}

func (c *Container) AddChildAt(index int, control gxui.Control) *gxui.Child {
	if control.Parent() != nil {
		panic("Child already has a parent")
	}
	if index < 0 || index > len(c.children) {
		panic(fmt.Errorf("Index %d is out of bounds. Acceptable range: [%d - %d]",
			index, 0, len(c.children)))
	}

	child := &gxui.Child{Control: control}

	c.children = append(c.children, nil)
	copy(c.children[index+1:], c.children[index:])
	c.children[index] = child

	control.SetParent(c.outer)
	if c.outer.Attached() {
		control.Attach()
	}
	if !c.relayoutSuspended {
		c.outer.Relayout()
	}
	return child
}

func (c *Container) RemoveChild(control gxui.Control) {
	for i := range c.children {
		if c.children[i].Control == control {
			c.outer.RemoveChildAt(i)
			return
		}
	}
	panic("Child not part of container")
}

func (c *Container) RemoveChildAt(index int) {
	child := c.children[index]
	c.children = append(c.children[:index], c.children[index+1:]...)
	child.Control.SetParent(nil)
	if c.outer.Attached() {
		child.Control.Detach()
	}
	if !c.relayoutSuspended {
		c.outer.Relayout()
	}
}

func (c *Container) RemoveAll() {
	for i := len(c.children) - 1; i >= 0; i-- {
		c.outer.RemoveChildAt(i)
	}
}

func (c *Container) ContainsPoint(p math.Point) bool {
	if !c.outer.IsVisible() || !c.outer.Size().Rect().Contains(p) {
		return false
	}
	for _, v := range c.children {
		if v.Control.ContainsPoint(p.Sub(v.Offset)) {
			return true
		}
	}
	if c.IsMouseEventTarget() {
		return true
	}
	return false
}
