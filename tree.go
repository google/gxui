// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type Tree interface {
	Focusable
	SetAdapter(TreeAdapter)
	Adapter() TreeAdapter

	// Show makes the item with the specified id visible, expanding the tree if
	// necessary.
	Show(AdapterItemId)

	// ExpandAll expands all tree nodes.
	ExpandAll()

	// CollapseAll collapses all tree nodes.
	CollapseAll()

	Selected() AdapterItemId
	Select(AdapterItemId)
	OnSelectionChanged(func(AdapterItemId)) EventSubscription
}
