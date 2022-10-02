package mock

import (
	"d7024e/internal/test/mock/util"
	"d7024e/kademlia/network"
	"d7024e/kademlia/network/routing"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/stretchr/testify/mock"
)

// A mock object for the INetwork interface
type NetworkMockObject struct {
	mock.Mock
}

func (net *NetworkMockObject) GetMe() *routing.Contact {
	args := net.Called()
	return util.GetPointerOrNil[routing.Contact](args, 0)
}

func (net *NetworkMockObject) GetRoutingTable() routing.IRoutingTable {
	args := net.Called()
	return util.GetPointerOrNil[RoutingTableMockObject](args, 0)
}

func (net *NetworkMockObject) GetDatastore() *cmap.ConcurrentMap[[]byte] {
	args := net.Called()
	return util.GetPointerOrNil[cmap.ConcurrentMap[[]byte]](args, 0)
}

func (net *NetworkMockObject) NewNetworkMessage(
	rpc int,
	sender *routing.Contact,
	target *routing.Contact,
	bodyDigest string,
	body string,
	contacts []routing.Contact,
) *network.NetworkMessage {
	args := net.Called(rpc, sender, target, bodyDigest, body, contacts)
	return util.GetPointerOrNil[network.NetworkMessage](args, 0)
}

func (net *NetworkMockObject) Listen() {}

func (net *NetworkMockObject) StopListen() {}

func (net *NetworkMockObject) SendMessageWithResponse(msg network.NetworkMessage) (response network.NetworkMessage, timeout bool) {
	args := net.Called(msg)
	return args.Get(0).(network.NetworkMessage), args.Bool(1)
}

func (net *NetworkMockObject) SendMessage(msg network.NetworkMessage) {}
