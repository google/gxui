package mixins

import (
	"github.com/google/gxui"
	"github.com/google/gxui/math"
	"github.com/google/gxui/mixins/base"
)

type Cell struct {
	x, y, w, h int
}

func (c Cell) AtColumn(x int) bool {
	return c.x <= x && c.x+c.w >= x
}

func (c Cell) AtRow(y int) bool {
	return c.y <= y && c.y+c.h >= y
}

type TableLayoutOuter interface {
	base.ContainerOuter
}

type TableLayout struct {
	base.Container

	outer TableLayoutOuter

	grid    map[gxui.Control]Cell
	rows    int
	columns int
}

func (l *TableLayout) Init(outer TableLayoutOuter, theme gxui.Theme) {
	l.Container.Init(outer, theme)
	l.outer = outer
	l.grid = make(map[gxui.Control]Cell)

	// Interface compliance test
	_ = gxui.TableLayout(l)
}

func (l *TableLayout) LayoutChildren() {
	s := l.outer.Size().Contract(l.outer.Padding())
	o := l.outer.Padding().LT()

	cw, ch := s.W/l.columns, s.H/l.rows

	var cr math.Rect

	for _, c := range l.outer.Children() {
		cm := c.Control.Margin()
		cell := l.grid[c.Control]

		x, y := cell.x*cw, cell.y*ch
		w, h := x+cell.w*cw, y+cell.h*ch

		cr = math.CreateRect(x+cm.L, y+cm.T, w-cm.R, h-cm.B)

		c.Layout(cr.Offset(o).Canon())
	}
}

func (l *TableLayout) DesiredSize(min, max math.Size) math.Size {
	return max
}

func (l *TableLayout) SetGrid(columns, rows int) {
	if l.columns != columns {
		if l.columns > columns {
			for c := l.columns; c > columns; c-- {
				for _, cell := range l.grid {
					if cell.AtColumn(c) {
						panic("Can't remove column with cells")
					}
				}
				l.columns--
			}
		} else {
			l.columns = columns
		}
	}

	if l.rows != rows {
		if l.rows > rows {
			for r := l.rows; r > rows; r-- {
				for _, cell := range l.grid {
					if cell.AtRow(r) {
						panic("Can't remove row with cells")
					}
				}
				l.rows--
			}
		} else {
			l.rows = rows
		}
	}

	if l.rows != rows || l.columns != columns {
		l.LayoutChildren()
	}
}

func (l *TableLayout) SetChildAt(x, y, w, h int, child gxui.Control) *gxui.Child {
	if x+w > l.columns || y+h > l.rows {
		panic("Cell is out of grid")
	}

	for _, c := range l.grid {
		if c.x+c.w > x && c.x < x+w && c.y+c.h > y && c.y < y+h {
			panic("Cell already has a child")
		}
	}

	l.grid[child] = Cell{x, y, w, h}
	return l.Container.AddChild(child)
}

func (l *TableLayout) RemoveChild(child gxui.Control) {
	delete(l.grid, child)
	l.Container.RemoveChild(child)
}
