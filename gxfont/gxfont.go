// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gxfont provides default fonts.
//
// Note that the Roboto fonts are owned by Google Inc. (one of the Go
// Authors) and released under the Apache 2 license. Any notices
// distributed with applications build with GXUI and using this package
// should include both the GXUI license and the Roboto font license.
package gxfont

//go:generate go run mkfont.go

// Default is the standard GXUI sans-serif font.
var Default []byte
