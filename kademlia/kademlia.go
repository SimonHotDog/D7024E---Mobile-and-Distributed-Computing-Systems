package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
)

type Kademlia struct {
	Routing *RoutingTable
	Me      *Contact
	Network *Network
	Data    map[string][]byte
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) []byte {
	return kademlia.Data[hash]
}

func (kademlia *Kademlia) Store(data []byte) {
	if kademlia.Data == nil {
		kademlia.Data = make(map[string][]byte)
	}

	kademlia.Data[Hash(data)] = data
}

// Hashes data and returns key
func Hash(data []byte) string {
	sha1 := sha1.Sum([]byte(data))
	key := hex.EncodeToString(sha1[:])

	return key
}

// Finds k closest nodes
func KClosest(key string) {
	// TODO
	// Implement a algorithm for finding
}
