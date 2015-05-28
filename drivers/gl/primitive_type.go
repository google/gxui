// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"
	"reflect"

	"github.com/goxjs/gl"
)

type primitiveType int

const (
	ptFloat  primitiveType = gl.FLOAT
	ptInt    primitiveType = gl.INT
	ptUint   primitiveType = gl.UNSIGNED_INT
	ptUshort primitiveType = gl.UNSIGNED_SHORT
	ptUbyte  primitiveType = gl.UNSIGNED_BYTE
)

func (p primitiveType) sizeInBytes() int {
	switch p {
	case ptFloat:
		return 4
	case ptInt:
		return 4
	case ptUint:
		return 4
	case ptUshort:
		return 2
	case ptUbyte:
		return 1
	default:
		panic(fmt.Errorf("Unknown primitiveType 0x%.4x", p))
	}
}

func (p primitiveType) isArrayOfType(array interface{}) bool {
	ty := reflect.TypeOf(array).Elem()
	switch p {
	case ptFloat:
		return ty.Name() == "float32"
	case ptInt:
		return ty.Name() == "int32"
	case ptUint:
		return ty.Name() == "uint32"
	case ptUshort:
		return ty.Name() == "uint16"
	case ptUbyte:
		return ty.Name() == "uint8"
	default:
		panic(fmt.Errorf("Unknown primitiveType 0x%.4x", p))
	}
}
