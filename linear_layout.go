// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type LinearLayout interface {
	Control
	Container
	Orientation() Orientation
	SetOrientation(Orientation)
	HorizontalAlignment() HorizontalAlignment
	SetHorizontalAlignment(HorizontalAlignment)
	VerticalAlignment() VerticalAlignment
	SetVerticalAlignment(VerticalAlignment)
	BorderPen() Pen
	SetBorderPen(Pen)
	BackgroundBrush() Brush
	SetBackgroundBrush(Brush)
}
