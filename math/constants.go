// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import (
	"math"
)

var ZeroSize Size
var ZeroSpacing Spacing
var ZeroPoint Point

var MaxSize = Size{0x40000000, 0x40000000} // Picked to avoid integer overflows

var Mat2Ident Mat2 = CreateMat2(1, 0, 0, 1)
var Mat3Ident Mat3 = CreateMat3(1, 0, 0, 0, 1, 0, 0, 0, 1)

const Pi = float32(math.Pi)
const TwoPi = Pi * 2.0

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -(MaxInt - 1)
