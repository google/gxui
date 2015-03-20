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
	outer.Bounds
	outer.IsVisibler
	outer.LayoutChildren
	outer.Parenter
}

type Container struct {
	outer              ContainerOuter
	children           []gxui.Control
	isMouseEventTarget bool
	relayoutSuspended  bool
}

func (c *Container) Init(outer ContainerOuter) {
	c.outer = outer
	c.children = []gxui.Control{}
	outer.OnAttach(func() {
		for _, v := range c.children {
			v.Attach()
		}
	})
	outer.OnDetach(func() {
		for _, v := range c.children {
			v.Detach()
		}
	})
}

func (c *Container) Children() []gxui.Control {
	return c.children
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

// gxui.Container compliance
func (c *Container) ChildCount() int {
	return len(c.children)
}

func (c *Container) ChildIndex(child gxui.Control) int {
	for i, v := range c.children {
		if v == child {
			return i
		}
	}
	return -1
}

func (c *Container) ChildAt(index int) gxui.Control {
	return c.children[index]
}

func (c *Container) AddChild(child gxui.Control) {
	c.outer.AddChildAt(len(c.children), child)
}

func (c *Container) AddChildAt(index int, child gxui.Control) {
	if child.Parent() != nil {
		panic("Child already has a parent")
	}
	if index < 0 || index > len(c.children) {
		panic(fmt.Errorf("Index %d is out of bounds. Acceptable range: [%d - %d]",
			index, 0, len(c.children)))
	}

	c.children = append(c.children, nil)
	copy(c.children[index+1:], c.children[index:])
	c.children[index] = child

	child.SetParent(c.outer)
	if c.outer.Attached() {
		child.Attach()
	}
	if !c.relayoutSuspended {
		c.outer.Relayout()
	}
}

func (c *Container) RemoveChild(child gxui.Control) {
	i := c.ChildIndex(child)
	if i >= 0 {
		c.outer.RemoveChildAt(i)
	} else {
		panic("Child not part of container")
	}
}

func (c *Container) RemoveChildAt(index int) {
	child := c.children[index]
	c.children = append(c.children[:index], c.children[index+1:]...)
	child.SetParent(nil)
	if c.outer.Attached() {
		child.Detach()
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
	if !c.outer.IsVisible() || !c.outer.Bounds().Size().Rect().Contains(p) {
		return false
	}
	for _, v := range c.children {
		if v.ContainsPoint(p.Sub(v.Bounds().Min)) {
			return true
		}
	}
	if c.IsMouseEventTarget() {
		return true
	}
	return false
}
