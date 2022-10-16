package network

import (
	"d7024e/kademlia/datastore"
	"d7024e/kademlia/network/routing"
	"d7024e/util"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

const (
	// RPC
	MESSAGE_RPC_PING         = 1
	MESSAGE_RPC_STORE        = 2
	MESSAGE_RPC_FIND_NODE    = 3
	MESSAGE_RPC_FIND_VALUE   = 4
	MESSAGE_RPC_DATA_REFRESH = 5
	MESSAGE_RPC_DATA_FORGET  = 6

	// RPC response
	MESSAGE_RESPONSE = 10
)

const (
	NETWORK_INCOMING_BUFFER        = 8192
	NETWORK_REQUEST_TIMEOUT        = 2 * time.Second
	NETWORK_REQUEST_TIMEOUT_STRING = "::timeout::"
)

type INetwork interface {
	GetMe() *routing.Contact
	GetRoutingTable() routing.IRoutingTable
	GetDatastore() datastore.IDataStore

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
	NewNetworkMessage(
		rpc int,
		sender *routing.Contact,
		target *routing.Contact,
		bodyDigest string,
		body string,
		contacts []routing.Contact,
	) *NetworkMessage

	// Listen for incoming UDP network messages.
	Listen()

	// Stop listening for incoming UDP network messages.
	StopListen()

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
	SendMessageWithResponse(msg NetworkMessage) (response NetworkMessage, timeout bool)

	// Send network message and don't wait for response
	// Parameters:
	//
	//	msg: The message to send
	SendMessage(msg NetworkMessage)
}

type Network struct {
	me           *routing.Contact
	routingtable routing.IRoutingTable
	datastore    datastore.IDataStore

	incomingData       chan []byte
	messageCounter     *util.Counter
	port               int
	quitListenSig      chan struct{}
	incomingDataLock   sync.Mutex
	incomingDataSocket *net.UDPConn
}

type NetworkMessage struct {
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
func NewNetwork(port int, datastore datastore.IDataStore) (*Network, *routing.Contact) {
	myAddress := fmt.Sprintf("%s:%d", GetOutboundIP(), port)
	me := routing.NewContact(routing.NewRandomKademliaID(), myAddress)

	net := Network{
		me:             &me,
		routingtable:   routing.NewRoutingTable(me),
		datastore:      datastore,
		incomingData:   make(chan []byte),
		messageCounter: util.MakeCounter(),
		port:           port,
		quitListenSig:  make(chan struct{}, 1),
	}
	return &net, &me
}

func (network *Network) GetMe() *routing.Contact {
	return network.me
}

func (network *Network) GetRoutingTable() routing.IRoutingTable {
	return network.routingtable
}

func (network *Network) GetDatastore() datastore.IDataStore {
	return network.datastore
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
		RPC:        rpc,
		Sender:     sender,
		Target:     target,
		BodyDigest: bodyDigest,
		Body:       body,
		Contacts:   contacts,
	}
}

func (network *Network) Listen() {
	addr := net.UDPAddr{
		Port: network.port,
		IP:   net.ParseIP(network.me.Address),
	}

	socket, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatal(err)
		return
	}

	network.incomingDataLock.Lock()
	network.incomingDataSocket = socket
	network.incomingDataLock.Unlock()

	defer socket.Close()

	for {
		buf := make([]byte, NETWORK_INCOMING_BUFFER)
		select {
		case <-network.quitListenSig:
			return
		default:
		}
		len, remote, udpError := socket.ReadFromUDP(buf)
		if udpError != nil {
			log.Println(udpError)
		}
		go network.incomingDataHandler(remote, buf[:len])
	}
}

func (network *Network) StopListen() {
	network.quitListenSig <- struct{}{}
	network.incomingDataLock.Lock()
	network.incomingDataSocket.Close()
	network.incomingDataLock.Unlock()
}

func (network *Network) SendMessageWithResponse(msg NetworkMessage) (response NetworkMessage, timeout bool) {
	udpAddr, _ := net.ResolveUDPAddr("udp", msg.Target.Address)
	res, err := network.sendRequest(udpAddr, msg, true)
	if err == nil {
		return *res, false
	}
	return *new(NetworkMessage), true
}

