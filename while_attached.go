// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"fmt"
	"reflect"
)

// WhileAttached binds the function callback to event for the duration of c
// being attached.
// event can either be:
//  • A event of the signature:    event(callback)
//  • A function of the signature: func(callback) EventSubscription
func WhileAttached(c Control, event, callback interface{}) {
	if err := verifyWhileAttachedSignature(event, callback); err != nil {
		panic(err)
	}
	var s EventSubscription
	bind := func() {
		if e, ok := event.(Event); ok {
			s = e.Listen(callback)
		} else {
			params := []reflect.Value{reflect.ValueOf(callback)}
			res := reflect.ValueOf(event).Call(params)[0]
			s = res.Interface().(EventSubscription)
		}
	}
	if c.Attached() {
		bind()
	}
	c.OnAttach(bind)
	c.OnDetach(func() { s.Unlisten() })
}

func verifyWhileAttachedSignature(event, callback interface{}) error {
	if _, ok := event.(Event); ok {
		return nil // Leave validation up to Event
	}
	e, c := reflect.TypeOf(event), reflect.TypeOf(callback)
	if e.Kind() != reflect.Func {
		return fmt.Errorf("event must be of type Event or func, got type %T", event)
	}
	if c.Kind() != reflect.Func {
		return fmt.Errorf("callback must be of type func, got type %T", callback)
	}
	if c := e.NumIn(); c != 1 {
		return fmt.Errorf("event as func must only take 1 parameter, got %d", c)
	}
	if got := e.In(0); got != c {
		return fmt.Errorf("event as func must only take 1 parameter of type callback, got type %s, callback type: %s",
			got, c)
	}
	if c := e.NumOut(); c != 1 {
		return fmt.Errorf("event as func must only return 1 value, got %d", c)
	}
	if got, expected := e.Out(0), reflect.TypeOf((*EventSubscription)(nil)).Elem(); got != expected {
		return fmt.Errorf("event as func must only return 1 value of type %s, got type %s",
			expected, got)
	}
	return nil
}
