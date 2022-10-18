package network

import (
	"d7024e/kademlia/network/routing"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlipSenderTarget(t *testing.T) {
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

	assert.Equal(t, expectedSender, actualSender)
	assert.Equal(t, expectedTarget, actualTarget)
}

func TestSendMessageWithResponse_OnTimeout_TargetNodeShouldBeRemoved(t *testing.T) {
	targetNode := routing.NewContact(routing.NewRandomKademliaID(), ":14041")
	networkA, _ := CreateTestNetwork(14041)
	networkA.GetRoutingTable().AddContact(targetNode)
	msgToSend := *networkA.NewNetworkMessage(MESSAGE_RPC_PING, networkA.GetMe(), &targetNode, "", "", nil)

	_, timeout := networkA.SendMessageWithResponse(msgToSend)
	nodesInRoutingTable := networkA.GetRoutingTable().Nodes()

	assert.True(t, timeout, "Expected timeout to be true")
	assert.NotContains(t, nodesInRoutingTable, targetNode)
}
