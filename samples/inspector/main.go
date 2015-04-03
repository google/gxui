// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/math"
	"github.com/google/gxui/samples/inspector/property"
	"github.com/google/gxui/themes/dark"
)

type updater struct {
	inner  property.Property
	update func()
}

func (p updater) Name() string            { return p.inner.Name() }
func (p updater) Type() reflect.Type      { return p.inner.Type() }
func (p updater) Get() reflect.Value      { return p.inner.Get() }
func (p updater) Set(value reflect.Value) { p.inner.Set(value); p.update() }
func (p updater) CanSet() bool            { return p.inner.CanSet() }

// TODO: Add EnumValues() method on to the enum type?
func enumValues(value interface{}) []gxui.AdapterItem {
	switch value.(type) {
	case bool:
		return []gxui.AdapterItem{false, true}

	case gxui.ButtonType:
		return []gxui.AdapterItem{
			gxui.PushButton,
			gxui.ToggleButton,
		}

	case gxui.Direction:
		return []gxui.AdapterItem{
			gxui.TopToBottom,
			gxui.LeftToRight,
			gxui.BottomToTop,
			gxui.RightToLeft,
		}

	case gxui.HorizontalAlignment:
		return []gxui.AdapterItem{
			gxui.AlignLeft,
			gxui.AlignCenter,
			gxui.AlignRight,
		}

	case gxui.Orientation:
		return []gxui.AdapterItem{gxui.Vertical, gxui.Horizontal}

	case gxui.SizeMode:
		return []gxui.AdapterItem{gxui.ExpandToContent, gxui.Fill}

	case gxui.VerticalAlignment:
		return []gxui.AdapterItem{
			gxui.AlignTop,
			gxui.AlignMiddle,
			gxui.AlignBottom,
		}

	default:
		return nil
	}
}

