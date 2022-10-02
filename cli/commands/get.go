package commands

import (
	"d7024e/kademlia"
	"errors"
)

func GetObjectByHash(context kademlia.IKademlia, args string) (string, error) {
	cleanHash := RemoveDoubleQuotes(args)
	value, _ := context.LookupData(cleanHash)
	if value != nil {
		return string(value), nil
	} else {
		return "", errors.New("data not found")
	}
}
