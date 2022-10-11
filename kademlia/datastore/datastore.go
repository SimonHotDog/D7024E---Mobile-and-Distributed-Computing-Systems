package datastore

import (
	"d7024e/util"
	"log"
	"sync"
	"time"
)

type IDataStore interface {
	// Get a dataobject from the datastore by a given key.
	//
	// Parameters:
	// 	`key` - The key to search for.
	//
	// Returns:
	// 	The datavalue associated with the key and exists will be true.
	// 	Otherwise value is nil and exists will be false.
	Get(key string) (value []byte, exists bool)

	// Add a new dataobject to the datastore. Only one dataobject can be
	// associated with a key and adding a new dataobject with an existing key
	// will not add the value.
	//
	// Parameters:
	// 	`key` - The key to add.
	// 	`value` - The value to add.
	//
	// Returns:
	// 	True if the key was added. Otherwise false.
	Set(key string, value []byte) (ok bool)

	// Remove a key/value pair from the datastore.
	//
	// Parameters:
	// 	`key` - The key to remove.
	//
	// Returns:
	// 	If an dataobject with the specified key exists, the datavalue is returned and ok is true.
	// 	Otherwise, datavalue is nil and ok is false.
	Remove(key string) (value []byte, ok bool)

	// Remove all expired dataobjects from the datastore.
	RemoveExpired()

	// Refresh the TTL of a dataobject associated with a given key.
	//
	// Parameters:
	// 	`key` - The key to refresh.
	//
	// Returns:
	// 	True if the dataobject was found and refreshed. Otherwise, false.
	Refresh(key string) (ok bool)
}

type keyValuePair struct {
	Key   string
	Value []byte
}

type dataObject struct {
	Expiration time.Time
	Value      []byte
}

type DataStore struct {
	janitor           *Janitor
	defaultExpiration time.Duration
	onExpired         func(key string, value []byte)
	dataobjects       map[string]dataObject
	lock              sync.Mutex
	time              util.ITimeProvider
}

// Create a new datastore.
//
// Parameters:
//
//	`ttl` - The default expiration time for dataobjects.
//	`onExpired` - A function that is called when a dataobject expires.
//	`timeprovider` - A timeprovider that is used to get the current time.
func NewDataStore(ttl time.Duration, onExpired func(key string, value []byte), timeprovider util.ITimeProvider) *DataStore {
	datastore := new(DataStore)
	datastore.defaultExpiration = ttl
	datastore.onExpired = onExpired
	datastore.dataobjects = make(map[string]dataObject)
	datastore.time = timeprovider

	runJanitor(datastore, ttl)

	return datastore
}

func (store *DataStore) Get(key string) (value []byte, exists bool) {
	log.Println("Datastore: served dataobject", key)

	store.lock.Lock()
	defer store.lock.Unlock()

	dataObject, exists := store.dataobjects[key]
	if !exists {
		return nil, false
	}

	if dataObject.IsExpired(store.time) {
		store._remove(key)
		return nil, false
	}

	dataObject.Refresh(store.defaultExpiration, store.time)
	store.dataobjects[key] = dataObject

	return dataObject.Value, exists
}

func (store *DataStore) Set(key string, value []byte) (ok bool) {
	log.Println("Datastore: added dataobject", key)

	store.lock.Lock()
	defer store.lock.Unlock()

	dataobject, exists := store.dataobjects[key]
	if exists && !dataobject.IsExpired(store.time) {
		return false
	}

	store.dataobjects[key] = dataObject{
		Expiration: store.time.Now().Add(store.defaultExpiration),
		Value:      value,
	}

	return true
}

func (store *DataStore) Remove(key string) (value []byte, ok bool) {
	log.Println("Datastore: removed dataobject", key)

	store.lock.Lock()
	defer store.lock.Unlock()

	return store._remove(key)
}

func (store *DataStore) _remove(key string) (value []byte, ok bool) {
	dataobject, exists := store.dataobjects[key]
	if exists {
		delete(store.dataobjects, key)
		return dataobject.Value, true
	}
	return nil, false
}

func (store *DataStore) RemoveExpired() {
	var removedItems []keyValuePair

	store.lock.Lock()
	for key, dataobject := range store.dataobjects {
		if dataobject.IsExpired(store.time) {
			value, evicted := store._remove(key)
			if evicted {
				removedItems = append(removedItems, keyValuePair{key, value})
			}
		}
	}
	store.lock.Unlock()

	if store.onExpired == nil {
		return
	}

	for _, item := range removedItems {
		store.onExpired(item.Key, item.Value)
	}
}

func (store *DataStore) Refresh(key string) (ok bool) {
	log.Println("Datastore: refreshed dataobject", key)

	store.lock.Lock()
	defer store.lock.Unlock()

	return store._refresh(key)
}

func (store *DataStore) _refresh(key string) (ok bool) {
	dataObject, exists := store.dataobjects[key]
	if !exists {
		return false
	}

	if dataObject.IsExpired(store.time) {
		store._remove(key)
		return false
	}

	dataObject.Refresh(store.defaultExpiration, store.time)
	store.dataobjects[key] = dataObject

	return true
}

func (object *dataObject) Refresh(expirationTime time.Duration, timeProvider util.ITimeProvider) {
	object.Expiration = timeProvider.Now().Add(expirationTime)
}

func (object *dataObject) IsExpired(timeProvider util.ITimeProvider) (isExpired bool) {
	return timeProvider.Now().After(object.Expiration)
}
