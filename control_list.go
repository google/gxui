// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type ControlList []Control

func (l ControlList) Contains(c Control) bool {
	for _, i := range l {
		if i == c {
			return true
		}
	}
	return false
}
