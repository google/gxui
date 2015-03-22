// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/google/gxui/math"
)

// AdapterItem is a user defined type that can be used to uniquely identify a
// single item in an adapter. The type must support equality and be hashable.
type AdapterItem interface{}

// Adapter is an interface used to visualize a set of items.
// Users of the Adapter should presume the data is unchanged until the
// OnDataChanged or OnDataReplaced events are fired.
type Adapter interface {
	// Count returns the number of items represented by this adapter.
	Count() int

	// ItemAt returns the AdapterItem for the item at index i. It is important
	// for the Adapter to return consistent AdapterItems for the same data item,
	// so that selections can be persisted, or re-ordering animations can be
	// played.
	// The AdapterItem returned must be equality-unique across all indices.
	ItemAt(index int) AdapterItem

	// ItemIndex returns the index of item.
	ItemIndex(item AdapterItem) int

	// Size returns the size that each of the item's controls will be displayed
	// at for the given theme.
	Size(Theme) math.Size

	// Create returns a Control visualizing the item at the specified index.
	Create(theme Theme, index int) Control

	// OnDataChanged registers f to be called when there is a partial change in
	// the items of the adapter.
	OnDataChanged(f func()) EventSubscription

	// OnDataReplaced registers f to be called when there is a complete
	// replacement of items in the adapter.
	OnDataReplaced(f func()) EventSubscription
}
