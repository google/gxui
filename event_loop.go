// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

func EventLoop(driver Driver) {
	for ev := range driver.Events() {
		ev()
	}
}
