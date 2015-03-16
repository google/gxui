// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package interval

import (
	test "github.com/google/gxui/testing"
	"math/rand"
)
import "testing"

func TestU64ListMergeFromEmpty(t *testing.T) {
	l := &U64List{}
	Merge(l, CreateU64Inc(0, 0))
	test.AssertEquals(t, &U64List{
		CreateU64Inc(0, 0),
	}, l)
}

func TestU64ListDuplicate1Len(t *testing.T) {
	l := &U64List{
		CreateU64Inc(10, 10),
	}
	Merge(l, CreateU64Inc(10, 10))
	test.AssertEquals(t, &U64List{
		CreateU64Inc(10, 10),
	}, l)
}

func TestU64ListDuplicate0Len(t *testing.T) {
	l := &U64List{
		U64{10, 0},
	}
	Merge(l, U64{10, 0})
	test.AssertEquals(t, &U64List{U64{10, 0}}, l)
}

// Test for adding a new interval that does not intersect existing intervals
func TestU64ListMergeInBetween(t *testing.T) {
	l := &U64List{
		CreateU64Inc(0, 10),
		CreateU64Inc(40, 50),
	}
	Merge(l, CreateU64Inc(20, 30))
	test.AssertEquals(t, &U64List{
		CreateU64Inc(0, 10),
		CreateU64Inc(20, 30),
		CreateU64Inc(40, 50),
	}, l)
}

func TestU64ListMergeSingleBefore(t *testing.T) {
	l := &U64List{CreateU64Inc(10, 20)}
	Merge(l, CreateU64Inc(0, 5))
	test.AssertEquals(t, &U64List{
		CreateU64Inc(0, 5),
		CreateU64Inc(10, 20),
	}, l)
}

func TestU64ListMergeSingleAfter(t *testing.T) {
	l := &U64List{CreateU64Inc(0, 5)}
	Merge(l, CreateU64Inc(10, 20))
	test.AssertEquals(t, &U64List{
		CreateU64Inc(0, 5),
		CreateU64Inc(10, 20),
	}, l)
}

// Extending tests
func TestU64ListMergeExtendSingleFront(t *testing.T) {
	l := &U64List{CreateU64Inc(3, 5)}
	Merge(l, CreateU64Inc(0, 3))
	test.AssertEquals(t, &U64List{
		CreateU64Inc(0, 5),
	}, l)
}

func TestU64ListMergeExtendSingleBack(t *testing.T) {
	l := &U64List{CreateU64Inc(3, 5)}
	Merge(l, CreateU64Inc(5, 7))
	test.AssertEquals(t, &U64List{
		CreateU64Inc(3, 7),
	}, l)
}

func TestU64ListMergeExtendMiddleFront(t *testing.T) {
	l := &U64List{
		CreateU64Inc(0, 1),
		CreateU64Inc(3, 3),
		CreateU64Inc(5, 6),
	}
	Merge(l, CreateU64Inc(2, 3))
	test.AssertEquals(t, &U64List{
		CreateU64Inc(0, 1),
		CreateU64Inc(2, 3),
		CreateU64Inc(5, 6),
	}, l)
}

func TestU64ListMergeExtendMiddleBack(t *testing.T) {
	l := &U64List{
		CreateU64Inc(0, 1),
		CreateU64Inc(3, 3),
		CreateU64Inc(5, 6),
	}
	Merge(l, CreateU64Inc(3, 4))
	test.AssertEquals(t, &U64List{
		CreateU64Inc(0, 1),
		CreateU64Inc(3, 4),
		CreateU64Inc(5, 6),
	}, l)
}

// Encompassed tests
func TestU64ListMergeEdgeFront(t *testing.T) {
	l := &U64List{
		CreateU64Inc(10, 20),
	}
	Merge(l, CreateU64Inc(10, 11))
	test.AssertEquals(t, &U64List{
		CreateU64Inc(10, 20),
	}, l)
}

func TestU64ListMergeEdgeBack(t *testing.T) {
	l := &U64List{
		CreateU64Inc(10, 20),
	}
	Merge(l, CreateU64Inc(19, 20))
	test.AssertEquals(t, &U64List{
		CreateU64Inc(10, 20),
	}, l)
}

func TestU64ListMergeFirstTwo(t *testing.T) {
	l := &U64List{
		CreateU64Inc(0, 1),
		CreateU64Inc(2, 4),
		CreateU64Inc(5, 6),
	}
	Merge(l, CreateU64Inc(1, 2))
	test.AssertEquals(t, &U64List{
		CreateU64Inc(0, 4),
		CreateU64Inc(5, 6),
	}, l)
}

