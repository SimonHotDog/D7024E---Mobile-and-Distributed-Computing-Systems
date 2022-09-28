package commands

import (
	"d7024e/kademlia"
	"fmt"
	"strings"
)

func GetAvaliableCommands(context *kademlia.Kademlia, args string) (string, error) {
	var commands = AllCommands()
	var sb strings.Builder

	// Calulate the longest command
	longestLength := 0
	for _, command := range commands {
		commandStructure := fmt.Sprintf("%s %s", command.Name, command.Args)
		if len(commandStructure) > longestLength {
			longestLength = len(commandStructure)
		}
	}

	// Generate string of avaliable commands and their description
	sb.WriteString("Available commands:\n")
	for _, command := range commands {
		commandStructure := fmt.Sprintf("%s %s", command.Name, command.Args)
		line := fmt.Sprintf("  %s%s   %s\n", commandStructure, strings.Repeat(" ", longestLength-len(commandStructure)), command.Description)
		sb.WriteString(line)
	}

	return sb.String(), nil
}
