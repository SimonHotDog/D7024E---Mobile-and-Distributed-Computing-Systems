package commands

import (
	"d7024e/internal/test/mock"
	"d7024e/kademlia/network/routing"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebugRoutingtable(t *testing.T) {
	lookupId := routing.NewKademliaID("0000000000000000000000000000000000000FFF")
	me := routing.NewContact(lookupId, "localhost:8000")
	expectedContactInRoutingtable := routing.NewContact(routing.NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002")

	networkMock := new(mock.NetworkMockObject)
	routingMock := new(mock.RoutingTableMockObject)
	kademliaMock := new(mock.KademliaMockObject)
	kademliaMock.On("GetMe").Return(&me)
	kademliaMock.On("GetNetwork").Return(networkMock)
	networkMock.On("GetRoutingTable").Return(routingMock)
	routingMock.On("GetNumberOfNodes").Return(3)
	routingMock.On("Nodes").Return([]routing.Contact{
		routing.NewContact(routing.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"),
		routing.NewContact(routing.NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"),
		expectedContactInRoutingtable,
	})

	actual, _ := Debug_routingTable(kademliaMock, "")

	assert.Contains(t, actual, me.String())
	assert.Contains(t, actual, " nodes in routingtable")
	assert.Contains(t, actual, expectedContactInRoutingtable.String())
}
