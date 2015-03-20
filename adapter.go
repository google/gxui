// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/google/gxui/math"
)

// AdapterItem is a user defined type that can be used to uniquely identify a
// single item in an adapter. The type must support equality.
type AdapterItem interface{}

type Adapter interface {
	ItemSize(theme Theme) math.Size
	Count() int
	ItemAt(index int) AdapterItem
	ItemIndex(item AdapterItem) int
	Create(theme Theme, index int) Control
	OnDataChanged(func()) EventSubscription
	OnDataReplaced(func()) EventSubscription
}
