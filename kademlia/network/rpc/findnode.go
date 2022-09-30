package rpc

import (
	"d7024e/kademlia/network"
	"d7024e/kademlia/network/routing"
	"log"
)

// Send a message to the specified contact.
//
// If the contact responds, the returned contacts will added to the contacts channel.
// Otherwise, an empty array will be added to the contacts channel.
func SendFindContactMessage(net *network.Network, contact *routing.Contact, id *routing.KademliaID, contacts chan []routing.Contact) {
	msg := net.NewNetworkMessage(network.MESSAGE_RPC_FIND_NODE, net.Me, contact, "", id.String(), nil)

	response, timeout := net.SendMessageWithResponse(*msg)

	if timeout {
		contacts <- make([]routing.Contact, 0)
		log.Printf("Find contact timeout: %s\n", contact.String())
	} else {
		contacts <- response.Contacts
	}
}
