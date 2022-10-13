package kademlia

import (
	"d7024e/kademlia/network/routing"
	"math"
	"sort"
	"sync"
)

type CandidateList struct {
	// Sorted list of candidates
	candidates []Candidate
	targetID   *routing.KademliaID
	lock       sync.RWMutex
	Limit      int
}

type Candidate struct {
	Contact  routing.Contact
	Checked  bool
	Distance routing.KademliaID
}

func NewCandidateList(targetID *routing.KademliaID, candidateLimit int) *CandidateList {
	// TODO: Remove candidate limit, or improve it
	cl := &CandidateList{
		targetID: targetID,
		Limit:    math.MaxInt,
	}

	return cl
}

func (cl *CandidateList) AddMultiple(contacts []routing.Contact) {
	for _, contact := range contacts {
		cl.Add(contact)
	}
}

func (cl *CandidateList) Add(contact routing.Contact) {
	if cl.Exists(contact.ID) {
		return
	}

	contact.CalcDistance(cl.targetID)
	candidate := Candidate{contact, false, *contact.Distance}

	cl.lock.Lock()
	defer cl.lock.Unlock()
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

func (cl *CandidateList) Get(id *routing.KademliaID) *Candidate {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	for i := 0; i < len(cl.candidates); i++ {
		if cl.candidates[i].Contact.ID.Equals(id) {
			return &cl.candidates[i]
		}
	}
	return nil
}

func (cl *CandidateList) GetAll() []Candidate {
	return cl.candidates
}

// removed candidate from candidate list
func (cl *CandidateList) Remove(id *routing.KademliaID) {
	cl.lock.Lock()
	defer cl.lock.Unlock()

	for i := 0; i < len(cl.candidates); i++ {
		if cl.candidates[i].Contact.ID.Equals(id) {
			cl.candidates = append(cl.candidates[:i], cl.candidates[i+1:]...)
		}
	}
}

// Checks if candidate exists
func (cl *CandidateList) Exists(id *routing.KademliaID) bool {
	cl.lock.RLock()
	defer cl.lock.RUnlock()
	for _, candidate := range cl.candidates {
		if candidate.Contact.ID.Equals(id) {
			return true
		}
	}
	return false
}

// Get length of candidate list
func (cl *CandidateList) Len() int {
	return len(cl.candidates)
}

// Mark candidate as checked
func (cl *CandidateList) Check(id *routing.KademliaID) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	for i, candidate := range cl.candidates {
		if candidate.Contact.ID.Equals(id) {
			cl.candidates[i].Checked = true
		}
	}
}
