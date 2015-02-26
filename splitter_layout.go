// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type SplitterLayout interface {
	Control
	Container
	ChildWeight(Control) float32
	SetChildWeight(Control, float32)
	Orientation() Orientation
	SetOrientation(Orientation)
}
