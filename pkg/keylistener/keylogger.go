package keylistener

import (
	"syscall"
	"unicode/utf8"
	"unsafe"

	"github.com/TheTitanrain/w32"
)

var (
	moduser32 = syscall.NewLazyDLL("user32.dll")

	procGetKeyboardLayout     = moduser32.NewProc("GetKeyboardLayout")
	procGetKeyboardState      = moduser32.NewProc("GetKeyboardState")
	procToUnicodeEx           = moduser32.NewProc("ToUnicodeEx")
	procGetKeyboardLayoutList = moduser32.NewProc("GetKeyboardLayoutList")
	procMapVirtualKeyEx       = moduser32.NewProc("MapVirtualKeyExW")
	procGetKeyState           = moduser32.NewProc("GetKeyState")
)

// NewKeylogger creates a new keylogger depending on
// the platform we are running on (currently only Windows
// is supported)
func NewKeylogger() Keylogger {
	kl := Keylogger{}

	return kl
}

// Keylogger represents the keylogger
type Keylogger struct {
	Macro    []int
	lastKeys []int
}

// Key is a single key entered by the user
type Key struct {
	Empty   bool
	Rune    rune
	Keycode int
	Name    string
}

// Check Macro returns true if the macro is complete
func (kl *Keylogger) CheckMacro() bool {
	keys := kl.GetKeys()

	if len(keys) != len(kl.Macro) {
		return false
	}

	for _, key := range keys {
		isIn := false
		for _, k := range kl.Macro {
			if k == key.Keycode {
				isIn = true
				break
			}
		}
		if !isIn {
			return false
		}
	}

	return true
}

// GetKey gets the current entered key by the user, if there is any
func (kl *Keylogger) GetKeys() (keys []Key) {
	var activeKeys []int
	var keyState uint16

	for i := 0; i < 256; i++ {
		keyState = w32.GetAsyncKeyState(i)

		// Check if the most significant bit is set (key is down)
		// And check if the key is not a non-char key (except for space, 0x20)
		if keyState&(1<<15) != 0 && (i < 160 || i > 165) && (i < 91 || i > 93) {
			activeKeys = append(activeKeys, i)
		}
	}

	isNew := false
	for _, key := range activeKeys {
		isIn := false
		for _, otherKey := range kl.lastKeys {
			if key == otherKey {
				isIn = true
				break
			}
		}
		if !isIn {
			isNew = true
			break
		}
	}

	if isNew || len(kl.lastKeys) != len(activeKeys) {
		for _, key := range activeKeys {
			keys = append(keys, kl.ParseKeycode(key, keyState))
		}
	}

	kl.lastKeys = activeKeys

	return keys
}

// ParseKeycode returns the correct Key struct for a key taking in account the current keyboard settings
// That struct contains the Rune for the key
func (kl Keylogger) ParseKeycode(keyCode int, keyState uint16) Key {
	key := Key{Empty: false, Keycode: keyCode, Name: CodeToName(keyCode)}

	// Only one rune has to fit in
	outBuf := make([]uint16, 1)

	// Buffer to store the keyboard state in
	kbState := make([]uint8, 256)

	// Get keyboard layout for this process (0)
	kbLayout, _, _ := procGetKeyboardLayout.Call(uintptr(0))

	// Put all key modifier keys inside the kbState list
	if w32.GetAsyncKeyState(w32.VK_SHIFT)&(1<<15) != 0 {
		kbState[w32.VK_SHIFT] = 0xFF
	}

	capitalState, _, _ := procGetKeyState.Call(uintptr(w32.VK_CAPITAL))
	if capitalState != 0 {
		kbState[w32.VK_CAPITAL] = 0xFF
	}

	if w32.GetAsyncKeyState(w32.VK_CONTROL)&(1<<15) != 0 {
		kbState[w32.VK_CONTROL] = 0xFF
	}

	if w32.GetAsyncKeyState(w32.VK_MENU)&(1<<15) != 0 {
		kbState[w32.VK_MENU] = 0xFF
	}

	_, _, _ = procToUnicodeEx.Call(
		uintptr(keyCode),
		uintptr(0),
		uintptr(unsafe.Pointer(&kbState[0])),
		uintptr(unsafe.Pointer(&outBuf[0])),
		uintptr(1),
		uintptr(1),
		uintptr(kbLayout))

	key.Rune, _ = utf8.DecodeRuneInString(syscall.UTF16ToString(outBuf))

	return key
}
