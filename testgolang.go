package main

import (
	"fmt"
	"time"
)

var updateCurrentLine string = "\033[2K\r"

func main() {
	fmt.Println("Hello World!")
	fmt.Printf(updateCurrentLine + "Nya~!")
	time.Sleep(time.Millisecond * 1000)
	fmt.Printf(updateCurrentLine + style.mycolor("Meow!") + "\n")
	fmt.Printf("style: %v\n", style)
	//os.ReadFile()
}
