package main

import (
	"autofunction/src/stuct"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var (
	moduser32       = syscall.NewLazyDLL("user32.dll")
	procPostMessage = moduser32.NewProc("PostMessageW")
	procFindWindowW = moduser32.NewProc("FindWindowW")
	keyChar         = map[string]int32{
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
)

func main() {
	var function1 = os.Args[1]
	fmt.Printf("%s : param function1 = %s", "hello test", function1)
	var config = readJsonFile(function1)
	//playground(config)
	auto(config)
	fmt.Printf("%s", "done program")
}

func auto(config stuct.ConfigModel) {

	var processId = getProcessId("Ragnarok")
	const WM_KEYDOWN = 0x100

	for true {
		for _, elm := range config.Config {
			var spec = elm.Spec
			var key = strings.ToUpper(elm.Key)
			var delay = elm.Delay

			if spec == "key" {
				var keyInt = keyChar[key]
				postMessage(processId, WM_KEYDOWN, keyInt, 0)
			} else if spec == "wait" {
				var yourTime = rand.Int31n(int32(delay))
				var duration = time.Duration(yourTime) * time.Millisecond
				time.Sleep(duration)
			}
		}
	}

}

func getProcessId(windowName string) uintptr {
	ret, _, err := procFindWindowW.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(windowName))), 0)
	log.Printf("getProcessId ret: %v error: %v", ret, err)

	return ret
}

func postMessage(processId uintptr, keyType uint32, key int32, param int32) {
	ret, _, err := procPostMessage.Call(processId, uintptr(keyType), uintptr(key), uintptr(param))
	log.Printf("postMessage ret: %v error: %v", ret, err)
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
