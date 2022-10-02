package rpc

import (
	"d7024e/kademlia/network"
	"d7024e/kademlia/network/routing"
	"log"
)

func SendPingMessage(net network.INetwork, contact *routing.Contact, alive chan bool) {
	msg := net.NewNetworkMessage(network.MESSAGE_RPC_PING, net.GetMe(), contact, "", "", nil)

	_, timeout := net.SendMessageWithResponse(*msg)

	if timeout {
		alive <- false
		log.Printf("Ping timeout: %s\n", contact.String())
	} else {
		alive <- true
	}
}
