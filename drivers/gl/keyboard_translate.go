// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"github.com/google/gxui"

	"github.com/goxjs/glfw"
)

func translateKeyboardKey(in glfw.Key) gxui.KeyboardKey {
	switch in {
	case glfw.KeySpace:
		return gxui.KeySpace
	case glfw.KeyApostrophe:
		return gxui.KeyApostrophe
	case glfw.KeyComma:
		return gxui.KeyComma
	case glfw.KeyMinus:
		return gxui.KeyMinus
	case glfw.KeyPeriod:
		return gxui.KeyPeriod
	case glfw.KeySlash:
		return gxui.KeySlash
	case glfw.Key0:
		return gxui.Key0
	case glfw.Key1:
		return gxui.Key1
	case glfw.Key2:
		return gxui.Key2
	case glfw.Key3:
		return gxui.Key3
	case glfw.Key4:
		return gxui.Key4
	case glfw.Key5:
		return gxui.Key5
	case glfw.Key6:
		return gxui.Key6
	case glfw.Key7:
		return gxui.Key7
	case glfw.Key8:
		return gxui.Key8
	case glfw.Key9:
		return gxui.Key9
	case glfw.KeySemicolon:
		return gxui.KeySemicolon
	case glfw.KeyEqual:
		return gxui.KeyEqual
	case glfw.KeyA:
		return gxui.KeyA
	case glfw.KeyB:
		return gxui.KeyB
	case glfw.KeyC:
		return gxui.KeyC
	case glfw.KeyD:
		return gxui.KeyD
	case glfw.KeyE:
		return gxui.KeyE
	case glfw.KeyF:
		return gxui.KeyF
	case glfw.KeyG:
		return gxui.KeyG
	case glfw.KeyH:
		return gxui.KeyH
	case glfw.KeyI:
		return gxui.KeyI
	case glfw.KeyJ:
		return gxui.KeyJ
	case glfw.KeyK:
		return gxui.KeyK
	case glfw.KeyL:
		return gxui.KeyL
	case glfw.KeyM:
		return gxui.KeyM
	case glfw.KeyN:
		return gxui.KeyN
	case glfw.KeyO:
		return gxui.KeyO
	case glfw.KeyP:
		return gxui.KeyP
	case glfw.KeyQ:
		return gxui.KeyQ
	case glfw.KeyR:
		return gxui.KeyR
	case glfw.KeyS:
		return gxui.KeyS
	case glfw.KeyT:
		return gxui.KeyT
	case glfw.KeyU:
		return gxui.KeyU
	case glfw.KeyV:
		return gxui.KeyV
	case glfw.KeyW:
		return gxui.KeyW
	case glfw.KeyX:
		return gxui.KeyX
	case glfw.KeyY:
		return gxui.KeyY
	case glfw.KeyZ:
		return gxui.KeyZ
	case glfw.KeyLeftBracket:
		return gxui.KeyLeftBracket
	case glfw.KeyBackslash:
		return gxui.KeyBackslash
	case glfw.KeyRightBracket:
		return gxui.KeyRightBracket
	case glfw.KeyGraveAccent:
		return gxui.KeyGraveAccent
	case glfw.KeyWorld1:
		return gxui.KeyWorld1
	case glfw.KeyWorld2:
		return gxui.KeyWorld2
	case glfw.KeyEscape:
		return gxui.KeyEscape
	case glfw.KeyEnter:
		return gxui.KeyEnter
	case glfw.KeyTab:
		return gxui.KeyTab
	case glfw.KeyBackspace:
		return gxui.KeyBackspace
	case glfw.KeyInsert:
		return gxui.KeyInsert
	case glfw.KeyDelete:
		return gxui.KeyDelete
	case glfw.KeyRight:
		return gxui.KeyRight
	case glfw.KeyLeft:
		return gxui.KeyLeft
	case glfw.KeyDown:
		return gxui.KeyDown
	case glfw.KeyUp:
		return gxui.KeyUp
	case glfw.KeyPageUp:
		return gxui.KeyPageUp
	case glfw.KeyPageDown:
		return gxui.KeyPageDown
	case glfw.KeyHome:
		return gxui.KeyHome
	case glfw.KeyEnd:
		return gxui.KeyEnd
	case glfw.KeyCapsLock:
		return gxui.KeyCapsLock
	case glfw.KeyScrollLock:
		return gxui.KeyScrollLock
	case glfw.KeyNumLock:
		return gxui.KeyNumLock
	case glfw.KeyPrintScreen:
		return gxui.KeyPrintScreen
	case glfw.KeyPause:
		return gxui.KeyPause
	case glfw.KeyF1:
		return gxui.KeyF1
	case glfw.KeyF2:
		return gxui.KeyF2
	case glfw.KeyF3:
		return gxui.KeyF3
	case glfw.KeyF4:
		return gxui.KeyF4
	case glfw.KeyF5:
		return gxui.KeyF5
	case glfw.KeyF6:
		return gxui.KeyF6
	case glfw.KeyF7:
		return gxui.KeyF7
	case glfw.KeyF8:
		return gxui.KeyF8
	case glfw.KeyF9:
		return gxui.KeyF9
	case glfw.KeyF10:
		return gxui.KeyF10
	case glfw.KeyF11:
		return gxui.KeyF11
	case glfw.KeyF12:
		return gxui.KeyF12
	case glfw.KeyKP0:
		return gxui.KeyKp0
	case glfw.KeyKP1:
		return gxui.KeyKp1
	case glfw.KeyKP2:
		return gxui.KeyKp2
	case glfw.KeyKP3:
		return gxui.KeyKp3
	case glfw.KeyKP4:
		return gxui.KeyKp4
	case glfw.KeyKP5:
		return gxui.KeyKp5
	case glfw.KeyKP6:
		return gxui.KeyKp6
	case glfw.KeyKP7:
		return gxui.KeyKp7
	case glfw.KeyKP8:
		return gxui.KeyKp8
	case glfw.KeyKP9:
		return gxui.KeyKp9
	case glfw.KeyKPDecimal:
		return gxui.KeyKpDecimal
	case glfw.KeyKPDivide:
		return gxui.KeyKpDivide
	case glfw.KeyKPMultiply:
		return gxui.KeyKpMultiply
	case glfw.KeyKPSubtract:
		return gxui.KeyKpSubtract
	case glfw.KeyKPAdd:
		return gxui.KeyKpAdd
	case glfw.KeyKPEnter:
		return gxui.KeyKpEnter
	case glfw.KeyKPEqual:
		return gxui.KeyKpEqual
	case glfw.KeyLeftShift:
		return gxui.KeyLeftShift
	case glfw.KeyLeftControl:
		return gxui.KeyLeftControl
	case glfw.KeyLeftAlt:
		return gxui.KeyLeftAlt
	case glfw.KeyLeftSuper:
		return gxui.KeyLeftSuper
	case glfw.KeyRightShift:
		return gxui.KeyRightShift
	case glfw.KeyRightControl:
		return gxui.KeyRightControl
	case glfw.KeyRightAlt:
		return gxui.KeyRightAlt
	case glfw.KeyRightSuper:
		return gxui.KeyRightSuper
	case glfw.KeyMenu:
		return gxui.KeyMenu
	default:
		return gxui.KeyUnknown
	}
}

func translateKeyboardModifier(in glfw.ModifierKey) gxui.KeyboardModifier {
	out := gxui.ModNone
	if in&glfw.ModShift != 0 {
		out |= gxui.ModShift
	}
	if in&glfw.ModControl != 0 {
		out |= gxui.ModControl
	}
	if in&glfw.ModAlt != 0 {
		out |= gxui.ModAlt
	}
	if in&glfw.ModSuper != 0 {
		out |= gxui.ModSuper
	}
	return out
}
