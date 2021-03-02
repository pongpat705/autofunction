package main

import (
	"autofunction/src/stuct"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var (
	moduser32               = syscall.NewLazyDLL("user32.dll")
	procPostMessage         = moduser32.NewProc("PostMessageW")
	procFindWindowW         = moduser32.NewProc("FindWindowW")
	procSendInput           = moduser32.NewProc("SendInput")
	procSetWindowsHookExW   = moduser32.NewProc("SetWindowsHookExW")
	procLowLevelKeyboard    = moduser32.NewProc("LowLevelKeyboardProc")
	procCallNextHookEx      = moduser32.NewProc("CallNextHookEx")
	procUnhookWindowsHookEx = moduser32.NewProc("UnhookWindowsHookEx")
	procGetMessage          = moduser32.NewProc("GetMessageW")
	procTranslateMessage    = moduser32.NewProc("TranslateMessage")
	procDispatchMessage     = moduser32.NewProc("DispatchMessageW")
	keyboardHook            HHOOK
	keyChar                 = map[string]DWORD{
		"A":   0x41,
		"B":   0x42,
		"C":   0x43,
		"D":   0x44,
		"E":   0x45,
		"F":   0x46,
		"G":   0x47,
		"H":   0x48,
		"I":   0x49,
		"J":   0x4a,
		"K":   0x4b,
		"L":   0x4c,
		"M":   0x4d,
		"N":   0x4e,
		"O":   0x4f,
		"P":   0x50,
		"Q":   0x51,
		"R":   0x52,
		"S":   0x53,
		"T":   0x54,
		"U":   0x55,
		"V":   0x56,
		"W":   0x57,
		"X":   0x58,
		"Y":   0x59,
		"Z":   0x5a,
		"0":   0x30,
		"1":   0x31,
		"2":   0x32,
		"3":   0x33,
		"4":   0x34,
		"5":   0x35,
		"6":   0x36,
		"7":   0x37,
		"8":   0x38,
		"9":   0x39,
		"F1":  0x01000030,
		"F2":  0x01000031,
		"F3":  0x01000032,
		"F4":  0x01000033,
		"F5":  0x01000034,
		"F6":  0x01000035,
		"F7":  0x01000036,
		"F8":  0x01000037,
		"F9":  0x01000038,
		"F10": 0x01000039,
		"F11": 0x0100003a,
		"F12": 0x0100003b,
	}
	sysAction = map[string]DWORD{
		"KEYDOWN":         0x100,
		"LEFT_BUTTONDOWN": 0x0201,
		"LEFT_BUTTONUP":   0x0202,
	}
)

func main() {
	var arg1 = os.Args[1]
	fmt.Printf("%s : param function1 = %s", "hello test", arg1)
	var config = readJsonFile(arg1)
	var processId = getProcessId("Ragnarok")
	if config.Mode == "LOOP" {
		for true {
			auto(config, processId)
		}
	} else if config.Mode == "LISTEN" {
		keyPressing(config, processId)
	}

	fmt.Printf("%s", "done program")
}

func auto(config stuct.ConfigModel, processId uintptr) {

	for _, elm := range config.Config {
		var spec = elm.Spec
		var key = strings.ToUpper(elm.Key)
		var delay = elm.Delay

		if spec == "key" {
			var keyInt = keyChar[key]
			postMessage(processId, sysAction["KEYDOWN"], keyInt, 0)
		} else if spec == "wait" {
			var duration = time.Duration(delay) * time.Millisecond
			time.Sleep(duration)
		} else if spec == "mouse" {
			var keyInt = sysAction[key]
			postMessage(processId, keyInt, keyInt, 0)
		}
	}

}

func keyPressing(config stuct.ConfigModel, processId uintptr) {
	keyboardHook = SetWindowsHookEx(WH_KEYBOARD_LL,
		func(nCode int, wparam WPARAM, lparam LPARAM) LRESULT {
			if nCode == 0 && wparam == WM_KEYDOWN {
				fmt.Println("key pressed:")
				kbdstruct := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lparam))
				code := byte(kbdstruct.VkCode)
				fmt.Printf("%q", code)
				var pressingKey = keyChar[strings.ToUpper(config.WhenPressKey)]
				if kbdstruct.VkCode == pressingKey {
					auto(config, processId)
				}
			}
			return CallNextHookEx(keyboardHook, nCode, wparam, lparam)
		}, 0, 0)
	var msg MSG
	for GetMessage(&msg, 0, 0, 0) != 0 {

	}

	UnhookWindowsHookEx(keyboardHook)
	keyboardHook = 0
}

