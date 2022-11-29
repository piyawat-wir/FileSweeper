package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"golang.org/x/exp/slices"
)

var updateCurrentLine string = "\033[2K\r"

var FILE_ERROR_RETRY = 3
var FILE_PATH = "./keyword/"
var DEFAULT_FILE_MODE os.FileMode = 0666

var verbose = false

type CompoundType struct {
	Directory []string
	File      []string
	Extension []string
}

var Compound = CompoundType{
	Directory: []string{},
	File:      []string{},
	Extension: []string{},
}

var compoundMap = map[string](*[]string){
	"Directory": &Compound.Directory,
	"File":      &Compound.File,
	"Extension": &Compound.Extension,
}

func warn(str string) {
	if verbose {
		fmt.Println(style.warn(str))
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(style.error("\nERROR: %v", err))
		os.Exit(1)
	}
}

func loadData() {

	// Check path existence
	var _, err = os.ReadDir(FILE_PATH)
	if os.IsNotExist(err) {
		if os.MkdirAll(FILE_PATH, DEFAULT_FILE_MODE) != nil {
			fmt.Println(style.error("Directory '" + FILE_PATH + "' cannot be created!"))
			os.Exit(1)
		}
		warn("Directory '" + FILE_PATH + "' is created.")
	}

	for compdName, cPtr := range compoundMap {

		var fpname = filepath.Join(FILE_PATH, compdName+".txt")
		var filePtr, err = os.Open(fpname)

		// Check file existence
		if os.IsNotExist(err) {

			// Create new file
			filePtr, err = os.Create(fpname)
			if err != nil {
				fmt.Println(style.error("File '" + fpname + "' cannot be created!"))
				os.Exit(1)
			}
			warn("File '" + fpname + "' is created.")

		}
		// Read file
		var count = 0
		var scanner = bufio.NewScanner(filePtr)
		for scanner.Scan() {
			var line = scanner.Text()
			color.New(color.Faint).Printf(updateCurrentLine+" > %v: %v", compdName, line)
			*cPtr = append(*cPtr, scanner.Text())
			count++
		}
		var suffixS = ""
		if count > 1 {
			suffixS = "s"
		}
		fmt.Println(updateCurrentLine + style.log(fmt.Sprintf(" + %v %v name%v", count, compdName, suffixS)))

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

var MAX_RETRY = 5
var client = http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConnsPerHost:   5,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

func getStatusCode(url string) int {
	var retry = 0
	var resp *http.Response
	var err error
	for {
		if retry >= MAX_RETRY {
			os.Exit(1)
		}
		resp, err = client.Get(url)
		if err == nil {
			break
		}
		fmt.Println(style.error("\nERROR: %v", err))
		retry++
		dotdotdotWait(" Retrying", 5)
	}

	return resp.StatusCode
}

var acceptHTTPStatus = []int{
	http.StatusOK,
	http.StatusForbidden,
	http.StatusMovedPermanently,
	http.StatusFound,
	http.StatusNotModified,
	http.StatusTemporaryRedirect,
}

func isHTTPStatusAccepted(status int) bool {
	return slices.Contains(acceptHTTPStatus, status)
}

var MapSingularToPlural = map[string]string{
	"directory": "directories",
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
	fmt.Print(style.slight(
		updateCurrentLine+"%s %s",
		progressPercentText(scanned, scanTarget),
		sweepURL,
	))
}
func printFound(respStatus int, sweepURL string) {
	fmt.Print(style.success(
		updateCurrentLine+"[%d] %s : %s\n",
		respStatus,
		http.StatusText(respStatus),
		sweepURL,
	))
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

	fmt.Printf(updateCurrentLine+"Found %d %s\n",
		len(foundPath),
		decidePlural("directory", len(foundPath)),
	)

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

	fmt.Printf(updateCurrentLine+"Found %d %s\n",
		len(foundFile),
		decidePlural("directory", len(foundFile)),
	)

	return foundFile
}

func randomPath(hostname string) {
	var mydir = Compound.Directory[rand.Intn(len(Compound.Directory))]
	var myfile = Compound.File[rand.Intn(len(Compound.File))]
	var myext = Compound.Extension[rand.Intn(len(Compound.Extension))]
	// var path, _ = url.JoinPath(hostname, mydir, myfile+"."+myext)
	var path = path.Join(hostname, mydir, myfile+"."+myext)

	fmt.Println(path)
}

func main() {
	fmt.Println(style.title("{Keywords}"))
	loadData()

	var hostname = "http://localhost"
	var searchDepth = 2

	fmt.Println(style.title("\n# Start %d-Level Path Sweeping @ %s", searchDepth, hostname))

	var allFoundPaths = []string{}
	var workingPath = []string{}
	var respStatus = getStatusCode(hostname)

	fmt.Println("Level 0")
	if isHTTPStatusAccepted(respStatus) {
		workingPath = append(workingPath, hostname)
		allFoundPaths = append(allFoundPaths, workingPath...)
		printFound(respStatus, hostname)
	}

	for depth := 1; depth < searchDepth; depth++ {
		fmt.Printf("Level %d\n", depth)
		for _, currentPath := range workingPath {
			if strings.Count(currentPath, "/")-1 >= depth {
				workingPath = sweepPath(currentPath)
				allFoundPaths = append(allFoundPaths, workingPath...)
			}
		}
	}

	fmt.Println(style.log("Total path found: %d\n", len(allFoundPaths)))

	fmt.Println(style.title("# Start File Sweeping"))

	for _, foundPath := range allFoundPaths {
		sweepFile(foundPath)
	}
}
