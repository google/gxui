// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import "testing"

func TestRound(t *testing.T) {
	check := []struct {
		value    float32
		expected int
	}{
		{-1.1, -1},
		{-0.9, -1},
		{-0.5, -1},
		{-0.1, 0},
		{0.1, 0},
		{0.5, 0},
		{0.9, 1},
		{1.1, 1},
	}

	for _, v := range check {
		got := Round(v.value)
		if got != v.expected {
			t.Errorf("Round(%v) returned unexpected value. Expected: %v, Got: %v", v.value, v.expected, got)
		}
	}
}
