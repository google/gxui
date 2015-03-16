// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import "github.com/google/gxui/interval"

type TextSelectionList []TextSelection

func (l TextSelectionList) Transform(from int, transform func(i int) int) TextSelectionList {
	res := TextSelectionList{}
	for _, s := range l {
		start := s.start
		end := s.end
		if start >= from {
			start = transform(start)
		}
		if end >= from {
			end = transform(end)
		}
		interval.Merge(&res, TextSelection{start, end, s.caretAtStart})
	}
	return res
}

func (l TextSelectionList) TransformCarets(from int, transform func(i int) int) TextSelectionList {
	res := TextSelectionList{}
	for _, s := range l {
		if s.caretAtStart && s.start >= from {
			s.start = transform(s.start)
		} else if s.end >= from {
			s.end = transform(s.end)
		}
		if s.start > s.end {
			tmp := s.start
			s.start = s.end
			s.end = tmp
			s.caretAtStart = !s.caretAtStart
		}
		interval.Merge(&res, s)
	}
	return res
}

func (l TextSelectionList) Len() int {
	return len(l)
}

func (l TextSelectionList) Cap() int {
	return cap(l)
}

func (l *TextSelectionList) SetLen(len int) {
	*l = (*l)[:len]
}

func (l *TextSelectionList) GrowTo(length, capacity int) {
	old := *l
	*l = make(TextSelectionList, length, capacity)
	copy(*l, old)
}

func (l TextSelectionList) Copy(to, from, count int) {
	copy(l[to:to+count], l[from:from+count])
}

func (l TextSelectionList) GetInterval(index int) (start, end uint64) {
	return l[index].Span()
}

func (l TextSelectionList) SetInterval(index int, start, end uint64) {
	l[index].start = int(start)
	l[index].end = int(end)
}

func (l TextSelectionList) MergeData(index int, i interval.Node) {
	l[index].caretAtStart = i.(TextSelection).caretAtStart
}
