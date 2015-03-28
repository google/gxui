// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"bytes"
	"fmt"
	"sync/atomic"
	"time"
)

const historySize = 100

type count int32

func (c *count) inc() {
	atomic.AddInt32((*int32)(c), 1)
}

func (c *count) dec() {
	if atomic.AddInt32((*int32)(c), -1) < 0 {
		panic("Count has gone negative")
	}
}

type globalDriverStats struct {
	canvasCount       count
	shapeCount        count
	vertexBufferCount count
	vertexStreamCount count
	indexBufferCount  count
}

func (s globalDriverStats) String() string {
	buffer := &bytes.Buffer{}
	fmt.Fprintf(buffer, "Canvas count: %d\n", s.canvasCount)
	fmt.Fprintf(buffer, "Shape count: %d\n", s.shapeCount)
	fmt.Fprintf(buffer, "Vertex buffer count: %d\n", s.vertexBufferCount)
	fmt.Fprintf(buffer, "Vertex stream count: %d\n", s.vertexStreamCount)
	fmt.Fprintf(buffer, "Index buffer count: %d\n", s.indexBufferCount)
	return buffer.String()
}

var globalStats globalDriverStats

type timer struct {
	name     string
	duration time.Duration
	total    time.Duration
	history  [historySize]time.Duration
	timer    time.Time
	current  int
}

func (t *timer) start() {
	t.timer = time.Now()
}

func (t *timer) stop() {
	duration := time.Since(t.timer)
	t.total -= t.history[t.current]
	t.history[t.current] = duration
	t.total += t.history[t.current]
	t.duration = t.total / historySize
	t.current++
	if t.current >= historySize {
		t.current = 0
	}
}

func (t timer) Format(f fmt.State, c rune) {
	ps := 1.0 / t.duration.Seconds()
	max := t.history[0]
	min := t.history[0]
	for _, h := range t.history {
		if max < h {
			max = h
		}
		if min > h {
			min = h
		}
	}
	fmt.Fprintf(f, "%s: %v [%.2f/s] %v<%v<%v", t.name, t.duration, ps, min, t.history[t.current], max)
}

type contextStats struct {
	textureCount              int
	framebufferUsedCount      int
	framebufferFreeCount      int
	framebufferBytesAllocated int
	vertexStreamCount         int
	indexBufferCount          int
	shaderProgramCount        int
	frameCount                int
	drawCallCount             int
	timers                    []timer
}

func (s *contextStats) timer(name string) *timer {
	for i := range s.timers {
		t := &s.timers[i]
		if t.name == name {
			return t
		}
	}
	s.timers = append(s.timers, timer{name: name})
	return &s.timers[len(s.timers)-1]
}

func (s contextStats) String() string {
	buffer := &bytes.Buffer{}
	for _, t := range s.timers {
		fmt.Fprintf(buffer, "%v\n", t)
	}
	fmt.Fprintf(buffer, "Draw calls per frame: %d\n", s.drawCallCount)
	fmt.Fprintf(buffer, "Frame count: %d\n", s.frameCount)
	fmt.Fprintf(buffer, "Textures: %d\n", s.textureCount)
	fmt.Fprintf(buffer, "Framebuffers Used: %d\n", s.framebufferUsedCount)
	fmt.Fprintf(buffer, "Framebuffers Free: %d\n", s.framebufferFreeCount)
	fmt.Fprintf(buffer, "Framebuffer bytes allocated: %d\n", s.framebufferBytesAllocated)
	fmt.Fprintf(buffer, "Vertex stream count: %d\n", s.vertexStreamCount)
	fmt.Fprintf(buffer, "Index buffer count: %d\n", s.indexBufferCount)
	fmt.Fprintf(buffer, "Shader program count: %d\n", s.shaderProgramCount)
	return buffer.String()
}
