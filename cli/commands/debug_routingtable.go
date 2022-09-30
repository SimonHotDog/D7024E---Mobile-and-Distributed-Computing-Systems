package commands

import (
	"d7024e/kademlia"
	"fmt"
	"strings"
)

func Debug_routingTable(context *kademlia.Kademlia, args string) (string, error) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("I am %s\n\n", context.Me.String()))
	sb.WriteString(fmt.Sprintf("%d nodes in routingtable:\n", context.Network.Routingtable.GetNumberOfNodes()))
	for _, contact := range context.Network.Routingtable.Nodes() {
		sb.WriteString(fmt.Sprintf("   %s\n", contact.String()))
	}
	return sb.String(), nil
}
