// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"reflect"
)

type EventSubscription interface {
	Unlisten()
}

type Event interface {
	Fire(args ...interface{})
	Listen(interface{}) EventSubscription
	ParameterTypes() []reflect.Type
}