func TestU64ListMergeLastTwo(t *testing.T) {
	l := &U64List{
		CreateU64Inc(0, 1),
		CreateU64Inc(2, 4),
		CreateU64Inc(5, 6),
	}
	Merge(l, CreateU64Inc(3, 6))
	test.AssertEquals(t, &U64List{
		CreateU64Inc(0, 1),
		CreateU64Inc(2, 6),
	}, l)
}

func TestU64ListMergeOverlap(t *testing.T) {
	l := &U64List{
		CreateU64Inc(10, 11),
		CreateU64Inc(12, 14),
		CreateU64Inc(15, 16),
	}
	Merge(l, CreateU64Inc(11, 15))
	test.AssertEquals(t, &U64List{
		CreateU64Inc(10, 16),
	}, l)
}

func TestU64ListMergeEncompass(t *testing.T) {
	l := &U64List{
		CreateU64Inc(10, 11),
		CreateU64Inc(12, 14),
		CreateU64Inc(15, 16),
	}
	Merge(l, CreateU64Inc(0, 20))
	test.AssertEquals(t, &U64List{
		CreateU64Inc(0, 20),
	}, l)
}

func TestU64ListReplaceFromEmpty(t *testing.T) {
	l := &U64List{}
	Replace(l, CreateU64Inc(0, 10))
	test.AssertEquals(t, &U64List{
		CreateU64Inc(0, 10),
	}, l)
}

func TestU64ListReplaceSingleWhole(t *testing.T) {
	l := &U64List{CreateU64Inc(5, 10)}
	Replace(l, CreateU64Inc(5, 10))
	test.AssertEquals(t, &U64List{
		CreateU64Inc(5, 10),
	}, l)
}

func TestU64ListReplaceSingleInner(t *testing.T) {
	l := &U64List{CreateU64Inc(5, 25)}
	Replace(l, CreateU64Inc(15, 20))
	test.AssertEquals(t, &U64List{
		CreateU64Inc(5, 14),
		CreateU64Inc(15, 20),
		CreateU64Inc(21, 25),
	}, l)
}

func TestU64ListReplaceInBetween(t *testing.T) {
	l := &U64List{
		CreateU64Inc(5, 10),
		CreateU64Inc(25, 30),
	}
	Replace(l, CreateU64Inc(15, 20))
	test.AssertEquals(t, &U64List{
		CreateU64Inc(5, 10),
		CreateU64Inc(15, 20),
		CreateU64Inc(25, 30),
	}, l)
}

func TestU64ListRemoveFromEmpty(t *testing.T) {
	l := &U64List{}
	Remove(l, CreateU64Inc(0, 0))
	test.AssertEquals(t, &U64List{}, l)
}

func TestU64ListRemoveSingle(t *testing.T) {
	l := &U64List{CreateU64Inc(3, 5)}
	Remove(l, CreateU64Inc(3, 5))
	test.AssertEquals(t, &U64List{}, l)
}

func TestU64ListRemoveBeforeSingle(t *testing.T) {
	l := &U64List{CreateU64Inc(3, 5)}
	Remove(l, CreateU64Inc(0, 2))
	test.AssertEquals(t, &U64List{
		CreateU64Inc(3, 5),
	}, l)
}

func TestU64ListRemoveAfterSingle(t *testing.T) {
	l := &U64List{CreateU64Inc(3, 5)}
	Remove(l, CreateU64Inc(6, 7))
	test.AssertEquals(t, &U64List{
		CreateU64Inc(3, 5),
	}, l)
}

func TestU64ListRemoveTrimFront(t *testing.T) {
	l := &U64List{CreateU64Inc(10, 20)}
	Remove(l, CreateU64Inc(5, 14))
	test.AssertEquals(t, &U64List{
		CreateU64Inc(15, 20),
	}, l)
}

func TestU64ListRemoveSplit0(t *testing.T) {
	l := &U64List{
		CreateU64Inc(10, 20),
		CreateU64Inc(30, 40),
		CreateU64Inc(50, 60),
	}
	Remove(l, CreateU64Inc(31, 39))
	test.AssertEquals(t, &U64List{
		CreateU64Inc(10, 20),
		CreateU64Inc(30, 30),
		CreateU64Inc(40, 40),
		CreateU64Inc(50, 60),
	}, l)
}

func TestU64ListRemoveSplit1(t *testing.T) {
	l := &U64List{
		CreateU64Inc(10, 20),
		CreateU64Inc(30, 40),
		CreateU64Inc(50, 60),
	}
	Remove(l, CreateU64Inc(35, 55))
	test.AssertEquals(t, &U64List{
		CreateU64Inc(10, 20),
		CreateU64Inc(30, 34),
		CreateU64Inc(56, 60),
	}, l)
}

func TestU64ListRemoveAll0(t *testing.T) {
	l := &U64List{
		CreateU64Inc(10, 20),
		CreateU64Inc(30, 40),
		CreateU64Inc(50, 60),
	}
	Remove(l, CreateU64Inc(10, 60))
	test.AssertEquals(t, &U64List{}, l)
}

