package commands

import (
	"d7024e/kademlia"
	"errors"
)

func PutObjectInStore(context kademlia.IKademlia, args string) (string, error) {
	if args == "" {
		return "", errors.New("expected 1 argument, but got 0")
	}

	cleanContent := RemoveDoubleQuotes(args)
	dataToSend := []byte(cleanContent)
	value, err := context.Store(dataToSend)
	if err == nil {
		return value, nil
	} else {
		return "", err
	}
}
