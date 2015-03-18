// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"fmt"
	test "github.com/google/gxui/testing"
	"testing"
)

func parseTBC(markup string) *TextBoxController {
	tbc := CreateTextBoxController()
	tbc.selections = TextSelectionList{}
	runes := make([]rune, 0, 32)
	sel := TextSelection{}
	for _, c := range markup {
		i := len(runes)
		switch c {
		case '|':
			tbc.AddCaret(i)
		case '{':
			sel.start = i
			sel.caretAtStart = false
		case '[':
			sel.start = i
			sel.caretAtStart = true
		case ']':
			sel.end = i
			if sel.CaretAtStart() {
				panic("Carat should be at end")
			}
			tbc.AddSelection(sel)
		case '}':
			sel.end = i
			if !sel.CaretAtStart() {
				panic("Carat should be at start")
			}
			tbc.AddSelection(sel)
		default:
			runes = append(runes, c)
		}
	}
	tbc.SetTextRunes(runes)
	return tbc
}

func assertTBCTextAndSelectionsEqual(t *testing.T, markup string, c *TextBoxController) {
	expected := parseTBC(markup)
	test.AssertEquals(t, expected.Text(), c.Text())
	test.AssertEquals(t, expected.selections, c.selections)
}

func TestTBCLineIndent(t *testing.T) {
	c := parseTBC("  ÀÁ\n    BB\nĆ\n      D\n   EE")
	test.AssertEquals(t, 2, c.LineIndent(0))
	test.AssertEquals(t, 4, c.LineIndent(1))
	test.AssertEquals(t, 0, c.LineIndent(2))
	test.AssertEquals(t, 6, c.LineIndent(3))
	test.AssertEquals(t, 3, c.LineIndent(4))
}

func TestParseTBCCarets(t *testing.T) {
	c := parseTBC("he|llo\n|wor|ld")
	test.AssertEquals(t, "hello\nworld", c.Text())
	test.AssertEquals(t, 2, c.LineCount())
	test.AssertEquals(t, 3, c.SelectionCount())
	test.AssertEquals(t, 2, c.Caret(0))
	test.AssertEquals(t, 6, c.Caret(1))
	test.AssertEquals(t, 9, c.Caret(2))
}

func TestParseTBCSelections(t *testing.T) {
	c := parseTBC("he[llo}\n|wor{ld]")
	test.AssertEquals(t, "hello\nworld", c.Text())
	test.AssertEquals(t, 2, c.LineCount())
	test.AssertEquals(t, 3, c.SelectionCount())
	test.AssertEquals(t, TextSelection{2, 5, true}, c.Selection(0))
	test.AssertEquals(t, TextSelection{6, 6, false}, c.Selection(1))
	test.AssertEquals(t, TextSelection{9, 11, false}, c.Selection(2))
}

func TestTBCWordAt(t *testing.T) {
	check := func(str, expected string) {
		c := parseTBC(str)
		s, e := c.WordAt(c.FirstCaret())
		c.SetSelection(TextSelection{s, e, false})
		assertTBCTextAndSelectionsEqual(t, expected, c)
	}
	check("abc.dE|f()", "abc.{dEf]()")
	check("dE|f10()", "{dEf10]()")
	check("hello_|world.foo", "{hello_world].foo")
}

func TestTBCReplaceAll(t *testing.T) {
	c := parseTBC("ħę|ľĺő\n|ŵōř|ŀď")
	c.ReplaceAll("_")
	assertTBCTextAndSelectionsEqual(t, "ħę{_]ľĺő\n{_]ŵōř{_]ŀď", c)
	c.Deselect(false)
	assertTBCTextAndSelectionsEqual(t, "ħę_|ľĺő\n_|ŵōř_|ŀď", c)
	c.ReplaceAll("¶§")
	assertTBCTextAndSelectionsEqual(t, "ħę_{¶§]ľĺő\n_{¶§]ŵōř_{¶§]ŀď", c)
	c.Deselect(false)
	assertTBCTextAndSelectionsEqual(t, "ħę_¶§|ľĺő\n_¶§|ŵōř_¶§|ŀď", c)
}

func TestTBCReplace(t *testing.T) {
	c := parseTBC("ħę|ľĺő\n|ŵōř|ŀď")
	c.Replace(func(s TextSelection) string {
		return fmt.Sprintf("%d", s.start)
	})
	assertTBCTextAndSelectionsEqual(t, "ħę{2]ľĺő\n{6]ŵōř{9]ŀď", c)
}

