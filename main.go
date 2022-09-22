package main

import (
	"d7024e/cli"
	"d7024e/kademlia"
	"flag"
	"fmt"
	"io"
	"log"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func main() {
	port, _, verbose := parseCommandlineFlags()
	if !*verbose {
		log.SetOutput(io.Discard)
	}

	myAddress := fmt.Sprintf("%s:%d", kademlia.GetOutboundIP(), *port)
	me := kademlia.NewContact(kademlia.NewRandomKademliaID(), myAddress)
	context := kademlia.Kademlia{Me: &me, Routing: kademlia.NewRoutingTable(me)}
	network := kademlia.CreateNewNetwork(&context, *port)

	context.Network = &network

	go network.Listen() // TODO: Notify it is actually listening
	go context.LookupContact(*me.ID)
	cli.Open(&context)
}

func parseCommandlineFlags() (port *int, bootstrapNode *bool, verbose *bool) {
	port = flag.Int("p", 14041, "Portnumber")
	bootstrapNode = flag.Bool("b", false, "Indicates whether the node is a bootstrap node")
	verbose = flag.Bool("v", false, "Indicates if a log should be created")

	flag.Parse()

	return port, bootstrapNode, verbose
}
