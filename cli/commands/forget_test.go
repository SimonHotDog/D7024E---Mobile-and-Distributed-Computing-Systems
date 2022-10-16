package commands

import (
	mocks "d7024e/internal/test/mock"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestForgetObjectInStore(t *testing.T) {
	expectedHash := "myhash"
	kademliaMock := new(mocks.KademliaMockObject)
	kademliaMock.On("LookupContact", mock.Anything).Return(nil)
	kademliaMock.On("ForgetData", mock.Anything).Return(nil)

	_, err := ForgetObjectInStore(kademliaMock, expectedHash)
	assert.Nil(t, err)
}

func TestForgetObjectInStore_WithNoArgs_ShouldReturnError(t *testing.T) {
	kademliaMock := new(mocks.KademliaMockObject)

	_, err := ForgetObjectInStore(kademliaMock, "")
	assert.NotNil(t, err)
}
