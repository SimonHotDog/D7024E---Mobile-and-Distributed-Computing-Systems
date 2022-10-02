package commands

import (
	mocks "d7024e/internal/test/mock"
	"d7024e/kademlia/network"
	"d7024e/kademlia/network/routing"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSendPingToAlive(t *testing.T) {
	me := &routing.Contact{}
	networkMock := new(mocks.NetworkMockObject)
	kademliaMock := new(mocks.KademliaMockObject)
	kademliaMock.On("GetNetwork").Return(networkMock)
	networkMock.On("GetMe").Return(me)
	networkMock.On("NewNetworkMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&network.NetworkMessage{})
	networkMock.On("SendMessageWithResponse", network.NetworkMessage{}).Return(network.NetworkMessage{}, false)

	actual, _ := Debug_sendPing(kademliaMock, "")
	assert.Contains(t, actual, "alive")
}

func TestSendPingToDead(t *testing.T) {
	me := &routing.Contact{}
	networkMock := new(mocks.NetworkMockObject)
	kademliaMock := new(mocks.KademliaMockObject)
	kademliaMock.On("GetNetwork").Return(networkMock)
	networkMock.On("GetMe").Return(me)
	networkMock.On("NewNetworkMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&network.NetworkMessage{})
	networkMock.On("SendMessageWithResponse", network.NetworkMessage{}).Return(network.NetworkMessage{}, true)

	actual, _ := Debug_sendPing(kademliaMock, "")
	assert.Contains(t, actual, "dead")
}
