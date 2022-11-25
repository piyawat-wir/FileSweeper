package main

import (
	"fmt"
	"os"
	"time"
)

var updateCurrentLine string = "\033[2K\r"

func main() {
	testprint()

	var filename = "mydoc.txt"
	var data, err = os.ReadFile(filename)
	if err != nil {
		fmt.Println(style.error("File not found"))
		os.WriteFile(filename, []byte("hello world!"), 0644)
	}

	fmt.Println(style.log("Data in file:"))
	fmt.Println(string(data))

}

func testprint() {
	fmt.Printf(updateCurrentLine + "Hello World!")
	time.Sleep(time.Millisecond * 1000)
	fmt.Printf(updateCurrentLine + "Nya~! OwO")
	time.Sleep(time.Millisecond * 1000)
	fmt.Printf(updateCurrentLine + style.mycolor("Meow! >w<") + "\n")
	time.Sleep(time.Millisecond * 1000)
}
