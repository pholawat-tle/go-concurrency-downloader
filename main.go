package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
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
			simpleDownload(params[0])
			_ = contentSize
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
