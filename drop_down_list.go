// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type DropDownList interface {
	Focusable
	Container
	SetBubbleOverlay(BubbleOverlay)
	BubbleOverlay() BubbleOverlay
	Adapter() Adapter
	SetAdapter(Adapter)
	BorderPen() Pen
	SetBorderPen(Pen)
	BackgroundBrush() Brush
	SetBackgroundBrush(Brush)
	Selected() AdapterItemId
	Select(AdapterItemId)
	OnSelectionChanged(func(AdapterItemId)) EventSubscription
	OnShowList(func()) EventSubscription
	OnHideList(func()) EventSubscription
}
