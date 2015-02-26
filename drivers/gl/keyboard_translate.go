// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"github.com/go-gl/glfw3"
	"gxui"
)

func translateKeyboardKey(in glfw3.Key) gxui.KeyboardKey {
	switch in {
	case glfw3.KeySpace:
		return gxui.KeySpace
	case glfw3.KeyApostrophe:
		return gxui.KeyApostrophe
	case glfw3.KeyComma:
		return gxui.KeyComma
	case glfw3.KeyMinus:
		return gxui.KeyMinus
	case glfw3.KeyPeriod:
		return gxui.KeyPeriod
	case glfw3.KeySlash:
		return gxui.KeySlash
	case glfw3.Key0:
		return gxui.Key0
	case glfw3.Key1:
		return gxui.Key1
	case glfw3.Key2:
		return gxui.Key2
	case glfw3.Key3:
		return gxui.Key3
	case glfw3.Key4:
		return gxui.Key4
	case glfw3.Key5:
		return gxui.Key5
	case glfw3.Key6:
		return gxui.Key6
	case glfw3.Key7:
		return gxui.Key7
	case glfw3.Key8:
		return gxui.Key8
	case glfw3.Key9:
		return gxui.Key9
	case glfw3.KeySemicolon:
		return gxui.KeySemicolon
	case glfw3.KeyEqual:
		return gxui.KeyEqual
	case glfw3.KeyA:
		return gxui.KeyA
	case glfw3.KeyB:
		return gxui.KeyB
	case glfw3.KeyC:
		return gxui.KeyC
	case glfw3.KeyD:
		return gxui.KeyD
	case glfw3.KeyE:
		return gxui.KeyE
	case glfw3.KeyF:
		return gxui.KeyF
	case glfw3.KeyG:
		return gxui.KeyG
	case glfw3.KeyH:
		return gxui.KeyH
	case glfw3.KeyI:
		return gxui.KeyI
	case glfw3.KeyJ:
		return gxui.KeyJ
	case glfw3.KeyK:
		return gxui.KeyK
	case glfw3.KeyL:
		return gxui.KeyL
	case glfw3.KeyM:
		return gxui.KeyM
	case glfw3.KeyN:
		return gxui.KeyN
	case glfw3.KeyO:
		return gxui.KeyO
	case glfw3.KeyP:
		return gxui.KeyP
	case glfw3.KeyQ:
		return gxui.KeyQ
	case glfw3.KeyR:
		return gxui.KeyR
	case glfw3.KeyS:
		return gxui.KeyS
	case glfw3.KeyT:
		return gxui.KeyT
	case glfw3.KeyU:
		return gxui.KeyU
	case glfw3.KeyV:
		return gxui.KeyV
	case glfw3.KeyW:
		return gxui.KeyW
	case glfw3.KeyX:
		return gxui.KeyX
	case glfw3.KeyY:
		return gxui.KeyY
	case glfw3.KeyZ:
		return gxui.KeyZ
	case glfw3.KeyLeftBracket:
		return gxui.KeyLeftBracket
	case glfw3.KeyBackslash:
		return gxui.KeyBackslash
	case glfw3.KeyRightBracket:
		return gxui.KeyRightBracket
	case glfw3.KeyGraveAccent:
		return gxui.KeyGraveAccent
	case glfw3.KeyWorld1:
		return gxui.KeyWorld1
	case glfw3.KeyWorld2:
		return gxui.KeyWorld2
	case glfw3.KeyEscape:
		return gxui.KeyEscape
	case glfw3.KeyEnter:
		return gxui.KeyEnter
	case glfw3.KeyTab:
		return gxui.KeyTab
	case glfw3.KeyBackspace:
		return gxui.KeyBackspace
	case glfw3.KeyInsert:
		return gxui.KeyInsert
	case glfw3.KeyDelete:
		return gxui.KeyDelete
	case glfw3.KeyRight:
		return gxui.KeyRight
	case glfw3.KeyLeft:
		return gxui.KeyLeft
	case glfw3.KeyDown:
		return gxui.KeyDown
	case glfw3.KeyUp:
		return gxui.KeyUp
	case glfw3.KeyPageUp:
		return gxui.KeyPageUp
	case glfw3.KeyPageDown:
		return gxui.KeyPageDown
	case glfw3.KeyHome:
		return gxui.KeyHome
	case glfw3.KeyEnd:
		return gxui.KeyEnd
	case glfw3.KeyCapsLock:
		return gxui.KeyCapsLock
	case glfw3.KeyScrollLock:
		return gxui.KeyScrollLock
	case glfw3.KeyNumLock:
		return gxui.KeyNumLock
	case glfw3.KeyPrintScreen:
		return gxui.KeyPrintScreen
	case glfw3.KeyPause:
		return gxui.KeyPause
	case glfw3.KeyF1:
		return gxui.KeyF1
	case glfw3.KeyF2:
		return gxui.KeyF2
	case glfw3.KeyF3:
		return gxui.KeyF3
	case glfw3.KeyF4:
		return gxui.KeyF4
	case glfw3.KeyF5:
		return gxui.KeyF5
	case glfw3.KeyF6:
		return gxui.KeyF6
	case glfw3.KeyF7:
		return gxui.KeyF7
	case glfw3.KeyF8:
		return gxui.KeyF8
	case glfw3.KeyF9:
		return gxui.KeyF9
	case glfw3.KeyF10:
		return gxui.KeyF10
	case glfw3.KeyF11:
		return gxui.KeyF11
	case glfw3.KeyF12:
		return gxui.KeyF12
	case glfw3.KeyKP0:
		return gxui.KeyKp0
	case glfw3.KeyKP1:
		return gxui.KeyKp1
	case glfw3.KeyKP2:
		return gxui.KeyKp2
	case glfw3.KeyKP3:
		return gxui.KeyKp3
	case glfw3.KeyKP4:
		return gxui.KeyKp4
	case glfw3.KeyKP5:
		return gxui.KeyKp5
	case glfw3.KeyKP6:
		return gxui.KeyKp6
	case glfw3.KeyKP7:
		return gxui.KeyKp7
	case glfw3.KeyKP8:
		return gxui.KeyKp8
	case glfw3.KeyKP9:
		return gxui.KeyKp9
	case glfw3.KeyKPDecimal:
		return gxui.KeyKpDecimal
	case glfw3.KeyKPDivide:
		return gxui.KeyKpDivide
	case glfw3.KeyKPMultiply:
		return gxui.KeyKpMultiply
	case glfw3.KeyKPSubtract:
		return gxui.KeyKpSubtract
	case glfw3.KeyKPAdd:
		return gxui.KeyKpAdd
	case glfw3.KeyKPEnter:
		return gxui.KeyKpEnter
	case glfw3.KeyKPEqual:
		return gxui.KeyKpEqual
	case glfw3.KeyLeftShift:
		return gxui.KeyLeftShift
	case glfw3.KeyLeftControl:
		return gxui.KeyLeftControl
	case glfw3.KeyLeftAlt:
		return gxui.KeyLeftAlt
	case glfw3.KeyLeftSuper:
		return gxui.KeyLeftSuper
	case glfw3.KeyRightShift:
		return gxui.KeyRightShift
	case glfw3.KeyRightControl:
		return gxui.KeyRightControl
	case glfw3.KeyRightAlt:
		return gxui.KeyRightAlt
	case glfw3.KeyRightSuper:
		return gxui.KeyRightSuper
	case glfw3.KeyMenu:
		return gxui.KeyMenu
	default:
		return gxui.KeyUnknown
	}
}

func translateKeyboardModifier(in glfw3.ModifierKey) gxui.KeyboardModifier {
	out := gxui.ModNone
	if in&glfw3.ModShift != 0 {
		out |= gxui.ModShift
	}
	if in&glfw3.ModControl != 0 {
		out |= gxui.ModControl
	}
	if in&glfw3.ModAlt != 0 {
		out |= gxui.ModAlt
	}
	if in&glfw3.ModSuper != 0 {
		out |= gxui.ModSuper
	}
	return out
}
