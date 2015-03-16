// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/google/gxui/math"
	"sort"
	"strings"
)

type FilteredListItem struct {
	Name string
	Data interface{}
}

type FilteredListAdapter struct {
	AdapterBase
	items  []FilteredListItem
	order  []int
	rorder []int
	scores []int
}

func (s *FilteredListAdapter) score(str, strLower, partial, partialLower string) int {
	l := len(partial)
	for i := l; i > 0; i-- {
		c := 0
		if strings.Contains(str, partial[:i]) {
			c = i*i + 1
		} else if strings.Contains(strLower, partialLower[:i]) {
			c = i * i
		}

		if c > 0 {
			return c + s.score(str, strLower, partial[i:], partialLower[i:])
		}
	}
	return 0
}

func (a *FilteredListAdapter) SetItems(items []FilteredListItem) {
	a.items = items
	a.order = make([]int, len(items))
	a.rorder = make([]int, len(items))
	a.scores = make([]int, len(items))
	for i, _ := range a.order {
		a.order[i] = i
		a.rorder[i] = i
	}
	a.onDataReplaced.Fire()
}

func (a *FilteredListAdapter) Item(id AdapterItemId) FilteredListItem {
	return a.items[id]
}

func (a *FilteredListAdapter) Sort(partial string) {
	for i, _ := range a.order {
		a.order[i] = i
	}
	partialLower := strings.ToLower(partial)
	for i, s := range a.items {
		a.scores[i] = a.score(s.Name, strings.ToLower(s.Name), partial, partialLower)
	}
	sort.Sort(a)
	for i, j := range a.order {
		a.rorder[j] = i
	}
	a.DataChanged()
}

// sort.Interface compliance
func (a *FilteredListAdapter) Len() int {
	return len(a.items)
}

func (a *FilteredListAdapter) Less(i, j int) bool {
	return a.scores[a.order[i]] > a.scores[a.order[j]]
}

func (a *FilteredListAdapter) Swap(i, j int) {
	t := a.order[i]
	a.order[i] = a.order[j]
	a.order[j] = t
}

// Adapter compliance
func (a *FilteredListAdapter) ItemSize(theme Theme) math.Size {
	return math.Size{W: 200, H: 16}
}

func (a *FilteredListAdapter) Count() int {
	return len(a.items)
}

func (a *FilteredListAdapter) ItemId(index int) AdapterItemId {
	return AdapterItemId(a.order[index])
}

func (a *FilteredListAdapter) ItemIndex(id AdapterItemId) int {
	return a.rorder[id]
}

func (a *FilteredListAdapter) Create(theme Theme, index int) Control {
	item := a.Item(a.ItemId(index))
	l := theme.CreateLabel()
	l.SetMargin(math.ZeroSpacing)
	l.SetMultiline(false)
	l.SetText(item.Name)
	return l
}
