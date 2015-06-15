// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"fmt"
	"github.com/google/gxui/math"
	"reflect"
	"runtime"
)

func indent(depth int) string {
	s := ""
	for i := 0; i < depth; i++ {
		s += "   |"
	}
	return s
}

func dump(c interface{}, depth int) string {
	s := ""
	switch t := c.(type) {
	case Control:
		s += fmt.Sprintf("(%p) %T Size: %+v Margin: %+v \n", &t, t, t.Size(), t.Margin())
	default:
		s += fmt.Sprintf("%T\n", t)
	}
	switch t := c.(type) {
	case Container:
		for i, c := range t.Children() {
			s += fmt.Sprintf("%s--- Child %d: ", indent(depth), i)
			s += dump(c.Control, depth+1)
		}
	}
	return s
}

func Dump(c interface{}) {
	fmt.Printf("%s\n", dump(c, 0))
}

func FunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func BreadcrumbsAt(p Container, pnt math.Point) string {
	s := reflect.TypeOf(p).String()
	for _, c := range p.Children() {
		b := c.Control.Size().Rect().Offset(c.Offset)
		if b.Contains(pnt) {
			switch t := c.Control.(type) {
			case Container:
				return s + " > " + BreadcrumbsAt(t, pnt.Sub(c.Offset))
			default:
				return s + " > " + reflect.TypeOf(c.Control).String()
			}
		}
	}
	return s
}

func Path(p interface{}) string {
	if p == nil {
		return "nil"
	}

	s := reflect.TypeOf(p).String()

	if c, _ := p.(Control); c != nil {
		if c.Parent() != nil {
			return Path(c.Parent()) + " > " + s
		}
	}

	return s
}
