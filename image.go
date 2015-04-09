// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/google/gxui/math"
)

type ScalingMode int

const (
	Scaling1to1 ScalingMode = iota
	ScalingExpandGreedy
	ScalingExplicitSize
)

type AspectMode int

const (
	AspectStretch = iota
	AspectCorrectLetterbox
	AspectCorrectCrop
)

type Image interface {
	Control
	Texture() Texture
	SetTexture(Texture)
	Canvas() Canvas
	SetCanvas(Canvas)
	BorderPen() Pen
	SetBorderPen(Pen)
	BackgroundBrush() Brush
	SetBackgroundBrush(Brush)
	ScalingMode() ScalingMode
	SetScalingMode(ScalingMode)
	SetExplicitSize(math.Size)
	AspectMode() AspectMode
	SetAspectMode(AspectMode)
	PixelAt(math.Point) (math.Point, bool) // TODO: Remove
}
