// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type DropDownList interface {
	Focusable
	Container
	SetBubbleOverlay(BubbleOverlay)
	BubbleOverlay() BubbleOverlay
	Adapter() ListAdapter
	SetAdapter(ListAdapter)
	BorderPen() Pen
	SetBorderPen(Pen)
	BackgroundBrush() Brush
	SetBackgroundBrush(Brush)
	Selected() AdapterItem
	Select(AdapterItem)
	OnSelectionChanged(func(AdapterItem)) EventSubscription
	OnShowList(func()) EventSubscription
	OnHideList(func()) EventSubscription
}
