package routing

const bucketSize = 20

type IRoutingTable interface {
	// AddContact add a new contact to the correct Bucket. Contact will not be added if it is me
	AddContact(contact Contact)

	// RemoveContact removes a contact from the correct Bucket if it exists
	RemoveContact(contactId *KademliaID)

	// FindClosestContacts finds the count closest Contacts to the target in the RoutingTable
	FindClosestContacts(target *KademliaID, count int) []Contact

	// Get number of nodes in the RoutingTable
	GetNumberOfNodes() int

	// Get all nodes in the RoutingTable
	Nodes() []Contact
}

// RoutingTable definition
// keeps a refrence contact of me and an array of buckets
type RoutingTable struct {
	me      Contact
	buckets [IDLength * 8]*bucket
}

// NewRoutingTable returns a new instance of a RoutingTable
func NewRoutingTable(me Contact) *RoutingTable {
	routingTable := &RoutingTable{}
	for i := 0; i < IDLength*8; i++ {
		routingTable.buckets[i] = newBucket()
	}
	routingTable.me = me
	return routingTable
}

func (routingTable *RoutingTable) AddContact(contact Contact) {
	if routingTable.me.ID.Equals(contact.ID) {
		return
	}
	bucketIndex := routingTable.getBucketIndex(contact.ID)
	bucket := routingTable.buckets[bucketIndex]
	bucket.AddContact(contact)
}

func (routingTable *RoutingTable) RemoveContact(contactId *KademliaID) {
	bucketIndex := routingTable.getBucketIndex(contactId)
	bucket := routingTable.buckets[bucketIndex]
	bucket.RemoveContact(contactId)
}

func (routingTable *RoutingTable) FindClosestContacts(target *KademliaID, count int) []Contact {
	var candidates ContactCandidates
	bucketIndex := routingTable.getBucketIndex(target)
	bucket := routingTable.buckets[bucketIndex]

	candidates.Append(bucket.GetContactAndCalcDistance(target))

	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < IDLength*8) && candidates.Len() < count; i++ {
		if bucketIndex-i >= 0 {
			bucket = routingTable.buckets[bucketIndex-i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
		if bucketIndex+i < IDLength*8 {
			bucket = routingTable.buckets[bucketIndex+i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
	}

	candidates.Sort()

	if count > candidates.Len() {
		count = candidates.Len()
	}

	return candidates.GetContacts(count)
}

// getBucketIndex get the correct Bucket index for the KademliaID
func (routingTable *RoutingTable) getBucketIndex(id *KademliaID) int {
	distance := id.CalcDistance(routingTable.me.ID)
	for i := 0; i < IDLength; i++ {
		for j := 0; j < 8; j++ {
			if (distance[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}

	return IDLength*8 - 1
}

func (routingTable *RoutingTable) GetNumberOfNodes() int {
	nodes := 0
	for _, bucket := range routingTable.buckets {
		nodes += bucket.Len()
	}
	return nodes
}

func (routingTable *RoutingTable) Nodes() []Contact {
	var contacts []Contact
	for _, bucket := range routingTable.buckets {
		for e := bucket.list.Front(); e != nil; e = e.Next() {
			x := e.Value.(Contact)
			contacts = append(contacts, x)
		}
	}
	return contacts
}
