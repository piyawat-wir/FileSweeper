package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

var updateCurrentLine string = "\033[2K\r"
var MapSingularToPlural = map[string]string{
	"directory": "directories",
	"file":      "files",
}

func warn(str string) {
	if VERBOSE_ENABLE {
		fmt.Println(style.warn(str))
	}
}
func checkErr(err error) {
	if err != nil {
		fmt.Println(style.error("\nERROR: %v", err))
		os.Exit(1)
	}
}
func dotdotdotWait(str string, second int) {
	fmt.Print(style.slight(str))
	for i := 0; i < second; i++ {
		time.Sleep(time.Second)
		fmt.Print(style.slight("."))
	}
	fmt.Println()
}
func decidePlural(word string, num int) string {
	if plural, found := MapSingularToPlural[word]; found {
		if num > 1 {
			return plural
		}
	}
	return word
}
func progressPercentText(now int, target int) string {
	return fmt.Sprintf("[%v%%] ", int(100*now/target))
}
func printScanning(scanned int, scanTarget int, sweepURL string) {
	if VERBOSE_ENABLE {
		fmt.Print(style.slight(
			updateCurrentLine+"  %s %s",
			progressPercentText(scanned, scanTarget),
			sweepURL,
		))
	}
}
func printFound(respStatus int, sweepURL string) {
	fmt.Print(style.success(
		updateCurrentLine+"  [%d] %s : %s\n",
		respStatus,
		http.StatusText(respStatus),
		sweepURL,
	))
}
