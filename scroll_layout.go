// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type ScrollLayout interface {
	Control
	Parent
	SetChild(Control)
	Child() Control
	SetScrollAxis(horizontal, vertical bool)
	ScrollAxis() (horizontal, vertical bool)
	BorderPen() Pen
	SetBorderPen(Pen)
	BackgroundBrush() Brush
	SetBackgroundBrush(Brush)
}