func TestU64ListRemoveAll1(t *testing.T) {
	l := &U64List{
		CreateU64Inc(10, 20),
		CreateU64Inc(30, 40),
		CreateU64Inc(50, 60),
	}
	Remove(l, CreateU64Inc(0, 100))
	test.AssertEquals(t, &U64List{}, l)
}

func TestU64ListIntersectEmpty(t *testing.T) {
	l := U64List{}
	first, count := Intersect(&l, CreateU64Inc(0, 10))
	test.AssertEquals(t, 0, first)
	test.AssertEquals(t, 0, count)
}

func TestU64ListIntersectSingle(t *testing.T) {
	l := U64List{
		CreateU64Inc(10, 20),
	}
	first, count := Intersect(&l, CreateU64Inc(10, 20))
	test.AssertEquals(t, 0, first)
	test.AssertEquals(t, 1, count)
}

func TestU64ListIntersectSingleFront(t *testing.T) {
	l := U64List{
		CreateU64Inc(10, 20),
	}
	first, count := Intersect(&l, CreateU64Inc(10, 15))
	test.AssertEquals(t, 0, first)
	test.AssertEquals(t, 1, count)
}

func TestU64ListIntersectSingleBack(t *testing.T) {
	l := U64List{
		CreateU64Inc(10, 20),
	}
	first, count := Intersect(&l, CreateU64Inc(15, 20))
	test.AssertEquals(t, 0, first)
	test.AssertEquals(t, 1, count)
}

func TestU64ListIntersectSingleMiddle(t *testing.T) {
	l := U64List{
		CreateU64Inc(10, 20),
	}
	first, count := Intersect(&l, CreateU64Inc(12, 18))
	test.AssertEquals(t, 0, first)
	test.AssertEquals(t, 1, count)
}

func TestU64ListIntersectOverlapThree(t *testing.T) {
	l := U64List{
		CreateU64Inc(10, 20),
		CreateU64Inc(30, 40),
		CreateU64Inc(50, 60),
	}
	first, count := Intersect(&l, CreateU64Inc(15, 55))
	test.AssertEquals(t, 0, first)
	test.AssertEquals(t, 3, count)
}

func TestU64ListIntersectOverlapFirstTwo(t *testing.T) {
	l := U64List{
		CreateU64Inc(10, 20),
		CreateU64Inc(30, 40),
		CreateU64Inc(50, 60),
	}
	first, count := Intersect(&l, CreateU64Inc(10, 35))
	test.AssertEquals(t, 0, first)
	test.AssertEquals(t, 2, count)
}

func TestU64ListIntersectOverlapLastTwo(t *testing.T) {
	l := U64List{
		CreateU64Inc(10, 20),
		CreateU64Inc(30, 40),
		CreateU64Inc(50, 60),
	}
	first, count := Intersect(&l, CreateU64Inc(35, 60))
	test.AssertEquals(t, 1, first)
	test.AssertEquals(t, 2, count)
}

func TestU64ListIndexOf(t *testing.T) {
	l := U64List{
		CreateU64Inc(10, 20),
		CreateU64Inc(30, 40),
		CreateU64Inc(50, 60),
	}
	test.AssertEquals(t, -1, IndexOf(&l, 0))
	test.AssertEquals(t, -1, IndexOf(&l, 9))
	test.AssertEquals(t, 0, IndexOf(&l, 10))
	test.AssertEquals(t, 0, IndexOf(&l, 15))
	test.AssertEquals(t, 0, IndexOf(&l, 20))
	test.AssertEquals(t, -1, IndexOf(&l, 21))
	test.AssertEquals(t, 1, IndexOf(&l, 32))
	test.AssertEquals(t, 2, IndexOf(&l, 60))
}

const maxIntervalValue = 100000
const maxIntervalRange = 10000

type iteration struct{ merge, replace U64 }

func buildRands(b *testing.B) []iteration {
	b.StopTimer()
	defer b.StartTimer()
	iterations := make([]iteration, b.N)
	rand.Seed(1)
	for i := range iterations {
		iterations[i].merge.first = uint64(rand.Intn(maxIntervalValue))
		iterations[i].merge.count = uint64(rand.Intn(maxIntervalRange))
		iterations[i].replace.first = uint64(rand.Intn(maxIntervalValue))
		iterations[i].replace.count = uint64(rand.Intn(maxIntervalRange))
	}
	return iterations
}

func BenchmarkGeneral(b *testing.B) {
	iterations := buildRands(b)
	l := U64List{}
	for _, iter := range iterations {
		Merge(&l, iter.merge)
		Replace(&l, iter.replace)
	}
}
