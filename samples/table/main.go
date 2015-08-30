package main

import (
	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/samples/flags"
)

func appMain(driver gxui.Driver) {
	theme := flags.CreateTheme(driver)

	label1 := theme.CreateLabel()
	label1.SetColor(gxui.White)
	label1.SetText("1x1")

	cell1x1 := theme.CreateLinearLayout()
	cell1x1.SetBackgroundBrush(gxui.CreateBrush(gxui.Blue40))
	cell1x1.SetHorizontalAlignment(gxui.AlignCenter)
	cell1x1.AddChild(label1)

	label2 := theme.CreateLabel()
	label2.SetColor(gxui.White)
	label2.SetText("2x1")

	cell2x1 := theme.CreateLinearLayout()
	cell2x1.SetBackgroundBrush(gxui.CreateBrush(gxui.Green40))
	cell2x1.SetHorizontalAlignment(gxui.AlignCenter)
	cell2x1.AddChild(label2)

	label3 := theme.CreateLabel()
	label3.SetColor(gxui.White)
	label3.SetText("1x2")

	cell1x2 := theme.CreateLinearLayout()
	cell1x2.SetBackgroundBrush(gxui.CreateBrush(gxui.Red40))
	cell1x2.SetHorizontalAlignment(gxui.AlignCenter)
	cell1x2.AddChild(label3)

	table := theme.CreateTableLayout()
	table.SetGrid(3, 2) // rows, columns

	// row, column, horizontal span, vertical span
	table.SetChildAt(0, 0, 1, 1, cell1x1)
	table.SetChildAt(0, 1, 2, 1, cell2x1)
	table.SetChildAt(2, 0, 1, 2, cell1x2)

	window := theme.CreateWindow(800, 600, "Table")
	window.AddChild(table)
	window.OnClose(driver.Terminate)
}

func main() {
	gl.StartDriver(appMain)
}
