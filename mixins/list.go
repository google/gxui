// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"gxui"
	"gxui/assert"
	"gxui/math"
	"gxui/mixins/base"
	"gxui/mixins/parts"
)

type ListOuter interface {
	base.ContainerOuter
	PaintBackground(c gxui.Canvas, r math.Rect)
	PaintMouseOverBackground(c gxui.Canvas, r math.Rect)
	PaintSelection(c gxui.Canvas, r math.Rect)
	PaintBorder(c gxui.Canvas, r math.Rect)
}

type ListItem struct {
	Control             gxui.Control
	Index               int
	Mark                int
	OnClickSubscription gxui.EventSubscription
}

type List struct {
	base.Container
	parts.BackgroundBorderPainter
	parts.Focusable

	outer ListOuter

	theme                    gxui.Theme
	adapter                  gxui.Adapter
	scrollBar                gxui.ScrollBar
	scrollBarEnabled         bool
	selectedId               gxui.AdapterItemId
	onSelectionChanged       gxui.Event
	items                    map[gxui.AdapterItemId]ListItem
	orientation              gxui.Orientation
	scrollOffset             int
	itemSize                 math.Size
	itemCount                int
	layoutMark               int
	mousePosition            math.Point
	itemMouseOver            gxui.Control
	onItemClicked            gxui.Event
	dataChangedSubscription  gxui.EventSubscription
	dataReplacedSubscription gxui.EventSubscription
}

func (l *List) Init(outer ListOuter, theme gxui.Theme) {
	l.outer = outer
	l.Container.Init(outer, theme)
	l.BackgroundBorderPainter.Init(outer)
	l.Focusable.Init(outer)

	l.theme = theme
	l.scrollBar = theme.CreateScrollBar()
	l.scrollBarEnabled = true
	l.scrollBar.OnScroll(func(from, to int) { l.SetScrollOffset(from) })
	l.selectedId = gxui.InvalidAdapterItemId
	l.SetOrientation(gxui.Vertical)
	l.SetBackgroundBrush(gxui.TransparentBrush)
	l.SetMouseEventTarget(true)

	l.items = make(map[gxui.AdapterItemId]ListItem)

	// Interface compliance test
	_ = gxui.List(l)
}

func (l *List) UpdateItemMouseOver() {
	if !l.IsMouseOver() {
		if l.itemMouseOver != nil {
			l.itemMouseOver = nil
			l.Redraw()
		}
		return
	}
	for _, item := range l.items {
		if item.Control.Bounds().Contains(l.mousePosition) {
			if l.itemMouseOver != item.Control {
				l.itemMouseOver = item.Control
				l.Redraw()
				return
			}
		}
	}
}

func (l *List) LayoutChildren() {
	if l.adapter == nil {
		l.outer.RemoveAll()
		return
	}

	s := l.outer.Bounds().Size().Contract(l.Padding())
	o := l.Padding().LT()

	var itemSize math.Size
	if l.orientation.Horizontal() {
		itemSize = math.Size{W: l.itemSize.W, H: s.H}
	} else {
		itemSize = math.Size{W: s.W, H: l.itemSize.H}
	}

	startIndex, endIndex := l.VisibleItemRange(true)
	majorAxisItemSize := l.MajorAxisItemSize()

	d := startIndex*majorAxisItemSize - l.scrollOffset

	mark := l.layoutMark
	l.layoutMark++

	for idx := startIndex; idx < endIndex; idx++ {
		id := l.adapter.ItemId(idx)

		item, found := l.items[id]
		if found {
			assert.False(item.Mark == mark, "Adapter returned duplicate id (%v) for indices %v and %v",
				id, item.Index, idx)
		} else {
			item.Control = l.adapter.Create(l.theme, idx)
			item.OnClickSubscription = item.Control.OnClick(func(ev gxui.MouseEvent) {
				l.ItemClicked(ev, id)
			})
			l.AddChildAt(0, item.Control)
		}
		item.Mark = mark
		item.Index = idx
		l.items[id] = item

		c := item.Control
		cm := c.Margin()
		cs := itemSize.Contract(cm).Max(math.ZeroSize)
		if l.orientation.Horizontal() {
			c.Layout(math.CreateRect(d, cm.T, d+cs.W, cm.T+cs.H).Offset(o))
		} else {
			c.Layout(math.CreateRect(cm.L, d, cm.L+cs.W, d+cs.H).Offset(o))
		}
		d += majorAxisItemSize
	}

	// Reap unused items
	for id, item := range l.items {
		if item.Mark != mark {
			item.OnClickSubscription.Unlisten()
			l.RemoveChild(item.Control)
			delete(l.items, id)
		}
	}

	if l.scrollBarEnabled {
		ss := l.scrollBar.DesiredSize(math.ZeroSize, s)
		if l.Orientation().Horizontal() {
			l.scrollBar.Layout(math.CreateRect(0, s.H-ss.H, s.W, s.H).Canon().Offset(o))
		} else {
			l.scrollBar.Layout(math.CreateRect(s.W-ss.W, 0, s.W, s.H).Canon().Offset(o))
		}

		// Only show the scroll bar if needed
		entireContentVisible := startIndex == 0 && endIndex == l.itemCount
		l.scrollBar.SetVisible(!entireContentVisible)
		if l.scrollBar.Parent() != l.outer {
			l.AddChild(l.scrollBar)
		}
	}

	l.UpdateItemMouseOver()
}

