// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/google/gxui/math"
	"time"
)

type ToolTipCreator func(math.Point) Control

type toolTipTracker struct {
	creator      ToolTipCreator
	control      Control
	onEnterES    EventSubscription
	onExitES     EventSubscription
	onMoveES     EventSubscription
	lastPosition math.Point
}

type ToolTipController struct {
	driver        Driver
	timer         *time.Timer
	bubbleOverlay BubbleOverlay
	trackers      []*toolTipTracker
	showing       *toolTipTracker
}

func (c *ToolTipController) beginTimer(tracker *toolTipTracker, timeout time.Duration) {
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
	if timeout > 0 {
		c.timer = time.AfterFunc(timeout, func() {
			c.driver.Call(func() { c.showToolTipForTracker(tracker) })
		})
	} else {
		c.showToolTipForTracker(tracker)
	}
}

func (c *ToolTipController) showToolTipForTracker(tracker *toolTipTracker) {
	toolTip := tracker.creator(tracker.lastPosition)
	if toolTip != nil {
		at := TransformCoordinate(tracker.lastPosition, tracker.control, c.bubbleOverlay)
		c.ShowToolTip(toolTip, at)
		c.showing = tracker
	} else {
		c.hideToolTipForTracker(tracker)
	}
}

func (c *ToolTipController) hideToolTipForTracker(tracker *toolTipTracker) {
	if c.showing == tracker {
		c.bubbleOverlay.Hide()
		c.showing = nil
	}
}

func CreateToolTipController(bubbleOverlay BubbleOverlay, driver Driver) *ToolTipController {
	return &ToolTipController{
		driver:        driver,
		bubbleOverlay: bubbleOverlay,
	}
}

func (c *ToolTipController) AddToolTip(control Control, delaySeconds float32, creator ToolTipCreator) {
	tracker := &toolTipTracker{
		control: control,
		creator: creator,
	}
	duration := time.Duration(delaySeconds * float32(time.Second))
	bind := func() {
		tracker.onEnterES = control.OnMouseEnter(func(ev MouseEvent) {
			tracker.lastPosition = ev.Point
			c.beginTimer(tracker, duration)
		})
		tracker.onExitES = control.OnMouseExit(func(ev MouseEvent) {
			if c.timer != nil {
				c.timer.Stop()
				c.timer = nil
			}
			c.hideToolTipForTracker(tracker)
		})
		tracker.onMoveES = control.OnMouseMove(func(ev MouseEvent) {
			tracker.lastPosition = ev.Point
			c.beginTimer(tracker, duration)
		})
	}
	control.OnAttach(bind)
	control.OnDetach(func() {
		if c.timer != nil {
			c.timer.Stop()
			c.timer = nil
		}
		tracker.onEnterES.Unlisten()
		tracker.onExitES.Unlisten()
		tracker.onMoveES.Unlisten()
	})
	if control.Attached() {
		bind()
	}
}

func (c *ToolTipController) ShowToolTip(toolTip Control, at math.Point) {
	c.bubbleOverlay.Show(toolTip, at)
}
