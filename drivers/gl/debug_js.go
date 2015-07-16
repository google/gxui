// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build js

package gl

func (d *driver) discoverUIGoRoutine() {
	println("discoverUIGoRoutine not yet implemented on js architecture")
}

func (d *driver) AssertUIGoroutine() {
	// AssertUIGoroutine not yet implemented on js architecture, so it never panics.
}
