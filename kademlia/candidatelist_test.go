package kademlia

import (
	"testing"
)

func TestNewCandidateList(t *testing.T) {
	testname := "Create new CandidateList"
	t.Run(testname, func(t *testing.T) {
		// Arrange
		contacts := make([]Contact, 10)
		for i := 0; i < len(contacts); i++ {
			contacts[i] = NewContact(NewRandomKademliaID(), "")
		}
		targetid := NewRandomKademliaID()

		// Act
		cl := NewCandidateList(targetid, contacts)

		// Assert
		if cl.targetID != targetid {
			t.Errorf("Expected targetid %v, got %v", targetid, cl.targetID)
		}
		for i := 0; i < len(contacts[:LIMIT]); i++ { //Check order
			if cl.candidates[i].contact != contacts[i] {
				t.Errorf("Expected contact %v, got %v", contacts[i], cl.candidates[i])
			}
		}
	})
}

func TestLen(t *testing.T) {
	testname := "Check length"
	t.Run(testname, func(t *testing.T) {
		//arrange
		contacts := make([]Contact, 7)
		for i := 0; i < len(contacts); i++ {
			contacts[i] = NewContact(NewRandomKademliaID(), "")
		}
		targetid := NewRandomKademliaID()

		cl := NewCandidateList(targetid, contacts)
		//act

		actual := cl.Len()
		//assert
		if actual != 7 {
			t.Errorf("Expected targetid %v, got %v", 7, actual)
		}
	})
}

func TestCandidateExists(t *testing.T) {
	testname := "Check if candidate exists in candidatelist"
	t.Run(testname, func(t *testing.T) {
		//arrange
		contacts := make([]Contact, 5)
		ids := make([]*KademliaID, 5)

		for i := 0; i < len(contacts); i++ {
			x := NewRandomKademliaID()
			contacts[i] = NewContact(x, "")
			ids[i] = x
		}

		targetid := NewRandomKademliaID()

		cl := NewCandidateList(targetid, contacts)

		//Act

		actual := cl.candidateExists(ids[2])

		if actual != true {
			t.Errorf("Expected targetid %v, got %v", true, actual)
		}

		x := NewRandomKademliaID()
		actual2 := cl.candidateExists(x)

		if actual2 != false {
			t.Errorf("Expected targetid %v, got %v", false, actual2)
		}

	})
}

func TestRemoveFromCandidateList(t *testing.T) {
	testname := "Remove from candidate list"
	t.Run(testname, func(t *testing.T) {
		//Arrange
		contacts := make([]Contact, 6)
		ids := make([]*KademliaID, 6)

		for i := 0; i < len(contacts); i++ {
			x := NewRandomKademliaID()
			contacts[i] = NewContact(x, "")
			ids[i] = x
		}

		targetid := NewRandomKademliaID()

		cl := NewCandidateList(targetid, contacts)

		//Act
		cl.removeFromCandidateList(ids[3])

		//Assert
		actual := cl.candidateExists(ids[3])

		if actual != false {
			t.Errorf("Expected targetid %v, got %v", false, actual)
		}
	})
}

func TestGetCandidateFromID(t *testing.T) {
	testname := "Gets Candidate from ID"
	t.Run(testname, func(t *testing.T) {
		//Arrange
		contacts := make([]Contact, 4)
		ids := make([]*KademliaID, 4)

		for i := 0; i < len(contacts); i++ {
			x := NewRandomKademliaID()
			contacts[i] = NewContact(x, "")
			ids[i] = x
		}

		targetid := NewRandomKademliaID()

		cl := NewCandidateList(targetid, contacts)

		//Act

		actual := cl.getCandidateFromID(ids[3]).contact

		x := NewRandomKademliaID()

		actual2 := cl.getCandidateFromID(x)

		//Assert

		if actual != contacts[3] {
			t.Errorf("Expected targetid %v, got %v", contacts[3], actual)
		}

		if actual2 != nil {
			t.Errorf("Expected targetid %v, got %v", nil, actual2)
		}

	})
}

func TestLess(t *testing.T) {
	testname := "Check less than"
	t.Run(testname, func(t *testing.T) {
		//Arrange
		contacts := make([]Contact, 2)

		for i := 0; i < len(contacts); i++ {
			x := NewRandomKademliaID()
			contacts[i] = NewContact(x, "")
		}

		contacts[0].distance = NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
		contacts[1].distance = NewKademliaID("0000000000000000000000000000000000000000")

		targetid := NewRandomKademliaID()

		cl := NewCandidateList(targetid, contacts)

		//Act

		actual := cl.Less(1, 0)

		//Assert

		if actual != true {
			t.Errorf("Expected targetid %v, got %v", true, actual)
		}

	})
}

//TODO Test add to candidate list
