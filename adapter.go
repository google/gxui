// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"gxui/math"
)

const InvalidAdapterItemId AdapterItemId = 0xFFFFFFFFFFFFFFFF

type AdapterItemId uint64

func (i AdapterItemId) IsValid() bool { return i != InvalidAdapterItemId }

type Adapter interface {
	ItemSize(theme Theme) math.Size
	Count() int
	ItemId(index int) AdapterItemId
	ItemIndex(id AdapterItemId) int
	Create(theme Theme, index int) Control
	OnDataChanged(func()) EventSubscription
	OnDataReplaced(func()) EventSubscription
}
