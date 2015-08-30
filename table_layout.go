package gxui

type TableLayout interface {
	Control

	Parent

	SetGrid(rows, columns int)
	// Add child at cell {x, y} with size of {w, h}
	SetChildAt(x, y, w, h int, child Control) *Child
	RemoveChild(child Control)
}
