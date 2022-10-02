package commands

import (
	"d7024e/internal/test/mock"
	"d7024e/kademlia/network/routing"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebugLookupMe(t *testing.T) {
	lookupId := routing.NewKademliaID("0000000000000000000000000000000000000FFF")
	expectedContact := routing.NewContact(lookupId, "localhost:8080")
	testObj := new(mock.KademliaMockObject)
	testObj.On("GetMe").Return(&expectedContact)
	testObj.On("LookupContact", lookupId).Return([]routing.Contact{expectedContact})

	actual, _ := Debug_lookupMe(testObj, "")

	assert.Contains(t, actual, "Recieved 1 nodes")
	assert.Contains(t, actual, lookupId.String())
}
