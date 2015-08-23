package gxui

type GridLayout interface {
  Control

  Container

  SetGrid(rows, columns int)
  SetChildAt(x, y, w, h int, child Control) *Child
}
