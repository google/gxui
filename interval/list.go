// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package interval

import "sort"

type Node interface {
	Span() (start, end uint64)
}

type RList interface {
	Len() int
	GetInterval(index int) (start, end uint64)
	SetInterval(index int, start, end uint64)
}

type List interface {
	RList
	Copy(to, from, count int)
	Cap() int
	SetLen(len int)
	GrowTo(length, capacity int)
}

type ExtendedList interface {
	MergeData(index int, i Node)
}

type intersection struct {
	overlap            int
	lowIndex           int
	lowStart, lowEnd   uint64
	intersectsLow      bool
	highIndex          int
	highStart, highEnd uint64
	intersectsHigh     bool
}

const (
	minSpace = 3 // Max growth of 2 plus one slot for temporary
	minCap   = 5
)

func Merge(l List, i Node) {
	start, end := i.Span()
	s := intersection{}
	s.intersect(l, start, end)
	adjust(l, s.lowIndex, 1-s.overlap)
	if s.intersectsLow {
		start = s.lowStart
	}
	if s.intersectsHigh {
		end = s.highEnd
	}
	l.SetInterval(s.lowIndex, start, end)
	if dl, ok := l.(ExtendedList); ok {
		dl.MergeData(s.lowIndex, i)
	}
}

func Replace(l List, i Node) {
	start, end := i.Span()
	index, start, end := replace(l, start, end, true)
	l.SetInterval(index, start, end)
	if dl, ok := l.(ExtendedList); ok {
		dl.MergeData(index, i)
	}
}

func Remove(l List, i Node) {
	start, end := i.Span()
	replace(l, start, end, false)
}

func Intersect(l RList, i Node) (first, count int) {
	start, end := i.Span()
	s := intersection{}
	s.intersect(l, start, end)
	return s.lowIndex, s.overlap
}

type Visitor func(start, end uint64, index int)

func Visit(l RList, i Node, v Visitor) {
	start, end := i.Span()
	s := intersection{}
	s.intersect(l, start, end)
	for index := s.lowIndex; index < s.lowIndex+s.overlap; index++ {
		s, e := l.GetInterval(index)
		if s < start {
			s = start
		}
		if e > end {
			e = end
		}
		v(s, e, index)
	}
}

func Contains(l RList, p uint64) bool {
	return IndexOf(l, p) >= 0
}

func IndexOf(l RList, p uint64) int {
	index := sort.Search(l.Len(), func(at int) bool {
		iStart, _ := l.GetInterval(at)
		return p < iStart
	})
	index--
	if index >= 0 {
		_, iEnd := l.GetInterval(index)
		if p < iEnd {
			return index
		}
	}
	return -1
}

func FindStart(l RList, at int, start uint64) bool {
	_, iEnd := l.GetInterval(at)
	return start < iEnd
}

func FindEnd(l RList, at int, end uint64) bool {
	iStart, _ := l.GetInterval(at)
	return end <= iStart
}

func Search(l RList, v uint64, f func(l RList, at int, v uint64) bool) int {
	i, j := 0, l.Len()
	for i < j {
		h := i + (j-i)/2
		if !f(l, h, v) {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}

func (s *intersection) intersect(l RList, start, end uint64) {
	beforeLen := Search(l, start, FindStart)
	afterIndex := Search(l, end, FindEnd)
	if afterIndex < beforeLen {
		afterIndex, beforeLen = beforeLen, afterIndex
	}
	s.lowIndex = beforeLen
	s.highIndex = afterIndex - 1
	s.overlap = afterIndex - beforeLen
	s.intersectsLow = false
	s.intersectsHigh = false
	if s.overlap > 0 {
		s.lowStart, s.lowEnd = l.GetInterval(s.lowIndex)
		s.intersectsLow = s.lowStart < start
		s.highStart, s.highEnd = l.GetInterval(s.highIndex)
		s.intersectsHigh = end < s.highEnd
	}
}

func adjust(l List, at, delta int) {
	if delta == 0 {
		return
	}
	oldLen := l.Len()
	newLen := oldLen + delta
	if delta > 0 {
		cap := l.Cap()
		if cap < newLen {
			newCap := newLen + (newLen >> 1)
			l.GrowTo(newLen, newCap)
		} else {
			l.SetLen(newLen)
		}
	}
	copyStart := at - delta
	copyTo := at
	if copyStart < 0 {
		copyTo -= copyStart
		copyStart = 0
	}
	l.Copy(copyTo, copyStart, newLen-copyTo)
	if delta < 0 {
		l.SetLen(newLen)
	}
}

func replace(l List, start, end uint64, add bool) (int, uint64, uint64) {
	s := intersection{}
	s.intersect(l, start, end)
	if s.overlap == 0 {
		if add {
			adjust(l, s.lowIndex, 1)
		}
		return s.lowIndex, start, end
	}

	insertLen := 0
	insertPoint := s.lowIndex
	if s.intersectsLow {
		s.lowEnd = start
		insertLen++
		insertPoint++
	}
	if add {
		insertLen++
	}
	if s.intersectsHigh {
		s.highStart = end
		insertLen++
	}
	delta := insertLen - s.overlap
	adjust(l, insertPoint, delta)
	if s.intersectsLow {
		l.SetInterval(s.lowIndex, s.lowStart, s.lowEnd)
	}
	if s.intersectsHigh {
		l.SetInterval(s.lowIndex+insertLen-1, s.highStart, s.highEnd)
	}
	return insertPoint, start, end
}
