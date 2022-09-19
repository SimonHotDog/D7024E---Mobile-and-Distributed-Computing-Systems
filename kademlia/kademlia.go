package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
)

type Kademlia struct {
	Routing *RoutingTable
	Me      *Contact
	Network *Network
}

// Hyperparameters
const K int = 3 //k closest
const A int = 1 //alpha, 1 is effectively no concurrency

//TODO use mutex for concurrency when needed

/*func (kademlia *Kademlia) LookupContact(target *Contact) {

	localKClosestList := kademlia.Routing.FindClosestContacts(target.ID, K) //List of contacts closest to node

	res := make(chan []Contact)

	toSearchList := initSearchList()

	panic("error")

	//alpha criterion
	if len(kClosestList) > A {
		toSearch = append(toSearch, kClosestList[0:A]...)
	}

	//Call recursive lookup function
	kClosestList = kademlia.LookupContactInner(toSearch, searched, kClosestList, target)

}*/

func (kademlia *Kademlia) LookupContact(target *Contact) {

	kClosestTemp := kademlia.Routing.FindClosestContacts(target.ID, K)
	cl := NewCandidateList(target.ID, kClosestTemp)

	channelList := make([]chan string, K)

	//Call recursive lookup

}

func (kademlia *Kademlia) LookupContactInner(target *Contact, cl *CandidateList, channelList *[]chan string, msg string) {

	nodesChecked := 0
	ids := []*KademliaID{}

	i := 0
	kademlia.LookupContactInnerHelper(target, cl, channelList, msg, ids, i, nodesChecked)

}

func (kademlia *Kademlia) LookupContactInnerHelper(target *Contact, cl *CandidateList, channelList *[]chan string, msg string, ids []*KademliaID, i int, nodesChecked int) (int, []*KademliaID) {
	if !(i < cl.Len() && nodesChecked < A) {
		return nodesChecked, ids
	}

	cl.candidates[i].checked = true
	rpc := kademlia.Network.SendFindContactMessage() //Todo check how to send an rpc with message code
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
