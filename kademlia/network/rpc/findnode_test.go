package rpc

import (
	"d7024e/kademlia/network"
	"d7024e/kademlia/network/routing"
	"fmt"
	"testing"
	"time"
)

func TestSendFindNodeMessage(t *testing.T) {
	tests := []struct {
		nContactsInRoutingTable int
		expectedNContacts       int
	}{
		// TODO: Get contact limit from kademlia parameter k
		{0, 1},
		{5, 6},
		{19, 20},
		{23, 20},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("Find %d nodes close to me", test.expectedNContacts)
		t.Run(testname, func(t *testing.T) {
			contacts := make(chan []routing.Contact, 1)
			networkA, _ := network.CreateTestNetwork(14041)
			networkB, _ := network.CreateTestNetwork(14048)
			for i := 0; i < test.nContactsInRoutingTable; i++ {
				networkB.Routingtable.AddContact(routing.NewContact(routing.NewRandomKademliaID(), ""))
			}

			go networkA.Listen()
			go networkB.Listen()
			defer networkA.StopListen()
			defer networkB.StopListen()

			time.Sleep(20 * time.Millisecond)

			SendFindContactMessage(networkA, networkB.Me, networkA.Me.ID, contacts)
			actual := <-contacts

			if len(actual) != test.expectedNContacts {
				t.Errorf("Expected %v, got %v", test.expectedNContacts, actual)
			}
		})
	}
}

func TestSendFindNodeMessageTimout(t *testing.T) {
	expectedNContacts := 0
	testname := "Timeout when waiting for find nodes response"
	t.Run(testname, func(t *testing.T) {
		contacts := make(chan []routing.Contact, 1)
		network, _ := network.CreateTestNetwork(14041)

		time.Sleep(20 * time.Millisecond)

		SendFindContactMessage(network, network.Me, network.Me.ID, contacts)
		actual := <-contacts

		if len(actual) != expectedNContacts {
			t.Errorf("Expected %v, got %v", expectedNContacts, actual)
		}
	})
}
