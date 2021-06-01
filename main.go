package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"sync"
)

var commandList = map[string]func(params []string){
	"help":     helpCommand,
	"test":     testCommand,
	"download": downloadCommand,
}

var helpList = map[string]string{
	"help": "	Definition of help command",
	"test": "	Test if the server supports partial request",
	"download": "Download file from a URL",
}

var helpText = "\nUsage: concurrent-downloader COMMAND [PARAMETERS]\n\nCommands:\n"

func main() {
	for key, val := range helpList {
		helpText += fmt.Sprintf("	%s	%s\n", key, val)
	}
	hasArgs, args := parseArgs()
	if hasArgs {
		parseCommand(args)
	} else {
		fmt.Print(helpText)
	}
}

func parseArgs() (bool, []string) {
	if len(os.Args) > 1 {
		return true, os.Args[1:]
	}

	return false, []string{}
}

func parseCommand(args []string) {
	var command string
	var params []string
	if len(args) > 1 {
		command = args[0]
		params = args[1:]
	} else {
		command = args[0]
		params = make([]string, 0)
	}

	if val, ok := commandList[command]; ok {
		val(params)
		return
	}
	fmt.Printf("The command %s doesn't exist\n", command)
}

func helpCommand(params []string) {
	if len(params) > 1 {
		fmt.Println("Invalid Input : This command expects no parameter!")
	} else if len(params) == 1 {
		if val, ok := helpList[params[0]]; ok {
			fmt.Println(val)
		} else {
			fmt.Printf("The command %s doesn't exist!\n", params[0])
		}
	} else {
		fmt.Print(helpText)
	}
}

func testCommand(params []string) {
	if len(params) != 1 {
		fmt.Println("Invalid Input : This command expects 1 parameter, which is the file's URL")
	} else {
		res, err := http.Head(params[0])
		if err != nil {
			log.Fatal(err)
		}
		if res.StatusCode == http.StatusOK && res.Header.Get("Accept-Ranges") == "bytes" {
			fmt.Printf("%s can be downloaded with concurrent-downloader\n", params[0])
		} else {
			fmt.Printf("%s can't be downloaded with concurrent-downloader\n", params[0])
		}
	}
}

func downloadCommand(params []string) {
	if len(params) != 1 {
		fmt.Println("Invalid Input : This command expects 1 parameter, which is the file's URL")
	} else {
		res, err := http.Head(params[0])
		if err != nil {
			log.Fatal(err)
		}

		if res.StatusCode == http.StatusOK && res.Header.Get("Accept-Ranges") == "bytes" {
			contentSize, err := strconv.Atoi(res.Header.Get("Content-Length"))
			if err != nil {
				log.Fatal(err)
			}
			multiPartDownload(params[0], contentSize)
		} else {
			simpleDownload(params[0])
		}
	}

}

func simpleDownload(url string) {

	// make a request
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// create the output file
	f, err := os.OpenFile(path.Base(url), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// copy to output file
	buffer := make([]byte, 1024)
	_, err = io.CopyBuffer(io.MultiWriter(f), res.Body, buffer)
	if err != nil {
		log.Fatal(err)
	}
}

func multiPartDownload(url string, contentSize int) {
	partSize := contentSize / 16

	startRange := 0
	wg := &sync.WaitGroup{}
	wg.Add(16)

	for i := 1; i <= 16; i++ {
		if i == 16 {
			go downloadPartial(startRange, contentSize, i, wg, url)
		} else {
			go downloadPartial(startRange, startRange+partSize, i, wg, url)
		}

		startRange += partSize + 1
	}

	wg.Wait()
	merge(url)
}

func downloadPartial(rangeStart int, rangeStop int, partNo int, wg *sync.WaitGroup, url string) {
	defer wg.Done()
	fmt.Printf("Downloading part %d from byte %d to %d\n", partNo, rangeStart, rangeStop)

	if rangeStart >= rangeStop {
		// nothing to download
		return
	}

	// create a request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", rangeStart, rangeStop))

	// make a request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// create the output file
	outputPath := getPartFilename(path.Base(url), partNo)
	flags := os.O_CREATE | os.O_WRONLY

	f, err := os.OpenFile(outputPath, flags, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// copy to output file
	for {
		_, err = io.CopyN(io.MultiWriter(f), res.Body, int64(1024))
		if err != nil {
			if err == io.EOF {
				return
			} else {
				log.Fatal(err)
			}
		}
	}
}

func merge(url string) {
	destination, err := os.OpenFile(path.Base(url), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer destination.Close()

	for i := 1; i <= 16; i++ {
		filename := getPartFilename(path.Base(url), i)
		source, err := os.OpenFile(filename, os.O_RDONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		io.Copy(destination, source)
		source.Close()
		os.Remove(filename)
	}
}

func getPartFilename(outName string, partNum int) string {
	return outName + ".part" + strconv.Itoa(partNum)
}
