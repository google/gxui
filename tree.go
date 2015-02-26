// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type Tree interface {
	Focusable
	SetAdapter(TreeAdapter)
	Adapter() TreeAdapter
	Selected() AdapterItemId
	Select(AdapterItemId)
	OnSelectionChanged(func(AdapterItemId)) EventSubscription
}
