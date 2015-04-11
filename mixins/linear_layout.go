// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/google/gxui"
	"github.com/google/gxui/mixins/base"
	"github.com/google/gxui/mixins/parts"
)

type LinearLayoutOuter interface {
	base.ContainerOuter
}

type LinearLayout struct {
	base.Container
	parts.LinearLayout
	parts.BackgroundBorderPainter
}

func (l *LinearLayout) Init(outer LinearLayoutOuter, theme gxui.Theme) {
	l.Container.Init(outer, theme)
	l.LinearLayout.Init(outer)
	l.BackgroundBorderPainter.Init(outer)
	l.SetMouseEventTarget(true)
	l.SetBackgroundBrush(gxui.TransparentBrush)
	l.SetBorderPen(gxui.TransparentPen)

	// Interface compliance test
	_ = gxui.LinearLayout(l)
}

func (l *LinearLayout) Paint(c gxui.Canvas) {
	r := l.Size().Rect()
	l.BackgroundBorderPainter.PaintBackground(c, r)
	l.PaintChildren.Paint(c)
	l.BackgroundBorderPainter.PaintBorder(c, r)
}
