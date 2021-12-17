package keyboard

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	user32                = syscall.NewLazyDLL("user32.dll")
	vkKeyScanExAProc      = user32.NewProc("VkKeyScanExA")
	getKeyboardLayoutProc = user32.NewProc("GetKeyboardLayout")
	mapVirtualKeyExAProc  = user32.NewProc("MapVirtualKeyExA")
	getKeyStateProc       = user32.NewProc("GetKeyState")
	keydbProc             = user32.NewProc("keybd_event")
	sendInputProc         = user32.NewProc("SendInput")
)

// INPUT winuser.h structure
type INPUT struct {
	inputType uint32
	ki        KEYBDINPUT
	padding   uint64
}

// KEYBDINPUT winuser.h structure
type KEYBDINPUT struct {
	wVk         uint16
	wScan       uint16
	dwFlags     uint32
	time        uint32
	dwExtraInfo uint64
}

type inputBuilder struct {
	inputs    []INPUT
	capitalOn bool
}

func (ib *inputBuilder) append(i INPUT) {
	ib.inputs = append(ib.inputs, i)
}

func (ib *inputBuilder) pressCapital() {
	capsDown := NewKeyboardInput(_VK_CAPITAL, 0, 0)
	capsUp := NewKeyboardInput(_VK_CAPITAL, 0, _KEYEVENTF_KEYUP)
	ib.append(capsDown)
	ib.append(capsUp)
	ib.capitalOn = !ib.capitalOn
}

// NewKeyboardInput INPUT constructor
func NewKeyboardInput(wVk uint16, wScan uint16, dwFlags uint32) (kbinput INPUT) {
	kbinput.inputType = 1
	kbinput.ki.wVk = wVk
	kbinput.ki.wScan = wScan
	kbinput.ki.dwFlags = dwFlags
	return
}

// SendString converts string chars into virtual key codes and feeds them into SendInput system call
func SendString(s string) {
	// get keyboard layout for VkKeyScanExA
	hkl, _, _ := getKeyboardLayoutProc.Call(uintptr(0))
	// get caps lock state 0x14 is VK_CAPITAL
	state, _, _ := getKeyStateProc.Call(uintptr(_VK_CAPITAL))

	var inputs inputBuilder
	isCapsOn := state&0x0001 != 0
	inputs.capitalOn = isCapsOn

	for _, c := range s {
		vkscan, _, _ := vkKeyScanExAProc.Call(uintptr(c), hkl)
		var vkc = vkscan & 0xFF
		var shiftState int16 = int16((vkscan >> 8) & 0xFF)

		if shiftState == 1 && !inputs.capitalOn {
			inputs.pressCapital()
		}

		// translate VK_X to VK_NUMPADX
		if vkc >= 0x30 && vkc <= 0x39 {
			vkc += 0x30
		}

		vsc, _, _ := mapVirtualKeyExAProc.Call(vkc, uintptr(0))
		inputDown := NewKeyboardInput(uint16(vkc), uint16(vsc), 0)
		inputUp := NewKeyboardInput(uint16(vkc), uint16(vsc), _KEYEVENTF_KEYUP)

		inputs.append(inputDown)
		inputs.append(inputUp)
	}

	if isCapsOn != inputs.capitalOn {
		inputs.pressCapital()
	}

	var dummy INPUT

	n, _, err := sendInputProc.Call(
		uintptr(len(inputs.inputs)),
		uintptr(unsafe.Pointer((*[1]INPUT)(inputs.inputs))), // get underlying array pointer from slice
		uintptr(unsafe.Sizeof(dummy)))

	fmt.Println(n, err, unsafe.Sizeof(dummy))
}

const (
	_VK_SHIFT        = 0x10
	_VK_CTRL         = 0x11
	_VK_ALT          = 0x12
	_VK_CAPITAL      = 0x14
	_KEYEVENTF_KEYUP = 0x0002
)
