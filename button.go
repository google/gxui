// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type ButtonType int

const (
	PushButton ButtonType = iota
	ToggleButton
)

type Button interface {
	LinearLayout
	Text() string
	SetText(string)
	Type() ButtonType
	SetType(ButtonType)
	IsChecked() bool
	SetChecked(bool)
}
