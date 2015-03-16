// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"
	"reflect"

	"github.com/go-gl-legacy/gl"
)

type PrimitiveType int

const (
	FLOAT  PrimitiveType = gl.FLOAT
	INT    PrimitiveType = gl.INT
	UINT   PrimitiveType = gl.UNSIGNED_INT
	USHORT PrimitiveType = gl.UNSIGNED_SHORT
	UBYTE  PrimitiveType = gl.UNSIGNED_BYTE
)

func (p PrimitiveType) SizeInBytes() int {
	switch p {
	case FLOAT:
		return 4
	case INT:
		return 4
	case UINT:
		return 4
	case USHORT:
		return 2
	case UBYTE:
		return 1
	default:
		panic(fmt.Errorf("Unknown PrimitiveType 0x%.4x", p))
	}
}

func (p PrimitiveType) IsArrayOfType(array interface{}) bool {
	ty := reflect.TypeOf(array).Elem()
	switch p {
	case FLOAT:
		return ty.Name() == "float32"
	case INT:
		return ty.Name() == "int32"
	case UINT:
		return ty.Name() == "uint32"
	case USHORT:
		return ty.Name() == "uint16"
	case UBYTE:
		return ty.Name() == "uint8"
	default:
		panic(fmt.Errorf("Unknown PrimitiveType 0x%.4x", p))
	}
}