func (l *List) Layout(rect math.Rect) {
	l.Layoutable.Layout(rect)
	// Ensure scroll offset is still valid
	l.SetScrollOffset(l.scrollOffset)
}

func (l *List) DesiredSize(min, max math.Size) math.Size {
	if l.adapter == nil {
		return min
	}
	count := math.Max(l.itemCount, 1)
	var s math.Size
	if l.orientation.Horizontal() {
		s = math.Size{W: l.itemSize.W * count, H: l.itemSize.H}
	} else {
		s = math.Size{W: l.itemSize.W, H: l.itemSize.H * count}
	}
	if l.scrollBarEnabled {
		if l.orientation.Horizontal() {
			s.H += l.scrollBar.DesiredSize(min, max).H
		} else {
			s.W += l.scrollBar.DesiredSize(min, max).W
		}
	}
	return s.Expand(l.outer.Padding()).Clamp(min, max)
}

func (l *List) ScrollBarEnabled(bool) bool {
	return l.scrollBarEnabled
}

func (l *List) SetScrollBarEnabled(enabled bool) {
	if l.scrollBarEnabled != enabled {
		l.scrollBarEnabled = enabled
		l.Relayout()
	}
}

func (l *List) SetScrollOffset(scrollOffset int) {
	if l.adapter == nil {
		return
	}
	b := l.outer.Bounds().Contract(l.outer.Padding())
	if l.orientation.Horizontal() {
		maxScroll := math.Max(l.itemSize.W*l.itemCount-b.W(), 0)
		scrollOffset = math.Clamp(scrollOffset, 0, maxScroll)
		l.scrollBar.SetScrollPosition(scrollOffset, scrollOffset+b.W())
	} else {
		maxScroll := math.Max(l.itemSize.H*l.itemCount-b.H(), 0)
		scrollOffset = math.Clamp(scrollOffset, 0, maxScroll)
		l.scrollBar.SetScrollPosition(scrollOffset, scrollOffset+b.H())
	}
	if l.scrollOffset != scrollOffset {
		l.scrollOffset = scrollOffset
		l.LayoutChildren()
	}
}

func (l *List) MajorAxisItemSize() int {
	return l.orientation.Major(l.itemSize.WH())
}

func (l *List) VisibleItemRange(includePartiallyVisible bool) (startIndex, endIndex int) {
	if l.itemCount == 0 {
		return 0, 0
	}
	s := l.outer.Bounds().Size()
	p := l.outer.Padding()
	majorAxisItemSize := l.MajorAxisItemSize()
	startIndex = l.scrollOffset
	if !includePartiallyVisible {
		startIndex += majorAxisItemSize - 1
	}
	if l.orientation.Horizontal() {
		endIndex = l.scrollOffset + s.W - p.W()
	} else {
		endIndex = l.scrollOffset + s.H - p.H()
	}
	if includePartiallyVisible {
		endIndex += majorAxisItemSize - 1
	}
	startIndex = math.Max(startIndex/majorAxisItemSize, 0)
	endIndex = math.Min(endIndex/majorAxisItemSize, l.itemCount)

	return startIndex, endIndex
}

