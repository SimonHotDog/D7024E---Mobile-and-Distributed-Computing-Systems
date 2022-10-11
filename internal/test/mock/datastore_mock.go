package mock

import (
	"d7024e/internal/test/mock/util"

	"github.com/stretchr/testify/mock"
)

// A mock object for the INetwork interface
type DataStoreMockObject struct {
	mock.Mock
}

func (store *DataStoreMockObject) Get(key string) (value []byte, exists bool) {
	args := store.Called()
	return util.GetArrayOrNil[byte](args, 0), args.Bool(1)
}

func (store *DataStoreMockObject) Set(key string, value []byte) (ok bool) {
	args := store.Called()
	return args.Bool(0)
}

func (store *DataStoreMockObject) Remove(key string) (value []byte, ok bool) {
	args := store.Called()
	return util.GetArrayOrNil[byte](args, 0), args.Bool(1)
}

func (store *DataStoreMockObject) RemoveExpired() {}

func (store *DataStoreMockObject) Refresh(key string) (ok bool) {
	args := store.Called()
	return args.Bool(0)
}
