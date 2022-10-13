package network

import (
	"d7024e/kademlia/network/routing"
	"testing"
)

func TestFlipSenderTarget(t *testing.T) {
	testname := "Flip sender and target in networkmessage"
	t.Run(testname, func(t *testing.T) {
		network, me := CreateTestNetwork(14041)
		expectedSender := *me
		expectedTarget := routing.NewContact(routing.NewRandomKademliaID(), "127.0.0.1")
		netmsg := NetworkMessage{
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
