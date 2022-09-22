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

	//Find the k closest
	closestList := kademlia.Routing.FindClosestContacts(target.ID, K)
	cl := NewCandidateList(target.ID, closestList)
	checkedcl := CandidateList{}

	var chans [A]chan []Contact //https://stackoverflow.com/questions/2893004/how-to-allocate-an-array-of-channels
	for i := range chans {
		chans[i] = make(chan []Contact)
	}

	j := 0
	for i := 0; i < cl.Len() && j < A; i++ {
		if !cl.candidates[i].checked {
			go kademlia.Network.SendFindContactMessage(&cl.candidates[i].contact, target.ID, chans[j])
			cl.candidates[i].checked = true
			checkedcl.addToCandidateList(&cl.candidates[i].contact)
			j++
		}
	}

	for i := 0; i < len(chans); i++ {
		temp := <-chans[i]
		checkedcl.candidates[i].connectedContacts = temp
	}

}

func (kademlia *Kademlia) LookupContactInner(target *Contact, cl *CandidateList) {

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
