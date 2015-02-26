// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type HorizontalAlignment int

const (
	AlignLeft HorizontalAlignment = iota
	AlignCenter
	AlignRight
)

type VerticalAlignment int

const (
	AlignTop VerticalAlignment = iota
	AlignMiddle
	AlignBottom
)
