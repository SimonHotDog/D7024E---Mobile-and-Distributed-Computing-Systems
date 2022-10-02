package network

import (
	"d7024e/kademlia/network/routing"

	cmap "github.com/orcaman/concurrent-map/v2"
)

func CreateTestNetwork(port int) (*Network, *routing.Contact) {
	datastore := cmap.New[[]byte]()
	network, me := NewNetwork(port, &datastore)
	return network, me
}
