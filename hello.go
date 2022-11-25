package main

import (
	"fmt"
	"time"

	color "github.com/fatih/color"
)

var updateCurrentLine string = "\033[2K\r"
var myColor = color.New(color.FgHiBlue, color.BlinkRapid, color.Bold).SprintFunc()

func main() {
	fmt.Println("Hello World!")
	fmt.Printf(updateCurrentLine + "Nya~!")
	time.Sleep(time.Millisecond * 1000)
	fmt.Printf(updateCurrentLine + myColor("Meow!"))
}
