package commands

import (
	"d7024e/kademlia"
	"errors"
)

func GetObjectByHash(context kademlia.IKademlia, args string) (string, error) {
	if args == "" {
		return "", errors.New("expected 1 argument, but got 0")
	}
	cleanHash := RemoveDoubleQuotes(args)
	value, _ := context.LookupData(cleanHash)
	if value != nil {
		return string(value), nil
	} else {
		return "", errors.New("data not found")
	}
}