func TestTBCReplaceWithNewline(t *testing.T) {
	c := parseTBC("   XX|\n  YY|YY\n    ZZZ|")
	c.ReplaceWithNewline()
	assertTBCTextAndSelectionsEqual(t, "   XX\n|\n  YY\n|YY\n    ZZZ\n|", c)
}

func ReplaceWithNewlineKeepIndent(t *testing.T) {
	c := parseTBC("   XX|\n  YY|YY\n    ZZZ|")
	c.ReplaceWithNewline()
	assertTBCTextAndSelectionsEqual(t, "   XX\n   |\n  YY\n  |YY\n    ZZZ\n    |", c)
}

func TestTBCReplaceSelection(t *testing.T) {
	c := parseTBC("ħ{ęľ]ĺő\nę[ľĺ}\n|\n{ŵōřŀď]")
	c.ReplaceAll("ŵōřŀď")
	assertTBCTextAndSelectionsEqual(t, "ħ{ŵōřŀď]ĺő\nę[ŵōřŀď}\n{ŵōřŀď]\n{ŵōřŀď]", c)
}

func TestTBCBackspaceMulti(t *testing.T) {
	c := parseTBC("ħęľ|ĺő\nŵōř|ŀď")
	c.Backspace()
	assertTBCTextAndSelectionsEqual(t, "ħę|ĺő\nŵō|ŀď", c)
	c.Backspace()
	assertTBCTextAndSelectionsEqual(t, "ħ|ĺő\nŵ|ŀď", c)
}

func TestTBCBackspaceCaretOverlap(t *testing.T) {
	c := parseTBC("ħęľ|ĺ|ő\nŵōřŀď")
	c.Backspace()
	assertTBCTextAndSelectionsEqual(t, "ħę|ő\nŵōřŀď", c)
	c.Backspace()
	assertTBCTextAndSelectionsEqual(t, "ħ|ő\nŵōřŀď", c)
}

func TestTBCBackspaceAtStart(t *testing.T) {
	c := parseTBC("|ħęľĺő\n|ŵōřŀď")
	c.Backspace()
	assertTBCTextAndSelectionsEqual(t, "|ħęľĺő|ŵōřŀď", c)
}

func TestTBCBackspaceSelection(t *testing.T) {
	c := parseTBC("ħ[ęľ}ĺ{ő]\n[ŵ}{ō]řŀď|")
	c.Backspace()
	assertTBCTextAndSelectionsEqual(t, "ħ|ĺ|\n|řŀ|", c)
}

func TestTBCDeleteAtEnd(t *testing.T) {
	c := parseTBC("ħęľĺő|\nŵōřŀď|")
	c.Delete()
	assertTBCTextAndSelectionsEqual(t, "ħęľĺő|ŵōřŀď|", c)
}

func TestTBCDeleteSelection(t *testing.T) {
	c := parseTBC("ħ[ęľ}ĺ{ő]\n[ŵ}{ō]řŀď|")
	c.Delete()
	assertTBCTextAndSelectionsEqual(t, "ħ|ĺ|\n|řŀď|", c)
}

func TestTBCAddCaretsUp(t *testing.T) {
	c := parseTBC("ÀÁAA|\nBBB|\nĆ|\nDDD|")
	c.AddCaretsUp()
	assertTBCTextAndSelectionsEqual(t, "|ÀÁA|A|\nB|BB|\nĆ|\nDDD|", c)
}

func TestTBCAddCaretsDown(t *testing.T) {
	c := parseTBC("ÀÁAA|\nBBB|\nĆ|\nDDD|")
	c.AddCaretsDown()
	assertTBCTextAndSelectionsEqual(t, "ÀÁAA|\nBBB|\nĆ|\nD|DD|", c)
}

func TestTBCMoveUp(t *testing.T) {
	c := parseTBC("ÀÁAA|\nBBB|\nĆ|\nDDD|")
	c.MoveUp()
	assertTBCTextAndSelectionsEqual(t, "|ÀÁA|A\nB|BB\nĆ|\nDDD", c)
}

func TestTBCMoveDown(t *testing.T) {
	c := parseTBC("ÀÁAA|\nBBB|\nĆ|\nDDD|")
	c.MoveDown()
	assertTBCTextAndSelectionsEqual(t, "ÀÁAA\nBBB|\nĆ|\nD|DD|", c)
}

