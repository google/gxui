// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

// AdapterItem is a user defined type that can be used to uniquely identify a
// single item in an adapter. The type must support equality and be hashable.
type AdapterItem interface{}
