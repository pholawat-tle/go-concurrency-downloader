package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var commandList = map[string]func(params []string){
	"help": helpCommand,
	"test": testCommand,
}

var helpList = map[string]string{
	"help": "Definition of help command",
	"test": "Test if the server supports partial request",
}

var helpText = "\nUsage: concurrent-downloader COMMAND [PARAMETERS]\n\nCommands:\n"



func main(){
	for key , val := range helpList {
		helpText += fmt.Sprintf("	%s	%s\n",key,val)
	}
	hasArgs, args := parseArgs()
	if hasArgs{
		parseCommand(args)
	}else{
		fmt.Print(helpText)
	}
}

func parseArgs() (bool,[]string) {
	if(len(os.Args) > 1) {
		return true, os.Args[1:]
	}
	
	return false,[]string{}
}

func parseCommand(args []string) {
	var command string
	var params []string
	if(len(args) > 1){
		command = args[0]
		params = args[1:]
	}else{
		command = args[0]
		params = make([]string, 0)
	}
	
	if val, ok := commandList[command]; ok{
		val(params)
		return
	}
	fmt.Printf("The command %s doesn't exist\n", command)
}

func helpCommand(params []string){
	if len(params) > 1 {
		fmt.Println("Invalid Input : This command expects no parameter!")
	} else if len(params) == 1 {
		if val, ok := helpList[params[0]]; ok {
			fmt.Println(val)
		}else{
			fmt.Printf("The command %s doesn't exist!\n", params[0])
		}
	} else{
		fmt.Print(helpText)
	}
}

func testCommand(params []string){
	if len(params) != 1 {
		fmt.Println("Invalid Input : This command expects 1 parameter, which is the file's URL")
	}else{
		res, err := http.Head(params[0])
		if err != nil {
   			log.Fatal(err)
		}
		if res.StatusCode == http.StatusOK && res.Header.Get("Accept-Ranges") == "bytes" {
			fmt.Printf("%s can be downloaded with concurrent-downloader\n", params[0])
		}else{
			fmt.Printf("%s can't be downloaded with concurrent-downloader\n", params[0])
		}
	}
}