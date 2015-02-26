// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"bytes"
	"fmt"
	"time"
)

const HistorySize = 100

type GlobalStats struct {
	CanvasCount       int
	ShapeCount        int
	VertexBufferCount int
	VertexStreamCount int
	IndexBufferCount  int
}

func (s GlobalStats) String() string {
	buffer := &bytes.Buffer{}
	fmt.Fprintf(buffer, "Canvas count: %d\n", s.CanvasCount)
	fmt.Fprintf(buffer, "Shape count: %d\n", s.ShapeCount)
	fmt.Fprintf(buffer, "Vertex buffer count: %d\n", s.VertexBufferCount)
	fmt.Fprintf(buffer, "Vertex stream count: %d\n", s.VertexStreamCount)
	fmt.Fprintf(buffer, "Index buffer count: %d\n", s.IndexBufferCount)
	return buffer.String()
}

var globalStats GlobalStats

type Timer struct {
	Name     string
	Duration time.Duration
	total    time.Duration
	history  [HistorySize]time.Duration
	timer    time.Time
	current  int
}

func (t *Timer) Start() {
	t.timer = time.Now()
}

func (t *Timer) Stop() {
	duration := time.Since(t.timer)
	t.total -= t.history[t.current]
	t.history[t.current] = duration
	t.total += t.history[t.current]
	t.Duration = t.total / HistorySize
	t.current++
	if t.current >= HistorySize {
		t.current = 0
	}
}

func (t Timer) Format(f fmt.State, c rune) {
	ps := 1.0 / t.Duration.Seconds()
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
	fmt.Fprintf(f, "%s: %v [%.2f/s] %v<%v<%v", t.Name, t.Duration, ps, min, t.history[t.current], max)
}

type Stats struct {
	TextureCount              int
	FramebufferUsedCount      int
	FramebufferFreeCount      int
	FramebufferBytesAllocated int
	VertexStreamCount         int
	IndexBufferCount          int
	ShaderProgramCount        int
	FrameCount                int
	DrawCallCount             int
	Timers                    []Timer
}

func (s *Stats) Timer(name string) *Timer {
	for i := range s.Timers {
		t := &s.Timers[i]
		if t.Name == name {
			return t
		}
	}
	s.Timers = append(s.Timers, Timer{Name: name})
	return &s.Timers[len(s.Timers)-1]
}

func (s Stats) String() string {
	buffer := &bytes.Buffer{}
	for _, t := range s.Timers {
		fmt.Fprintf(buffer, "%v\n", t)
	}
	fmt.Fprintf(buffer, "Draw calls per frame: %d\n", s.DrawCallCount)
	fmt.Fprintf(buffer, "Frame count: %d\n", s.FrameCount)
	fmt.Fprintf(buffer, "Textures: %d\n", s.TextureCount)
	fmt.Fprintf(buffer, "Framebuffers Used: %d\n", s.FramebufferUsedCount)
	fmt.Fprintf(buffer, "Framebuffers Free: %d\n", s.FramebufferFreeCount)
	fmt.Fprintf(buffer, "Framebuffer bytes allocated: %d\n", s.FramebufferBytesAllocated)
	fmt.Fprintf(buffer, "Vertex stream count: %d\n", s.VertexStreamCount)
	fmt.Fprintf(buffer, "Index buffer count: %d\n", s.IndexBufferCount)
	fmt.Fprintf(buffer, "Shader program count: %d\n", s.ShaderProgramCount)
	return buffer.String()
}
