// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import test "github.com/google/gxui/testing"
import "testing"

func TestRectConstrainWithin(t *testing.T) {
	r1 := CreateRect(0, 0, 100, 100)
	r2 := CreateRect(40, 40, 60, 60)
	test.AssertEquals(t, r2, r2.Constrain(r1))
}

func TestRectConstrainOutMin(t *testing.T) {
	r1 := CreateRect(0, 0, 100, 100)
	r2 := CreateRect(-20, -20, 20, 20)
	test.AssertEquals(t, CreateRect(0, 0, 40, 40), r2.Constrain(r1))
}

func TestRectConstrainOutMax(t *testing.T) {
	r1 := CreateRect(0, 0, 100, 100)
	r2 := CreateRect(80, 80, 120, 120)
	test.AssertEquals(t, CreateRect(60, 60, 100, 100), r2.Constrain(r1))
}
