// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import test "github.com/google/gxui/testing"
import "testing"

func TestSizeEdgeAlignedFitTopEdge(t *testing.T) {
	outer := CreateRect(0, 0, 100, 100)
	s := Size{10, 10}
	p := Point{50, 50}
	test.AssertEquals(t, CreateRect(45, 50, 55, 60), s.EdgeAlignedFit(outer, p))
}

func TestSizeEdgeAlignedFitBottomEdge(t *testing.T) {
	outer := CreateRect(0, 0, 100, 100)
	s := Size{10, 10}
	p := Point{50, 95}
	test.AssertEquals(t, CreateRect(45, 85, 55, 95), s.EdgeAlignedFit(outer, p))
}

func TestSizeEdgeAlignedFitLeftEdge(t *testing.T) {
	outer := CreateRect(0, 0, 100, 100)
	s := Size{10, 80}
	p := Point{5, 50}
	test.AssertEquals(t, CreateRect(5, 10, 15, 90), s.EdgeAlignedFit(outer, p))
}

func TestSizeEdgeAlignedFitRightEdge(t *testing.T) {
	outer := CreateRect(0, 0, 100, 100)
	s := Size{10, 80}
	p := Point{95, 50}
	test.AssertEquals(t, CreateRect(85, 10, 95, 90), s.EdgeAlignedFit(outer, p))
}
