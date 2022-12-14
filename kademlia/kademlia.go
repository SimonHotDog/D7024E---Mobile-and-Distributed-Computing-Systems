package kademlia

import (
	"d7024e/kademlia/datastore"
	"d7024e/kademlia/network"
	"d7024e/kademlia/network/routing"
	"d7024e/kademlia/network/rpc"
	"d7024e/util"
	"errors"
	"log"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type IKademlia interface {
	// Get me
	GetMe() *routing.Contact
	GetNetwork() network.INetwork
	GetDataStore() datastore.IDataStore

	LookupContact(targetID *routing.KademliaID) []routing.Contact
	LookupData(hash string) ([]byte, *routing.Contact)
	Store(data []byte) (string, error)
	ForgetData(hash string, contacts []routing.Contact) error
	JoinNetwork(contact *routing.Contact, retries int) bool
}

type Kademlia struct {
	me        *routing.Contact
	network   network.INetwork
	dataStore datastore.IDataStore
}

// Hyperparameters
const K int = 20 //k closest
const A int = 3  //alpha, 1 is effectively no concurrency

func NewKademlia(me *routing.Contact, network network.INetwork, datastore datastore.IDataStore) *Kademlia {
	return &Kademlia{me, network, datastore}
}

// Getters
func (kademlia *Kademlia) GetMe() *routing.Contact            { return kademlia.me }
func (kademlia *Kademlia) GetNetwork() network.INetwork       { return kademlia.network }
func (kademlia *Kademlia) GetDataStore() datastore.IDataStore { return kademlia.dataStore }

// Lookup contacts
func (kademlia *Kademlia) LookupContact(targetID *routing.KademliaID) []routing.Contact {
	candidateList := NewCandidateList(targetID, K)
	kClosestContacts := kademlia.network.GetRoutingTable().FindClosestContacts(targetID, K)

	candidateList.AddMultiple(kClosestContacts)
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
	if kademliaIdFromHash == nil {
		return nil, nil
	}
	contacts := kademlia.LookupContact(kademliaIdFromHash)
	noResponseArray := []*routing.Contact{}

	for _, contact := range contacts {
		go rpc.SendRefreshDataMessage(kademlia.network, &contact, hash)
	}

	for _, contact := range contacts {
		valuechannel := make(chan string, 1)
		go rpc.SendLookupMessage(kademlia.network, &contact, hash, valuechannel) //send FindLocally to each
		value := <-valuechannel
		if value == "" {
			noResponseArray = append(noResponseArray, &contact)
		} else {
			if len(noResponseArray) > 0 {
				lastContact := noResponseArray[len(noResponseArray)-1]
				log.Printf("Storing at last empty contact with ID:  %v \n", lastContact.ID)
				go rpc.SendStoreMessage(kademlia.network, lastContact, hash, []byte(value))
			}

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

// Send forget message to specified contacts
func (kademlia *Kademlia) ForgetData(hash string, contacts []routing.Contact) error {
	kademliaIdFromHash := routing.NewKademliaID(hash)
	if kademliaIdFromHash == nil {
		return errors.New("invalid hash")
	}

	for _, contact := range contacts {
		rpc.SendForgetDataMessage(kademlia.network, &contact, hash)
	}

	return nil
}

// Join a kademlia network by through a known node
func (kademlia *Kademlia) JoinNetwork(knownNode *routing.Contact, retries int) bool {
	log.Printf("Joining network via %v...", knownNode)

	contacts, deadContacts := kademlia.joinNetworkAux(knownNode, 0, retries)

	if contacts == 0 {
		log.Printf("Failed to join network, no contacts received")
		return false
	} else if contacts != 0 && contacts == deadContacts {
		log.Printf("Failed to join network, no contacts responded in time")
		return false
	}
	log.Printf("Succesfully joined network, recieved %d (%d dead) nodes from %v\n", contacts, deadContacts, knownNode.Address)
	return true
}

func (kademlia *Kademlia) joinNetworkAux(knownNode *routing.Contact, numberOfRetries int, maxRestries int) (numberOfContacts, deadContacts int) {
	// Limit number of attempts to join network
	if numberOfRetries > maxRestries {
		return 0, 0
	}

	repononseChannel := make(chan []routing.Contact)
	go rpc.SendFindContactMessage(kademlia.network, knownNode, kademlia.me.ID, repononseChannel)

	// Ping all recieved contacts and add them to routing-table if they respond
	contacts := <-repononseChannel
	backoffTime := getExponentialBackoffTime(numberOfRetries)
	if len(contacts) == 0 {
		log.Printf("No contacts recieved from %v, trying again in %v\n", knownNode.Address, backoffTime)
		time.Sleep(backoffTime)
		return kademlia.joinNetworkAux(knownNode, numberOfRetries+1, maxRestries)
	}

	var deadNodes uint32
	var wg sync.WaitGroup
	for _, contact := range contacts {
		wg.Add(1)
		go func(contact routing.Contact) {
			aliveChannel := make(chan bool)
			go rpc.SendPingMessage(kademlia.network, &contact, aliveChannel)
			if <-aliveChannel {
				kademlia.network.GetRoutingTable().AddContact(contact)
			} else {
				atomic.AddUint32(&deadNodes, 1)
			}
			wg.Done()
		}(contact)
	}

	wg.Wait()

	// If all nodes are dead, try again
	if len(contacts) == int(deadNodes) {
		time.Sleep(backoffTime)
		return kademlia.joinNetworkAux(knownNode, numberOfRetries+1, maxRestries)
	}

	return len(contacts), int(deadNodes)
}

func getExponentialBackoffTime(attemptNumber int) time.Duration {
	// Inspiration from https://cloud.google.com/iot/docs/how-tos/exponential-backoff
	wait := int(math.Pow(2, float64(attemptNumber)))
	randomTime := rand.Intn(20)
	proposedBackoffTime := wait + randomTime
	maxWaitTime := 1000
	return time.Duration(min(proposedBackoffTime, maxWaitTime)) * time.Millisecond
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
