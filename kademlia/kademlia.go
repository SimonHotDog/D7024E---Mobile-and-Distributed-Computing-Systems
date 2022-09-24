package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"sync"
)

type Kademlia struct {
	Routing *RoutingTable
	Me      *Contact
	Network *Network
	Data    map[string][]byte
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) []byte {
	return kademlia.Data[hash]
}

func (kademlia *Kademlia) Store(data []byte) {
	if kademlia.Data == nil {
		kademlia.Data = make(map[string][]byte)
	}

	kademlia.Data[Hash(data)] = data
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

// Finds k closest nodes
func KClosest(key string) {
	// TODO
	// Implement a algorithm for finding
}
