// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import test "github.com/google/gxui/testing"
import (
	"github.com/google/gxui/interval"
	"testing"
)

func TestTextSelectionMergeOne(t *testing.T) {
	s := TextSelection{5, 10, true}
	l := TextSelectionList{}
	interval.Merge(&l, s)
	test.AssertEquals(t, TextSelectionList{s}, l)
}

func TestTextSelectionMergeInner(t *testing.T) {
	s1 := TextSelection{5, 10, true}
	s2 := TextSelection{6, 9, false}
	l := TextSelectionList{s1}
	interval.Merge(&l, s2)
	test.AssertEquals(t, TextSelectionList{
		TextSelection{5, 10, false},
	}, l)
}

func TestTextSelectionMergeAtStart(t *testing.T) {
	s1 := TextSelection{6, 9, true}
	s2 := TextSelection{6, 7, false}
	l := TextSelectionList{s1}
	interval.Merge(&l, s2)
	test.AssertEquals(t, TextSelectionList{
		TextSelection{6, 9, false},
	}, l)
}

func TestTextSelectionMergeAtEnd(t *testing.T) {
	s1 := TextSelection{6, 9, true}
	s2 := TextSelection{8, 9, false}
	l := TextSelectionList{s1}
	interval.Merge(&l, s2)
	test.AssertEquals(t, TextSelectionList{
		TextSelection{6, 9, false},
	}, l)
}

func TestTextSelectionMergeEncompass(t *testing.T) {
	s1 := TextSelection{6, 9, false}
	s2 := TextSelection{5, 10, true}
	l := TextSelectionList{s1}
	interval.Merge(&l, s2)
	test.AssertEquals(t, TextSelectionList{
		TextSelection{5, 10, true},
	}, l)
}

func TestTextSelectionMergeDuplicate(t *testing.T) {
	s1 := TextSelection{2, 6, false}
	s2 := TextSelection{2, 6, true}
	l := TextSelectionList{s1}
	interval.Merge(&l, s2)
	test.AssertEquals(t, TextSelectionList{
		TextSelection{2, 6, true},
	}, l)
}

func TestTextSelectionMergeDuplicate0Len(t *testing.T) {
	s1 := TextSelection{2, 2, false}
	s2 := TextSelection{2, 2, true}
	l := TextSelectionList{s1}
	interval.Merge(&l, s2)
	test.AssertEquals(t, TextSelectionList{
		TextSelection{2, 2, true},
	}, l)
}

func TestTextSelectionMergeExtendStart(t *testing.T) {
	s1 := TextSelection{6, 9, false}
	s2 := TextSelection{1, 7, true}
	l := TextSelectionList{s1}
	interval.Merge(&l, s2)
	test.AssertEquals(t, TextSelectionList{
		TextSelection{1, 9, true},
	}, l)
}

func TestTextSelectionMergeExtendEnd(t *testing.T) {
	s1 := TextSelection{6, 9, true}
	s2 := TextSelection{8, 15, false}
	l := TextSelectionList{s1}
	interval.Merge(&l, s2)
	test.AssertEquals(t, TextSelectionList{
		TextSelection{6, 15, false},
	}, l)
}

func TestTextSelectionMergeBeforeStart(t *testing.T) {
	s1 := TextSelection{6, 9, true}
	s2 := TextSelection{2, 6, false}
	l := TextSelectionList{s1}
	interval.Merge(&l, s2)
	test.AssertEquals(t, TextSelectionList{
		TextSelection{2, 6, false},
		TextSelection{6, 9, true},
	}, l)
}

func TestTextSelectionMergeAfterEnd(t *testing.T) {
	s1 := TextSelection{2, 6, false}
	s2 := TextSelection{6, 9, true}
	l := TextSelectionList{s1}
	interval.Merge(&l, s2)
	test.AssertEquals(t, TextSelectionList{
		TextSelection{2, 6, false},
		TextSelection{6, 9, true},
	}, l)
}
