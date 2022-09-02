package cli

import (
	"bufio"
	"d7024e/kademlia"
	"fmt"
	"os"
	"strings"
)

func PrintHello() {
	fmt.Println("Hello! This is your CLI speaking!")
	fmt.Println("Type 'help' for a list of commands.")
}

var commands = [][]string{
	{"get", "Takes a hash and downloads the file from the network."},
	{"help", "Help on commands."},
	{"put", "Uploads a file to the network and returns the hash if succesful."},
	{"stat", "Displays the status of the network."},
	{"exit", "Exit the CLI."},
}

func Open(context *kademlia.Kademlia) {
	for {
		fmt.Print(">>> ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("Error reading input.")
		}

		input = strings.TrimRight(input, "\r\n")
		cmd := strings.SplitN(input, " ", 2)

		switch cmd[0] {
		case "exit":
			exitApplication()
		case "get":
			if len(cmd) != 2 {
				fmt.Println("Invalid arguments.")
				continue
			}
			getObjectByHash(context, &cmd[1])
		case "help":
			printAvaliableCommands()
		case "put":
			if len(cmd) != 2 {
				fmt.Println("Invalid arguments.")
				continue
			}
			putObjectInStore(context, &cmd[1])
		case "stat":
			fmt.Println("Not implemented yet") // TODO: Should this be a thing?
		default:
			fmt.Println("Command not found. Type 'help' for a list of commands.")
		}
	}
}

func getObjectByHash(context *kademlia.Kademlia, hash *string) (string, error) {
	// TODO: Return error if object not found
	// TODO: Return the requested object if found
	fmt.Println("You asked for the object with hash", *hash)
	context.LookupData(*hash)
	return "", nil
}

func putObjectInStore(context *kademlia.Kademlia, content *string) (string, error) {
	// TODO: Return error if object could not be stored
	// TODO: Return the hash of the data
	fmt.Println("You asked to put the object", *content)
	dataToSend := []byte(*content)
	context.Store(dataToSend)
	return "", nil
}

func exitApplication() {
	fmt.Println("Goodbye!")
	os.Exit(0)
}

func printAvaliableCommands() {
	fmt.Println("Available commands:")

	// Calulate the longest command
	longestLength := 0
	for _, command := range commands {
		if len(command[0]) > longestLength {
			longestLength = len(command[0])
		}
	}

	// Print avaliable commands and their description
	for i := 0; i < len(commands); i++ {
		line := fmt.Sprintf("  %s%s   %s\n", commands[i][0], strings.Repeat(" ", longestLength-len(commands[i][0])), commands[i][1])
		fmt.Print(line)
	}

	fmt.Println()
}
