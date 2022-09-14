package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
)

type Kademlia struct {
	routing *RoutingTable
}

type CandidateList struct {
	l []Candidate
}

type Candidate struct {
	contact Contact
	checked bool
}

// Hyperparameters
const K int = 3 //k closest
const A int = 1 //alpha, 1 is effectively no concurrency

//TODO use mutex for concurrency when needed

func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {

	//List setups.
	localKClosest := kademlia.routing.FindClosestContacts(target.ID, K) //List of contacts closest to node
	kClosestList := make([]Contact, 0)                                  //List of k closest
	kClosestList = append(kClosestList, localKClosest...)
	toSearch := make([]Contact, 0) //List of nodes to be looked up
	toSearch = append(toSearch, kClosestList...)
	searched := make([]Contact, 0) //List of looked up nodes

	//alpha criterion
	if len(kClosestList) > A {
		toSearch = append(toSearch, kClosestList[0:A]...)
	}

	//Call recursive lookup function
	kClosestList = kademlia.LookupContactInner(toSearch, searched, kClosestList, target)

	return kClosestList
}

// Recursive inner fucntion for node lookup
func (Kademlia *Kademlia) LookupContactInner(toSearch []Contact, searched []Contact, kClosestList []Contact, target *Contact) []Contact {
	panic("Not yet implemented")
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	key := Hash(data)

	target := NewKademliaID(key)
	_ = target // To supress error while unfinished
}

// Hashes data and returns key
func Hash(data []byte) string {
	sha1 := sha1.Sum([]byte(data))
	key := hex.EncodeToString(sha1[:])

	return key
}
