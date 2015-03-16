// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import test "github.com/google/gxui/testing"
import "testing"

func TestMat2InvertIdent(t *testing.T) {
	m := Mat2Ident.Invert()
	test.AssertEquals(t, Mat2Ident, m)
}

func TestMat2InvertSimple(t *testing.T) {
	m := CreateMat2(4, 3, 3, 2)
	a := m.Invert()
	e := CreateMat2(-2, 3, 3, -4)
	test.AssertEquals(t, e, a)
}
