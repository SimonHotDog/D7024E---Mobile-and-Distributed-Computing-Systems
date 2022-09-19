package kademlia

import "sort"

const LIMIT = 8 //size of the list

type CandidateList struct {
	candidates       [LIMIT]*Candidate //limit indicates the list limit
	closestCandidate *Contact
	targetID         *KademliaID
}

type Candidate struct {
	contact Contact
	checked bool
}

func (cl *CandidateList) addToCandidateList(c *Contact) {
	if cl.candidateExists(c) {
		return
	}

	c.CalcDistance(cl.targetID)

	if cl.Len() == LIMIT {
		if c.Less(&cl.candidates[LIMIT-1].contact) {
			cl.candidates[LIMIT-1] = &Candidate{*c, true}
		}
	} else {
		for i := 0; i < len(cl.candidates); i++ {
			if cl.candidates[i] == nil {
				cl.candidates[i] = &Candidate{*c, true}
			}
		}
	}

	sort.Sort(cl)

}

//removed candidate from candidate list
func (cl *CandidateList) removeFromCandidateList(c *Contact) {
	for i := 0; i < len(cl.candidates); i++ {
		if (cl.candidates[i] != nil) && cl.candidates[i].contact.ID.Equals(c.ID) {
			cl.candidates[i] = nil
		}
	}
}

// Checks if candidate exists
func (cl *CandidateList) candidateExists(c *Contact) bool {
	for _, candidate := range cl.candidates {
		if (candidate != nil) && candidate.contact.ID.Equals(c.ID) {
			return true
		}
	}

	return false
}

func initCandidateList(targetID *KademliaID, c []Contact) *CandidateList {
	cl := &CandidateList{}
	cl.closestCandidate = &c[0]
	cl.targetID = targetID

	for i, co := range c {
		cl.candidates[i] = &Candidate{co, false}
	}

	return cl
}

//Below Required functions for sort interfacing

//gets length of list
func (cl *CandidateList) Len() int {
	l := 0

	for _, candidate := range cl.candidates {
		if candidate != nil {
			l++
		}
	}

	return l
}

//Checks less contacts
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

//Swaps 2 elements
func (cl *CandidateList) Swap(a, b int) {
	cl.candidates[a], cl.candidates[b] = cl.candidates[b], cl.candidates[a]
}
