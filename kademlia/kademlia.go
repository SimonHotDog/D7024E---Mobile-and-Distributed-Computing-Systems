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

func (kademlia *Kademlia) LookupContact(targetID KademliaID) *Contact {
	//Find the k closest
	closestList := kademlia.Routing.FindClosestContacts(&targetID, K)
	cl := NewCandidateList(&targetID, closestList)

	if cl.candidateExists(&targetID) {
		return &cl.getCandidateFromID(&targetID).contact
	}

	// TODO: check closest

	checkedcl := CandidateList{}

	var chans [A]chan []Contact //https://stackoverflow.com/questions/2893004/how-to-allocate-an-array-of-channels
	for i := range chans {
		chans[i] = make(chan []Contact)
	}

	j := 0
	for i := 0; i < cl.Len() && j < A; i++ {
		if !cl.candidates[i].checked {
			go kademlia.Network.SendFindContactMessage(&cl.candidates[i].contact, &targetID, chans[j])
			cl.candidates[i].checked = true
			checkedcl.addToCandidateList(&cl.candidates[i].contact)
			j++
		}
	}

	for i := 0; i < len(chans); i++ {
		temp := <-chans[i]
		temp2 := NewCandidateList(&targetID, temp)
		checkedcl.candidates[i].connectedContacts = temp2
	}

	return nil //???
}

// Recursive step
func (kademlia *Kademlia) LookupContactInner(targetID *KademliaID, c *Candidate) {
	if c.connectedContacts.candidateExists(targetID) {
		//Return whatever
	}

	var chans [A]chan []Contact
	for i := range chans {
		chans[i] = make(chan []Contact)
	}

	checkedcl := CandidateList{}

	j := 0
	for i := 0; i < c.connectedContacts.Len() && j < A; i++ {
		if !c.connectedContacts.candidates[i].checked {
			go kademlia.Network.SendFindContactMessage(&c.connectedContacts.candidates[i].contact, targetID, chans[j])
			c.connectedContacts.candidates[i].checked = true
			checkedcl.addToCandidateList(&c.connectedContacts.candidates[i].contact)
			j++
		}
	}

	for i := 0; i < len(chans); i++ {
		temp := <-chans[i]
		temp2 := NewCandidateList(targetID, temp)
		checkedcl.candidates[i].connectedContacts = temp2
	}

	for i := 0; i < checkedcl.Len(); i++ {
		kademlia.LookupContactInner(targetID, checkedcl.candidates[i])
	}
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
