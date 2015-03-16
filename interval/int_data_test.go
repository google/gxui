// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package interval

import test "github.com/google/gxui/testing"

import "testing"

func TestIntDataListReplace(t *testing.T) {
	l := &IntDataList{
		CreateIntData(10, 30, 1),
		CreateIntData(40, 80, 2),
	}
	Replace(l, CreateIntData(00, 15, 3))
	Replace(l, CreateIntData(25, 45, 4))
	Replace(l, CreateIntData(55, 65, 5))
	Replace(l, CreateIntData(75, 90, 6))
	test.AssertEquals(t, &IntDataList{
		CreateIntData(00, 15, 3),
		CreateIntData(15, 25, 1),
		CreateIntData(25, 45, 4),
		CreateIntData(45, 55, 2),
		CreateIntData(55, 65, 5),
		CreateIntData(65, 75, 2),
		CreateIntData(75, 90, 6),
	}, l)
}

func TestIntDataListMerge(t *testing.T) {
	l := &IntDataList{
		CreateIntData(10, 30, 1),
		CreateIntData(40, 80, 2),
	}
	Merge(l, CreateIntData(00, 15, 3))
	Merge(l, CreateIntData(25, 45, 4))
	Merge(l, CreateIntData(55, 65, 5))
	Merge(l, CreateIntData(75, 90, 6))
	test.AssertEquals(t, &IntDataList{
		CreateIntData(00, 90, 6),
	}, l)
}
