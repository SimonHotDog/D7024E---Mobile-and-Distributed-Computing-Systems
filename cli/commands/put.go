package commands

import (
	"d7024e/kademlia"
)

func PutObjectInStore(context kademlia.IKademlia, args string) (string, error) {
	cleanContent := RemoveDoubleQuotes(args)
	dataToSend := []byte(cleanContent)

	value, err := context.Store(dataToSend)

	if err == nil {
		return value, nil
	} else {
		return "", err
	}
}
