// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type ScrollBar interface {
	Control

	OnScroll(func(from, to int)) EventSubscription
	ScrollPosition() (from, to int)
	SetScrollPosition(from, to int)
	ScrollLimit() int
	SetScrollLimit(l int)
	AutoHide() bool
	SetAutoHide(l bool)
	Orientation() Orientation
	SetOrientation(Orientation)
}
