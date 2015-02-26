// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type Focusable interface {
	Control

	IsFocusable() bool
	HasFocus() bool

	GainedFocus()
	LostFocus()

	OnGainedFocus(func()) EventSubscription
	OnLostFocus(func()) EventSubscription
}
