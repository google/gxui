package mixins

import (
  "github.com/google/gxui"
  "github.com/google/gxui/math"
  "github.com/google/gxui/mixins/base"
)

type Cell struct {
  x, y, w, h int
}

type GridLayoutOuter interface {
  base.ContainerOuter
}

type GridLayout struct {
  base.Container

  outer GridLayoutOuter

  grid []Cell
  rows int
  columns int
}

func (l *GridLayout) Init(outer GridLayoutOuter, theme gxui.Theme) {
  l.Container.Init(outer, theme)
  l.outer = outer

  _ = gxui.GridLayout(l)
}

func (l *GridLayout) LayoutChildren() {
  s := l.outer.Size().Contract(l.outer.Padding())
	o := l.outer.Padding().LT()

  children := l.outer.Children()

  cw, ch := s.W / l.columns, s.H / l.rows

  var cr math.Rect

  for i, cell := range l.grid {
    c := children[i]
    cm := c.Control.Margin()

    x, y := cell.x * cw, cell.y * ch
    w, h := x + cell.w * cw, y + cell.h * ch

    cr = math.CreateRect(x+cm.L, y+cm.T, w-cm.R, h-cm.B)

    c.Layout(cr.Offset(o).Canon())
  }
}

func (l *GridLayout) DesiredSize(min, max math.Size) math.Size {
	return max
}

func (l *GridLayout) SetGrid(columns, rows int) {
  l.columns = columns
  l.rows = rows
}

func (l *GridLayout) SetChildAt(x, y, w, h int, child gxui.Control) *gxui.Child {
  l.grid = append(l.grid, Cell{x, y, w, h})
  return l.Container.AddChild(child)
}

func (l *GridLayout) RemoveChild(child gxui.Control) {
  for i, c := range l.Container.Children() {
		if c.Control == child {
      l.grid = append(l.grid[:i], l.grid[i+1:]...)
			l.Container.RemoveChildAt(i)
      return
		}
	}
  panic("Child not part of container")
}
