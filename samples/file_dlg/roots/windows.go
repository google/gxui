// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package roots

import "os"

// Roots returns the list of drives avaliable on this machine.
func Roots() []string {
	roots := []string{}
	for drive := 'A'; drive <= 'Z'; drive++ {
		path := string(drive) + ":"
		if _, err := os.Stat(path); err == nil {
			roots = append(roots, path)
		}
	}
	return roots
}
