package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

var FILE_ERROR_RETRY = 3
var FILE_PATH = "./keyword/"
var DEFAULT_FILE_MODE os.FileMode = 0666
var MAX_RETRY = 3
var VERBOSE_ENABLE = false

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
func readFlags() (string, int) {
	var defaultport = 80
	var hostname = flag.String("host", "", "target hostname")
	var port = flag.Int("port", defaultport, "target port")
	var searchDepth = flag.Int("depth", 2, "maximum nested directory")
	var verbose = flag.Bool("verbose", true, "display realtime logging")
	flag.Parse()

	VERBOSE_ENABLE = *verbose

	if _, err := url.ParseRequestURI(*hostname); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var hostnameport = *hostname
	if *port != defaultport {
		hostnameport += ":" + strconv.Itoa(*port)
	}
	return hostnameport, *searchDepth
}

func main() {

	// Initialization
	var hostname, searchDepth = readFlags()

	fmt.Println(style.title("{FileSweep}"))
	fmt.Println(style.log("  Hostname: " + hostname + "\n"))

	fmt.Println(style.title("{Keywords}"))
	loadData()

	var starttime = time.Now()
	fmt.Println("Start time:", starttime.Format("2006-01-02 15:04:05"))

	// Path sweep
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

	// File sweep
	fmt.Println(style.title("# Start File Sweeping"))
	for _, foundPath := range allFoundPaths {
		sweepFile(foundPath)
	}

	var elapsed = int(time.Now().Unix() - starttime.Unix())
	fmt.Println("Total time elapsed:", elapsed, "sec")
}
