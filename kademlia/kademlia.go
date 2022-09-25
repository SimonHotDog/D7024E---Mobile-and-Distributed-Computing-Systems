package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"log"
	"sync"
)

type Kademlia struct {
	Routing *RoutingTable
	Me      *Contact
	Network *Network
	Data    map[string][]byte
}

// Hyperparameters
const K int = 20 //k closest
const A int = 3  //alpha, 1 is effectively no concurrency

// Lookup contacts
func (kademlia *Kademlia) LookupContact(targetID *KademliaID) []Contact {
	candidateList := NewCandidateList(targetID, K)
	kClosestContacts := kademlia.Routing.FindClosestContacts(targetID, K)

	kademlia.lookupContactAux(targetID, kClosestContacts, candidateList)

	contacts := make([]Contact, candidateList.Len())
	for i, candidate := range candidateList.GetAll() {
		contacts[i] = candidate.Contact
	}

	return contacts
}

func (kademlia *Kademlia) lookupContactAux(targetID *KademliaID, contacts []Contact, cl *CandidateList) {
	var wg sync.WaitGroup

	for i, contact := range contacts {
		if i > A {
			break
		}
		wg.Add(1)
		go func(contact Contact, targetId *KademliaID, cl *CandidateList, wg *sync.WaitGroup) {
			defer wg.Done()

			candidate := cl.Get(contact.ID)
			if candidate != nil && candidate.Checked {
				// Already checked
				return
			}

			channel := make(chan []Contact, 1)
			go kademlia.Network.SendFindContactMessage(&contact, targetID, channel)
			contacts := <-channel
			cl.Check(contact.ID)

			if len(contacts) == 0 {
				// No contacts recieved
				return
			}
			if contact.ID.CalcDistance(targetId).Less(contacts[0].ID.CalcDistance(targetId)) {
				// Recieved contacts not closer than current
				return
			}

			cl.AddMultiple(contacts)
			kademlia.lookupContactAux(targetId, contacts, cl)
		}(contact, targetID, cl, &wg)
	}

	wg.Wait()
}

func (kademlia *Kademlia) LookupData(hash string) []byte {
	return kademlia.Data[hash]
}

// find send store message to closest nodes
func (kademlia *Kademlia) Store(data []byte) {

	hashed := Hash(data)
	stringToByte := []byte(hashed)

	contacts := kademlia.LookupContact((*KademliaID)(stringToByte))

	if len(contacts) == 0 {
		errors.New("No suitable contacts found for storage")
	} else {
		for contact := range contacts { // for each of the <=5 contacts found...
			log.Println("Storing at... ", contact.ID)
			go kademlia.Network.SendStore(&contact, targetID, channel) //send StoreLocally to each
		}
	}

}

// store data locally on node
func (kademlia *Kademlia) StoreLocally(data []byte) {

	hashed := Hash(data)

	if kademlia.Data == nil {
		kademlia.Data = make(map[string][]byte)
	}

	kademlia.Data[hashed] = data

}

// Join a kademlia network by through a known node
func (kademlia *Kademlia) JoinNetwork(knownNode *Contact) {
	log.Printf("Joining network via %v...", knownNode)
	repononseChannel := make(chan []Contact)
	go kademlia.Network.SendFindContactMessage(knownNode, kademlia.Me.ID, repononseChannel)

	// Ping all recieved contacts and add them to routing-table if they respond
	contacts := <-repononseChannel
	var wg sync.WaitGroup
	c := MakeCounter()
	for _, contact := range contacts {
		wg.Add(1)
		go func(contact Contact) {
			aliveChannel := make(chan bool)
			go kademlia.Network.SendPingMessage(&contact, aliveChannel)
			if <-aliveChannel {
				kademlia.Routing.AddContact(contact)
				c.Increase()
			}
			wg.Done()
		}(contact)
	}

	wg.Wait()
	log.Printf("Joined network and recieved %d (%d alive) nodes close to me from %v\n", len(contacts), c.GetNext()-1, knownNode.Address)
}

// Hashes data and returns key
func Hash(data []byte) string {
	sha1 := sha1.Sum([]byte(data))
	key := hex.EncodeToString(sha1[:])

	return key
}
