// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type MouseState int

func (s MouseState) IsDown(b MouseButton) bool {
	return s&(1<<uint(b)) != 0
}
