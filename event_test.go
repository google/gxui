// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import test "github.com/google/gxui/testing"
import "testing"

func TestEventNoArgs(t *testing.T) {
	e := CreateEvent(func() {})

	fired := false
	e.Listen(func() { fired = true })
	test.AssertEquals(t, false, fired)

	e.Fire()
	test.AssertEquals(t, true, fired)
}

type TestEvent func(i int, s string, b bool)

func TestEventExactArgs(t *testing.T) {
	e := CreateEvent(func(int, string, bool, int, int, bool) {})

	fired := false
	e.Listen(func(i1 int, s string, b1 bool, i2, i3 int, b2 bool) {
		test.AssertEquals(t, 1, i1)
		test.AssertEquals(t, "hello", s)
		test.AssertEquals(t, false, b1)
		test.AssertEquals(t, 2, i2)
		test.AssertEquals(t, 3, i3)
		test.AssertEquals(t, true, b2)
		fired = true
	})
	test.AssertEquals(t, false, fired)

	e.Fire(1, "hello", false, 2, 3, true)
	test.AssertEquals(t, true, fired)
}

func TestEventNilArgs(t *testing.T) {
	e := CreateEvent(func(chan int, func(), interface{}, map[int]int, *int, []int) {})

	fired := false
	e.Listen(func(c chan int, f func(), i interface{}, m map[int]int, p *int, s []int) {
		test.AssertEquals(t, true, nil == c)
		test.AssertEquals(t, true, nil == f)
		test.AssertEquals(t, true, nil == i)
		test.AssertEquals(t, true, nil == m)
		test.AssertEquals(t, true, nil == p)
		test.AssertEquals(t, true, nil == s)
		fired = true
	})
	test.AssertEquals(t, false, fired)

	e.Fire(nil, nil, nil, nil, nil, nil)
	test.AssertEquals(t, true, fired)
}

func TestEventMixedVariadic(t *testing.T) {
	e := CreateEvent(func(int, int, ...int) {})

	fired := false
	e.Listen(func(a, b int, cde ...int) {
		test.AssertEquals(t, 3, len(cde))

		test.AssertEquals(t, 0, a)
		test.AssertEquals(t, 1, b)
		test.AssertEquals(t, 2, cde[0])
		test.AssertEquals(t, 3, cde[1])
		test.AssertEquals(t, 4, cde[2])
		fired = true
	})
	e.Fire(0, 1, 2, 3, 4)
	test.AssertEquals(t, true, fired)
}

func TestEventSingleVariadic(t *testing.T) {
	e := CreateEvent(func(...int) {})

	fired := false
	e.Listen(func(va ...int) {
		test.AssertEquals(t, 3, len(va))

		test.AssertEquals(t, 2, va[0])
		test.AssertEquals(t, 3, va[1])
		test.AssertEquals(t, 4, va[2])
		fired = true
	})
	e.Fire(2, 3, 4)
	test.AssertEquals(t, true, fired)
}

func TestEventEmptyVariadic(t *testing.T) {
	e := CreateEvent(func(...int) {})

	fired := false
	e.Listen(func(va ...int) {
		test.AssertEquals(t, 0, len(va))
		fired = true
	})
	e.Fire()
	test.AssertEquals(t, true, fired)
}

func TestEventChaining(t *testing.T) {
	e1 := CreateEvent(func(int, string, bool, int, int, bool) {})
	e2 := CreateEvent(func(int, string, bool, int, int, bool) {})

	e1.Listen(e2)

	fired := false
	e2.Listen(func(i1 int, s string, b1 bool, i2, i3 int, b2 bool) {
		test.AssertEquals(t, 1, i1)
		test.AssertEquals(t, "hello", s)
		test.AssertEquals(t, false, b1)
		test.AssertEquals(t, 2, i2)
		test.AssertEquals(t, 3, i3)
		test.AssertEquals(t, true, b2)
		fired = true
	})
	test.AssertEquals(t, false, fired)

	e1.Fire(1, "hello", false, 2, 3, true)
	test.AssertEquals(t, true, fired)
}

func TestEventUnlisten(t *testing.T) {
	eI := CreateEvent(func() {})
	eJ := CreateEvent(func() {})
	eK := CreateEvent(func() {})

	e := CreateEvent(func() {})
	e.Listen(eI)
	e.Listen(eJ)
	e.Listen(eK)

	i, j, k := 0, 0, 0
	subI := eI.Listen(func() { i++ })
	subJ := eJ.Listen(func() { j++ })
	subK := eK.Listen(func() { k++ })

	test.AssertEquals(t, 0, i)
	test.AssertEquals(t, 0, j)
	test.AssertEquals(t, 0, k)

	e.Fire()
	test.AssertEquals(t, 1, i)
	test.AssertEquals(t, 1, j)
	test.AssertEquals(t, 1, k)

	subJ.Unlisten()

	e.Fire()
	test.AssertEquals(t, 2, i)
	test.AssertEquals(t, 1, j)
	test.AssertEquals(t, 2, k)

	subK.Unlisten()

	e.Fire()
	test.AssertEquals(t, 3, i)
	test.AssertEquals(t, 1, j)
	test.AssertEquals(t, 2, k)

	subI.Unlisten()

	e.Fire()
	test.AssertEquals(t, 3, i)
	test.AssertEquals(t, 1, j)
	test.AssertEquals(t, 2, k)
}

// TODO: Add tests for early signature mismatch failures
