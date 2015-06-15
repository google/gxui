// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"bytes"
	"fmt"
	"unicode/utf8"

	"github.com/google/gxui/math"
)

type ParentPoint struct {
	C Parent
	P math.Point
}

type ControlPoint struct {
	C Control
	P math.Point
}

type ControlPointList []ControlPoint

func (l ControlPointList) Contains(c Control) bool {
	_, found := l.Find(c)
	return found
}

func (l ControlPointList) Find(c Control) (math.Point, bool) {
	for _, i := range l {
		if i.C == c {
			return i.P, true
		}
	}
	return math.Point{}, false
}

func ValidateHierarchy(p Parent) {
	for _, c := range p.Children() {
		if p != c.Control.Parent() {
			panic(fmt.Errorf("Child's parent is not as expected.\nChild: %s\nExpected parent: %s",
				Path(c.Control), Path(p)))
		}
		if cp, ok := c.Control.(Parent); ok {
			ValidateHierarchy(cp)
		}
	}
}

func CommonAncestor(a, b Control) Parent {
	seen := make(map[Parent]bool)
	if c, _ := a.(Parent); c != nil {
		seen[c] = true
	}
	for a != nil {
		p := a.Parent()
		seen[p] = true
		a, _ = p.(Control)
	}
	if c, _ := b.(Parent); c != nil {
		if seen[c] {
			return c
		}
	}
	for b != nil {
		p := b.Parent()
		if seen[p] {
			return p
		}
		b, _ = p.(Control)
	}
	return nil
}

func TopControlsUnder(p math.Point, c Parent) ControlPointList {
	children := c.Children()
	for i := len(children) - 1; i >= 0; i-- {
		child := children[i]
		cp := p.Sub(child.Offset)
		if child.Control.ContainsPoint(cp) {
			l := ControlPointList{ControlPoint{child.Control, cp}}
			if cc, ok := child.Control.(Parent); ok {
				l = append(l, TopControlsUnder(cp, cc)...)
			}
			return l
		}
	}
	return ControlPointList{}
}

func ControlsUnder(p math.Point, c Parent) ControlPointList {
	toVisit := []ParentPoint{ParentPoint{c, p}}
	l := ControlPointList{}
	for len(toVisit) > 0 {
		c = toVisit[0].C
		p = toVisit[0].P
		toVisit = toVisit[1:]
		for _, child := range c.Children() {
			cp := p.Sub(child.Offset)
			if child.Control.ContainsPoint(cp) {
				l = append(l, ControlPoint{child.Control, cp})
				if cc, ok := child.Control.(Parent); ok {
					toVisit = append(toVisit, ParentPoint{cc, cp})
				}
			}
		}
	}
	return l
}

func WindowToChild(coord math.Point, to Control) math.Point {
	c := to
	for {
		p := c.Parent()
		if p == nil {
			panic("Control's parent was nil")
		}
		child := p.Children().Find(c)
		if child == nil {
			Dump(p)
			panic(fmt.Errorf("Control's parent (%p %T) did not contain control (%p %T).", &p, p, &c, c))
		}
		coord = coord.Sub(child.Offset)
		if _, ok := p.(Window); ok {
			return coord
		}
		c = p.(Control)
	}
}

func ChildToParent(coord math.Point, from Control, to Parent) math.Point {
	c := from
	for {
		p := c.Parent()
		if p == nil {
			panic(fmt.Errorf("Control detached: %s", Path(c)))
		}
		child := p.Children().Find(c)
		if child == nil {
			Dump(p)
			panic(fmt.Errorf("Control's parent (%p %T) did not contain control (%p %T).", &p, p, &c, c))
		}
		coord = coord.Add(child.Offset)
		if p == to {
			return coord
		}

		if control, ok := p.(Control); ok {
			c = control
		} else {
			Dump(p)
			panic(fmt.Errorf("ChildToParent (%p %T) -> (%p %T) reached non-control parent (%p %T).",
				&from, from, &to, to, &p, p))
		}
	}
}

func ParentToChild(coord math.Point, from Parent, to Control) math.Point {
	return coord.Sub(ChildToParent(math.ZeroPoint, to, from))
}

func TransformCoordinate(coord math.Point, from, to Control) math.Point {
	if from == to {
		return coord
	}

	ancestor := CommonAncestor(from, to)
	if ancestor == nil {
		panic(fmt.Errorf("No common ancestor between %s and %s", Path(from), Path(to)))
	}

	if parent, ok := ancestor.(Control); !ok || parent != from {
		coord = ChildToParent(coord, from, ancestor)
	}
	if parent, ok := ancestor.(Control); !ok || parent != to {
		coord = ParentToChild(coord, ancestor, to)
	}
	return coord
}

// FindControl performs a depth-first search of the controls starting from root,
// calling test with each visited control. If test returns true then the search
// is stopped and FindControl returns the Control passed to test. If no call to
// test returns true then FindControl returns nil.
func FindControl(root Parent, test func(Control) (found bool)) Control {
	if c, ok := root.(Control); ok && test(c) {
		return c
	}

	for _, child := range root.Children() {
		if test(child.Control) {
			return child.Control
		}
		if parent, ok := child.Control.(Parent); ok {
			if c := FindControl(parent, test); c != nil {
				return c
			}
		}
	}
	return nil
}

func WindowContaining(c Control) Window {
	for {
		p := c.Parent()
		if p == nil {
			panic("Control's parent was nil")
		}
		if window, ok := p.(Window); ok {
			return window
		}
		c = p.(Control)
	}
}

func SetFocus(focusable Focusable) {
	wnd := WindowContaining(focusable)
	wnd.SetFocus(focusable)
}

func StringToRuneArray(str string) []rune {
	return bytes.Runes([]byte(str))
}

func RuneArrayToString(arr []rune) string {
	tmp := make([]byte, 8)
	enc := make([]byte, 0, len(arr))
	offset := 0
	for _, r := range arr {
		size := utf8.EncodeRune(tmp, r)
		enc = append(enc, tmp[:size]...)
		offset += size
	}
	return string(enc)
}