func (l *List) ItemSizeChanged() {
	l.itemSize = l.adapter.ItemSize(l.theme)
	l.scrollBar.SetScrollLimit(l.itemCount * l.MajorAxisItemSize())
	l.SetScrollOffset(l.scrollOffset)
	l.outer.Relayout()
}

func (l *List) DataChanged() {
	l.itemCount = l.adapter.Count()
	l.ItemSizeChanged()
}

func (l *List) DataReplaced() {
	l.selectedId = gxui.InvalidAdapterItemId
	for id, item := range l.items {
		item.OnClickSubscription.Unlisten()
		l.RemoveChild(item.Control)
		delete(l.items, id)
	}
	l.DataChanged()
}

func (l *List) Paint(c gxui.Canvas) {
	r := l.outer.Bounds().Size().Rect()
	l.outer.PaintBackground(c, r)
	l.Container.Paint(c)
	l.outer.PaintBorder(c, r)
}

func (l *List) PaintSelection(c gxui.Canvas, r math.Rect) {
	c.DrawRoundedRect(r, 2.0, 2.0, 2.0, 2.0, gxui.WhitePen, gxui.TransparentBrush)
}

func (l *List) PaintMouseOverBackground(c gxui.Canvas, r math.Rect) {
	c.DrawRoundedRect(r, 2.0, 2.0, 2.0, 2.0, gxui.TransparentPen, gxui.CreateBrush(gxui.Gray90))
}

func (l *List) SelectPrevious() {
	if l.selectedId != gxui.InvalidAdapterItemId {
		selectedIndex := l.adapter.ItemIndex(l.selectedId)
		l.Select(l.adapter.ItemId(math.Mod(selectedIndex-1, l.itemCount)))
	} else {
		l.Select(l.adapter.ItemId(0))
	}
}

func (l *List) SelectNext() {
	if l.selectedId != gxui.InvalidAdapterItemId {
		selectedIndex := l.adapter.ItemIndex(l.selectedId)
		l.Select(l.adapter.ItemId(math.Mod(selectedIndex+1, l.itemCount)))
	} else {
		l.Select(l.adapter.ItemId(0))
	}
}

// PaintChildren overrides
func (l *List) PaintChild(c gxui.Canvas, child gxui.Control, idx int) {
	if child == l.itemMouseOver {
		b := child.Bounds().Expand(child.Margin())
		l.outer.PaintMouseOverBackground(c, b)
	}
	l.PaintChildren.PaintChild(c, child, idx)
	if selected, found := l.items[l.selectedId]; found {
		if child == selected.Control {
			b := child.Bounds().Expand(child.Margin())
			l.outer.PaintSelection(c, b)
		}
	}
}

// InputEventHandler override
func (l *List) MouseMove(ev gxui.MouseEvent) {
	l.InputEventHandler.MouseMove(ev)
	l.mousePosition = ev.Point
	l.UpdateItemMouseOver()
}

func (l *List) MouseExit(ev gxui.MouseEvent) {
	l.InputEventHandler.MouseExit(ev)
	l.itemMouseOver = nil
}

func (l *List) MouseScroll(ev gxui.MouseEvent) (consume bool) {
	if ev.ScrollY == 0 {
		return l.InputEventHandler.MouseScroll(ev)
	}
	prevOffset := l.scrollOffset
	if l.orientation.Horizontal() {
		delta := ev.ScrollY * l.itemSize.W / 8
		l.SetScrollOffset(l.scrollOffset - delta)
	} else {
		delta := ev.ScrollY * l.itemSize.H / 8
		l.SetScrollOffset(l.scrollOffset - delta)
	}
	return prevOffset != l.scrollOffset
}

