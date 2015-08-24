// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"fmt"
	"strings"

	"github.com/google/gxui/math"
)

type Child struct {
	Control Control
	Offset  math.Point
}

type Parent interface {
	Children() Children
	Relayout()
	Redraw()
}

type Container interface {
	Parent
	AddChild(child Control) *Child
	AddChildAt(index int, child Control) *Child
	RemoveChild(child Control)
	RemoveChildAt(index int)
	RemoveAll()
	Padding() math.Spacing
	SetPadding(math.Spacing)
}

// String returns a string describing the child type and bounds.
func (c *Child) String() string {
	return fmt.Sprintf("Type: %T, Bounds: %v", c.Control, c.Bounds())
}

// Bounds returns the Child bounds relative to the parent.
func (c *Child) Bounds() math.Rect {
	return c.Control.Size().Rect().Offset(c.Offset)
}

// Layout sets the Child size and offset relative to the parent.
// Layout should only be called by the Child's parent.
func (c *Child) Layout(rect math.Rect) {
	c.Offset = rect.Min
	c.Control.SetSize(rect.Size())
}

// Children is a list of Child pointers.
type Children []*Child

// String returns a string describing the child type and bounds.
func (c Children) String() string {
	s := make([]string, len(c))
	for i, c := range c {
		s[i] = fmt.Sprintf("%d: %s", i, c.String())
	}
	return strings.Join(s, "\n")
}

// IndexOf returns and returns the index of the child control, or -1 if the
// child is not in this Children list.
func (c Children) IndexOf(control Control) int {
	for i, child := range c {
		if child.Control == control {
			return i
		}
	}
	return -1
}

// Find returns and returns the Child pointer for the given Control, or nil
// if the child is not in this Children list.
func (c Children) Find(control Control) *Child {
	for _, child := range c {
		if child.Control == control {
			return child
		}
	}
	return nil
}
