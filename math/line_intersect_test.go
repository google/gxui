// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import test "github.com/google/gxui/testing"
import (
	"testing"
)

func v(a, b float32) Vec2 {
	return Vec2{X: a, Y: b}
}

func TestLinesIntersectCrossing(t *testing.T) {
	test.AssertEquals(t, true, linesIntersect(v(5, 0), v(5, 10), v(0, 5), v(10, 5)))
}

func TestLinesIntersectCrossingEnds(t *testing.T) {
	test.AssertEquals(t, false, linesIntersect(v(0, 0), v(0, 5), v(0, 0), v(5, 0)))
	test.AssertEquals(t, false, linesIntersect(v(0, 0), v(0, 5), v(0, 5), v(5, 5)))
}
func TestLinesIntersectColinearOverlap(t *testing.T) {
	test.AssertEquals(t, true, linesIntersect(v(0, 0), v(5, 0), v(0, 0), v(5, 0)))
	test.AssertEquals(t, true, linesIntersect(v(0, 0), v(5, 0), v(1, 0), v(4, 0)))
	test.AssertEquals(t, true, linesIntersect(v(0, 0), v(5, 0), v(0, 0), v(7, 0)))
}

func TestLinesIntersectParallel(t *testing.T) {
	test.AssertEquals(t, false, linesIntersect(v(5, 0), v(5, 10), v(6, 0), v(6, 10)))
}

func TestLinesIntersectColinearNonOverlap(t *testing.T) {
	test.AssertEquals(t, false, linesIntersect(v(0, 0), v(5, 0), v(5, 0), v(10, 0)))
	test.AssertEquals(t, false, linesIntersect(v(0, 0), v(4, 0), v(6, 0), v(10, 0)))
}
