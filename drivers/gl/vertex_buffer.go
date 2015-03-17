// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import "fmt"

type VertexBuffer struct {
	refCounted
	Streams     map[string]*VertexStream
	VertexCount int
}

func CreateVertexBuffer(streams ...*VertexStream) *VertexBuffer {
	vb := &VertexBuffer{
		Streams: map[string]*VertexStream{},
	}
	vb.init()
	for i, s := range streams {
		if i == 0 {
			vb.VertexCount = s.VertexCount()
		} else {
			if vb.VertexCount != s.VertexCount() {
				panic(fmt.Errorf("Inconsistent vertex count in vertex buffer. %s has %d vertices, %s has %d",
					streams[i-1].Name(), streams[i-1].VertexCount(), s.Name(), s.VertexCount()))
			}
		}
		vb.Streams[s.Name()] = s
	}
	globalStats.VertexBufferCount++
	return vb
}

func (vb *VertexBuffer) Release() {
	if vb.release() {
		for _, s := range vb.Streams {
			s.Release()
		}
		globalStats.VertexBufferCount--
	}
}
