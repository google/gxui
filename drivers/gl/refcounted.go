// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

/*
import (
	"fmt"
	"runtime"
	"strings"
	"sync/atomic"
)

const debugTrackReferences = false

type refCounted struct {
	refCount int32
	history  []string
}

func verifyRefCountIsZero(r *refCounted) {
	if r.alive() {
		panic(fmt.Errorf("RefCounted object was garbage collected with a reference count of %d.\n%s",
			r.refCount, strings.Join(r.history, "\n")))
	}
}

func (r *refCounted) init() {
	r.refCount = 1

	if debugTrackReferences {
		_, file, line, _ := runtime.Caller(1)
		r.history = append(r.history, fmt.Sprintf("0 -> 1: %s:%d", file, line))
		runtime.SetFinalizer(r, verifyRefCountIsZero)
	}
}

func (r *refCounted) addRef() {
	r.assertAlive("AddRef")
	count := atomic.AddInt32(&r.refCount, 1)

	if debugTrackReferences {
		_, file, line, _ := runtime.Caller(1)
		r.history = append(r.history, fmt.Sprintf("%d -> %d: %s:%d",
			count-1, count, file, line))
	}
}

func (r *refCounted) release() bool {
	r.assertAlive("Release")
	count := atomic.AddInt32(&r.refCount, -1)

	if debugTrackReferences {
		_, file, line, _ := runtime.Caller(2)
		r.history = append(r.history, fmt.Sprintf("%d -> %d: %s:%d",
			count+1, count, file, line))
	}
	return count == 0
}

func (r *refCounted) alive() bool {
	return atomic.LoadInt32(&r.refCount) > 0
}

func (r *refCounted) assertAlive(funcName string) {
	if !r.alive() {
		if debugTrackReferences {
			panic(fmt.Errorf("Attempting to call %s()) on a fully released object.\n%s",
				funcName, strings.Join(r.history, "\n")))
		} else {
			panic(fmt.Errorf("Attempting to call %s() on a fully released object. Enable debugTrackReferences for more info.", funcName))
		}
	}
}
*/
