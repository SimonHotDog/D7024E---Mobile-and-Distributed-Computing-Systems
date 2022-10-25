package commands

import (
	"d7024e/kademlia"
	"d7024e/kademlia/network/routing"
	"errors"
	"fmt"
)

func ForgetObjectInStore(context kademlia.IKademlia, args string) (string, error) {
	if args == "" {
		return "", errors.New("expected 1 argument, but got 0")
	}

	cleanHash := RemoveDoubleQuotes(args)
	closestContacts := context.LookupContact(routing.NewKademliaID(cleanHash))
	if closestContacts == nil {
		return "", errors.New("Object not found")
	}
	err := context.ForgetData(cleanHash, closestContacts)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Successfully forgot dataobject with hash %s", cleanHash), nil
}
