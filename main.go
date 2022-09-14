package main

import (
	"d7024e/cli"
	"d7024e/kademlia"
)

func main() {

	me := kademlia.NewContact(kademlia.NewRandomKademliaID(), kademlia.GetOutboundIP())
	context := kademlia.Kademlia{Me: &me}
	network := kademlia.Network{Kademlia: &context}

	context.Network = &network

	go network.Listen() // TODO: Notify it is actually listening
	go context.LookupContact(&me)
	cli.Open(&context)
}
