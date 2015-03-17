// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package interval

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

type U64 struct {
	first uint64
	count uint64
}

type U64List []U64

func CreateU64(first, count uint64) U64 {
	return U64{first, count}
}

func CreateU64Inc(first, last uint64) U64 {
	if last < first {
		panic(fmt.Errorf("Attempting to set last before first! [0x%.16x-0x%.16x]", first, last))
	}
	return U64{first, 1 + last - first}
}

func (i U64) Expand(value uint64) U64 {
	if i.first > value {
		i.first = value
	}
	if i.Last() < value {
		i.count += value - i.Last()
	}
	return i
}

func (i U64) Contains(value uint64) bool {
	return i.first <= value && value <= i.Last()
}

func (i U64) Range() (start, end uint64) {
	return i.first, i.first + i.count
}

func (i U64) First() uint64 {
	return i.first
}

func (i U64) Last() uint64 {
	return i.first + i.count - 1
}

func (i U64) Count() uint64 {
	return i.count
}

func (i U64) String() string {
	return fmt.Sprintf("[0x%.16x-0x%.16x]", i.first, i.Last())
}

func (i U64) Span() (start, end uint64) {
	return i.first, i.first + i.count
}

// encoding.BinaryMarshaler compliance
func (i U64) MarshalBinary() ([]byte, error) {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.LittleEndian, i.first)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, i.count)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// encoding.BinaryUnmarshaler compliance
func (i *U64) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	err := binary.Read(buf, binary.LittleEndian, &i.first)
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.LittleEndian, &i.count)
	return err
}

func (l U64List) Len() int {
	return len(l)
}

func (l U64List) Cap() int {
	return cap(l)
}

func (l *U64List) SetLen(len int) {
	*l = (*l)[:len]
}

func (l *U64List) GrowTo(length, capacity int) {
	old := *l
	*l = make(U64List, length, capacity)
	copy(*l, old)
}

func (l U64List) Copy(to, from, count int) {
	copy(l[to:to+count], l[from:from+count])
}

func (l U64List) GetInterval(index int) (start, end uint64) {
	return l[index].Span()
}

func (l U64List) SetInterval(index int, start, end uint64) {
	l[index].first = start
	l[index].count = end - start
}

func (l U64List) Overlaps(i IntData) U64List {
	first, count := Intersect(l, i)
	return l[first : first+count]
}

func (l U64List) String() string {
	s := make([]string, len(l))
	for i, v := range l {
		s[i] = fmt.Sprintf("%d%s", i, v)
	}
	return "{" + strings.Join(s, ",") + "}"
}
