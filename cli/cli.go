package cli

import (
	"bufio"
	"d7024e/cli/commands"
	"d7024e/kademlia"
	"fmt"
	"io"
	"strings"
)

type Cli struct {
	out     io.Writer
	in      io.Reader
	context kademlia.IKademlia
}

// Create a new CLI
//
// Parameters:
//
//	`out`: The output stream all commands will write to
//	`in`: The input stream the CLI will read from
//	`context`: The Kadmlia context the CLI will be using
func NewCli(out io.Writer, in io.Reader, context kademlia.IKademlia) *Cli {
	return &Cli{
		out:     out,
		in:      in,
		context: context,
	}
}

// Open the CLI and start reading input
//
// Parameters:
//
//	`isUserFriendly`: If true, the CLI will print user friendly messages such as greetings. Otherwise only output from commands will be printed.
func (cli *Cli) Open(isUserFriendly bool) {
	if isUserFriendly {
		fmt.Fprintln(cli.out, "Hello! This is your CLI speaking!")
		fmt.Fprintln(cli.out, "Type 'help' for a list of commands.")
	}

	commands := commands.AllCommands()
	reader := bufio.NewReader(cli.in)

	for {
		if isUserFriendly {
			fmt.Fprint(cli.out, ">>> ")
		}

		input, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			panic(fmt.Sprintf("Error reading input: %v\n", err))
		}

		commandName, commandArgs := splitCommand(input)
		if commandName == "" {
			continue
		}

		commandFound := false
		for _, command := range commands {
			if command.Name == commandName {
				result, err := command.Action(cli.context, commandArgs)
				if err != nil {
					fmt.Fprintf(cli.out, "Error: %s\n", err)
				} else {
					fmt.Fprintln(cli.out, result)
				}
				commandFound = true
				break
			}
		}
		if err == io.EOF {
			return
		}

		if !commandFound {
			fmt.Fprintln(cli.out, "Command not found. Type 'help' for a list of commands")
		}
	}
}

func splitCommand(input string) (name, args string) {
	input = strings.TrimRight(input, "\r\n")
	cmd := strings.SplitN(input, " ", 2)

	if len(cmd) == 2 {
		return cmd[0], cmd[1]
	}
	return cmd[0], ""
}
