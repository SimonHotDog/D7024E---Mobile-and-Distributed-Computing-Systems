package commands

import (
	"d7024e/kademlia"
	"fmt"
	"strings"
)

func Debug_routingTable(context kademlia.IKademlia, args string) (string, error) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("I am %s\n\n", context.GetMe().String()))
	sb.WriteString(fmt.Sprintf("%d nodes in routingtable:\n", context.GetNetwork().GetRoutingTable().GetNumberOfNodes()))
	for _, contact := range context.GetNetwork().GetRoutingTable().Nodes() {
		sb.WriteString(fmt.Sprintf("   %s\n", contact.String()))
	}
	return sb.String(), nil
}
