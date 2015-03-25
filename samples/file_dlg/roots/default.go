// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows

package roots

import (
	"os"
	"path/filepath"
)

// Roots returns the list of root directories avaliable on this machine.
func Roots() []string {
	roots := []string{}
	filepath.Walk("/", func(subpath string, info os.FileInfo, err error) error {
		if err == nil && "/" != subpath {
			roots = append(roots, subpath)
			if info.IsDir() {
				return filepath.SkipDir
			}
		}
		return nil
	})
	return roots
}
