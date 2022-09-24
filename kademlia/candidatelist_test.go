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
		cl := NewCandidateList(targetid, 8)
		cl.AddMultiple(contacts)

		// Assert
		if cl.targetID != targetid {
			t.Errorf("Expected targetid %v, got %v", targetid, cl.targetID)
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

		cl := NewCandidateList(targetid, 8)
		cl.AddMultiple(contacts)

		//Act

		actual := cl.Exists(ids[2])

		if actual != true {
			t.Errorf("Expected targetid %v, got %v", true, actual)
		}

		x := NewRandomKademliaID()
		actual2 := cl.Exists(x)

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

		cl := NewCandidateList(targetid, 8)
		cl.AddMultiple(contacts)

		//Act
		cl.Remove(ids[3])

		//Assert
		actual := cl.Exists(ids[3])

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
		for i := 0; i < len(contacts); i++ {
			contacts[i] = NewContact(NewRandomKademliaID(), "")
		}

		targetid := NewRandomKademliaID()
		cl := NewCandidateList(targetid, 8)
		cl.AddMultiple(contacts)

		//Act
		actual := cl.Get(contacts[3].ID).Contact

		//Assert
		if !actual.ID.Equals(contacts[3].ID) {
			t.Errorf("Expected targetid %v, got %v", contacts[3].String(), actual.String())
		}
	})
}

func TestAddWhenListIsNotFull(t *testing.T) {
	testname := "Add candidate when list is empty"
	t.Run(testname, func(t *testing.T) {
		//Arrange
		var contacts []Contact
		targetid := NewRandomKademliaID()
		contacts = append(contacts, NewContact(NewRandomKademliaID(), ""))
		contacts = append(contacts, NewContact(NewRandomKademliaID(), ""))

		cl := NewCandidateList(targetid, 8)
		cl.AddMultiple(contacts)

		contactToAdd := NewContact(NewRandomKademliaID(), "")

		//Act
		cl.Add(contactToAdd)
		actual := cl.Get(contactToAdd.ID)

		//Assert
		if actual != nil && actual.Contact.ID.Equals(contactToAdd.ID) == false {
			t.Errorf("Expected contact %v, got %v", contactToAdd, actual.Contact)
		}
	})
}

func TestAddWhenListIsFullAndReplace(t *testing.T) {
	testname := "Add candidate when list is full and no replace"
	t.Run(testname, func(t *testing.T) {
		//Arrange
		var contacts []Contact
		targetid := NewKademliaID("0000000000000000000000000000000000000000")

		contacts = append(contacts, NewContact(NewKademliaID("0000000000000000000000000000000000000002"), ""))
		contacts = append(contacts, NewContact(NewKademliaID("0000000000000000000000000000000000000003"), ""))
		contacts = append(contacts, NewContact(NewKademliaID("0000000000000000000000000000000000000004"), ""))
		contacts = append(contacts, NewContact(NewKademliaID("0000000000000000000000000000000000000005"), ""))
		contacts = append(contacts, NewContact(NewKademliaID("0000000000000000000000000000000000000006"), ""))
		contacts = append(contacts, NewContact(NewKademliaID("0000000000000000000000000000000000000007"), ""))
		contacts = append(contacts, NewContact(NewKademliaID("0000000000000000000000000000000000000008"), ""))
		contacts = append(contacts, NewContact(NewKademliaID("00000000000000000000000000000000000000F0"), ""))

		cl := NewCandidateList(targetid, 8)
		cl.AddMultiple(contacts)

		contactToAdd := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "")

		//Act
		cl.Add(contactToAdd)
		actual := cl.Get(contactToAdd.ID)

		//Assert
		if actual != nil && actual.Contact.ID.Equals(contactToAdd.ID) == false {
			t.Errorf("Expected contact %v, got %v", contactToAdd, actual.Contact)
		}
	})
}

