// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/google/gxui/math"
)

type ProgressBar interface {
	Control
	SetDesiredSize(math.Size)
	SetProgress(int)
	Progress() int
	SetTarget(int)
	Target() int
}
