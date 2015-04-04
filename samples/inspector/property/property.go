// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package property

import (
	"fmt"
	"go/ast"
	"reflect"
	"strings"
)

type Property interface {
	Name() string
	Type() reflect.Type
	Get() reflect.Value
	Set(reflect.Value)
	CanSet() bool
}

func Properties(v reflect.Value) []Property {
	properties := []Property{}

	t := v.Type()

	switch t.Kind() {
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return properties
		}

	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			properties = append(properties, value{
				name:  t.Field(i).Name,
				value: v.Field(i),
			})
		}

	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			properties = append(properties, value{
				name:  fmt.Sprintf("[%d]", i),
				value: v.Index(i),
			})
		}
	}

	for i := 0; i < v.NumMethod(); i++ {
		getter := v.Method(i)
		if !ast.IsExported(t.Method(i).Name) {
			continue
		}

		ty := getterType(getter)
		if ty == nil {
			continue
		}

		name := strings.TrimPrefix(t.Method(i).Name, "Is")

		setter := v.MethodByName("Set" + name)
		if setter.IsValid() && ty != setterType(setter) {
			setter = reflect.Value{}
		}

		properties = append(properties, getterSetter{
			name:   name,
			ty:     ty,
			getter: getter,
			setter: setter,
		})
	}
	return properties
}

func getterType(method reflect.Value) reflect.Type {
	ty := method.Type()

	// A getter is a method that takes no arguments and returns a single value.
	if ty.NumIn() != 0 || ty.NumOut() != 1 {
		return nil
	}

	return ty.Out(0)
}

func setterType(method reflect.Value) reflect.Type {
	ty := method.Type()

	// A setter is a method that takes one argument and returns no values.
	if ty.NumIn() != 1 || ty.NumOut() != 0 {
		return nil
	}

	return ty.In(0)
}

func underlying(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Interface {
		elem := v.Elem()
		if elem.IsValid() {
			return v.Elem()
		}
	}
	return v
}
