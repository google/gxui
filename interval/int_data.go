// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package interval

type IntData struct {
	start, end int
	data       interface{}
}

type IntDataList []IntData

func CreateIntData(start, end int, data interface{}) IntData {
	return IntData{
		start: start,
		end:   end,
		data:  data,
	}
}

func (t IntData) Range() (start, end int) {
	return t.start, t.end
}

func (t IntData) Data() interface{} {
	return t.data
}

func (t IntData) Span() (start, end uint64) {
	return uint64(t.start), uint64(t.end)
}

func (t IntData) Contains(v int) bool {
	return v >= t.start && v < t.end
}

func (l IntDataList) Len() int {
	return len(l)
}

func (l IntDataList) Cap() int {
	return cap(l)
}

func (l *IntDataList) SetLen(len int) {
	*l = (*l)[:len]
}

func (l *IntDataList) GrowTo(length, capacity int) {
	old := *l
	*l = make(IntDataList, length, capacity)
	copy(*l, old)
}

func (l IntDataList) Copy(to, from, count int) {
	copy(l[to:to+count], l[from:from+count])
}

func (l IntDataList) GetInterval(index int) (start, end uint64) {
	return l[index].Span()
}

func (l IntDataList) SetInterval(index int, start, end uint64) {
	l[index].start = int(start)
	l[index].end = int(end)
}

func (l IntDataList) MergeData(index int, i Node) {
	l[index].data = i.(IntData).data
}

func (l IntDataList) Overlaps(i IntData) IntDataList {
	first, count := Intersect(l, i)
	return l[first : first+count]
}
