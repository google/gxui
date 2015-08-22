package mixins

import (
  "github.com/google/gxui"
  "github.com/google/gxui/math"
  "github.com/google/gxui/mixins/base"
)

type GridLayoutOuter interface {
  base.ContainerOuter
}

type GridLayout struct {
  base.Container

  outer GridLayoutOuter

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
  for row := 0; row < l.rows; row++ {
    for column := 0; column < l.columns; column++ {
      c := children[row*l.columns+column]

      cm := c.Control.Margin()

      w, h := column*cw, row*ch

      cr = math.CreateRect(w+cm.L, h+cm.T, w+cw-cm.R, h+ch-cm.B)

      c.Layout(cr.Offset(o).Canon())
    }
  }
}

func (l *GridLayout) DesiredSize(min, max math.Size) math.Size {
	return max
}

func (l *GridLayout) SetGrid(columns, rows int, children ...gxui.Control) {
  l.columns = columns
  l.rows = rows

  for _, c := range children {
    l.Container.AddChild(c)
  }
}
