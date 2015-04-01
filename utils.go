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
		if p != c.Parent() {
			panic(fmt.Errorf("Child's parent is not as expected.\nChild: %s\nExpected parent: %s",
				Path(c), Path(p)))
		}
		if cp, ok := c.(Parent); ok {
			ValidateHierarchy(cp)
		}
	}
}

func CommonAncestor(a, b Control) Container {
	seen := make(map[Container]bool)
	if c, _ := a.(Container); c != nil {
		seen[c] = true
	}
	for a != nil {
		p := a.Parent()
		seen[p] = true
		a, _ = p.(Control)
	}
	if c, _ := b.(Container); c != nil {
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
		cp := p.Sub(child.Bounds().Min)
		if child.ContainsPoint(cp) {
			l := ControlPointList{ControlPoint{child, cp}}
			if cc, ok := child.(Parent); ok {
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
			cp := p.Sub(child.Bounds().Min)
			if child.ContainsPoint(cp) {
				l = append(l, ControlPoint{child, cp})
				if cc, ok := child.(Parent); ok {
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
		coord = coord.Sub(c.Bounds().Min)
		p := c.Parent()
		if p == nil {
			panic("Control's parent was nil")
		}
		if _, ok := p.(Window); ok {
			return coord
		}
		c = p.(Control)
	}
}

func ChildToParent(coord math.Point, from Control, to Parent) math.Point {
	c := from
	for {
		coord = coord.Add(c.Bounds().Min)
		p := c.Parent()
		if p == nil {
			panic(fmt.Errorf("Control detached: %s", Path(c)))
		}
		if p == to {
			return coord
		}
		c = p.(Control)
	}
}

func ParentToChild(coord math.Point, from Parent, to Control) math.Point {
	return coord.Sub(ChildToParent(math.ZeroPoint, to, from))
}

func TransformCoordinate(coord math.Point, from, to Control) math.Point {
	ancestor := CommonAncestor(from, to)
	if ancestor == nil {
		panic(fmt.Errorf("No common ancestor between %s and %s", Path(from), Path(to)))
	}

	coord = ChildToParent(coord, from, ancestor)
	coord = ParentToChild(coord, ancestor, to)
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
		if test(child) {
			return child
		}
		if parent, ok := child.(Parent); ok {
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
