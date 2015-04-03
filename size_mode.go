// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import "fmt"

type SizeMode int

const (
	ExpandToContent SizeMode = iota
	Fill
)

func (s SizeMode) ExpandToContent() bool {
	return s == ExpandToContent
}

func (s SizeMode) Fill() bool {
	return s == Fill
}

func (s SizeMode) String() string {
	switch s {
	case ExpandToContent:
		return "Expand To Content"
	case Fill:
		return "Fill"
	default:
		return fmt.Sprintf("SizeMode<%d>", s)
	}
}
