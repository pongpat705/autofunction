package main

import (
	"autofunction/src/stuct"
	"encoding/json"
	"fmt"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	var function1 = os.Getenv("jsonPath")
	fmt.Printf("%s : param function1 = %s", "hello test", function1)
	var config = readJsonFile(function1)
	playground(config)
	fmt.Printf("%s", "done program")

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

func playground(config stuct.ConfigModel) {

	evChan := hook.Start()
	defer hook.End()

	var pressCtrl = false
	var pressEscape = false
	for ev := range evChan {
		fmt.Println("hook: ", ev)

		var color = getPixelAtMousePointerColor(config.IndicatorPositionX, config.IndicatorPositionY)
		isColorWhiteDo(color, config.KeyToPress)

		//break program
		if ev.Kind == hook.KeyHold && ev.Keychar == 65535 {
			pressCtrl = true
		}
		if ev.Kind == hook.KeyUp && ev.Keychar == 65535 {
			pressCtrl = false
		}
		if pressCtrl && ev.Keychar == 27 {
			pressEscape = true
		}
		fmt.Printf("pressCtrl %t pressEscape %t ", pressCtrl, pressEscape)
		if pressCtrl && pressEscape {
			break
		}

	}
}

func getPixelAtMousePointerColor(x int, y int) string {
	//x, y := robotgo.GetMousePos()
	color := robotgo.GetPixelColor(x, y)
	fmt.Printf("Color of pixel at (%d, %d) is 0x%s\n", x, y, color)
	return color
}

func isColorWhiteDo(color string, keyToPress string) {
	if color == strings.ToLower("FFFFFF") {
		fmt.Printf("tapping %s\n", keyToPress)
		err := robotgo.KeyTap(keyToPress)
		if err != "" {
			fmt.Println("robotgo.KeyTap run error is: ", err)
		}
	}
}
