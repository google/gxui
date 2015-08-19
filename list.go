// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import "github.com/google/gxui/math"

type List interface {
	Focusable
	Parent
	Adapter() ListAdapter
	SetAdapter(ListAdapter)
	SetOrientation(Orientation)
	Orientation() Orientation
	BorderPen() Pen
	SetBorderPen(Pen)
	BackgroundBrush() Brush
	SetBackgroundBrush(Brush)
	ScrollTo(AdapterItem)
	IsItemVisible(AdapterItem) bool
	ItemControl(AdapterItem) Control
	Selected() AdapterItem
	Select(AdapterItem) bool
	OnSelectionChanged(func(AdapterItem)) EventSubscription
	OnItemClicked(func(MouseEvent, AdapterItem)) EventSubscription
}

// ListAdapter is an interface used to visualize a flat set of items.
// Users of the ListAdapter should presume the data is unchanged until the
// OnDataChanged or OnDataReplaced events are fired.
type ListAdapter interface {
	// Count returns the total number of items.
	Count() int

	// ItemAt returns the AdapterItem for the item at index i. It is important
	// for the Adapter to return consistent AdapterItems for the same data, so
	// that selections can be persisted, or re-ordering animations can be played
	// when the dataset changes.
	// The AdapterItem returned must be equality-unique across all indices.
	ItemAt(index int) AdapterItem

	// ItemIndex returns the index of item, or -1 if the adapter does not contain
	// item.
	ItemIndex(item AdapterItem) int

	// Create returns a Control visualizing the item at the specified index.
	Create(theme Theme, index int) Control

	// Size returns the size that each of the item's controls will be displayed
	// at for the given theme.
	Size(Theme) math.Size

	// OnDataChanged registers f to be called when there is a partial change in
	// the items of the adapter. Scroll positions and selections should be
	// preserved if possible.
	// If recreateControls is true then each of the visible controls should be
	// recreated by re-calling Create().
	OnDataChanged(f func(recreateControls bool)) EventSubscription

	// OnDataReplaced registers f to be called when there is a complete
	// replacement of items in the adapter.
	OnDataReplaced(f func()) EventSubscription
}
