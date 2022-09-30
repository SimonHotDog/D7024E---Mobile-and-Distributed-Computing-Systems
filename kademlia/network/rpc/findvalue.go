package rpc

import (
	"d7024e/kademlia/network"
	"d7024e/kademlia/network/routing"
	"log"
)

// Send Lookup command to find data
//
// If lookup is succesful, the found value will be added to the value channel.
// Otherwise, a "::timeout::" string-value will be added to the value channel.
func SendLookupMessage(net *network.Network, contact *routing.Contact, hash string, value chan string) {
	msg := net.NewNetworkMessage(network.MESSAGE_RPC_FIND_VALUE, net.Me, contact, hash, "", nil)

	response, timeout := net.SendMessageWithResponse(*msg)

	if timeout {
		value <- network.NETWORK_REQUEST_TIMEOUT_STRING
		log.Printf("Lookup timeout: %s\n", contact.String())
	} else {
		value <- response.Body
	}
}
