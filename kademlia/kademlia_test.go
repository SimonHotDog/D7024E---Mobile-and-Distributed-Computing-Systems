package kademlia

import (
	mocks "d7024e/internal/test/mock"
	"d7024e/kademlia/datastore"
	"d7024e/kademlia/network"
	"d7024e/kademlia/network/routing"
	"d7024e/util"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetters(t *testing.T) {
	me := routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000000"), "node0")
	networkMock := new(mocks.NetworkMockObject)
	datastore := datastore.NewDataStore(3600)
	kademlia := NewKademlia(&me, networkMock, datastore)

	assert.Equal(t, kademlia.GetMe(), &me)
	assert.Equal(t, kademlia.GetNetwork(), networkMock)
	assert.Equal(t, kademlia.GetDataStore(), datastore)
}

func TestLookupContact(t *testing.T) {
	/*
		Routes:
		- node A -> [node B]
		- node B -> [node D, node C]
		- node C -> [node D]
		- node D -> []

		Expected result is: [node A, node B, node C, node D]
	*/

	targetId := routing.NewKademliaID("0000000000000000000000000000000000000000")
	me := routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000000"), "node0")

	// Create nodes
	nodeA := routing.NewContact(routing.NewKademliaID("000000000000000000000000000000000000000F"), "nodeA")
	nodeB := routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000007"), "nodeB")
	nodeC := routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000003"), "nodeC")
	nodeD := routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000001"), "nodeD")
	nodeAContacts := []routing.Contact{nodeB}
	nodeBContacts := []routing.Contact{nodeD, nodeC}
	nodeCContacts := []routing.Contact{nodeB}
	nodeDContacts := []routing.Contact{}

	// Create network message to send from nodeA to nodeB
	nodeA_Request := network.NetworkMessage{ID: 2}
	nodeA_Response := network.NetworkMessage{ID: 3, Contacts: nodeAContacts}
	nodeB_Request := network.NetworkMessage{ID: 4}
	nodeB_Response := network.NetworkMessage{ID: 5, Contacts: nodeBContacts}
	nodeC_Request := network.NetworkMessage{ID: 6}
	nodeC_Response := network.NetworkMessage{ID: 7, Contacts: nodeCContacts}
	nodeD_Request := network.NetworkMessage{ID: 8}
	nodeD_Response := network.NetworkMessage{ID: 9, Contacts: nodeDContacts}

	// Setup mocks
	networkMock := new(mocks.NetworkMockObject)
	routingMock := new(mocks.RoutingTableMockObject)
	networkMock.On("GetMe").Return(&me)
	networkMock.On("GetRoutingTable").Return(routingMock)
	routingMock.On("FindClosestContacts", targetId, mock.Anything).Return([]routing.Contact{nodeA})
	networkMock.On("NewNetworkMessage", mock.Anything, mock.Anything, &nodeA, mock.Anything, mock.Anything, mock.Anything).Return(&nodeA_Request)
	networkMock.On("NewNetworkMessage", mock.Anything, mock.Anything, &nodeB, mock.Anything, mock.Anything, mock.Anything).Return(&nodeB_Request)
	networkMock.On("NewNetworkMessage", mock.Anything, mock.Anything, &nodeC, mock.Anything, mock.Anything, mock.Anything).Return(&nodeC_Request)
	networkMock.On("NewNetworkMessage", mock.Anything, mock.Anything, &nodeD, mock.Anything, mock.Anything, mock.Anything).Return(&nodeD_Request)
	networkMock.On("SendMessageWithResponse", nodeA_Request).Return(nodeA_Response, false)
	networkMock.On("SendMessageWithResponse", nodeB_Request).Return(nodeB_Response, false)
	networkMock.On("SendMessageWithResponse", nodeC_Request).Return(nodeC_Response, false)
	networkMock.On("SendMessageWithResponse", nodeD_Request).Return(nodeD_Response, false)

	expected := []routing.Contact{nodeD, nodeC, nodeB, nodeA}

	// Run test
	kademlia := NewKademlia(&me, networkMock, nil)
	actual := kademlia.LookupContact(targetId)

	assert.Equal(t, len(expected), len(actual))
	for i := 0; i < len(expected); i++ {
		assert.Equal(t, expected[i].ID, actual[i].ID)
	}
}

func TestLookupContactAux(t *testing.T) {
	/*
		Routes:
		- node A -> [node B]
		- node B -> []

		Node A is already checked, so we expect node A to be returned
	*/

	targetId := routing.NewKademliaID("0000000000000000000000000000000000000000")
	me := routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000000"), "node0")

	// Create nodes
	nodeA := routing.NewContact(routing.NewKademliaID("000000000000000000000000000000000000001F"), "nodeA")
	nodeB := routing.NewContact(routing.NewKademliaID("000000000000000000000000000000000000000F"), "nodeB")
	nodeAContacts := []routing.Contact{nodeB}
	nodeBContacts := []routing.Contact{}

	// Create and populate candidatelist
	candidateList := NewCandidateList(targetId, math.MaxInt)
	candidateList.Add(nodeA)
	candidateList.Check(nodeA.ID)

	// Create network message to send from nodeA to nodeB
	nodeA_Request := network.NetworkMessage{ID: 2}
	nodeA_Response := network.NetworkMessage{ID: 3, Contacts: nodeAContacts}
	nodeB_Request := network.NetworkMessage{ID: 4}
	nodeB_Response := network.NetworkMessage{ID: 5, Contacts: nodeBContacts}

	// Setup mocks
	networkMock := new(mocks.NetworkMockObject)
	routingMock := new(mocks.RoutingTableMockObject)
	networkMock.On("GetMe").Return(&me)
	networkMock.On("GetRoutingTable").Return(routingMock)
	routingMock.On("FindClosestContacts", targetId, mock.Anything).Return([]routing.Contact{nodeA})
	networkMock.On("NewNetworkMessage", mock.Anything, mock.Anything, &nodeA, mock.Anything, mock.Anything, mock.Anything).Return(&nodeA_Request)
	networkMock.On("NewNetworkMessage", mock.Anything, mock.Anything, &nodeB, mock.Anything, mock.Anything, mock.Anything).Return(&nodeB_Request)
	networkMock.On("SendMessageWithResponse", nodeA_Request).Return(nodeA_Response, false)
	networkMock.On("SendMessageWithResponse", nodeB_Request).Return(nodeB_Response, false)

	expected := []routing.Contact{nodeA}

	// Run test
	kademlia := NewKademlia(&me, networkMock, nil)
	kademlia.lookupContactAux(targetId, []routing.Contact{nodeA}, candidateList)
	actual := candidateList.GetAll()

	assert.Equal(t, len(expected), len(actual))
	for i := 0; i < len(expected); i++ {
		assert.Equal(t, expected[i].ID, actual[i].Contact.ID)
	}
}

func TestLookupDataSucces(t *testing.T) {
	expectedData := "data"
	expectedDataHash := util.Hash([]byte(expectedData))
	me := routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000000"), "node0")

	// Create nodes
	nodeA := routing.NewContact(routing.NewKademliaID("000000000000000000000000000000000000000F"), "nodeA")

	// Create network messages
	nodeAFindNode_Request := network.NetworkMessage{ID: 1}
	nodeAFindNode_Response := network.NetworkMessage{ID: 2, Contacts: []routing.Contact{}}
	nodeAFindValue_Request := network.NetworkMessage{ID: 3}
	nodeAFindValue_Response := network.NetworkMessage{ID: 4, Body: expectedData}

	// Setup mocks
	networkMock := new(mocks.NetworkMockObject)
	routingMock := new(mocks.RoutingTableMockObject)
	networkMock.On("GetMe").Return(&me)
	networkMock.On("GetRoutingTable").Return(routingMock)
	routingMock.On("FindClosestContacts", mock.Anything, mock.Anything).Return([]routing.Contact{nodeA})
	networkMock.On("NewNetworkMessage", network.MESSAGE_RPC_FIND_NODE, mock.Anything, mock.Anything, mock.Anything, expectedDataHash, mock.Anything).Return(&nodeAFindNode_Request)   // FindNode
	networkMock.On("NewNetworkMessage", network.MESSAGE_RPC_FIND_VALUE, mock.Anything, mock.Anything, expectedDataHash, mock.Anything, mock.Anything).Return(&nodeAFindValue_Request) // FindValue
	networkMock.On("SendMessageWithResponse", nodeAFindNode_Request).Return(nodeAFindNode_Response, false)
	networkMock.On("SendMessageWithResponse", nodeAFindValue_Request).Return(nodeAFindValue_Response, false)

	// Run test
	kademlia := NewKademlia(&me, networkMock, nil)
	actualData, actualContact := kademlia.LookupData(expectedDataHash)

	assert.Equal(t, expectedData, string(actualData))
	assert.Equal(t, nodeA.ID, actualContact.ID)
}

func TestLookupDataTimeout(t *testing.T) {
	t.Skip("Not implemented")
}

func TestLookupDataNotFound(t *testing.T) {
	var expectedData []byte = nil
	var expectedContact *routing.Contact = nil

	requestedData := "data"
	requestedDataHash := util.Hash([]byte(requestedData))
	me := routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000000"), "node0")

	// Create nodes
	nodeA := routing.NewContact(routing.NewKademliaID("000000000000000000000000000000000000000F"), "nodeA")

	// Create network messages
	nodeAFindNode_Request := network.NetworkMessage{ID: 1}
	nodeAFindNode_Response := network.NetworkMessage{ID: 2, Contacts: []routing.Contact{}}
	nodeAFindValue_Request := network.NetworkMessage{ID: 3}
	nodeAFindValue_Response := network.NetworkMessage{ID: 4, Body: ""}

	// Setup mocks
	networkMock := new(mocks.NetworkMockObject)
	routingMock := new(mocks.RoutingTableMockObject)
	networkMock.On("GetMe").Return(&me)
	networkMock.On("GetRoutingTable").Return(routingMock)
	routingMock.On("FindClosestContacts", mock.Anything, mock.Anything).Return([]routing.Contact{nodeA})
	networkMock.On("NewNetworkMessage", network.MESSAGE_RPC_FIND_NODE, mock.Anything, mock.Anything, mock.Anything, requestedDataHash, mock.Anything).Return(&nodeAFindNode_Request)   // FindNode
	networkMock.On("NewNetworkMessage", network.MESSAGE_RPC_FIND_VALUE, mock.Anything, mock.Anything, requestedDataHash, mock.Anything, mock.Anything).Return(&nodeAFindValue_Request) // FindValue
	networkMock.On("SendMessageWithResponse", nodeAFindNode_Request).Return(nodeAFindNode_Response, false)
	networkMock.On("SendMessageWithResponse", nodeAFindValue_Request).Return(nodeAFindValue_Response, false)

	// Run test
	kademlia := NewKademlia(&me, networkMock, nil)
	actualData, actualContact := kademlia.LookupData(requestedDataHash)

	assert.Equal(t, expectedData, actualData)
	assert.Equal(t, expectedContact, actualContact)
}

func TestStore(t *testing.T) {
	t.Skip("Not implemented")
}

func TestJoinNetwork(t *testing.T) {
	t.Skip("Not implemented")
}
