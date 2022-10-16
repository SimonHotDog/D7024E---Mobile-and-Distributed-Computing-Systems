package commands

import "d7024e/kademlia"

type Command struct {
	Name        string
	Args        string
	Description string
	Action      func(context kademlia.IKademlia, args string) (string, error)
}

func AllCommands() []Command {
	return []Command{
		{"forget", "[hash]", "Takes a hash and forgets any dataobject associated with it", ForgetObjectInStore},
		{"get", "[hash]", "Takes a hash and downloads the file from the network.", GetObjectByHash},
		{"help", "", "Help on ", GetAvaliableCommands},
		{"put", "[text]", "Uploads a file to the network and returns the hash if succesful.", PutObjectInStore},
		{"stat", "", "Displays the status of the network.", GetStatus},
		{"exit", "", "Exit the CLI.", ExitApplication},
		{"ping", "[address]", "DEBUG: Send a ping RPC to the target client", Debug_sendPing},
		{"whoami", "", "DEBUG: Lookup myself", Debug_lookupMe},
		{"routes", "", "DEBUG: Print routingtable", Debug_routingTable},
		{"lookup", "", "DEBUG: Lookup", Debug_lookupMe},
	}
}

func RemoveDoubleQuotes(str string) string {
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	return str
}
