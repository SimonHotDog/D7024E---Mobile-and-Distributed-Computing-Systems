package commands

import (
	mocks "d7024e/internal/test/mock"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetObjectByHash(t *testing.T) {
	testObj := new(mocks.KademliaMockObject)
	testObj.On("LookupData", mock.Anything).Return([]byte("test"), nil)

	str, _ := GetObjectByHash(testObj, "test")
	assert.Equal(t, "test", str)
}

func TestGetObjectByHash_DataNoteFound(t *testing.T) {
	testObj := new(mocks.KademliaMockObject)
	testObj.On("LookupData", mock.Anything).Return(nil, nil)

	_, err := GetObjectByHash(testObj, "test")
	assert.NotNil(t, err)
}
