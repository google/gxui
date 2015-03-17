// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"fmt"
	"github.com/google/gxui/math"
	"reflect"
)

type Viewer interface {
	View(t Theme) Control
}

type Stringer interface {
	String() string
}

type DefaultAdapter struct {
	AdapterBase
	data       reflect.Value
	itemSize   math.Size
	styleLabel func(Theme, Label)
}

func CreateDefaultAdapter() *DefaultAdapter {
	l := &DefaultAdapter{
		itemSize: math.Size{W: 200, H: 16},
	}
	return l
}

func (a *DefaultAdapter) SetItemSizeAsLargest(theme Theme) {
	s := math.Size{}
	font := theme.DefaultFont()
	for i := 0; i < a.Count(); i++ {
		e := a.data.Index(i).Interface()
		switch t := e.(type) {
		case Viewer:
			s = s.Max(t.View(theme).DesiredSize(math.ZeroSize, math.MaxSize))
		case Stringer:
			s = s.Max(font.Measure(t.String()))
		default:
			str := fmt.Sprintf("%+v", e)
			s = s.Max(font.Measure(str))
		}
	}
	a.SetItemSize(s)
}

func (a *DefaultAdapter) SetItemSize(s math.Size) {
	a.itemSize = s
	a.DataChanged()
}

func (a *DefaultAdapter) SetStyleLabel(f func(Theme, Label)) {
	a.styleLabel = f
	a.DataChanged()
}

func (a *DefaultAdapter) Count() int {
	if a.data.IsValid() {
		return a.data.Len()
	} else {
		return 0
	}
}

func (a *DefaultAdapter) IdOf(data interface{}) AdapterItemId {
	for i := 0; i < a.Count(); i++ {
		e := a.data.Index(i).Interface()
		if e == data {
			return a.ItemId(i)
		}
	}
	return InvalidAdapterItemId
}

func (a *DefaultAdapter) ValueOf(id AdapterItemId) interface{} {
	index := a.ItemIndex(id)
	return a.data.Index(index).Interface()
}

func (a *DefaultAdapter) ItemId(index int) AdapterItemId {
	return AdapterItemId(index)
}

func (a *DefaultAdapter) ItemIndex(id AdapterItemId) int {
	return int(id)
}

func (a *DefaultAdapter) ItemSize(theme Theme) math.Size {
	return a.itemSize
}

func (a *DefaultAdapter) Create(theme Theme, index int) Control {
	e := a.data.Index(index).Interface()
	switch t := e.(type) {
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
		l.SetText(fmt.Sprintf("%+v", e))
		if a.styleLabel != nil {
			a.styleLabel(theme, l)
		}
		return l
	}
}

func (a *DefaultAdapter) Data() interface{} {
	return a.data.Interface()
}

func (a *DefaultAdapter) SetData(data interface{}) {
	a.data = reflect.ValueOf(data)
	a.DataReplaced()
}