func getProcessId(windowName string) uintptr {
	ret, _, err := procFindWindowW.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(windowName))), 0)
	log.Printf("getProcessId ret: %v error: %v", ret, err)

	return ret
}

func postMessage(processId uintptr, keyType DWORD, key DWORD, param DWORD) {
	ret, _, err := procPostMessage.Call(processId, uintptr(keyType), uintptr(key), uintptr(param))
	log.Printf("postMessage ret: %v error: %v", ret, err)
}

func sendInput(processId uintptr, keyType DWORD, key DWORD, param DWORD) {
	ret, _, err := procPostMessage.Call(processId, uintptr(keyType), uintptr(key), uintptr(param))
	log.Printf("sendInput ret: %v error: %v", ret, err)
}

func readJsonFile(jsonfile string) stuct.ConfigModel {
	jsonFile, err := os.Open(jsonfile)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Successfully Opened %s", jsonFile)
	// defer the closing of our jsonFile so that we can parse it later on

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config stuct.ConfigModel
	unmarshallErr := json.Unmarshal(byteValue, &config)
	if unmarshallErr != nil {
		fmt.Println(unmarshallErr)
	}

	defer jsonFile.Close()

	return config
}

const (
	WH_KEYBOARD_LL = 13
	WH_KEYBOARD    = 2
	WM_KEYDOWN     = 256
	WM_SYSKEYDOWN  = 260
	WM_KEYUP       = 257
	WM_SYSKEYUP    = 261
	WM_KEYFIRST    = 256
	WM_KEYLAST     = 264
	PM_NOREMOVE    = 0x000
	PM_REMOVE      = 0x001
	PM_NOYIELD     = 0x002
	WM_LBUTTONDOWN = 513
	WM_RBUTTONDOWN = 516
	NULL           = 0
)

type (
	DWORD     uint32
	WPARAM    uintptr
	LPARAM    uintptr
	LRESULT   uintptr
	HANDLE    uintptr
	HINSTANCE HANDLE
	HHOOK     HANDLE
	HWND      HANDLE
)

type HOOKPROC func(int, WPARAM, LPARAM) LRESULT

type KBDLLHOOKSTRUCT struct {
	VkCode      DWORD
	ScanCode    DWORD
	Flags       DWORD
	Time        DWORD
	DwExtraInfo uintptr
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd162805.aspx
type POINT struct {
	X, Y int32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms644958.aspx
type MSG struct {
	Hwnd    HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

func SetWindowsHookEx(idHook int, lpfn HOOKPROC, hMod HINSTANCE, dwThreadId DWORD) HHOOK {
	ret, _, _ := procSetWindowsHookExW.Call(
		uintptr(idHook),
		uintptr(syscall.NewCallback(lpfn)),
		uintptr(hMod),
		uintptr(dwThreadId),
	)
	return HHOOK(ret)
}

func CallNextHookEx(hhk HHOOK, nCode int, wParam WPARAM, lParam LPARAM) LRESULT {
	ret, _, _ := procCallNextHookEx.Call(
		uintptr(hhk),
		uintptr(nCode),
		uintptr(wParam),
		uintptr(lParam),
	)
	return LRESULT(ret)
}

func UnhookWindowsHookEx(hhk HHOOK) bool {
	ret, _, _ := procUnhookWindowsHookEx.Call(
		uintptr(hhk),
	)
	return ret != 0
}

func GetMessage(msg *MSG, hwnd HWND, msgFilterMin uint32, msgFilterMax uint32) int {
	ret, _, _ := procGetMessage.Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax))
	return int(ret)
}

func TranslateMessage(msg *MSG) bool {
	ret, _, _ := procTranslateMessage.Call(
		uintptr(unsafe.Pointer(msg)))
	return ret != 0
}

func DispatchMessage(msg *MSG) uintptr {
	ret, _, _ := procDispatchMessage.Call(
		uintptr(unsafe.Pointer(msg)))
	return ret
}

func LowLevelKeyboardProc(nCode int, wParam WPARAM, lParam LPARAM) LRESULT {
	ret, _, _ := procLowLevelKeyboard.Call(
		uintptr(nCode),
		uintptr(wParam),
		uintptr(lParam),
	)
	return LRESULT(ret)
}
