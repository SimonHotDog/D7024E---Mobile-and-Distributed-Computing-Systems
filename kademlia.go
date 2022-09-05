package d7024e

import (
	"crypto/sha1"
	"encoding/hex"
)

type Kademlia struct {
	routing *RoutingTable
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {

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
