package rpc

import (
	"d7024e/kademlia/network"
	"d7024e/kademlia/network/routing"
)

// Sends a message to the specified contact, instructing it to refresh a
// dataobject with the hash value.
func SendRefreshDataMessage(net network.INetwork, node *routing.Contact, hash string) {
	msg := net.NewNetworkMessage(network.MESSAGE_RPC_DATA_REFRESH, net.GetMe(), node, hash, "", nil)
	net.SendMessage(*msg)
}
