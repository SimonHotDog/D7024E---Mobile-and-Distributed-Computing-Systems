package kademlia

import (
	"d7024e/kademlia/network"
	"d7024e/kademlia/network/routing"
	"d7024e/kademlia/network/rpc"
	"d7024e/util"
	"errors"
	"log"
	"sync"

	cmap "github.com/orcaman/concurrent-map/v2"
)

type IKademlia interface {
	// Get me
	GetMe() *routing.Contact
	GetNetwork() network.INetwork
	GetDataStore() *cmap.ConcurrentMap[[]byte]

	LookupContact(targetID *routing.KademliaID) []routing.Contact
	LookupData(hash string) ([]byte, *routing.Contact)
	Store(data []byte) (string, error)
	JoinNetwork(contact *routing.Contact)
}

type Kademlia struct {
	me        *routing.Contact
	network   network.INetwork
	dataStore *cmap.ConcurrentMap[[]byte]
}

// Hyperparameters
const K int = 20 //k closest
const A int = 3  //alpha, 1 is effectively no concurrency

func NewKademlia(me *routing.Contact, network network.INetwork, datastore *cmap.ConcurrentMap[[]byte]) *Kademlia {
	return &Kademlia{me, network, datastore}
}

// Getters
func (kademlia *Kademlia) GetMe() *routing.Contact                   { return kademlia.me }
func (kademlia *Kademlia) GetNetwork() network.INetwork              { return kademlia.network }
func (kademlia *Kademlia) GetDataStore() *cmap.ConcurrentMap[[]byte] { return kademlia.dataStore }

// Lookup contacts
func (kademlia *Kademlia) LookupContact(targetID *routing.KademliaID) []routing.Contact {
	candidateList := NewCandidateList(targetID, K)
	kClosestContacts := kademlia.network.GetRoutingTable().FindClosestContacts(targetID, K)

	kademlia.lookupContactAux(targetID, kClosestContacts, candidateList)

	contacts := make([]routing.Contact, candidateList.Len())
	for i, candidate := range candidateList.GetAll() {
		contacts[i] = candidate.Contact
	}

	return contacts
}

func (kademlia *Kademlia) lookupContactAux(targetID *routing.KademliaID, contacts []routing.Contact, cl *CandidateList) {
	var wg sync.WaitGroup

	for i, contact := range contacts {
		if i > A {
			break
		}
		wg.Add(1)
		go func(contact routing.Contact, targetId *routing.KademliaID, cl *CandidateList, wg *sync.WaitGroup) {
			defer wg.Done()

			candidate := cl.Get(contact.ID)
			if candidate != nil && candidate.Checked {
				// Already checked
				return
			} else {
				cl.Add(contact)
			}

			channel := make(chan []routing.Contact, 1)
			go rpc.SendFindContactMessage(kademlia.network, &contact, targetID, channel)
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

// send lookup message to closest nodes
func (kademlia *Kademlia) LookupData(hash string) ([]byte, *routing.Contact) {
	kademliaIdFromHash := routing.NewKademliaID(hash)
	contacts := kademlia.LookupContact(kademliaIdFromHash)

	for _, contact := range contacts { // for each of the <=5 contacts found...
		//fmt.Println(" trying to find data on node ", contact.ID)

		//kademlia.Network.SendLookup(&contact, hash) //send FindLocally to each

		valuechannel := make(chan string, 1)
		go rpc.SendLookupMessage(kademlia.network, &contact, hash, valuechannel) //send FindLocally to each
		value := <-valuechannel
		//fmt.Printf("Recieved value %v from node %v", value, contact.String())
		if value != "" {
			return []byte(value), &contact
		}
	}

	return nil, nil
}

// send store message to closest nodes
func (kademlia *Kademlia) Store(data []byte) (string, error) {
	hashed := util.Hash(data)
	stringToByte := []byte(hashed)

	contacts := kademlia.LookupContact((*routing.KademliaID)(stringToByte))

	if len(contacts) == 0 {
		err := errors.New("no suitable contacts found for storage")
		return "", err
	} else {
		for _, contact := range contacts { // for each of the <=5 contacts found...
			log.Printf("Storing message with hash %s at node %s\n", hashed, contact.String())
			// TODO: Make this concurrent
			ok := rpc.SendStoreMessage(kademlia.network, &contact, hashed, data) //send StoreLocally to each
			if !ok {
				log.Println("Could not store message at node " + contact.String())
			}
		}
	}
	return hashed, nil

}

// Join a kademlia network by through a known node
func (kademlia *Kademlia) JoinNetwork(knownNode *routing.Contact) {
	log.Printf("Joining network via %v...", knownNode)
	repononseChannel := make(chan []routing.Contact)
	go rpc.SendFindContactMessage(kademlia.network, knownNode, kademlia.me.ID, repononseChannel)

	// Ping all recieved contacts and add them to routing-table if they respond
	contacts := <-repononseChannel
	var wg sync.WaitGroup
	c := util.MakeCounter()
	for _, contact := range contacts {
		wg.Add(1)
		go func(contact routing.Contact) {
			aliveChannel := make(chan bool)
			go rpc.SendPingMessage(kademlia.network, &contact, aliveChannel)
			if <-aliveChannel {
				kademlia.network.GetRoutingTable().AddContact(contact)
				c.Increase()
			}
			wg.Done()
		}(contact)
	}

	wg.Wait()
	log.Printf("Joined network and recieved %d (%d alive) nodes close to me from %v\n", len(contacts), c.GetNext()-1, knownNode.Address)
}
