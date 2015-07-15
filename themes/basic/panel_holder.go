// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package basic

import (
	"github.com/google/gxui"
	"github.com/google/gxui/math"
	"github.com/google/gxui/mixins"
)

type PanelHolder struct {
	mixins.PanelHolder
	theme *Theme
}

func CreatePanelHolder(theme *Theme) gxui.PanelHolder {
	p := &PanelHolder{}
	p.PanelHolder.Init(p, theme)
	p.theme = theme
	p.SetMargin(math.Spacing{L: 0, T: 2, R: 0, B: 0})
	return p
}

func (p *PanelHolder) CreatePanelTab() mixins.PanelTab {
	return CreatePanelTab(p.theme)
}

func (p *PanelHolder) Paint(c gxui.Canvas) {
	panel := p.SelectedPanel()
	if panel != nil {
		bounds := p.Children().Find(panel).Bounds()
		c.DrawRoundedRect(bounds, 0.0, 0.0, 3.0, 3.0, p.theme.PanelBackgroundStyle.Pen, p.theme.PanelBackgroundStyle.Brush)
	}
	p.PanelHolder.Paint(c)
}
