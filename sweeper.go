package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/exp/slices"
)

func getStatusCode(url string) int {
	var retry = 0
	var resp *http.Response
	var err error
	for {
		resp, err = client.Head(url)
		if err == nil {
			break
		}
		fmt.Println(style.error("\nERROR: %v", err))
		retry++
		if retry >= MAX_RETRY {
			os.Exit(1)
		}
		dotdotdotWait(" Retrying", 5)
	}

	return resp.StatusCode
}
func isHTTPStatusAccepted(status int) bool {
	return slices.Contains(acceptHTTPStatus, status)
}
func sweepPath(startURL string) []string {
	var scanTarget = len(Compound.Directory)
	var scanned = 0
	var foundPath = []string{}
	for _, dirname := range Compound.Directory {
		var sweepURL, err = url.JoinPath(startURL, dirname)
		checkErr(err)

		scanned++
		printScanning(scanned, scanTarget, sweepURL)

		var respStatus = getStatusCode(sweepURL)
		if isHTTPStatusAccepted(respStatus) {
			foundPath = append(foundPath, sweepURL)
			printFound(respStatus, sweepURL)
		}
	}

	if len(foundPath) > 0 {
		fmt.Printf(updateCurrentLine+"Found %d %s @ %s\n",
			len(foundPath),
			decidePlural("directory", len(foundPath)),
			startURL,
		)
	}

	fmt.Print(updateCurrentLine)

	return foundPath
}
func sweepFile(startURL string) []string {
	var scanTarget = len(Compound.File) * len(Compound.Extension)
	var scanned = 0
	var foundFile = []string{}
	for _, fname := range Compound.File {
		for _, ext := range Compound.Extension {
			var expectFile = fname
			if ext != "" {
				expectFile += "." + ext
			}
			var sweepURL, err = url.JoinPath(startURL, expectFile)
			checkErr(err)

			scanned++
			printScanning(scanned, scanTarget, sweepURL)

			var respStatus = getStatusCode(sweepURL)
			if isHTTPStatusAccepted(respStatus) {
				foundFile = append(foundFile, sweepURL)
				printFound(respStatus, sweepURL)
			}
		}
	}

	if len(foundFile) > 0 {
		fmt.Printf(updateCurrentLine+"Found %d %s @ %s\n",
			len(foundFile),
			decidePlural("file", len(foundFile)),
			startURL,
		)
	}

	fmt.Print(updateCurrentLine)

	return foundFile
}
