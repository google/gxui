// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !js

package gl

import (
	"runtime"
	"strings"
)

// discoverUIGoRoutine finds and stores the program counter of the
// function 'applicationLoop' that must be in the callstack. The
// PC is stored so that AssertUIGoroutine can verify that the call
// came from the application loop (the UI go-routine).
func (d *driver) discoverUIGoRoutine() {
	for _, pc := range d.pcs[:runtime.Callers(2, d.pcs)] {
		name := runtime.FuncForPC(pc).Name()
		if strings.HasSuffix(name, "applicationLoop") {
			d.uiPC = pc
			return
		}
	}
	panic("applicationLoop was not found in the callstack")
}

func (d *driver) AssertUIGoroutine() {
	for _, pc := range d.pcs[:runtime.Callers(2, d.pcs)] {
		if pc == d.uiPC {
			return
		}
	}
	panic("AssertUIGoroutine called on a go-routine that was not the UI go-routine")
}
