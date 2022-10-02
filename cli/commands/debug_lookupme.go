package commands

import (
	"d7024e/kademlia"
	"fmt"
	"strings"
)

func Debug_lookupMe(context kademlia.IKademlia, args string) (string, error) {
	var sb strings.Builder

	contacts := context.LookupContact(context.GetMe().ID)
	sb.WriteString(fmt.Sprintf("Recieved %d nodes:\n", len(contacts)))
	for _, contact := range contacts {
		sb.WriteString(fmt.Sprintf("   %s\n", contact.String()))
	}
	return sb.String(), nil
}
