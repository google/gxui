// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

func EventLoop(driver Driver) {
	events := driver.Events()
	for {
		select {
		case ev, open := <-events:
			if open {
				ev()
			} else {
				return // closed channel represents driver shutdown
			}
		}
	}
}
