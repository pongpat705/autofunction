package main

import (
	"autofunction/src/stuct"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"syscall"
	"unsafe"
)

var (
	moduser32       = syscall.NewLazyDLL("user32.dll")
	procSendInput   = moduser32.NewProc("SendInput")
	procPostMessage = moduser32.NewProc("PostMessageW")
	procFindWindowW = moduser32.NewProc("FindWindowW")
)

func main() {
	//var function1 = os.Getenv("jsonPath")
	//fmt.Printf("%s : param function1 = %s", "hello test", function1)
	//var config = readJsonFile(function1)
	//playground(config)
	//fmt.Printf("%s", "done program")
	auto()
}

func auto() {

	var processId = getProcessId("Ragnarok")
	const WM_KEYDOWN = 0x100
	postMessage(processId, WM_KEYDOWN, 0x41, 0)

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
