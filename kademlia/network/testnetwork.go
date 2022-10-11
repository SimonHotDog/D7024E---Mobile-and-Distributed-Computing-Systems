package network

import (
	"d7024e/kademlia/datastore"
	"d7024e/kademlia/network/routing"
	"d7024e/util"
	"time"
)

func CreateTestNetwork(port int) (*Network, *routing.Contact) {
	timeprovider := &util.TimeProvider{}
	datastore := datastore.NewDataStore(time.Hour, nil, timeprovider)
	network, me := NewNetwork(port, datastore)
	return network, me
}
