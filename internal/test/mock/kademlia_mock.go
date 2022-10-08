package mock

import (
	"d7024e/internal/test/mock/util"
	"d7024e/kademlia/datastore"
	"d7024e/kademlia/network"
	"d7024e/kademlia/network/routing"

	"github.com/stretchr/testify/mock"
)

// A mock object for the IKademlia interface
type KademliaMockObject struct {
	mock.Mock
}

func (k KademliaMockObject) GetMe() *routing.Contact {
	args := k.Called()
	return util.GetPointerOrNil[routing.Contact](args, 0)
}
func (k KademliaMockObject) GetNetwork() network.INetwork {
	args := k.Called()
	return util.GetPointerOrNil[NetworkMockObject](args, 0)
}
func (k KademliaMockObject) GetDataStore() datastore.IDataStore {
	args := k.Called()
	return util.GetPointerOrNil[DataStoreMockObject](args, 0)
}

func (k KademliaMockObject) LookupContact(targetID *routing.KademliaID) []routing.Contact {
	args := k.Called(targetID)
	return util.GetArrayOrNil[routing.Contact](args, 0)
}

func (k KademliaMockObject) LookupData(hash string) ([]byte, *routing.Contact) {
	args := k.Called(hash)

	data := util.GetArrayOrNil[byte](args, 0)
	contact := util.GetPointerOrNil[routing.Contact](args, 1)
	return data, contact
}

func (k KademliaMockObject) Store(data []byte) (string, error) {
	args := k.Called(data)

	return args.String(0), args.Error(1)
}

func (k KademliaMockObject) JoinNetwork(contact *routing.Contact) {}
