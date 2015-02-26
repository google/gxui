// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testing

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)
import "reflect"

const kDeepCompareMaxDepth = 10

func AssertEquals(t *testing.T, expected interface{}, actual interface{}) {
	if diff, equal := DeepCompare(expected, actual); !equal {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d ASSERT: Equality assertion:\n%s\n", file, line, strings.Join(diff, "\n"))
		t.Fail()
	}
}

func AssertEqual(t *testing.T, name string, expected interface{}, actual interface{}) {
	if diff, equal := DeepCompare(expected, actual); !equal {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d ASSERT: %s not equal:\n%s\n", file, line, name, strings.Join(diff, "\n"))
		t.Fail()
	}
}

func DeepCompare(expected, got interface{}) (result []string, equal bool) {
	equal = true
	if expected == nil && got != nil {
		result = append(result, fmt.Sprintf("Expected <nil>, got %T", got))
		equal = false
		return
	}
	if got == nil {
		if expected != nil {
			result = append(result, fmt.Sprintf("Expected %T, got <nil>", expected))
			equal = false
		}
		return
	}
	return deepCompareD(reflect.ValueOf(expected), reflect.ValueOf(got), "v", 0)
}

func deepCompareD(e, g reflect.Value, trail string, depth int) (result []string, equal bool) {
	eIsValid, gIsValid := e.IsValid(), g.IsValid()
	switch {
	case !gIsValid && !eIsValid:
		equal = true
		return
	case !eIsValid && gIsValid:
		result = append(result, fmt.Sprintf("%s: Expected <invalid> got %s", trail, g.Type()))
		return
	case eIsValid && !gIsValid:
		result = append(result, fmt.Sprintf("%s: Expected %s got <invalid>", trail, e.Type()))
		return
	}

	if e.Type() != g.Type() {
		result = append(result, fmt.Sprintf("%s: Expected type %s, got type %s", trail, e.Type(), g.Type()))
		return
	}

	eIsNil, gIsNil := e.Type() == nil, g.Type() == nil
	switch e.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		eIsNil, gIsNil = e.IsNil(), g.IsNil()
	}
	switch {
	case eIsNil && gIsNil:
		equal = true
		return
	case eIsNil && !gIsNil:
		result = append(result, fmt.Sprintf("%s: Expected <nil> got %s", trail, g.Type()))
		return
	case !eIsNil && gIsNil:
		result = append(result, fmt.Sprintf("%s: Expected %s got <nil>", trail, e.Type()))
		return
	}

	switch e.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Slice, reflect.Interface:
		if depth >= kDeepCompareMaxDepth {
			// result = append(result, fmt.Sprintf("%s: Reached max depth for compare while inspecting expected types %s", trail, e.Type()))
			equal = true
			return
		}
	}

	switch e.Kind() {
	case reflect.Bool:
		if e.Bool() == g.Bool() {
			equal = true
		} else {
			result = append(result, fmt.Sprintf("%s: Expected %v, got %v", trail, e.Bool(), g.Bool()))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if e.Int() == g.Int() {
			equal = true
		} else {
			result = append(result, fmt.Sprintf("%s: Expected %v, got %v", trail, e.Int(), g.Int()))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if e.Uint() == g.Uint() {
			equal = true
		} else {
			result = append(result, fmt.Sprintf("%s: Expected %v, got %v", trail, e.Uint(), g.Uint()))
		}
	case reflect.Float32, reflect.Float64:
		if e.Float() == g.Float() {
			equal = true
		} else {
			result = append(result, fmt.Sprintf("%s: Expected %v, got %v", trail, e.Float(), g.Float()))
		}
	case reflect.Complex64, reflect.Complex128:
		if e.Complex() == g.Complex() {
			equal = true
		} else {
			result = append(result, fmt.Sprintf("%s: Expected %v, got %v", trail, e.Complex(), g.Complex()))
		}
	case reflect.Array, reflect.Slice:
		equal = true
		for i := 0; i < e.Len(); i++ {
			if i < g.Len() {
				res, eq := deepCompareD(e.Index(i), g.Index(i), trail+fmt.Sprintf("[%d]", i), depth+1)
				result = append(result, res...)
				equal = equal && eq
			} else {
				result = append(result, fmt.Sprintf("%s[%d]: Missing %v", trail, i, e.Index(i)))
				equal = false
			}
		}
		for i := e.Len(); i < g.Len(); i++ {
			result = append(result, fmt.Sprintf("%s[%d]: Unexpected %v", trail, i, g.Index(i)))
			equal = false
		}
	case reflect.Chan:
		//??
		result = append(result, fmt.Sprintf("%s: Cannot compare Chans", trail))
	case reflect.Func:
		//??
		result = append(result, fmt.Sprintf("%s: Cannot compare Funcs", trail))
	case reflect.Interface, reflect.Ptr:
		res, eq := deepCompareD(e.Elem(), g.Elem(), trail, depth+1)
		result = append(result, res...)
		equal = eq
	case reflect.Map:
		if e.Len() == g.Len() {
			// TODO: Iterate both maps.
			equal = true
			for _, k := range e.MapKeys() {
				res, eq := deepCompareD(e.MapIndex(k), g.MapIndex(k), trail+fmt.Sprintf("[%v]", k), depth+1)
				result = append(result, res...)
				equal = equal && eq
			}
		} else {
			result = append(result, fmt.Sprintf("%s: Expected length %d, got %d", trail, e.Len(), g.Len()))
		}
	case reflect.String:
		if e.String() == g.String() {
			equal = true
		} else {
			result = append(result, fmt.Sprintf("%s: Expected %v, got %v", trail, e.String(), g.String()))
		}
	case reflect.Struct:
		t := e.Type()
		equal = true
		for i := 0; i < t.NumField(); i++ {
			res, eq := deepCompareD(e.Field(i), g.Field(i), trail+fmt.Sprintf(".%s", t.Field(i).Name), depth+1)
			result = append(result, res...)
			equal = equal && eq
		}
	case reflect.UnsafePointer:
		//??
		result = append(result, fmt.Sprintf("%s: Cannot compare UnsafePointers", trail))
	default:
		// Should not be possible
		panic(fmt.Errorf("Unexpected types! expected: (type:%v, kind:%v, valid: %v) got:(type:%v, kind:%v, valid: %v)", e.Type(), e.Kind(), e.IsValid(), g.Type(), g.Kind(), g.IsValid()))
	}
	return
}
