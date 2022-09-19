package main

import (
	"d7024e/cli"
	"d7024e/kademlia"
	"log"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)

	me := kademlia.NewContact(kademlia.NewRandomKademliaID(), kademlia.GetOutboundIP())
	context := kademlia.Kademlia{Me: &me, Routing: kademlia.NewRoutingTable(me)}
	network := kademlia.CreateNewNetwork(&context)

	context.Network = &network

	go network.Listen() // TODO: Notify it is actually listening
	go context.LookupContact(&me)
	cli.Open(&context)
}
