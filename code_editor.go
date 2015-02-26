// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import "gxui/interval"

type CodeSuggestion interface {
	Name() string
	Code() string
}

type CodeSuggestionProvider interface {
	SuggestionsAt(runeIndex int) []CodeSuggestion
}

type CodeEditor interface {
	TextBox
	SetSyntaxLayer(index int, layer CodeSyntaxLayer)
	ClearSyntaxLayers()
	TabWidth() int
	SetTabWidth(int)
	SpanAt(layerIdx, runeIdx int) *interval.IntData
	SpansAt(runeIdx int) []interval.IntData
	SuggestionProvider() CodeSuggestionProvider
	SetSuggestionProvider(CodeSuggestionProvider)
	ShowSuggestionList()
	HideSuggestionList()
}
