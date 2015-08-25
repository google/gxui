// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"fmt"
	"reflect"

	"github.com/google/gxui/math"
)

type Viewer interface {
	View(t Theme) Control
}

type Stringer interface {
	String() string
}

type DefaultAdapter struct {
	AdapterBase
	items       reflect.Value
	itemToIndex map[AdapterItem]int
	size        math.Size
	styleLabel  func(Theme, Label)
}

func CreateDefaultAdapter() *DefaultAdapter {
	l := &DefaultAdapter{
		size: math.Size{W: 200, H: 16},
	}
	return l
}

func (a *DefaultAdapter) SetSizeAsLargest(theme Theme) {
	s := math.Size{}
	font := theme.DefaultFont()
	for i := 0; i < a.Count(); i++ {
		switch t := a.ItemAt(i).(type) {
		case Viewer:
			s = s.Max(t.View(theme).DesiredSize(math.ZeroSize, math.MaxSize))

		case Stringer:
			s = s.Max(font.Measure(&TextBlock{
				Runes: []rune(t.String()),
			}))

		default:
			s = s.Max(font.Measure(&TextBlock{
				Runes: []rune(fmt.Sprintf("%+v", t)),
			}))
		}
	}
	a.SetSize(s)
}

func (a *DefaultAdapter) SetStyleLabel(f func(Theme, Label)) {
	a.styleLabel = f
	a.DataChanged(true)
}

func (a *DefaultAdapter) Count() int {
	if !a.items.IsValid() {
		return 0
	}

	switch a.items.Kind() {
	case reflect.Slice, reflect.Array:
		return a.items.Len()

	default:
		return 1
	}
}

func (a *DefaultAdapter) ItemAt(index int) AdapterItem {
	count := a.Count()
	if index < 0 || index >= count {
		panic(fmt.Errorf("ItemAt index %d is out of bounds [%d, %d]",
			index, 0, count-1))
	}

	switch a.items.Kind() {
	case reflect.Slice, reflect.Array:
		return a.items.Index(index).Interface()

	default:
		return a.items.Interface()
	}

}

func (a *DefaultAdapter) ItemIndex(item AdapterItem) int {
	return a.itemToIndex[item]
}

func (a *DefaultAdapter) Size(theme Theme) math.Size {
	return a.size
}

func (a *DefaultAdapter) SetSize(s math.Size) {
	a.size = s
	a.DataChanged(true)
}

func (a *DefaultAdapter) Create(theme Theme, index int) Control {
	switch t := a.ItemAt(index).(type) {
	case Viewer:
		return t.View(theme)

	case Stringer:
		l := theme.CreateLabel()
		l.SetMargin(math.ZeroSpacing)
		l.SetMultiline(false)
		l.SetText(t.String())
		if a.styleLabel != nil {
			a.styleLabel(theme, l)
		}
		return l

	default:
		l := theme.CreateLabel()
		l.SetMargin(math.ZeroSpacing)
		l.SetMultiline(false)
		l.SetText(fmt.Sprintf("%+v", t))
		if a.styleLabel != nil {
			a.styleLabel(theme, l)
		}
		return l
	}
}

func (a *DefaultAdapter) Items() interface{} {
	return a.items.Interface()
}

func (a *DefaultAdapter) SetItems(items interface{}) {
	a.items = reflect.ValueOf(items)
	a.itemToIndex = make(map[AdapterItem]int)
	for idx := 0; idx < a.Count(); idx++ {
		a.itemToIndex[a.ItemAt(idx)] = idx
	}
	a.DataReplaced()
}
