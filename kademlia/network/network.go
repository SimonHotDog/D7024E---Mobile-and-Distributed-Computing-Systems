package network

import (
	"d7024e/kademlia/network/routing"
	"d7024e/util"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
)

const (
	// RPC
	MESSAGE_RPC_PING       = 1
	MESSAGE_RPC_STORE      = 2
	MESSAGE_RPC_FIND_NODE  = 3
	MESSAGE_RPC_FIND_VALUE = 4

	// RPC response
	// SUGGESTION Replace specific response types with a general MESSAGE_RESPONSE
	MESSAGE_PONG      = 5
	MESSAGE_NODE_LIST = 6
	MESSAGE_VALUE     = 7
)

const (
	NETWORK_INCOMING_BUFFER        = 8192
	NETWORK_REQUEST_TIMEOUT        = 2 * time.Second
	NETWORK_REQUEST_TIMEOUT_STRING = "::timeout::"
)

type Network struct {
	Me           *routing.Contact
	Routingtable *routing.RoutingTable
	Datastore    *cmap.ConcurrentMap[[]byte]

	incomingData       chan []byte
	outgoingMsg        chan NetworkMessage
	waiters            cmap.ConcurrentMap[chan NetworkMessage]
	messageCounter     *util.Counter
	port               int
	quitListenSig      chan struct{}
	incomingDataSocket *net.UDPConn
}

type NetworkMessage struct {
	ID         int64
	RPC        int
	Sender     *routing.Contact
	Target     *routing.Contact
	BodyDigest string
	Body       string
	Contacts   []routing.Contact
}

// Create a new network instance.
//
// Parameters:
//
//	port: The port to listen on.
//	datastore: The datastore to use.
//
// Returns:
//
//	A new network instance and a contact that will be used when communicating
//	with other nodes.
func NewNetwork(port int, datastore *cmap.ConcurrentMap[[]byte]) (*Network, *routing.Contact) {
	myAddress := fmt.Sprintf("%s:%d", GetOutboundIP(), port)
	me := routing.NewContact(routing.NewRandomKademliaID(), myAddress)

	net := Network{
		Me:           &me,
		Routingtable: routing.NewRoutingTable(me),

		Datastore:      datastore,
		incomingData:   make(chan []byte),
		outgoingMsg:    make(chan NetworkMessage),
		messageCounter: util.MakeCounter(),
		port:           port,
		quitListenSig:  make(chan struct{}, 1),
		waiters:        cmap.New[chan NetworkMessage](),
	}
	go net.messageSender()
	return &net, &me
}

func (network *Network) NewNetworkMessage(
	rpc int,
	sender *routing.Contact,
	target *routing.Contact,
	bodyDigest string,
	body string,
	contacts []routing.Contact,
) *NetworkMessage {
	return &NetworkMessage{
		ID:         network.messageCounter.GetNext(),
		RPC:        rpc,
		Sender:     sender,
		Target:     target,
		BodyDigest: bodyDigest,
		Body:       body,
		Contacts:   contacts,
	}
}

// Listen for incoming UDP network messages.
func (network *Network) Listen() {
	addr := net.UDPAddr{
		Port: network.port,
		IP:   net.ParseIP(network.Me.Address),
	}

	socket, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatal(err)
		return
	}
	network.incomingDataSocket = socket

	defer socket.Close()
	go network.incomingDataHandler()

	for {
		buf := make([]byte, NETWORK_INCOMING_BUFFER)
		select {
		case <-network.quitListenSig:
			return
		default:
		}
		// len, _, udpError := socket.ReadFromUDP(buf)
		len, udpError := socket.Read(buf)
		if udpError != nil {
			log.Println(udpError)
		}
		network.incomingData <- buf[:len]
	}
}

func (network *Network) StopListen() {
	network.quitListenSig <- struct{}{}
	network.incomingDataSocket.Close()
}

