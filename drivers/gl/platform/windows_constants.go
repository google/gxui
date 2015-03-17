// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package platform

// #cgo LDFLAGS: -lShell32
// #include <Shlobj.h>
//
// void getWindowsFontDirectory(char* path) {
//   SHGetFolderPath(0, CSIDL_FONTS, 0, 0, path);
// }
import "C"

const ScrollSpeed = 20.0

// Paths to try when a font is not found in the local data directory
var FontPaths = []string{
	getWindowsFontDirectory(),
}

const max_path = 260

func getWindowsFontDirectory() string {
	buf := make([]C.char, max_path)
	C.getWindowsFontDirectory(&buf[0])
	return C.GoString(&buf[0])
}
