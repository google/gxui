// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package assert

import (
	"fmt"
	"reflect"
	"runtime"
)

const debugAssertsEnabled = true

func err(msg string, args ...interface{}) {
	if args != nil {
		msg = fmt.Sprintf(msg, args...)
	}
	_, file, line, _ := runtime.Caller(2)
	panic(fmt.Errorf("%s:%d %s\n", file, line, msg))
}

func Assert(msg string, args ...interface{}) {
	if debugAssertsEnabled && args != nil {
		msg = fmt.Sprintf(msg, args...)
	}
	err("ASSERT: %s", msg)
}

func NoError(e error) {
	if debugAssertsEnabled && e != nil {
		err("ASSERT: Error '%s' returned", e.Error())
	}
}

func True(v bool, msg string, args ...interface{}) {
	if debugAssertsEnabled && !v {
		if args != nil {
			msg = fmt.Sprintf(msg, args...)
		}
		err("ASSERT: %s", msg)
	}
}

func False(v bool, msg string, args ...interface{}) {
	if debugAssertsEnabled && v {
		if args != nil {
			msg = fmt.Sprintf(msg, args...)
		}
		err("ASSERT: %s", msg)
	}
}

func safeIsNil(x interface{}) bool {
	v := reflect.ValueOf(x)
	switch v.Type().Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

func NotNil(v interface{}, name string, args ...interface{}) {
	if debugAssertsEnabled && v == nil || safeIsNil(v) {
		if args != nil {
			name = fmt.Sprintf(name, args...)
		}
		err("ASSERT: %s was nil\n", name)
	}
}

func Nil(v interface{}, name string, args ...interface{}) {
	if debugAssertsEnabled && v != nil && !safeIsNil(v) {
		if args != nil {
			name = fmt.Sprintf(name, args...)
		}
		err("ASSERT: %s was not nil. Got: %+v (type: %T)\n", name, v, v)
	}
}

func Equals(expected, actual interface{}, name string, args ...interface{}) {
	if debugAssertsEnabled && expected != actual {
		if args != nil {
			name = fmt.Sprintf(name, args...)
		}
		err("ASSERT: %s was not the expected value.\nExpected: %+v (type: %T)\nGot: %+v (type: %T)\n",
			name, expected, expected, actual, actual)
	}
}
