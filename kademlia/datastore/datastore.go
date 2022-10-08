package datastore

import "sync"

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

	// Add a new dataobject to the datastore.
	//
	// Parameters:
	// 	`key` - The key to add.
	// 	`value` - The value to add.
	//
	// Returns:
	// 	True if the key was added, false if the key already exists.
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

	// Refresh the TTL of a dataobject associated with a given key.
	//
	// Parameters:
	// 	`key` - The key to refresh.
	//
	// Returns:
	// 	True if the dataobject was found and refreshed. Otherwise, false.
	Refresh(key string) (ok bool)
}

type DataStore struct {
	// Time in seconds before removing a dataobject
	ttl int

	dataobjects map[string]string
	lock        sync.RWMutex
}

func NewDataStore(ttl int) *DataStore {
	return &DataStore{
		ttl:         ttl,
		dataobjects: make(map[string]string),
	}
}

func (store *DataStore) Get(key string) (value []byte, exists bool) {
	// TODO: Implement
	panic("not implemented")
}

func (store *DataStore) Set(key string, value []byte) (ok bool) {
	// TODO: Implement
	panic("not implemented")
}

func (store *DataStore) Remove(key string) (value []byte, ok bool) {
	// TODO: Implement
	panic("not implemented")
}

func (store *DataStore) Refresh(key string) (ok bool) {
	// TODO: Implement
	panic("not implemented")
}
