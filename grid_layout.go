package gxui

type GridLayout interface {
  Control

  Container

  SetGrid(rows, columns int, control ...Control)
}
