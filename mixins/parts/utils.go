// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parts

import (
	"github.com/google/gxui/mixins/outer"
)

func callLayoutChildrenIfSupported(i interface{}) {
	switch ty := i.(type) {
	case outer.LayoutChildren:
		ty.LayoutChildren()
	}
}