func TestAddWhenListIsFullAndNotReplace(t *testing.T) {
	testname := "Add candidate when list is full and no replace"
	t.Run(testname, func(t *testing.T) {
		//Arrange
		var contacts []Contact
		targetid := NewKademliaID("0000000000000000000000000000000000000000")

		contacts = append(contacts, NewContact(NewKademliaID("0000000000000000000000000000000000000001"), ""))
		contacts = append(contacts, NewContact(NewKademliaID("0000000000000000000000000000000000000002"), ""))
		contacts = append(contacts, NewContact(NewKademliaID("0000000000000000000000000000000000000003"), ""))
		contacts = append(contacts, NewContact(NewKademliaID("0000000000000000000000000000000000000004"), ""))
		contacts = append(contacts, NewContact(NewKademliaID("0000000000000000000000000000000000000005"), ""))
		contacts = append(contacts, NewContact(NewKademliaID("0000000000000000000000000000000000000006"), ""))
		contacts = append(contacts, NewContact(NewKademliaID("0000000000000000000000000000000000000007"), ""))
		contacts = append(contacts, NewContact(NewKademliaID("0000000000000000000000000000000000000008"), ""))

		cl := NewCandidateList(targetid, 8)
		cl.AddMultiple(contacts)

		contactToAdd := NewContact(NewKademliaID("000000000000000000000000000000000000000F"), "")

		//Act
		cl.Add(contactToAdd)
		actual := cl.Get(contactToAdd.ID)

		//Assert
		if actual != nil {
			t.Errorf("Expected %v, got %v", nil, actual.Contact)
		}
	})
}

// func TestLen(t *testing.T) {
// 	testname := "Check length"
// 	t.Run(testname, func(t *testing.T) {
// 		//arrange
// 		contacts := make([]Contact, 7)
// 		for i := 0; i < len(contacts); i++ {
// 			contacts[i] = NewContact(NewRandomKademliaID(), "")
// 		}
// 		targetid := NewRandomKademliaID()

// 		cl := NewCandidateList(targetid, contacts)
// 		//act

// 		actual := cl.Len()
// 		//assert
// 		if actual != 7 {
// 			t.Errorf("Expected targetid %v, got %v", 7, actual)
// 		}
// 	})
// }

// func TestLess(t *testing.T) {
// 	testname := "Check less than"
// 	t.Run(testname, func(t *testing.T) {
// 		//Arrange
// 		contacts := make([]Contact, 2)

// 		for i := 0; i < len(contacts); i++ {
// 			x := NewRandomKademliaID()
// 			contacts[i] = NewContact(x, "")
// 		}

// 		contacts[0].distance = NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
// 		contacts[1].distance = NewKademliaID("0000000000000000000000000000000000000000")

// 		targetid := NewRandomKademliaID()

// 		cl := NewCandidateList(targetid, contacts)

// 		//Act

// 		actual := cl.Less(1, 0)

// 		//Assert

// 		if actual != true {
// 			t.Errorf("Expected targetid %v, got %v", true, actual)
// 		}

// 	})
// }

// func TestSwap(t *testing.T) {
// 	testname := "Swap candidates"
// 	t.Run(testname, func(t *testing.T) {
// 		//Arrange

// 		contacts := make([]Contact, 4)

// 		for i := 0; i < len(contacts); i++ {
// 			x := NewRandomKademliaID()
// 			contacts[i] = NewContact(x, "")
// 		}

// 		id0 := contacts[0].ID
// 		id1 := contacts[1].ID

// 		targetid := NewRandomKademliaID()

// 		cl := NewCandidateList(targetid, contacts)

// 		//Act

// 		cl.Swap(0, 1)

// 		//Assert

// 		actual := cl.candidates[0].Contact.ID == id1 && cl.candidates[1].Contact.ID == id0

// 		if actual != true {
// 			t.Errorf("Expected targetid %v, got %v", true, actual)
// 		}
// 	})
// }