func TestTBCMoveLeft(t *testing.T) {
	c := parseTBC("ÀÁAA|\nBBB|\nĆ|\nDDD|")
	c.MoveLeft()
	assertTBCTextAndSelectionsEqual(t, "ÀÁA|A\nBB|B\n|Ć\nDD|D", c)
}

func TestTBCMoveRight(t *testing.T) {
	c := parseTBC("ÀÁAA|\nBBB|\nĆ|\nDDD|")
	c.MoveRight()
	assertTBCTextAndSelectionsEqual(t, "ÀÁAA\n|BBB\n|Ć\n|DDD|", c)
}

func TestTBCSelectUp(t *testing.T) {
	c := parseTBC("ÀÁAA|\nBBB|\nĆ|\nDDD|")
	c.SelectUp()
	assertTBCTextAndSelectionsEqual(t, "[ÀÁAA\nBBB\nĆ}[\nDDD}", c)
}

func TestTBCSelectDown(t *testing.T) {
	c := parseTBC("ÀÁAA|\nBBB|\nĆ|\nDDD|")
	c.SelectDown()
	assertTBCTextAndSelectionsEqual(t, "ÀÁAA{\nBBB]{\nĆ]{\nD]DD|", c)
}

func TestTBCSelectLeft(t *testing.T) {
	c := parseTBC("ÀÁAA|\nBBB|\nĆ|\nDDD|")
	c.SelectLeft()
	assertTBCTextAndSelectionsEqual(t, "ÀÁA[A}\nBB[B}\n[Ć}\nDD[D}", c)
}

func TestTBCSelectRight(t *testing.T) {
	c := parseTBC("ÀÁAA|\nBBB|\nĆ|\nDDD|")
	c.SelectRight()
	assertTBCTextAndSelectionsEqual(t, "ÀÁAA{\n]BBB{\n]Ć{\n]DDD|", c)
}

func TestTBCSelectToWordLeft(t *testing.T) {
	c := parseTBC("ÀÁ|AA\nBB B|B BB\nf|oo() Ć|\n|DDD")
	c.SelectLeftByWord()
	assertTBCTextAndSelectionsEqual(t, "[ÀÁ}AA\nBB [B}B BB\n[f}oo() [Ć}[\n}DDD", c)
}

func TestTBCSelectToWordRight(t *testing.T) {
	c := parseTBC("ÀÁ|AA\nBB B|B BB\nf|oo() Ć|\n|DDD")
	c.SelectRightByWord()
	assertTBCTextAndSelectionsEqual(t, "ÀÁ{AA]\nBB B{B] BB\nf{oo]() Ć{\n]{DDD]", c)
}

func TestTBCSelectHome(t *testing.T) {
	c := parseTBC("ÀÁ|AA\n  BB|B\nĆ|\n    |DDD\n   | EE")
	c.SelectHome()
	assertTBCTextAndSelectionsEqual(t, "[ÀÁ}AA\n  [BB}B\n[Ć}\n[    }DDD\n[   } EE", c)
}

func TestTBCSelectEnd(t *testing.T) {
	c := parseTBC("ÀÁ|AA\nBB|B\nĆ|\nD|DD")
	c.SelectEnd()
	assertTBCTextAndSelectionsEqual(t, "ÀÁ{AA]\nBB{B]\nĆ|\nD{DD]", c)
}

func TestTBCUnicode(t *testing.T) {
	c := parseTBC("|1£2£3")
	c.MoveRight()
	assertTBCTextAndSelectionsEqual(t, "1|£2£3", c)
	c.MoveRight()
	assertTBCTextAndSelectionsEqual(t, "1£|2£3", c)
	c.MoveRight()
	assertTBCTextAndSelectionsEqual(t, "1£2|£3", c)
}

func TestTBCIndentSelection(t *testing.T) {
	c := parseTBC("a{aa\n  b]bb|bb\n    [cc}\nddd\ne{e][e}e\n")
	c.IndentSelection(2)
	assertTBCTextAndSelectionsEqual(t, "  a{aa\n    b]bb|bb\n      [cc}\nddd\n  e{e][e}e\n", c)
}

func TestTBCUnindentSelection(t *testing.T) {
	c := parseTBC("  a{aa\n    b]bb|bb\n      [cc}\nddd\n  e{e][e}e\n")
	c.UnindentSelection(2)
	assertTBCTextAndSelectionsEqual(t, "a{aa\n  b]bb|bb\n    [cc}\nddd\ne{e][e}e\n", c)
}
