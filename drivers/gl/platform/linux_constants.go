// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux

package platform

const ScrollSpeed = 20.0

// Paths to try when a font is not found in the local data directory
var FontPaths = []string{
	"/usr/share/fonts",
	"/usr/local/share/fonts",
	"~/.fonts",
}
