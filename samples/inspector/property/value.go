// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package property

import "reflect"

type value struct {
	name  string
	value reflect.Value
}

func (f value) Name() string {
	return f.name
}

func (f value) Type() reflect.Type {
	return f.Get().Type()
}

func (f value) Get() reflect.Value {
	return underlying(f.value)
}

func (f value) Set(value reflect.Value) {
	f.value.Set(value.Convert(f.value.Type()))
}

func (f value) CanSet() bool {
	return f.value.CanSet()
}
