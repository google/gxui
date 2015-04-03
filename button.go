// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import "fmt"

type Button interface {
	LinearLayout
	Text() string
	SetText(string)
	Type() ButtonType
	SetType(ButtonType)
	IsChecked() bool
	SetChecked(bool)
}

type ButtonType int

const (
	PushButton ButtonType = iota
	ToggleButton
)

func (t ButtonType) String() string {
	switch t {
	case PushButton:
		return "Push Button"
	case ToggleButton:
		return "Toggle Button"
	default:
		return fmt.Sprintf("ButtonType<%d>", t)
	}
}
