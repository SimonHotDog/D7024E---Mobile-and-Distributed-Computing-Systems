package network

import (
	"d7024e/kademlia/datastore"
	"d7024e/kademlia/network/routing"
)

func CreateTestNetwork(port int) (*Network, *routing.Contact) {
	datastore := datastore.NewDataStore(3600)
	network, me := NewNetwork(port, datastore)
	return network, me
}
