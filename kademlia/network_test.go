package kademlia

import (
	"fmt"
	"strings"
	"testing"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
)

func createTestNetwork(port int) (*Network, *Contact) {
	address := fmt.Sprintf("127.0.0.1:%d", port)
	me := NewContact(NewRandomKademliaID(), address)
	context := Kademlia{Me: &me, Routing: NewRoutingTable(me), DataStore: cmap.New[[]byte]()}
	network := CreateNewNetwork(&context, port)
	return &network, &me
}

func TestFlipSenderTarget(t *testing.T) {
	testname := "Flip sender and target in networkmessage"
	t.Run(testname, func(t *testing.T) {
		network, me := createTestNetwork(14041)
		expectedSender := *me
		expectedTarget := NewContact(NewRandomKademliaID(), "127.0.0.1")
		netmsg := NetworkMessage{
			ID:     0,
			RPC:    1,
			Sender: &expectedTarget,
			Target: &expectedSender,
		}

		network.generateReturnMessage(&netmsg)
		actualSender := *netmsg.Sender
		actualTarget := *netmsg.Target

		if actualSender != expectedSender || actualTarget != expectedTarget {
			t.Errorf("Expected from %s to %s, got from %s to %s", &expectedSender, &expectedTarget, &actualSender, &actualTarget)
		}
	})
}

func TestPingMessage(t *testing.T) {
	testname := "Ping myself"
	t.Run(testname, func(t *testing.T) {
		expected := true
		alive := make(chan bool, 1)
		networkA, _ := createTestNetwork(14041)
		networkB, _ := createTestNetwork(14048)

		go networkA.Listen()
		go networkB.Listen()
		defer networkA.StopListen()
		defer networkB.StopListen()
		time.Sleep(20 * time.Millisecond)
		networkA.SendPingMessage(networkB.Kademlia.Me, alive)
		actual := <-alive

		if actual != expected {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
	})
}

func TestPingMessageFailure(t *testing.T) {
	testname := "Fail to ping myself"
	t.Run(testname, func(t *testing.T) {
		expected := false
		alive := make(chan bool, 1)
		network, me := createTestNetwork(14041)

		network.SendPingMessage(me, alive)
		actual := <-alive

		if actual != expected {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
	})
}

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
			contacts := make(chan []Contact, 1)
			networkA, _ := createTestNetwork(14041)
			networkB, _ := createTestNetwork(14048)
			for i := 0; i < test.nContactsInRoutingTable; i++ {
				networkB.Kademlia.Routing.AddContact(NewContact(NewRandomKademliaID(), ""))
			}

			go networkA.Listen()
			go networkB.Listen()
			defer networkA.StopListen()
			defer networkB.StopListen()

			time.Sleep(100 * time.Millisecond)

			networkA.SendFindContactMessage(networkB.Kademlia.Me, networkA.Kademlia.Me.ID, contacts)
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
		contacts := make(chan []Contact, 1)
		network, _ := createTestNetwork(14041)

		time.Sleep(100 * time.Millisecond)

		network.SendFindContactMessage(network.Kademlia.Me, network.Kademlia.Me.ID, contacts)
		actual := <-contacts

		if len(actual) != expectedNContacts {
			t.Errorf("Expected %v, got %v", expectedNContacts, actual)
		}
	})
}

func TestStoreMessage(t *testing.T) {
	testname := "Send store command"
	t.Run(testname, func(t *testing.T) {
		expected := "My Message"
		networkA, _ := createTestNetwork(14041)
		networkB, _ := createTestNetwork(14048)
		messageBytes := []byte(expected)
		messageHash := Hash([]byte(messageBytes))

		go networkA.Listen()
		go networkB.Listen()
		defer networkA.StopListen()
		defer networkB.StopListen()
		time.Sleep(20 * time.Millisecond)

		networkA.SendStoreMessage(networkB.Kademlia.Me, messageHash, messageBytes)
		storedMessage, _ := networkB.Kademlia.DataStore.Get(messageHash)
		actual := string(storedMessage)

		if strings.Compare(actual, expected) != 0 {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
	})
}

func TestStoreMessageTimeout(t *testing.T) {
	testname := "Send store command and wait for timeout"
	t.Run(testname, func(t *testing.T) {
		expected := false
		networkA, _ := createTestNetwork(14041)
		networkB, _ := createTestNetwork(14048)
		messageBytes := []byte("My Message")
		messageHash := Hash([]byte(messageBytes))

		go networkA.Listen()
		defer networkA.StopListen()
		time.Sleep(20 * time.Millisecond)

		actual := networkA.SendStoreMessage(networkB.Kademlia.Me, messageHash, messageBytes)

		if actual != expected {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
	})
}

func TestLookupMessage(t *testing.T) {
	testname := "Send lookup command"
	t.Run(testname, func(t *testing.T) {
		expected := "My Message"
		valueChannel := make(chan string, 1)
		networkA, _ := createTestNetwork(14041)
		networkB, _ := createTestNetwork(14048)
		messageBytes := []byte(expected)
		messageHash := Hash([]byte(messageBytes))
		networkB.Kademlia.DataStore.Set(messageHash, messageBytes)

		go networkA.Listen()
		go networkB.Listen()
		defer networkA.StopListen()
		defer networkB.StopListen()
		time.Sleep(20 * time.Millisecond)

		networkA.SendLookupMessage(networkB.Kademlia.Me, messageHash, valueChannel)
		actual := <-valueChannel

		if strings.Compare(actual, expected) != 0 {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
	})
}

func TestLookupMessageTimeout(t *testing.T) {
	testname := "Send lookup command and wait for timeout"
	t.Run(testname, func(t *testing.T) {
		expected := NETWORK_REQUEST_TIMEOUT_STRING
		valueChannel := make(chan string, 1)
		networkA, _ := createTestNetwork(14041)
		networkB, _ := createTestNetwork(14048)
		messageBytes := []byte("Test")
		messageHash := Hash([]byte(messageBytes))
		networkB.Kademlia.DataStore.Set(messageHash, messageBytes)

		go networkA.Listen()
		defer networkA.StopListen()
		time.Sleep(20 * time.Millisecond)

		networkA.SendLookupMessage(networkB.Kademlia.Me, messageHash, valueChannel)
		actual := <-valueChannel

		if strings.Compare(actual, expected) != 0 {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
	})
}

func TestStoreBeforeLookupMessage(t *testing.T) {
	testname := "Send store message before lookup"
	t.Run(testname, func(t *testing.T) {
		expected := "My Message"
		valueChannel := make(chan string, 1)
		networkA, _ := createTestNetwork(14041)
		networkB, _ := createTestNetwork(14048)
		messageBytes := []byte(expected)
		messageHash := Hash([]byte(messageBytes))

		go networkA.Listen()
		go networkB.Listen()
		defer networkA.StopListen()
		defer networkB.StopListen()
		time.Sleep(20 * time.Millisecond)

		storeOk := networkA.SendStoreMessage(networkB.Kademlia.Me, messageHash, messageBytes)
		if !storeOk {
			t.Errorf("Expected store to succeed")
			return
		}
		networkA.SendLookupMessage(networkB.Kademlia.Me, messageHash, valueChannel)
		actual := <-valueChannel

		if strings.Compare(actual, expected) != 0 {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
	})
}
