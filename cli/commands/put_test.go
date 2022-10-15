package commands

import (
	mocks "d7024e/internal/test/mock"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPutObjectInStore(t *testing.T) {
	expectedHash := "myhash"
	kademliaMock := new(mocks.KademliaMockObject)
	kademliaMock.On("Store", mock.Anything).Return(expectedHash, nil)

	actual, err := PutObjectInStore(kademliaMock, "myhash")
	assert.Equal(t, expectedHash, actual)
	assert.Nil(t, err)
}

func TestPutObjectInStoreShouldReturnError(t *testing.T) {
	kademliaMock := new(mocks.KademliaMockObject)
	kademliaMock.On("Store", mock.Anything).Return("", errors.New("error"))

	_, err := PutObjectInStore(kademliaMock, "")
	assert.NotNil(t, err)
}
