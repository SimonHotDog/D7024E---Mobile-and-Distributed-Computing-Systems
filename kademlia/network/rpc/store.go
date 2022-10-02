package rpc

import (
	"d7024e/kademlia/network"
	"d7024e/kademlia/network/routing"
	"log"
	"strconv"
)

// Send a store command to store data
//
// If store is succesful, success will be true. Otherwise, false.
func SendStoreMessage(net network.INetwork, contact *routing.Contact, hash string, data []byte) (succes bool) {

	msg := net.NewNetworkMessage(network.MESSAGE_RPC_STORE, net.GetMe(), contact, hash, string(data), nil)

	response, timeout := net.SendMessageWithResponse(*msg)

	if timeout {
		log.Printf("Store timeout: %s\n", contact.String())
		return false
	} else {
		succes, _ := strconv.ParseBool(response.Body)
		return succes
	}
}
