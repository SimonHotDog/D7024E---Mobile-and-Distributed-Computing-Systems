package kademlia

import (
	mocks "d7024e/internal/test/mock"
	"d7024e/kademlia/datastore"
	"d7024e/kademlia/network"
	"d7024e/kademlia/network/routing"
	"d7024e/util"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetters(t *testing.T) {
	me := routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000000"), "node0")
	networkMock := new(mocks.NetworkMockObject)
	datastore := datastore.NewDataStore(time.Hour, nil, nil)
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
	nodeA_Request := network.NetworkMessage{BodyDigest: "2"}
	nodeA_Response := network.NetworkMessage{BodyDigest: "3", Contacts: nodeAContacts}
	nodeB_Request := network.NetworkMessage{BodyDigest: "4"}
	nodeB_Response := network.NetworkMessage{BodyDigest: "5", Contacts: nodeBContacts}
	nodeC_Request := network.NetworkMessage{BodyDigest: "6"}
	nodeC_Response := network.NetworkMessage{BodyDigest: "7", Contacts: nodeCContacts}
	nodeD_Request := network.NetworkMessage{BodyDigest: "8"}
	nodeD_Response := network.NetworkMessage{BodyDigest: "9", Contacts: nodeDContacts}

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
	nodeA_Request := network.NetworkMessage{BodyDigest: "2"}
	nodeA_Response := network.NetworkMessage{BodyDigest: "3", Contacts: nodeAContacts}
	nodeB_Request := network.NetworkMessage{BodyDigest: "4"}
	nodeB_Response := network.NetworkMessage{BodyDigest: "5", Contacts: nodeBContacts}

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
	nodeAFindNode_Request := network.NetworkMessage{BodyDigest: "1"}
	nodeAFindNode_Response := network.NetworkMessage{BodyDigest: "2", Contacts: []routing.Contact{}}
	nodeAFindValue_Request := network.NetworkMessage{BodyDigest: "3"}
	nodeAFindValue_Response := network.NetworkMessage{BodyDigest: "4", Body: expectedData}
	rpcRefresh_Request := network.NetworkMessage{RPC: network.MESSAGE_RPC_DATA_REFRESH}

	// Setup mocks
	networkMock := new(mocks.NetworkMockObject)
	routingMock := new(mocks.RoutingTableMockObject)
	networkMock.On("GetMe").Return(&me)
	networkMock.On("GetRoutingTable").Return(routingMock)
	routingMock.On("FindClosestContacts", mock.Anything, mock.Anything).Return([]routing.Contact{nodeA})
	networkMock.On("NewNetworkMessage", network.MESSAGE_RPC_FIND_NODE, mock.Anything, mock.Anything, mock.Anything, expectedDataHash, mock.Anything).Return(&nodeAFindNode_Request)   // FindNode
	networkMock.On("NewNetworkMessage", network.MESSAGE_RPC_FIND_VALUE, mock.Anything, mock.Anything, expectedDataHash, mock.Anything, mock.Anything).Return(&nodeAFindValue_Request) // FindValue
	networkMock.On("NewNetworkMessage", network.MESSAGE_RPC_DATA_REFRESH, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&rpcRefresh_Request)
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
	nodeAFindNode_Request := network.NetworkMessage{BodyDigest: "1"}
	nodeAFindNode_Response := network.NetworkMessage{BodyDigest: "2", Contacts: []routing.Contact{}}
	nodeAFindValue_Request := network.NetworkMessage{BodyDigest: "3"}
	nodeAFindValue_Response := network.NetworkMessage{BodyDigest: "4", Body: ""}
	rpcRefresh_Request := network.NetworkMessage{RPC: network.MESSAGE_RPC_DATA_REFRESH}

	// Setup mocks
	networkMock := new(mocks.NetworkMockObject)
	routingMock := new(mocks.RoutingTableMockObject)
	networkMock.On("GetMe").Return(&me)
	networkMock.On("GetRoutingTable").Return(routingMock)
	routingMock.On("FindClosestContacts", mock.Anything, mock.Anything).Return([]routing.Contact{nodeA})
	networkMock.On("NewNetworkMessage", network.MESSAGE_RPC_FIND_NODE, mock.Anything, mock.Anything, mock.Anything, requestedDataHash, mock.Anything).Return(&nodeAFindNode_Request)   // FindNode
	networkMock.On("NewNetworkMessage", network.MESSAGE_RPC_FIND_VALUE, mock.Anything, mock.Anything, requestedDataHash, mock.Anything, mock.Anything).Return(&nodeAFindValue_Request) // FindValue
	networkMock.On("NewNetworkMessage", network.MESSAGE_RPC_DATA_REFRESH, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&rpcRefresh_Request)
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

func TestForgetData_WithHash_ShouldForgetData(t *testing.T) {
	dataHash := util.Hash([]byte("data"))
	me := routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000000"), "node0")
	nodeA := routing.NewContact(routing.NewKademliaID("000000000000000000000000000000000000000F"), "nodeA")
	forgetData_Request := network.NetworkMessage{BodyDigest: "3"}

	// Setup mocks
	networkMock := new(mocks.NetworkMockObject)
	networkMock.On("GetMe").Return(&me)
	networkMock.On("NewNetworkMessage", network.MESSAGE_RPC_DATA_FORGET, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&forgetData_Request)

	// Run test
	kademlia := NewKademlia(&me, networkMock, nil)
	contacts := []routing.Contact{nodeA}
	err := kademlia.ForgetData(dataHash, contacts)
	assert.Nil(t, err)
}

func TestForgetData_WhenHashIsInvalid_ShouldReturnError(t *testing.T) {
	t.Skip("Not implemented")
}

func TestJoinNetwork_WhenNetworkIsEmpty(t *testing.T) {
	me := routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000000"), "me")
	knownNode := routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000002"), "bootstrap")
	noContactsResponse := network.NetworkMessage{Contacts: []routing.Contact{}}

	networkMock := new(mocks.NetworkMockObject)
	networkMock.On("GetMe").Return(&me)
	networkMock.On("NewNetworkMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(new(network.NetworkMessage))
	networkMock.On("SendMessageWithResponse", mock.Anything).Return(noContactsResponse, false)

	kademlia := NewKademlia(&me, networkMock, nil)
	actual := kademlia.JoinNetwork(&knownNode, 1)

	assert.False(t, actual)
}

func TestJoinNetwork_WhenNetworkIsNonEmpty(t *testing.T) {
	me := routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000000"), "me")
	knownNode := routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000002"), "bootstrap")
	nodeA := routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000002"), "nodeA")

	networkMessageResponse := network.NetworkMessage{Contacts: []routing.Contact{nodeA}}

	networkMock := new(mocks.NetworkMockObject)
	routingMock := new(mocks.RoutingTableMockObject)
	networkMock.On("GetMe").Return(&me)
	networkMock.On("GetRoutingTable").Return(routingMock)
	networkMock.On("NewNetworkMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(new(network.NetworkMessage))
	networkMock.On("SendMessageWithResponse", mock.Anything).Return(networkMessageResponse, false)

	kademlia := NewKademlia(&me, networkMock, nil)
	actual := kademlia.JoinNetwork(&knownNode, 1)

	assert.True(t, actual)
}

func TestJoinNetworkAux_WhenSomeNodesRespond(t *testing.T) {
	expectedNumberOfNodes := 2
	expectedNumberOfDeadNodes := 1

	me := routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000000"), "me")
	knownNode := routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000002"), "bootstrap")
	nodeA := routing.NewContact(routing.NewKademliaID("000000000000000000000000000000000000000A"), "nodeA")
	nodeB := routing.NewContact(routing.NewKademliaID("000000000000000000000000000000000000000B"), "nodeB")

	findnode_response := network.NetworkMessage{Contacts: []routing.Contact{nodeA, nodeB}}
	findnode_request := network.NetworkMessage{RPC: network.MESSAGE_RPC_FIND_NODE}
	nodeA_ping_request := network.NetworkMessage{RPC: network.MESSAGE_RPC_PING, Target: &nodeA}
	nodeB_ping_request := network.NetworkMessage{RPC: network.MESSAGE_RPC_PING, Target: &nodeB}

	networkMock := new(mocks.NetworkMockObject)
	routingMock := new(mocks.RoutingTableMockObject)
	networkMock.On("GetMe").Return(&me)
	networkMock.On("GetRoutingTable").Return(routingMock)
	networkMock.On("NewNetworkMessage", network.MESSAGE_RPC_FIND_NODE, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&findnode_request)
	networkMock.On("NewNetworkMessage", network.MESSAGE_RPC_PING, mock.Anything, &nodeA, mock.Anything, mock.Anything, mock.Anything).Return(&nodeA_ping_request)
	networkMock.On("NewNetworkMessage", network.MESSAGE_RPC_PING, mock.Anything, &nodeB, mock.Anything, mock.Anything, mock.Anything).Return(&nodeB_ping_request)
	networkMock.On("SendMessageWithResponse", findnode_request).Return(findnode_response, false)
	networkMock.On("SendMessageWithResponse", nodeA_ping_request).Return(*new(network.NetworkMessage), true)
	networkMock.On("SendMessageWithResponse", nodeB_ping_request).Return(*new(network.NetworkMessage), false)

	kademlia := NewKademlia(&me, networkMock, nil)
	actualNodesRecieved, actualNodesDead := kademlia.joinNetworkAux(&knownNode, 0, 1)

	assert.Equal(t, expectedNumberOfNodes, actualNodesRecieved)
	assert.Equal(t, expectedNumberOfDeadNodes, actualNodesDead)
}

func TestJoinNetworkAux(t *testing.T) {
	t.Skip("Not implemented")
}

func Test_min(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name     string
		args     args
		expected int
	}{
		{name: "a < b", args: args{1, 2}, expected: 1},
		{name: "a == b", args: args{2, 2}, expected: 2},
		{name: "a > b", args: args{3, 2}, expected: 2},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := min(test.args.a, test.args.b)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func Test_getExponentialBackoffTime(t *testing.T) {
	attemptNumber := 3
	expectedAtLeast := 8 * time.Millisecond

	actual := getExponentialBackoffTime(attemptNumber)

	assert.Greater(t, actual, expectedAtLeast) // We can't assert exact value because of randomness
}
