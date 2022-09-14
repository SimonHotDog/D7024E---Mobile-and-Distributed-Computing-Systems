package kademlia

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
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
)

type Network struct {
	Kademlia *Kademlia
}

type NetworkMessage struct {
	RPC        int
	Sender     *Contact
	Target     *Contact
	BodyDigest []byte
	Body       string
}

// Listen for incoming UDP network messages.
func (network *Network) Listen() {
	// https://stackoverflow.com/a/27176812/6474775

	// TODO: Do all packet processing in a go routine
	// src: https://stackoverflow.com/a/23576771/6474775

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

	var buf [NETWORK_INCOMING_BUFFER]byte
	for {
		len, _, udpError := socket.ReadFromUDP(buf[:])
		if udpError != nil {
			log.Println(udpError)
		}
		var msg NetworkMessage
		encodeError := json.Unmarshal(buf[:len], &msg)
		if encodeError != nil {
			log.Println(encodeError)
		}

		go network.networkMessageHandler(msg)
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	fmt.Println("Me:", network.Kademlia.Me.String())
	msg := NetworkMessage{
		RPC:    MESSAGE_RPC_PING,
		Sender: network.Kademlia.Me,
		Target: contact,
	}
	sendMessage(msg)
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
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

// Take actions on a network message
func (network *Network) networkMessageHandler(msg NetworkMessage) {
	log.Printf("Message (%d) from %s", msg.RPC, msg.Sender.String())

	//network.Kademlia.Routing.AddContact(*msg.Sender)

	switch msg.RPC {
	case MESSAGE_RPC_PING:
		msg.RPC = MESSAGE_PONG
		flipSenderTarget(&msg)
		sendMessage(msg)
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

func flipSenderTarget(msg *NetworkMessage) {
	msg.Sender, msg.Target = msg.Target, msg.Sender
}

func sendMessage(msg NetworkMessage) {
	address := fmt.Sprintf("%s:%d", msg.Target.Address, NETWORK_PORT)
	conn, err := net.Dial("udp", address)
	if err != nil {
		log.Printf("UDP send error: %v", err)
		return
	}
	bytes, _ := json.Marshal(msg)
	conn.Write(bytes)
}
