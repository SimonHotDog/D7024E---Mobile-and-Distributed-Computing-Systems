package kademlia

import (
	"d7024e/kademlia/network/routing"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCandidateList(t *testing.T) {
	testname := "Create new CandidateList"
	t.Run(testname, func(t *testing.T) {
		// Arrange
		contacts := make([]routing.Contact, 10)
		for i := 0; i < len(contacts); i++ {
			contacts[i] = routing.NewContact(routing.NewRandomKademliaID(), "")
		}
		targetid := routing.NewRandomKademliaID()

		// Act
		cl := NewCandidateList(targetid, 8)
		cl.AddMultiple(contacts)

		// Assert
		assert.Equal(t, cl.targetID, targetid)
	})
}

func TestCandidateExists(t *testing.T) {
	testname := "Check if candidate exists in candidatelist"
	t.Run(testname, func(t *testing.T) {
		//arrange
		contacts := make([]routing.Contact, 5)
		ids := make([]*routing.KademliaID, 5)

		for i := 0; i < len(contacts); i++ {
			x := routing.NewRandomKademliaID()
			contacts[i] = routing.NewContact(x, "")
			ids[i] = x
		}

		targetid := routing.NewRandomKademliaID()

		cl := NewCandidateList(targetid, 8)
		cl.AddMultiple(contacts)

		//Act

		actual := cl.Exists(ids[2])

		assert.True(t, actual)

		x := routing.NewRandomKademliaID()
		actual2 := cl.Exists(x)

		assert.False(t, actual2)
	})
}

func TestRemoveFromCandidateList(t *testing.T) {
	testname := "Remove from candidate list"
	t.Run(testname, func(t *testing.T) {
		//Arrange
		contacts := make([]routing.Contact, 6)
		ids := make([]*routing.KademliaID, 6)

		for i := 0; i < len(contacts); i++ {
			x := routing.NewRandomKademliaID()
			contacts[i] = routing.NewContact(x, "")
			ids[i] = x
		}

		targetid := routing.NewRandomKademliaID()

		cl := NewCandidateList(targetid, 8)
		cl.AddMultiple(contacts)

		//Act
		cl.Remove(ids[3])

		//Assert
		actual := cl.Exists(ids[3])

		assert.False(t, actual)
	})
}

func TestGetCandidateFromID(t *testing.T) {
	testname := "Gets Candidate from ID"
	t.Run(testname, func(t *testing.T) {
		//Arrange
		wantedId := "F000000000000000000000000000000000000000"
		wantedContact := routing.NewContact(routing.NewKademliaID(wantedId), "")
		targetid := routing.NewRandomKademliaID()
		cl := NewCandidateList(targetid, 8)
		cl.Add(routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000001"), ""))
		cl.Add(routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000002"), ""))
		cl.Add(wantedContact)

		//Act
		actual := cl.Get(routing.NewKademliaID(wantedId))

		//Assert
		assert.NotNil(t, actual)
		assert.Equal(t, wantedContact.ID, actual.Contact.ID)
	})
}

func TestGetWhenNotFound(t *testing.T) {
	//Arrange
	nonExistingId := "F000000000000000000000000000000000000000"
	targetid := routing.NewRandomKademliaID()
	cl := NewCandidateList(targetid, 8)
	cl.Add(routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000001"), ""))
	cl.Add(routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000002"), ""))

	//Act
	actual := cl.Get(routing.NewKademliaID(nonExistingId))

	//Assert
	assert.Nil(t, actual)
}

func TestAddWhenListIsNotFull(t *testing.T) {
	testname := "Add candidate when list is empty"
	t.Run(testname, func(t *testing.T) {
		//Arrange
		var contacts []routing.Contact
		targetid := routing.NewRandomKademliaID()
		contacts = append(contacts, routing.NewContact(routing.NewRandomKademliaID(), ""))
		contacts = append(contacts, routing.NewContact(routing.NewRandomKademliaID(), ""))

		cl := NewCandidateList(targetid, 8)
		cl.AddMultiple(contacts)

		contactToAdd := routing.NewContact(routing.NewRandomKademliaID(), "")

		//Act
		cl.Add(contactToAdd)
		actual := cl.Get(contactToAdd.ID)

		//Assert
		assert.NotNil(t, actual)
		assert.Equal(t, contactToAdd.ID, actual.Contact.ID)
	})
}

func TestAddWhenListIsFullAndReplace(t *testing.T) {
	testname := "Add candidate when list is full and no replace"
	t.Run(testname, func(t *testing.T) {
		//Arrange
		var contacts []routing.Contact
		targetid := routing.NewKademliaID("0000000000000000000000000000000000000000")

		contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000002"), ""))
		contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000003"), ""))
		contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000004"), ""))
		contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000005"), ""))
		contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000006"), ""))
		contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000007"), ""))
		contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000008"), ""))
		contacts = append(contacts, routing.NewContact(routing.NewKademliaID("00000000000000000000000000000000000000F0"), ""))

		cl := NewCandidateList(targetid, 8)
		cl.AddMultiple(contacts)

		contactToAdd := routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000001"), "")

		//Act
		cl.Add(contactToAdd)
		actual := cl.Get(contactToAdd.ID)

		//Assert
		assert.NotNil(t, actual)
		assert.Equal(t, contactToAdd.ID, actual.Contact.ID)
	})
}

func TestAddWhenListIsFullAndNotReplace(t *testing.T) { // Disabled
	testname := "Add candidate when list is full and no replace"
	t.Run(testname, func(t *testing.T) {
		t.Skip()

		//Arrange
		var contacts []routing.Contact
		targetid := routing.NewKademliaID("0000000000000000000000000000000000000000")

		contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000001"), ""))
		contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000002"), ""))
		contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000003"), ""))
		contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000004"), ""))
		contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000005"), ""))
		contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000006"), ""))
		contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000007"), ""))
		contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000008"), ""))

		cl := NewCandidateList(targetid, 8)
		cl.AddMultiple(contacts)

		contactToAdd := routing.NewContact(routing.NewKademliaID("000000000000000000000000000000000000000F"), "")

		//Act
		cl.Add(contactToAdd)
		actual := cl.Get(contactToAdd.ID)

		//Assert
		assert.Nil(t, actual)
	})
}

func TestCheckedStatusWhenAddDuplicate(t *testing.T) {
	testname := "Preserve cghecked status when adding duplicate"
	t.Run(testname, func(t *testing.T) {
		//Arrange
		expected := true
		var contacts []routing.Contact
		targetid := routing.NewRandomKademliaID()
		contacts = append(contacts, routing.NewContact(routing.NewRandomKademliaID(), ""))
		contacts = append(contacts, routing.NewContact(routing.NewRandomKademliaID(), ""))

		cl := NewCandidateList(targetid, 8)
		cl.AddMultiple(contacts)

		contactToAdd := routing.NewContact(routing.NewRandomKademliaID(), "")

		//Act
		cl.Add(contactToAdd)
		cl.Check(contactToAdd.ID)
		cl.Add(contactToAdd)

		actual := cl.Get(contactToAdd.ID).Checked

		//Assert
		assert.Equal(t, expected, actual)
	})
}

func TestLen(t *testing.T) {
	expected := 8
	var contacts []routing.Contact
	contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000001"), ""))
	contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000002"), ""))
	contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000003"), ""))
	contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000004"), ""))
	contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000005"), ""))
	contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000006"), ""))
	contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000007"), ""))
	contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000008"), ""))

	cl := NewCandidateList(routing.NewKademliaID("0000000000000000000000000000000000000000"), 8)
	cl.AddMultiple(contacts)

	actual := cl.Len()

	assert.Equal(t, expected, actual)
}

func TestGetAll(t *testing.T) {
	expected := 8
	var contacts []routing.Contact
	contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000001"), ""))
	contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000002"), ""))
	contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000003"), ""))
	contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000004"), ""))
	contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000005"), ""))
	contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000006"), ""))
	contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000007"), ""))
	contacts = append(contacts, routing.NewContact(routing.NewKademliaID("0000000000000000000000000000000000000008"), ""))

	cl := NewCandidateList(routing.NewKademliaID("0000000000000000000000000000000000000000"), 8)
	cl.AddMultiple(contacts)

	actual := cl.GetAll()

	assert.Equal(t, expected, len(actual))
}
