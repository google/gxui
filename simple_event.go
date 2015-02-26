// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type SimpleEvent struct {
	EventBase
}

func CreateEvent(signature interface{}) Event {
	e := &SimpleEvent{}
	e.init(signature)
	return e
}
