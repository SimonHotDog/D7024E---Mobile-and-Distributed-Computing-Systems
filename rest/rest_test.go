package rest

import (
	mocks "d7024e/internal/test/mock"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPutHandle_ShoudlReturnSucces_WhenRequestIsCorrect(t *testing.T) {
	expectedHash := "myhash"
	valueToSend := "my message"
	reqBody := strings.NewReader(fmt.Sprintf("message=%s", valueToSend))
	req := httptest.NewRequest(http.MethodPut, "/objects", reqBody)
	w := httptest.NewRecorder()

	kademliaMock := new(mocks.KademliaMockObject)
	kademliaMock.On("Store", mock.Anything).Return(expectedHash, nil)

	context = kademliaMock
	putHandle(w, req)

	res := w.Result()
	defer res.Body.Close()
	resBody, _ := io.ReadAll(res.Body)

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Contains(t, string(resBody), "Data")
	assert.Contains(t, string(resBody), valueToSend)
	assert.Contains(t, string(resBody), "Location")
	assert.Contains(t, string(resBody), expectedHash)
}

func TestPutHandle_ShouldReturnError_WhenPutObjectInStoreFails(t *testing.T) {
	reqBody := strings.NewReader("message=")
	req := httptest.NewRequest(http.MethodPut, "/objects", reqBody)
	w := httptest.NewRecorder()

	kademliaMock := new(mocks.KademliaMockObject)
	kademliaMock.On("Store", mock.Anything).Return("", errors.New(""))

	context = kademliaMock
	putHandle(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestPutHandle_ShouldReturnError_WhenHttpMethodIsNotPut(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/objects", nil)
	w := httptest.NewRecorder()

	putHandle(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
}

func TestGetHandle_ShouldReturnSuccess_WhenObjectIsFound(t *testing.T) {
	expectedData := "my message"
	dataHash := "myhash"
	url := fmt.Sprintf("/objects/%s", dataHash)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	kademliaMock := new(mocks.KademliaMockObject)
	kademliaMock.On("LookupData", mock.Anything).Return([]byte(expectedData), nil)

	context = kademliaMock
	getHandle(w, req)

	res := w.Result()
	defer res.Body.Close()
	resBody, _ := io.ReadAll(res.Body)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, expectedData, string(resBody))
}

func TestGetHandle_ShouldReturnError_WhenObjectLookupFailed(t *testing.T) {
	dataHash := "myhash"
	url := fmt.Sprintf("/objects/%s", dataHash)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	kademliaMock := new(mocks.KademliaMockObject)
	kademliaMock.On("LookupData", mock.Anything).Return(nil, nil)

	context = kademliaMock
	getHandle(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestGetHandle_ShouldReturnError_WhenHttpMethodIsNotGet(t *testing.T) {
	dataHash := "myhash"
	url := fmt.Sprintf("/objects/%s", dataHash)
	req := httptest.NewRequest(http.MethodPost, url, nil)
	w := httptest.NewRecorder()

	getHandle(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
}

func TestHomepage_ShouldReturnSuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	homePage(w, req)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
