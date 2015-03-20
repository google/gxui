// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parts

import (
	"github.com/google/gxui"
	"github.com/google/gxui/mixins/outer"
)

type AttachableOuter interface {
	outer.Relayouter
}

type Attachable struct {
	outer    AttachableOuter
	onAttach gxui.Event
	onDetach gxui.Event
	attached bool
}

func (a *Attachable) Init(outer AttachableOuter) {
	a.outer = outer
}

func (a *Attachable) Attached() bool {
	return a.attached
}

func (a *Attachable) Attach() {
	if a.attached {
		panic("Control already attached")
	}
	a.attached = true
	if a.onAttach != nil {
		a.onAttach.Fire()
	}
}

func (a *Attachable) Detach() {
	if !a.attached {
		panic("Control already detached")
	}
	a.attached = false
	if a.onDetach != nil {
		a.onDetach.Fire()
	}
}

func (a *Attachable) OnAttach(f func()) gxui.EventSubscription {
	if a.onAttach == nil {
		a.onAttach = gxui.CreateEvent(func() {})
	}
	return a.onAttach.Listen(f)
}

func (a *Attachable) OnDetach(f func()) gxui.EventSubscription {
	if a.onDetach == nil {
		a.onDetach = gxui.CreateEvent(func() {})
	}
	return a.onDetach.Listen(f)
}
