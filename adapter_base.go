// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type AdapterBase struct {
	onDataChanged, onDataReplaced Event
}

func (a *AdapterBase) DataChanged() {
	if a.onDataChanged != nil {
		a.onDataChanged.Fire()
	}
}

func (a *AdapterBase) DataReplaced() {
	if a.onDataReplaced != nil {
		a.onDataReplaced.Fire()
	}
}

func (a *AdapterBase) OnDataChanged(f func()) EventSubscription {
	if a.onDataChanged == nil {
		a.onDataChanged = CreateEvent(func() {})
	}
	return a.onDataChanged.Listen(f)
}

func (a *AdapterBase) OnDataReplaced(f func()) EventSubscription {
	if a.onDataReplaced == nil {
		a.onDataReplaced = CreateEvent(func() {})
	}
	return a.onDataReplaced.Listen(f)
}
