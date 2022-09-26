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
	{"whoami", "DEBUG: Lookup myself"},
	{"routes", "DEBUG: Print routingtable"},
	{"lookup", "DEBUG: Lookup"},
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
	case "whoami":
		debug_lookupMe(context, cmd.arg)
		return "", nil
	case "routes":
		debug_routingTable(context)
		return "", nil
	default:
		return "", errors.New("command not found. Type 'help' for a list of commands")
	}
}

func getObjectByHash(context *kademlia.Kademlia, hash *string) (string, error) {
	//TODO: No return values yet, only print
	fmt.Println("You asked for the object with hash", *hash)
	context.LookupData(*hash)
	/*value := context.LookupData(*hash)
	if value == nil {
		return "", errors.New("data not found")
	} else {
		return string(value), nil
	}*/
	return "", nil
}

func putObjectInStore(context *kademlia.Kademlia, content *string) (string, error) {
	fmt.Println("You asked to put the object", *content)
	dataToSend := []byte(*content)

	value, error := context.Store(dataToSend)

	if error == nil {
		return value, nil
	} else {
		return "", error
	}
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
	contact := kademlia.Contact{Address: args, ID: nil}
	alive := make(chan bool)
	go context.Network.SendPingMessage(&contact, alive)

	if <-alive {
		fmt.Println("Node is alive")
	} else {
		fmt.Println("Node is dead")
	}
}

func debug_lookupMe(context *kademlia.Kademlia, args string) {
	contacts := context.LookupContact(context.Me.ID)
	fmt.Printf("Recieved %d nodes:\n", len(contacts))
	for _, contact := range contacts {
		fmt.Printf("   %s\n", contact.String())
	}
}

func debug_routingTable(context *kademlia.Kademlia) {
	out := fmt.Sprintf("I am %s\n\n", context.Me.String())
	out += fmt.Sprintf("%d nodes in routingtable:\n", context.Routing.GetNumberOfNodes())
	for _, contact := range context.Routing.Nodes() {
		out += fmt.Sprintf("   %s\n", contact.String())
	}
	fmt.Println(out)
}
