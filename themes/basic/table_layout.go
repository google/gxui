package basic

import (
	"github.com/google/gxui"
	"github.com/google/gxui/mixins"
)

func CreateTableLayout(theme *Theme) gxui.TableLayout {
	l := &mixins.TableLayout{}
	l.Init(l, theme)
	return l
}
