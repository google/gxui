// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"sort"
	"strings"
)

type FilteredListItem struct {
	Name string
	Data interface{}
}

type FilteredListAdapter struct {
	DefaultAdapter
	items []FilteredListItem
}

func (a *FilteredListAdapter) SetItems(items []FilteredListItem) {
	// Clone, as the order can be mutated
	a.items = append([]FilteredListItem{}, items...)
	a.DefaultAdapter.SetItems(a.items)
}

func (a *FilteredListAdapter) Sort(partial string) {
	partialLower := strings.ToLower(partial)
	sorter := flaSorter{items: a.items, scores: make([]int, len(a.items))}
	for i, s := range a.items {
		sorter.scores[i] = flaScore(s.Name, strings.ToLower(s.Name), partial, partialLower)
	}
	sort.Sort(sorter)
	a.DefaultAdapter.SetItems(a.items)
}

type flaSorter struct {
	items  []FilteredListItem
	scores []int
}

func (s flaSorter) Len() int {
	return len(s.items)
}

func (s flaSorter) Less(i, j int) bool {
	return s.scores[i] > s.scores[j]
}

func (s flaSorter) Swap(i, j int) {
	items, scores := s.items, s.scores
	items[i], items[j] = items[j], items[i]
	scores[i], scores[j] = scores[j], scores[i]
}

func flaScore(str, strLower, partial, partialLower string) int {
	l := len(partial)
	for i := l; i > 0; i-- {
		c := 0
		if strings.Contains(str, partial[:i]) {
			c = i*i + 1
		} else if strings.Contains(strLower, partialLower[:i]) {
			c = i * i
		}

		if c > 0 {
			return c + flaScore(str, strLower, partial[i:], partialLower[i:])
		}
	}
	return 0
}