// Send network message and wait on response.
//
// If the contact responds, the response will be returned and `timeout` be false.
//
// Otherwise, after time to respond exceeds `network.NETWORK_REQUEST_TIMEOUT`,
// timeout occured and `timeout` will be true.
//
// Parameters:
//
//	msg: The message to send
func (network *Network) SendMessageWithResponse(msg NetworkMessage) (response NetworkMessage, timeout bool) {
	network.waiters.Set(fmt.Sprint(msg.ID), make(chan NetworkMessage, 1))
	defer network.waiters.Remove(fmt.Sprint(msg.ID))

	network.outgoingMsg <- msg

	waitchannel, _ := network.waiters.Get(fmt.Sprint(msg.ID))
	select {
	case response := <-waitchannel:
		return response, false
	case <-time.After(NETWORK_REQUEST_TIMEOUT):
		return NetworkMessage{}, true
	}
}

// Send network message and don't wait for response
// Parameters:
//
//	msg: The message to send
func (network *Network) SendMessage(msg NetworkMessage) {
	network.outgoingMsg <- msg
}

// Process incoming network data in the 'network.incomingData' channel
func (network *Network) incomingDataHandler() {
	for data := range network.incomingData {
		msg, err := deserializeMessage(data)
		if err != nil {
			log.Printf("Deserialize error: %s\n", err)
			continue
		}

		log.Printf("Message (%d) from %s\n", msg.RPC, msg.Sender.String())

		network.messageHandler(msg)
	}
}

// Take actions on a network message
func (network *Network) messageHandler(msg *NetworkMessage) {
	network.Routingtable.AddContact(*msg.Sender)

	if c, ok := network.waiters.Get(fmt.Sprint(msg.ID)); ok {
		c <- *msg
		network.waiters.Remove(fmt.Sprint(msg.ID))
		return
	}

	switch msg.RPC {
	case MESSAGE_RPC_PING:
		msg.RPC = MESSAGE_PONG
		network.generateReturnMessage(msg)
		network.SendMessage(*msg)
	case MESSAGE_RPC_STORE:
		network.Datastore.Set(msg.BodyDigest, []byte(msg.Body))
		_, ok := network.Datastore.Get(msg.BodyDigest)
		if ok {
			msg.Body = "1"
		} else {
			msg.Body = "0"
		}
		msg.RPC = MESSAGE_VALUE
		network.generateReturnMessage(msg)
		network.SendMessage(*msg)
	case MESSAGE_RPC_FIND_NODE:
		contactId := routing.NewKademliaID(msg.Body)
		nodes := network.Routingtable.FindClosestContacts(contactId, 20) // TODO: Get count value (20) from some parameter

		msg.RPC = MESSAGE_NODE_LIST
		msg.Contacts = nodes
		network.generateReturnMessage(msg)
		network.SendMessage(*msg)
	case MESSAGE_RPC_FIND_VALUE:
		value, exists := network.Datastore.Get(msg.BodyDigest)
		if exists {
			log.Printf("Data: %v (%s) found on node %s\n", value, string(value), msg.Target.String())
		} else {
			log.Printf("Data not found on node %s\n", msg.Target.String())
		}

		msg.RPC = MESSAGE_VALUE
		msg.Body = string(value)
		network.generateReturnMessage(msg)
		network.SendMessage(*msg)
	case MESSAGE_PONG:
	case MESSAGE_NODE_LIST:
	case MESSAGE_VALUE:
		// TODO: Implement message response
	}
}

// Deserialize a byte array to a networkMessage
func deserializeMessage(data []byte) (*NetworkMessage, error) {
	var msg NetworkMessage
	encodeError := json.Unmarshal(data, &msg)
	if encodeError != nil {
		return nil, encodeError
	}
	return &msg, nil
}

// Serialize a networkMessage to a byte array
func serializeMessage(msg NetworkMessage) []byte {
	bytes, _ := json.Marshal(msg)
	return bytes
}

// Flip sender and target in a network message
func (network *Network) generateReturnMessage(msg *NetworkMessage) {
	returnContact := *msg.Sender
	msg.Target = &returnContact
	msg.Sender = network.Me
}

func (network *Network) messageSender() {
	for msg := range network.outgoingMsg {
		conn, err := net.Dial("udp", msg.Target.Address)
		if err != nil {
			log.Printf("UDP send error: %v", err)
			continue
		}
		bytes := serializeMessage(msg)
		conn.Write(bytes)
		conn.Close()
		log.Printf("Message (%d) sent to %s\n", msg.RPC, msg.Target.String())
	}
}
