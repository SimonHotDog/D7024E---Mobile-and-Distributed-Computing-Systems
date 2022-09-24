package kademlia

import (
	"sort"
)

const LIMIT = 8 //size of the list

type CandidateList struct {
	candidates [8]*Candidate //limit indicates the list limit. HARDCODED! FIX LATER!
	targetID   *KademliaID
}

type Candidate struct {
	contact           Contact //itself
	checked           bool
	connectedContacts *CandidateList
}

func (cl *CandidateList) addToCandidateList(c *Contact) {
	if cl.candidateExists(c.ID) {
		return
	}

	c.CalcDistance(cl.targetID)

	if cl.Len() == LIMIT {
		if c.Less(&cl.candidates[LIMIT-1].contact) {
			cl.candidates[LIMIT-1] = &Candidate{*c, true, nil}
		}
	} else {
		for i := 0; i < len(cl.candidates); i++ {
			if cl.candidates[i] == nil {
				cl.candidates[i] = &Candidate{*c, true, nil}
			}
		}
	}

	sort.Sort(cl)

}

func (cl *CandidateList) getCandidateFromID(id *KademliaID) *Candidate {
	for i := 0; i < cl.Len(); i++ {
		if cl.candidates[i].contact.ID == id {
			return cl.candidates[i]
		}
	}

	return nil //TODO. Is this really ok?
}

// removed candidate from candidate list
func (cl *CandidateList) removeFromCandidateList(id *KademliaID) {
	for i := 0; i < len(cl.candidates); i++ {
		if (cl.candidates[i] != nil) && cl.candidates[i].contact.ID.Equals(id) {
			cl.candidates[i] = nil
		}
	}
}

// Checks if candidate exists
func (cl *CandidateList) candidateExists(id *KademliaID) bool {
	for _, candidate := range cl.candidates {
		if (candidate != nil) && candidate.contact.ID.Equals(id) {
			return true
		}
	}

	return false
}

func NewCandidateList(targetID *KademliaID, candidates []Contact) *CandidateList {
	cl := &CandidateList{}
	//cl.closestCandidate = &candidates[0]
	cl.targetID = targetID

	if len(candidates) > LIMIT {
		candidates = candidates[:LIMIT]
	}

	for i, contact := range candidates {
		cl.candidates[i] = &Candidate{contact, false, nil}
	}
	return cl
}

//Below Required functions for sort interfacing

// gets length of list
func (cl *CandidateList) Len() int {
	l := 0

	for _, candidate := range cl.candidates {
		if candidate != nil {
			l++
		}
	}

	return l
}

// Checks less contacts
// a and b are indexes of candidates in candidate list
func (cl *CandidateList) Less(a, b int) bool {

	//Check if either element is null
	if cl.candidates[b] == nil {
		return true
	}

	if cl.candidates[a] == nil {
		return false
	}

	//Call contacts Less function
	return cl.candidates[a].contact.Less(&cl.candidates[b].contact)
}

// Swaps 2 elements
func (cl *CandidateList) Swap(a, b int) {
	cl.candidates[a], cl.candidates[b] = cl.candidates[b], cl.candidates[a]
}

func (cl *CandidateList) GetLimit() int {
	return LIMIT
}
