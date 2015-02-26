// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"
	"github.com/go-gl/gl"
)

func CheckError() {
	v := gl.GetError()
	if v != gl.GLenum(0) {
		err := fmt.Errorf("GL returned error 0x%.4x", v)
		panic(err)
	}
}