func controlFor(theme gxui.Theme, overlay gxui.BubbleOverlay, property property.Property) gxui.Control {
	value := property.Get()

	if enumVals := enumValues(value.Interface()); enumVals != nil {
		if property.CanSet() {
			adapter := gxui.CreateDefaultAdapter()
			adapter.SetItems(enumVals)
			control := theme.CreateDropDownList()
			control.SetBubbleOverlay(overlay)
			control.SetAdapter(adapter)
			control.Select(value.Interface())
			control.OnSelectionChanged(func(item gxui.AdapterItem) {
				property.Set(reflect.ValueOf(item))
			})
			return control
		} else {
			control := theme.CreateLabel()
			control.SetText(fmt.Sprintf("%v", value.Interface()))
			return control
		}
	}

	ty := value.Type()
	switch ty.Kind() {
	case reflect.String:
		if property.CanSet() {
			control := theme.CreateTextBox()
			control.SetPadding(math.ZeroSpacing)
			control.SetText(value.Interface().(string))
			control.OnTextChanged(func([]gxui.TextBoxEdit) {
				property.Set(reflect.ValueOf(control.Text()))
			})
			return control
		} else {
			control := theme.CreateLabel()
			control.SetText(fmt.Sprintf("%v", value.Interface()))
			return control
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if property.CanSet() {
			control := theme.CreateTextBox()
			control.SetPadding(math.ZeroSpacing)
			control.SetText(fmt.Sprintf("%d", value.Interface()))
			control.OnTextChanged(func([]gxui.TextBoxEdit) {
				i, err := strconv.ParseInt(control.Text(), 10, ty.Bits())
				if err == nil {
					property.Set(reflect.ValueOf(i))
				}
			})
			return control
		} else {
			control := theme.CreateLabel()
			control.SetText(fmt.Sprintf("%v", value.Interface()))
			return control
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if property.CanSet() {
			control := theme.CreateTextBox()
			control.SetPadding(math.ZeroSpacing)
			control.SetText(fmt.Sprintf("%d", value.Interface()))
			control.OnTextChanged(func([]gxui.TextBoxEdit) {
				i, err := strconv.ParseUint(control.Text(), 10, ty.Bits())
				if err == nil {
					property.Set(reflect.ValueOf(i))
				}
			})
			return control
		} else {
			control := theme.CreateLabel()
			control.SetText(fmt.Sprintf("%v", value.Interface()))
			return control
		}

	case reflect.Float32, reflect.Float64:
		if property.CanSet() {
			control := theme.CreateTextBox()
			control.SetPadding(math.ZeroSpacing)
			control.SetText(fmt.Sprintf("%f", value.Interface()))
			control.OnTextChanged(func([]gxui.TextBoxEdit) {
				f, err := strconv.ParseFloat(control.Text(), ty.Bits())
				if err == nil {
					property.Set(reflect.ValueOf(f))
				}
			})
			return control
		} else {
			control := theme.CreateLabel()
			control.SetText(fmt.Sprintf("%v", value.Interface()))
			return control
		}

	default:
		return nil
	}
}

type node struct {
	gxui.AdapterBase
	path       string
	target     interface{}
	properties []property.Property
	overlay    gxui.BubbleOverlay
}

func (n *node) Count() int {
	return len(n.properties)
}

func (n *node) set(target reflect.Value) {
	n.target = target

	// Wrap each of the settable properties in an updater property so that sets
	// trigger updates of the parent.
	n.properties = property.Properties(target)
	for i := range n.properties {
		if n.properties[i].CanSet() {
			n.properties[i] = updater{
				inner: n.properties[i],
				update: func() {
					n.DataChanged()
					n.set(target)
				},
			}
		}
	}
}

func (n *node) NodeAt(index int) gxui.TreeNode {
	property := n.properties[index]
	value := property.Get()

	if property.CanSet() && !value.CanAddr() {
		// The property can be set, but the returned value is non-addressable so
		// cannot be mutated.
		// Clone the value into something that can be mutated.
		clone := reflect.New(value.Type()).Elem()
		clone.Set(value)
		value = clone
	}

	child := &node{
		path:    n.ItemAt(index).(string),
		overlay: n.overlay,
	}
	child.set(value)
	child.OnDataChanged(func() {
		if property.CanSet() {
			property.Set(value)
		}
		n.DataChanged()
	})

	return child
}

func (n *node) ItemAt(index int) gxui.AdapterItem {
	return n.path + "." + n.properties[index].Name()
}

func (n *node) ItemIndex(item gxui.AdapterItem) int {
	for i := range n.properties {
		if n.ItemAt(i) == item {
			return i
		}
	}
	return -1
}

func (n *node) Create(theme gxui.Theme, index int) gxui.Control {
	property := n.properties[index]

	name := theme.CreateLabel()
	name.SetText(property.Name())

	l2r := theme.CreateLinearLayout()
	l2r.SetSizeMode(gxui.Fill)
	l2r.SetDirection(gxui.LeftToRight)
	l2r.AddChild(name)

	r2l := theme.CreateLinearLayout()
	r2l.SetDirection(gxui.RightToLeft)
	if c := controlFor(theme, n.overlay, property); c != nil {
		r2l.AddChild(c)
	} else {
		ty := theme.CreateLabel()
		ty.SetFont(theme.DefaultMonospaceFont())
		ty.SetColor(gxui.Gray60)
		ty.SetText(property.Type().String())
		r2l.AddChild(ty)
	}
	r2l.AddChild(l2r)

	// Pad enough so that the scroll bar doesn't obstruct the controls
	r2l.SetPadding(math.Spacing{R: 10})

	return r2l
}

type adapter struct{ node }

func (a *adapter) Size(gxui.Theme) math.Size {
	return math.Size{W: math.MaxSize.W, H: 20}
}

func appMain(driver gxui.Driver) {
	theme := dark.CreateTheme(driver)

	window := theme.CreateWindow(800, 600, "Inspector")
	window.OnClose(driver.Terminate)

	tree := theme.CreateTree()

	layout := theme.CreateLinearLayout()
	layout.SetDirection(gxui.TopToBottom)

	button := theme.CreateButton()
	button.SetText("Button")
	layout.AddChild(button)

	label := theme.CreateLabel()
	label.SetText("Label")
	layout.AddChild(label)

	progressbar := theme.CreateProgressBar()
	progressbar.SetDesiredSize(math.Size{W: 100, H: 20})
	layout.AddChild(progressbar)

	splitter := theme.CreateSplitterLayout()
	splitter.SetOrientation(gxui.Horizontal)
	splitter.AddChild(tree)
	splitter.AddChild(layout)

	overlay := theme.CreateBubbleOverlay()

	window.AddChild(splitter)
	window.AddChild(overlay)

	tree.SetAdapter(&adapter{
		node: node{
			properties: property.Properties(reflect.ValueOf(window)),
			overlay:    overlay,
		},
	})
}

func main() {
	gl.StartDriver(appMain)
}
