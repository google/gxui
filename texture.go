// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"image"

	"github.com/google/gxui/math"
)

type Texture interface {
	Image() image.Image
	Size() math.Size
	SizePixels() math.Size
	FlipY() bool
	SetFlipY(bool)
}