func (network *Network) SendMessage(msg NetworkMessage) {
	udpAddr, _ := net.ResolveUDPAddr("udp", msg.Target.Address)
	network.sendRequest(udpAddr, msg, false)
}

// Process incoming network data
func (network *Network) incomingDataHandler(senderAddr *net.UDPAddr, data []byte) {
	msg, err := deserializeMessage(data)
	if err != nil {
		log.Printf("Deserialize error: %s\n", err)
		return
	}

	log.Printf("Message (%d) from %s\n", msg.RPC, msg.Sender.String())

	network.messageHandler(senderAddr, msg)
}

// Take actions on a network message
func (network *Network) messageHandler(senderAddr *net.UDPAddr, msg *NetworkMessage) {
	network.routingtable.AddContact(*msg.Sender)

	switch msg.RPC {
	case MESSAGE_RPC_PING:
		network.generateReturnMessage(msg)
		network.sendResponse(senderAddr, *msg)
	case MESSAGE_RPC_STORE:
		ok := network.datastore.Set(msg.BodyDigest, []byte(msg.Body))

		msg.Body = strconv.FormatBool(ok)
		network.generateReturnMessage(msg)
		network.sendResponse(senderAddr, *msg)
	case MESSAGE_RPC_FIND_NODE:
		contactId := routing.NewKademliaID(msg.Body)
		nodes := network.routingtable.FindClosestContacts(contactId, 20) // TODO: Get count value (20) from some parameter

		msg.Contacts = nodes
		network.generateReturnMessage(msg)
		network.sendResponse(senderAddr, *msg)
	case MESSAGE_RPC_FIND_VALUE:
		value, exists := network.datastore.Get(msg.BodyDigest)
		if exists {
			log.Printf("Data: %v (%s) found on node %s\n", value, string(value), msg.Target.String())
		} else {
			log.Printf("Data not found on node %s\n", msg.Target.String())
		}

		msg.Body = string(value)
		network.generateReturnMessage(msg)
		network.sendResponse(senderAddr, *msg)
	case MESSAGE_RPC_DATA_REFRESH:
		keyToRefresh := msg.BodyDigest
		network.datastore.Refresh(keyToRefresh)
	case MESSAGE_RPC_DATA_FORGET:
		keyToForget := msg.BodyDigest
		network.datastore.Remove(keyToForget)
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
	msg.Sender = network.me
	msg.RPC = MESSAGE_RESPONSE
}

func (network *Network) sendResponse(addr *net.UDPAddr, msg NetworkMessage) {
	msg_bytes := serializeMessage(msg)
	_, err := network.incomingDataSocket.WriteToUDP(msg_bytes, addr)
	if err != nil {
		log.Printf("Send response error: %v\n", err)
	}
}

func (network *Network) sendRequest(recipient *net.UDPAddr, msg NetworkMessage, waitResponse bool) (*NetworkMessage, error) {
	log.Printf("Message sent to %s\n", recipient.String())
	conn, err := net.Dial("udp", recipient.String())
	if err != nil {
		log.Printf("UDP connection error: %v", err)
		return nil, err
	}
	defer conn.Close()

	// Send message
	bytes := serializeMessage(msg)
	_, err = conn.Write(bytes)
	if err != nil {
		log.Printf("Send request error: %v\n", err)
		return nil, err
	}
	if !waitResponse {
		return nil, nil
	}

	// Wait for response
	conn.SetReadDeadline(time.Now().Add(NETWORK_REQUEST_TIMEOUT))
	response_buffer := make([]byte, NETWORK_INCOMING_BUFFER)
	len, err := conn.Read(response_buffer)
	if err != nil {
		log.Printf("UDP read error: %v\n", err)
		return nil, err
	}

	// Deserialize response
	response, err := deserializeMessage(response_buffer[:len])
	if err != nil {
		log.Printf("Deserialize error: %s\n", err)
		return nil, err
	}

	return response, nil
}
