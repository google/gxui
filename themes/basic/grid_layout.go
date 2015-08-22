package basic

import (
	"github.com/google/gxui"
	"github.com/google/gxui/mixins"
)

func CreateGridLayout(theme *Theme) gxui.GridLayout {
  l := &mixins.GridLayout{}
  l.Init(l, theme)
  return l
}
