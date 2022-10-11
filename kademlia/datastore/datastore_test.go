package datastore

import (
	"d7024e/util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func emptyOnExpired(key string, value []byte) {}

func createNewDatastore(ttl time.Duration, currentTime time.Time) *DataStore {
	timeProvider := &util.FakeTimeProvider{InternalTime: currentTime}
	return NewDataStore(ttl, emptyOnExpired, timeProvider)
}

func TestNewDataStore(t *testing.T) {
	ttl := time.Hour
	timeProvider := &util.FakeTimeProvider{}
	datastore := NewDataStore(ttl, emptyOnExpired, timeProvider)

	assert.NotNil(t, datastore)
	assert.NotNil(t, datastore.janitor)
	assert.Equal(t, ttl, datastore.defaultExpiration)
	assert.NotNil(t, datastore.onExpired)
	assert.NotNil(t, datastore.dataobjects)
	assert.NotNil(t, datastore.time)
}

func TestDataStore_Get(t *testing.T) {
	expiredDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	currentDate := time.Date(1971, 1, 1, 0, 0, 0, 0, time.UTC)
	futureDate := time.Date(1972, 1, 1, 0, 0, 0, 0, time.UTC)

	dataObjects := map[string]dataObject{
		"key1": {Value: []byte("value1"), Expiration: futureDate},
		"key2": {Value: []byte("value1"), Expiration: futureDate},
		"key3": {Value: []byte("value1"), Expiration: currentDate},
		"key4": {Value: []byte("value1"), Expiration: expiredDate},
	}

	expectedValidKeys := []string{"key1", "key2", "key3"}
	expectedExpiredKeys := []string{"key4"}
	dataStore := createNewDatastore(time.Hour, currentDate)
	dataStore.dataobjects = dataObjects

	for _, key := range expectedValidKeys {
		value, ok := dataStore.Get(key)
		assert.True(t, ok)
		assert.Equal(t, dataObjects[key].Value, value)
	}
	for _, key := range expectedExpiredKeys {
		value, ok := dataStore.Get(key)
		assert.False(t, ok)
		assert.Nil(t, value)
	}
}

func TestDataStore_Set(t *testing.T) {
	currentDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedExpirationDate := currentDate.Add(time.Hour)
	expectedValidKeys := []string{"key1", "key2", "key3"}
	expectedString := "test"
	expectedObjectValue := []byte(expectedString)
	dataStore := createNewDatastore(time.Hour, currentDate)

	for _, key := range expectedValidKeys {
		dataStore.Set(key, expectedObjectValue)
	}

	for key, dataobject := range dataStore.dataobjects {
		assert.Contains(t, expectedValidKeys, key)
		assert.Equal(t, expectedObjectValue, dataobject.Value)
		assert.Equal(t, expectedExpirationDate, dataobject.Expiration)
	}
}

func TestDataStore_Set_WithExistingKey_WhenNotExpired_ShouldNotReplaceExisting(t *testing.T) {
	currentDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	futureDate := currentDate.Add(time.Minute)
	dataStore := createNewDatastore(time.Hour, currentDate)

	keyToAdd := "key1"
	dataObjects := map[string]dataObject{
		keyToAdd: {Value: []byte("value1"), Expiration: futureDate},
	}
	dataStore.dataobjects = dataObjects

	actualSetReturn := dataStore.Set(keyToAdd, []byte("test"))

	assert.False(t, actualSetReturn)
}

func TestDataStore_Set_WithExistingKey_WhenExpired_ShouldReplaceExisting(t *testing.T) {
	currentDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	expiredDate := currentDate.Add(-1 * time.Hour)
	dataStore := createNewDatastore(time.Hour, currentDate)

	keyToAdd := "key1"
	dataObjects := map[string]dataObject{
		keyToAdd: {Value: []byte("value1"), Expiration: expiredDate},
	}
	dataStore.dataobjects = dataObjects

	actualSetReturn := dataStore.Set(keyToAdd, []byte(""))

	assert.True(t, actualSetReturn)
}

func TestDataStore_Remove(t *testing.T) {
	var tests = []struct {
		keyToRemove         string
		expectedOutputOk    bool
		expectedOutputValue []byte
		expectedKeysInStore []string
	}{
		{"key1", true, []byte("value1"), []string{"key2"}},
		{"key4", false, nil, []string{"key1", "key2"}},
	}

	for _, test := range tests {
		t.Run(test.keyToRemove, func(t *testing.T) {
			currentDate := time.Date(1971, 1, 1, 0, 0, 0, 0, time.UTC)
			futureDate := time.Date(1972, 1, 1, 0, 0, 0, 0, time.UTC)

			dataStore := createNewDatastore(time.Hour, currentDate)
			dataStore.dataobjects = map[string]dataObject{
				"key1": {Value: []byte("value1"), Expiration: futureDate},
				"key2": {Value: []byte("value1"), Expiration: currentDate},
			}

			actualValue, actualOk := dataStore.Remove(test.keyToRemove)

			assert.Equal(t, test.expectedOutputOk, actualOk)
			assert.Equal(t, test.expectedOutputValue, actualValue)
			for key, _ := range dataStore.dataobjects {
				assert.Contains(t, test.expectedKeysInStore, key)
			}
		})
	}
}

func TestDataStore_RemoveExpired(t *testing.T) {
	expiredDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	currentDate := time.Date(1971, 1, 1, 0, 0, 0, 0, time.UTC)
	futureDate := time.Date(1972, 1, 1, 0, 0, 0, 0, time.UTC)

	dataObjects := map[string]dataObject{
		"key1": {Value: []byte("value1"), Expiration: futureDate},
		"key2": {Value: []byte("value1"), Expiration: futureDate},
		"key3": {Value: []byte("value1"), Expiration: currentDate},
		"key4": {Value: []byte("value1"), Expiration: expiredDate},
	}
	dataStore := createNewDatastore(time.Hour, currentDate)
	dataStore.dataobjects = dataObjects

	expectedValidKeys := []string{"key1", "key2", "key3"}

	dataStore.RemoveExpired()

	for key, _ := range dataStore.dataobjects {
		assert.Contains(t, expectedValidKeys, key)
	}
}

func TestDataStore_Refresh(t *testing.T) {
	ttl := time.Hour
	expiredDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	currentDate := time.Date(1971, 1, 1, 0, 0, 0, 0, time.UTC)
	futureDate := currentDate.Add(ttl)

	var tests = []struct {
		name                string
		keyToRefresh        string
		expectedKeysInStore []string
		expectedNewDate     time.Time
		expectedOutput      bool
	}{
		{
			name:                "Refresh not expired key",
			keyToRefresh:        "key1",
			expectedKeysInStore: []string{"key1", "key2", "key3", "key4"},
			expectedNewDate:     futureDate,
			expectedOutput:      true,
		},
		{
			name:                "Refresh expired key",
			keyToRefresh:        "key4",
			expectedKeysInStore: []string{"key1", "key2", "key3"},
			expectedNewDate:     futureDate,
			expectedOutput:      false,
		},
		{
			name:                "Refresh non-existing key",
			keyToRefresh:        "key5",
			expectedKeysInStore: []string{"key1", "key2", "key3", "key4"},
			expectedNewDate:     currentDate,
			expectedOutput:      false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dataStore := createNewDatastore(ttl, currentDate)
			dataStore.dataobjects = map[string]dataObject{
				"key1": {Value: []byte("value1"), Expiration: futureDate},
				"key2": {Value: []byte("value1"), Expiration: futureDate},
				"key3": {Value: []byte("value1"), Expiration: currentDate},
				"key4": {Value: []byte("value1"), Expiration: expiredDate},
			}

			actual := dataStore.Refresh(test.keyToRefresh)
			assert.Equal(t, actual, test.expectedOutput)

			refreshedDataObject, exists := dataStore.dataobjects[test.keyToRefresh]
			if exists {
				assert.Equal(t, test.expectedNewDate, refreshedDataObject.Expiration)
			}

			for _, key := range test.expectedKeysInStore {
				_, ok := dataStore.dataobjects[key]
				assert.True(t, ok)
			}
		})
	}
}

func Test_dataObject_Refresh(t *testing.T) {
	expectedTime := time.Date(1970, 1, 1, 1, 0, 0, 0, time.UTC)
	expectedValue := []byte("test")
	timeNow := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	timeProvider := &util.FakeTimeProvider{InternalTime: timeNow}
	dataObject := &dataObject{
		Expiration: timeNow,
		Value:      expectedValue,
	}

	dataObject.Refresh(time.Hour, timeProvider)

	assert.Equal(t, expectedTime, dataObject.Expiration)
	assert.Equal(t, expectedValue, dataObject.Value)
}