func (l *List) KeyPress(ev gxui.KeyboardEvent) (consume bool) {
	if l.itemCount > 0 {
		if l.orientation.Horizontal() {
			switch ev.Key {
			case gxui.KeyLeft:
				l.SelectPrevious()
				return true
			case gxui.KeyRight:
				l.SelectNext()
				return true
			case gxui.KeyPageUp:
				l.SetScrollOffset(l.scrollOffset - l.Bounds().W())
				return true
			case gxui.KeyPageDown:
				l.SetScrollOffset(l.scrollOffset + l.Bounds().W())
				return true
			}
		} else {
			switch ev.Key {
			case gxui.KeyUp:
				l.SelectPrevious()
				return true
			case gxui.KeyDown:
				l.SelectNext()
				return true
			case gxui.KeyPageUp:
				l.SetScrollOffset(l.scrollOffset - l.Bounds().H())
				return true
			case gxui.KeyPageDown:
				l.SetScrollOffset(l.scrollOffset + l.Bounds().H())
				return true
			}
		}
	}
	return l.Container.KeyPress(ev)
}

// gxui.List compliance
func (l *List) Adapter() gxui.Adapter {
	return l.adapter
}

func (l *List) SetAdapter(adapter gxui.Adapter) {
	if l.adapter != adapter {
		if l.adapter != nil {
			l.dataChangedSubscription.Unlisten()
			l.dataReplacedSubscription.Unlisten()
		}
		l.adapter = adapter
		if l.adapter != nil {
			l.dataChangedSubscription = l.adapter.OnDataChanged(l.DataChanged)
			l.dataReplacedSubscription = l.adapter.OnDataReplaced(l.DataReplaced)
		}
		l.DataReplaced()
	}
}

func (l *List) Orientation() gxui.Orientation {
	return l.orientation
}

func (l *List) SetOrientation(o gxui.Orientation) {
	l.scrollBar.SetOrientation(o)
	if l.orientation != o {
		l.orientation = o
		l.Relayout()
	}
}

func (l *List) ScrollTo(id gxui.AdapterItemId) {
	idx := l.adapter.ItemIndex(id)
	startIndex, endIndex := l.VisibleItemRange(false)
	if idx < startIndex {
		if l.Orientation().Horizontal() {
			l.SetScrollOffset(l.itemSize.W * idx)
		} else {
			l.SetScrollOffset(l.itemSize.H * idx)
		}
	} else if idx >= endIndex {
		count := endIndex - startIndex
		if l.Orientation().Horizontal() {
			l.SetScrollOffset(l.itemSize.W * (idx - count + 1))
		} else {
			l.SetScrollOffset(l.itemSize.H * (idx - count + 1))
		}
	}
}

func (l *List) IsItemVisible(id gxui.AdapterItemId) bool {
	_, found := l.items[id]
	return found
}

func (l *List) Item(id gxui.AdapterItemId) gxui.Control {
	if item, found := l.items[id]; found {
		return item.Control
	}
	return nil
}

func (l *List) ItemClicked(ev gxui.MouseEvent, id gxui.AdapterItemId) {
	if l.onItemClicked != nil {
		l.onItemClicked.Fire(ev, id)
	}
	l.Select(id)
}

func (l *List) OnItemClicked(f func(gxui.MouseEvent, gxui.AdapterItemId)) gxui.EventSubscription {
	if l.onItemClicked == nil {
		l.onItemClicked = gxui.CreateEvent(f)
	}
	return l.onItemClicked.Listen(f)
}

func (l *List) Selected() gxui.AdapterItemId {
	return l.selectedId
}

func (l *List) Select(id gxui.AdapterItemId) {
	if l.selectedId != id {
		l.selectedId = id
		if l.onSelectionChanged != nil {
			l.onSelectionChanged.Fire(id)
		}
		l.Redraw()
	}
	l.ScrollTo(id)
}

func (l *List) OnSelectionChanged(f func(gxui.AdapterItemId)) gxui.EventSubscription {
	if l.onItemClicked == nil {
		l.onSelectionChanged = gxui.CreateEvent(f)
	}
	return l.onSelectionChanged.Listen(f)
}
