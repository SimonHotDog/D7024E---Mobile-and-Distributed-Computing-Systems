package cli

import (
	"bufio"
	"d7024e/kademlia"
	"errors"
	"fmt"
	"os"
	"strings"
)

type command struct {
	name string
	arg  string
}

func PrintHello() {
	fmt.Println("Hello! This is your CLI speaking!")
	fmt.Println("Type 'help' for a list of commands.")
}

var commands = [][]string{
	{"get [hash]", "Takes a hash and downloads the file from the network."},
	{"help", "Help on commands."},
	{"put [text]", "Uploads a file to the network and returns the hash if succesful."},
	{"stat", "Displays the status of the network."},
	{"exit", "Exit the CLI."},
	{"ping [address]", "DEBUG: Send a ping RPC to the target client"},
}

func Open(context *kademlia.Kademlia) {
	PrintHello()
	for {
		fmt.Print(">>> ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("Error reading input.")
		}

		cmd := parseCommand(input)
		if cmd.name == "" {
			continue
		}

		result, err := performCommand(context, &cmd)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		} else {
			fmt.Println(result)
		}
	}

}

func parseCommand(input string) command {
	input = strings.TrimRight(input, "\r\n")
	cmd := strings.SplitN(input, " ", 2)

	if len(cmd) == 2 {
		return command{cmd[0], cmd[1]}
	}
	return command{cmd[0], ""}
}

func performCommand(context *kademlia.Kademlia, cmd *command) (string, error) {
	switch cmd.name {
	case "exit":
		exitApplication()
		return "", nil
	case "get":
		if cmd.arg == "" {
			return "", errors.New("expected 1 argument, but got 0")
		}
		getObjectByHash(context, &cmd.arg)
		return "", nil // TODO: Return object
	case "help":
		return getAvaliableCommands(), nil
	case "put":
		if cmd.arg == "" {
			return "", errors.New("expected 1 argument, but got 0")
		}
		putObjectInStore(context, &cmd.arg)
		return "", nil // TODO: Return hash
	case "stat":
		return "", errors.New("not implemented yet") // TODO: Should this be a thing?
	case "ping":
		debug_sendPing(context, cmd.arg)
		return "", nil
	default:
		return "", errors.New("command not found. Type 'help' for a list of commands")
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

func getAvaliableCommands() string {
	var sb strings.Builder

	// Calulate the longest command
	longestLength := 0
	for _, command := range commands {
		if len(command[0]) > longestLength {
			longestLength = len(command[0])
		}
	}

	// Generate string of avaliable commands and their description
	sb.WriteString("Available commands:\n")
	for i := 0; i < len(commands); i++ {
		line := fmt.Sprintf("  %s%s   %s\n", commands[i][0], strings.Repeat(" ", longestLength-len(commands[i][0])), commands[i][1])
		sb.WriteString(line)
	}

	return sb.String()
}

func debug_sendPing(context *kademlia.Kademlia, args string) {
	contact := kademlia.Contact{Address: args}
	context.Network.SendPingMessage(&contact)
}
