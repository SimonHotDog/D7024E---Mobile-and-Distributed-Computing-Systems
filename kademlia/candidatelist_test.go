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
		wantedId := "F000000000000000000000000000000000000000"
		wamtedContact := NewContact(NewKademliaID(wantedId), "")
		targetid := NewRandomKademliaID()
		cl := NewCandidateList(targetid, 8)
		cl.Add(NewContact(NewKademliaID("0000000000000000000000000000000000000001"), ""))
		cl.Add(NewContact(NewKademliaID("0000000000000000000000000000000000000002"), ""))
		cl.Add(wamtedContact)

		//Act
		actual := cl.Get(NewKademliaID(wantedId))
		if actual == nil {
			t.Errorf("Expected a candidate, got %v", actual)
		}

		//Assert
		if !actual.Contact.ID.Equals(wamtedContact.ID) {
			t.Errorf("Expected targetid %v, got %v", wamtedContact.String(), actual.Contact.String())
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

func TestAddWhenListIsFullAndNotReplace(t *testing.T) { // Disabled
	testname := "Add candidate when list is full and no replace"
	t.Run(testname, func(t *testing.T) {
		t.Skip()

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

func TestCheckedStatusWhenAddDuplicate(t *testing.T) {
	testname := "Preserve cghecked status when adding duplicate"
	t.Run(testname, func(t *testing.T) {
		//Arrange
		expected := true
		var contacts []Contact
		targetid := NewRandomKademliaID()
		contacts = append(contacts, NewContact(NewRandomKademliaID(), ""))
		contacts = append(contacts, NewContact(NewRandomKademliaID(), ""))

		cl := NewCandidateList(targetid, 8)
		cl.AddMultiple(contacts)

		contactToAdd := NewContact(NewRandomKademliaID(), "")

		//Act
		cl.Add(contactToAdd)
		cl.Check(contactToAdd.ID)
		cl.Add(contactToAdd)

		actual := cl.Get(contactToAdd.ID).Checked

		//Assert
		if actual != expected {
			t.Errorf("Expected contact %v, got %v", true, actual)
		}
	})
}
