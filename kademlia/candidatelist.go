package kademlia

import (
	"sort"
	"sync"
)

type CandidateList struct {
	// Sorted list of candidates
	candidates []Candidate
	targetID   *KademliaID
	lock       sync.RWMutex
	Limit      int
}

type Candidate struct {
	Contact  Contact
	Checked  bool
	Distance KademliaID
}

func NewCandidateList(targetID *KademliaID, candidateLimit int) *CandidateList {
	cl := &CandidateList{
		targetID: targetID,
		Limit:    candidateLimit,
	}

	return cl
}

func (cl *CandidateList) AddMultiple(contacts []Contact) {
	for _, contact := range contacts {
		cl.Add(contact)
	}
}

func (cl *CandidateList) Add(contact Contact) {
	if cl.Exists(contact.ID) {
		return
	}

	contact.CalcDistance(cl.targetID)
	candidate := Candidate{contact, true, *contact.distance}

	cl.lock.RLock()
	defer cl.lock.RUnlock()
	if len(cl.candidates) == cl.Limit {
		if contact.Less(&cl.candidates[cl.Limit-1].Contact) {
			cl.candidates[cl.Limit-1] = candidate
		}
	} else {
		cl.candidates = append(cl.candidates, candidate)
	}

	sort.Slice(cl.candidates, func(i, j int) bool {
		return cl.candidates[i].Distance.Less(&cl.candidates[j].Distance)
	})
}

func (cl *CandidateList) Get(id *KademliaID) *Candidate {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	for i := 0; i < len(cl.candidates); i++ {
		if cl.candidates[i].Contact.ID == id {
			return &cl.candidates[i]
		}
	}
	return nil
}

func (cl *CandidateList) GetAll() []Candidate {
	return cl.candidates
}

// removed candidate from candidate list
func (cl *CandidateList) Remove(id *KademliaID) {
	cl.lock.RLock()
	for i := 0; i < len(cl.candidates); i++ {
		if cl.candidates[i].Contact.ID.Equals(id) {
			cl.candidates = append(cl.candidates[:i], cl.candidates[i+1:]...)
		}
	}
	cl.lock.RUnlock()
}

// Checks if candidate exists
func (cl *CandidateList) Exists(id *KademliaID) bool {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	for _, candidate := range cl.candidates {
		if candidate.Contact.ID.Equals(id) {
			return true
		}
	}
	return false
}

func (cl *CandidateList) Len() int {
	return len(cl.candidates)
}

func (cl *CandidateList) Check(id *KademliaID) {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	for i, candidate := range cl.candidates {
		if candidate.Contact.ID.Equals(id) {
			cl.candidates[i].Checked = true
		}
	}
}
