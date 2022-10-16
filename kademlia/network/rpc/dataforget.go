package rpc

import (
	"d7024e/kademlia/network"
	"d7024e/kademlia/network/routing"
)

// Sends a message to the specified contact, instructing it to forget a
// dataobject with the hash value.
func SendForgetDataMessage(net network.INetwork, node *routing.Contact, hash string) {
	msg := net.NewNetworkMessage(network.MESSAGE_RPC_DATA_FORGET, net.GetMe(), node, hash, "", nil)
	net.SendMessage(*msg)
}
