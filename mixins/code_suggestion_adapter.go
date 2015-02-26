// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"gxui"
)

type SuggestionAdapter struct {
	gxui.FilteredListAdapter
}

func (a *SuggestionAdapter) SetSuggestions(suggestions []gxui.CodeSuggestion) {
	items := make([]gxui.FilteredListItem, len(suggestions))
	for i, s := range suggestions {
		items[i].Name = s.Name()
		items[i].Data = s
	}
	a.SetItems(items)
}

func (a *SuggestionAdapter) Suggestion(id gxui.AdapterItemId) gxui.CodeSuggestion {
	return a.Item(id).Data.(gxui.CodeSuggestion)
}
