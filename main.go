package main

import (
	"fmt"
	"os"
)

var commandList = map[string]func(params []string){
	"help": helpCommand,
}

var helpList = map[string]string{
	"help": "Definition of help command",
}

var helpText = "\nUsage: concurrent-downloader COMMAND [PARAMETERS]\n\nCommands:\n"



func main(){
	
	for key , val := range helpList {
		helpText += fmt.Sprintf("	%s		%s\n",key,val)
	}
	hasArgs, args := checkForArgs()
	if hasArgs{
		checkCommand(args)
	}else{
		fmt.Print(helpText)
	}
}

func checkForArgs() (bool,[]string) {
	if(len(os.Args) > 1) {
		return true, os.Args[1:]
	}
	
	return false,[]string{}
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
		fmt.Println(helpText)
	}
}

func checkCommand(args []string) {
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

