package kademlia

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	// RPC
	MESSAGE_RPC_PING       = 1
	MESSAGE_RPC_STORE      = 2
	MESSAGE_RPC_FIND_NODE  = 3
	MESSAGE_RPC_FIND_VALUE = 4

	// RPC response
	MESSAGE_PONG      = 5
	MESSAGE_NODE_LIST = 6
	MESSAGE_VALUE     = 7
)

const (
	NETWORK_PORT            = 14041
	NETWORK_INCOMING_BUFFER = 8192
	NETWORK_REQUEST_TIMEOUT = 2 * time.Second
)

type Network struct {
	Kademlia       *Kademlia
	incomingData   chan []byte
	outgoingMsg    chan NetworkMessage
	waiters        map[int64](chan NetworkMessage)
	messageCounter *Counter
}

type NetworkMessage struct {
	ID         int64
	RPC        int
	Sender     *Contact
	Target     *Contact
	BodyDigest []byte
	Body       string
}

func CreateNewNetwork(kademlia *Kademlia) Network {
	net := Network{
		Kademlia:       kademlia,
		incomingData:   make(chan []byte),
		outgoingMsg:    make(chan NetworkMessage),
		waiters:        make(map[int64]chan NetworkMessage),
		messageCounter: MakeCounter(),
	}
	return net
}

// Listen for incoming UDP network messages.
func (network *Network) Listen() {
	addr := net.UDPAddr{
		Port: NETWORK_PORT,
		IP:   net.ParseIP(network.Kademlia.Me.Address),
	}

	socket, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer socket.Close()
	go network.incomingDataHandler()

	var buf [NETWORK_INCOMING_BUFFER]byte
	for {
		len, _, udpError := socket.ReadFromUDP(buf[:])
		if udpError != nil {
			log.Println(udpError)
		}
		network.incomingData <- buf[:len]
	}
}

func (network *Network) SendPingMessage(contact *Contact, alive chan bool) {
	msg := NetworkMessage{
		ID:     network.messageCounter.GetNext(),
		RPC:    MESSAGE_RPC_PING,
		Sender: network.Kademlia.Me,
		Target: contact,
	}

	network.waiters[msg.ID] = make(chan NetworkMessage, 1)

	defer delete(network.waiters, msg.ID)
	go sendMessage(msg)

	select {
	case <-network.waiters[msg.ID]:
		alive <- true
	case <-time.After(NETWORK_REQUEST_TIMEOUT):
		alive <- false
	}
}

// Take actions on a network message
func (network *Network) incomingDataHandler() {
	for data := range network.incomingData {
		msg, err := deserializeMessage(data)
		if err != nil {
			log.Printf("Deserialize error: %s\n", err)
			// continue
		}

		log.Printf("Message (%d) from %s\n", msg.RPC, msg.Sender.String())

		network.messageHandler(msg)
	}
}

func (network *Network) messageHandler(msg *NetworkMessage) {
	network.Kademlia.Routing.AddContact(*msg.Sender)

	if c, ok := network.waiters[msg.ID]; ok {
		c <- *msg
		delete(network.waiters, msg.ID)
		return
	}

	switch msg.RPC {
	case MESSAGE_RPC_PING:
		msg.RPC = MESSAGE_PONG
		flipSenderTarget(msg)
		sendMessage(*msg)
	case MESSAGE_RPC_STORE:
		// TODO: Implement RPC
	case MESSAGE_RPC_FIND_NODE:
		// TODO: Implement RPC
	case MESSAGE_RPC_FIND_VALUE:
		// TODO: Implement RPC
	case MESSAGE_PONG:
	case MESSAGE_NODE_LIST:
		// TODO: Implement message response
	case MESSAGE_VALUE:
		// TODO: Implement message response
	}
}

// Get IP-address of this computer
func GetOutboundIP() string {
	// https://stackoverflow.com/a/37382208
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
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
func flipSenderTarget(msg *NetworkMessage) {
	msg.Sender, msg.Target = msg.Target, msg.Sender
}

// Send network message
func sendMessage(msg NetworkMessage) {
	address := fmt.Sprintf("%s:%d", msg.Target.Address, NETWORK_PORT)
	conn, err := net.Dial("udp", address)
	if err != nil {
		log.Printf("UDP send error: %v", err)
		return
	}
	bytes := serializeMessage(msg)
	conn.Write(bytes)
}
