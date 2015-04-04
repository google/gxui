// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package property

import "reflect"

type getterSetter struct {
	name   string
	ty     reflect.Type
	getter reflect.Value
	setter reflect.Value
}

func (p getterSetter) Name() string {
	return p.name
}

func (p getterSetter) Type() reflect.Type {
	return p.ty
}

func (p getterSetter) Get() reflect.Value {
	return underlying(p.getter.Call([]reflect.Value{})[0])
}

func (p getterSetter) Set(value reflect.Value) {
	p.setter.Call([]reflect.Value{value.Convert(p.ty)})
}

func (p getterSetter) CanSet() bool {
	return p.setter.IsValid()
}
