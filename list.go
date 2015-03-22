// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type List interface {
	Focusable
	Parent
	Adapter() Adapter
	SetAdapter(Adapter)
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
	Select(AdapterItem)
	OnSelectionChanged(func(AdapterItem)) EventSubscription
	OnItemClicked(func(MouseEvent, AdapterItem)) EventSubscription
}
