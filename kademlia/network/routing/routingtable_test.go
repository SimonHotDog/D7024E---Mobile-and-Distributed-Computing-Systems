package routing

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutingTable(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002"))

	contacts := rt.FindClosestContacts(NewKademliaID("2111111400000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())
	}
}

func TestAddMeContact(t *testing.T) {
	testname := "Me should not be added to routingtable"
	t.Run(testname, func(t *testing.T) {
		expected := 0
		me := NewContact(NewRandomKademliaID(), "")
		rt := NewRoutingTable(me)

		rt.AddContact(me)

		actual := rt.GetNumberOfNodes()

		if actual != expected {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
	})
}

func TestNodes(t *testing.T) {
	me := NewContact(NewKademliaID("ABC0000000000000000000000000000000000000"), "me")
	nodeA := NewContact(NewKademliaID("000000000000000000000000000000000000000F"), "nodeA")
	nodeB := NewContact(NewKademliaID("0000000000000000000000000000000000000007"), "nodeB")
	nodeC := NewContact(NewKademliaID("0000000000000000000000000000000000000003"), "nodeC")
	nodeD := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "nodeD")

	expected := []Contact{nodeD, nodeC, nodeB, nodeA}

	rt := NewRoutingTable(me)
	rt.AddContact(nodeA)
	rt.AddContact(nodeB)
	rt.AddContact(nodeC)
	rt.AddContact(nodeD)
	actual := rt.Nodes()

	for _, contact := range expected {
		assert.Contains(t, actual, contact)
	}
}

func TestNodesWhenEmpty(t *testing.T) {
	me := NewContact(NewKademliaID("ABC0000000000000000000000000000000000000"), "me")

	expectedNodesLen := 0

	rt := NewRoutingTable(me)
	actualNodesLen := len(rt.Nodes())

	assert.Equal(t, expectedNodesLen, actualNodesLen)
}
