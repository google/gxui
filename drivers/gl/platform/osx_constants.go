// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin

package platform

const ScrollSpeed = 4.0

// Paths to try when a font is not found in the local data directory
var FontPaths = []string{
	"/Library/Fonts/",
	"/System/Library/Fonts/",
}
