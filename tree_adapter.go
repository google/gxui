// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"gxui/math"
)

type TreeAdapterNode interface {
	Count() int
	ItemId(index int) AdapterItemId
	ItemIndex(id AdapterItemId) int
	Create(theme Theme, index int) Control
	CreateNode(index int) TreeAdapterNode
}

type TreeAdapter interface {
	TreeAdapterNode
	ItemSize(theme Theme) math.Size
	OnDataChanged(func()) EventSubscription
	OnDataReplaced(func()) EventSubscription
}
