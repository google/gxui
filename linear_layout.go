// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

// LinearLayout is a Container that lays out its child Controls into a column or
// row. The layout will always start by positioning the first (0'th) child, and
// then depending on the direction, will position each successive child either
// to the left, top, right or bottom of the preceding child Control.
// LinearLayout makes no effort to distribute remaining space evenly between the
// children - an child control that is laid out before others will reduce the
// remaining space given to the later children, even to the point that there is
// zero space remaining.
type LinearLayout interface {
	// LinearLayout extends the Control interface.
	Control

	// LinearLayout extends the Container interface.
	Container

	// Direction returns the direction of layout for this LinearLayout.
	Direction() Direction

	// Direction sets the direction of layout for this LinearLayout.
	SetDirection(Direction)

	// SizeMode returns the desired size behaviour for this LinearLayout.
	SizeMode() SizeMode

	// SetSizeMode sets the desired size behaviour for this LinearLayout.
	SetSizeMode(SizeMode)

	// HorizontalAlignment returns the alignment of the child Controls when laying
	// out TopToBottom or BottomToTop. It has no effect when the layout direction
	// is LeftToRight or RightToLeft.
	HorizontalAlignment() HorizontalAlignment

	// SetHorizontalAlignment sets the alignment of the child Controls when laying
	// out TopToBottom or BottomToTop. It has no effect when the layout direction
	// is LeftToRight or RightToLeft.
	SetHorizontalAlignment(HorizontalAlignment)

	// VerticalAlignment returns the alignment of the child Controls when laying
	// out LeftToRight or RightToLeft. It has no effect when the layout direction
	// is TopToBottom or BottomToTop.
	VerticalAlignment() VerticalAlignment

	// SetVerticalAlignment returns the alignment of the child Controls when
	// laying out LeftToRight or RightToLeft. It has no effect when the layout
	// direction is TopToBottom or BottomToTop.
	SetVerticalAlignment(VerticalAlignment)

	// BorderPen returns the Pen used to draw the LinearLayout's border.
	BorderPen() Pen

	// SetBorderPen sets the Pen used to draw the LinearLayout's border.
	SetBorderPen(Pen)

	// BackgroundBrush returns the Brush used to fill the LinearLayout's
	// background.
	BackgroundBrush() Brush

	// SetBackgroundBrush sets the Brush used to fill the LinearLayout's
	// background.
	SetBackgroundBrush(Brush)
}
