// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gxfont provides default fonts.
//
// Note that the Roboto and Droid Sans Mono fonts are owned by
// Google Inc. (one of the Go Authors) and released under the Apache 2
// license. Any notices distributed with applications build with GXUI
// and using this package should include both the GXUI license and the
// fonts license.
package gxfont

//go:generate go run mkfont.go

import (
	"bytes"
	"compress/flate"
	"io/ioutil"
)

var (
	// Default is the standard GXUI sans-serif font.
	Default []byte = inflate(roboto_regular)

	// Monospace is the standard GXUI fixed-width font.
	Monospace []byte = inflate(droid_sans_mono)
)

func inflate(src []byte) []byte {
	r := bytes.NewReader(src)
	b, err := ioutil.ReadAll(flate.NewReader(r))
	if err != nil {
		panic(err)
	}
	return b
}
