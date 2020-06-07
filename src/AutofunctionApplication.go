package main

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"os"
	"strings"
)

func main() {
	var function1 = os.Getenv("function1")
	fmt.Printf("%s : param function1 = %s", "hello test", function1)

	for true {
		x, y := robotgo.GetMousePos()
		color := robotgo.GetPixelColor(x, y)
		fmt.Printf("Color of pixel at (%d, %d) is 0x%s\n", x, y, color)

		if color == strings.ToLower("FFFFFF") {
			fmt.Printf("tapping %s", function1)
			robotgo.KeyTap(function1)
		}

	}



}
