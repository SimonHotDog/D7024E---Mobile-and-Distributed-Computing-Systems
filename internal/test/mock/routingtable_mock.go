package mock

import (
	"d7024e/internal/test/mock/util"
	"d7024e/kademlia/network/routing"

	"github.com/stretchr/testify/mock"
)

// A mock object for the IRoutingTable interface
type RoutingTableMockObject struct {
	mock.Mock
}

func (rt *RoutingTableMockObject) AddContact(contact routing.Contact) {}

func (rt *RoutingTableMockObject) RemoveContact(contactId *routing.KademliaID) {}

func (rt *RoutingTableMockObject) FindClosestContacts(target *routing.KademliaID, count int) []routing.Contact {
	args := rt.Called(target, count)
	return util.GetArrayOrNil[routing.Contact](args, 0)
}

func (rt *RoutingTableMockObject) GetNumberOfNodes() int {
	args := rt.Called()
	return args.Int(0)
}

func (rt *RoutingTableMockObject) Nodes() []routing.Contact {
	args := rt.Called()
	return util.GetArrayOrNil[routing.Contact](args, 0)
}
