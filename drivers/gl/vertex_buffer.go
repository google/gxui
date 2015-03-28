// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import "fmt"

type vertexBuffer struct {
	refCounted
	streams map[string]*vertexStream
	count   int
}

func newVertexBuffer(streams ...*vertexStream) *vertexBuffer {
	vb := &vertexBuffer{
		streams: map[string]*vertexStream{},
	}
	vb.init()
	for i, s := range streams {
		if i == 0 {
			vb.count = s.count
		} else {
			if vb.count != s.count {
				panic(fmt.Errorf("Inconsistent vertex count in vertex buffer. %s has %d vertices, %s has %d",
					streams[i-1].name, streams[i-1].count, s.name, s.count))
			}
		}
		vb.streams[s.name] = s
	}
	globalStats.vertexBufferCount.inc()
	return vb
}

func (vb *vertexBuffer) release() bool {
	if !vb.refCounted.release() {
		return false
	}
	for _, s := range vb.streams {
		s.release()
	}
	globalStats.vertexBufferCount.dec()
	return true
}
