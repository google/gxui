// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type TextSelection struct {
	start, end   int
	caretAtStart bool
}

func CreateTextSelection(start, end int, caretAtStart bool) TextSelection {
	if start < end {
		return TextSelection{start, end, caretAtStart}
	} else {
		return TextSelection{end, start, !caretAtStart}
	}
}

func (i TextSelection) Length() int             { return i.end - i.start }
func (i TextSelection) Range() (start, end int) { return i.start, i.end }
func (i TextSelection) Start() int              { return i.start }
func (i TextSelection) End() int                { return i.end }
func (i TextSelection) First() int              { return i.start }
func (i TextSelection) Last() int               { return i.end - 1 }
func (i TextSelection) CaretAtStart() bool      { return i.caretAtStart }

func (t TextSelection) Offset(i int) TextSelection {
	return TextSelection{
		start:        t.start + i,
		end:          t.end + i,
		caretAtStart: t.caretAtStart,
	}
}
func (i TextSelection) Caret() int {
	if i.caretAtStart {
		return i.start
	} else {
		return i.end
	}
}

func (i TextSelection) From() int { // TODO: Think of a better name for this function
	if i.caretAtStart {
		return i.end
	} else {
		return i.start
	}
}

func (i TextSelection) Span() (start, end uint64) {
	return uint64(i.start), uint64(i.end)
}
